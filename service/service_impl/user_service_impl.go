package service_impl

import (
	"fmt"
	"time"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/helpers"
	"api.mijkomp.com/models/entity"
	"api.mijkomp.com/models/request"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	db             *gorm.DB
	Validation     *validator.Validate
}

func NewUserService(userRepostitory repository.UserRepository, validation *validator.Validate, db *gorm.DB) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository: userRepostitory,
		Validation:     validation,
		db:             db,
	}
}

func (service *UserServiceImpl) Create(request request.RegisterUserPayload) response.User {

	err := service.Validation.Struct(request)
	exception.PanicIfNeeded(err)

	tx := service.db.Begin()
	defer helpers.CommitOrRollback(tx)

	var count int64
	service.db.Find(&entity.User{}, "email = ?", request.Email).Count(&count)
	if count > 0 {
		exception.PanicIfNeeded(exception.NewValidationError(fmt.Sprintf("Email %s telah terdaftar.", request.Email)))
	}

	service.db.Find(&entity.User{}, "user_name = ?", request.UserName).Count(&count)
	if count > 0 {
		exception.PanicIfNeeded(exception.NewValidationError(fmt.Sprintf("Username %s telah terdaftar.", request.UserName)))
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	hashString := string(hashedPassword)
	userModel := entity.User{
		FullName:    request.FullName,
		Email:       request.Email,
		UserName:    request.UserName,
		Password:    &hashString,
		CreatedById: 0,
		CreatedAt:   time.Now().UTC(),
	}

	user, err := service.UserRepository.Save(tx, userModel)
	exception.PanicIfNeeded(err)

	// create user verification
	userVerificationModel := entity.UserVerification{
		UserId:       user.Id,
		Code:         helpers.NewRandom(6),
		AlreadyUsed:  false,
		CreatedById:  user.Id,
		CreatedAt:    time.Now().UTC(),
		ModifiedById: user.Id,
		ModifiedAt:   time.Now().UTC(),
	}

	err = service.UserRepository.CreateUserVerification(tx, userVerificationModel)
	exception.PanicIfNeeded(err)

	return service.mapUserResponse(user)
}

func (service *UserServiceImpl) GetById(userId uint) response.User {
	userRes, err := service.UserRepository.GetById(service.db, userId)
	exception.PanicIfNeeded(err)

	return service.mapUserResponse(userRes)
}

func (service *UserServiceImpl) GetByEmail(email string) response.User {
	res, err := service.UserRepository.GetByEmail(service.db, email)
	if err == gorm.ErrRecordNotFound {
		panic(exception.NewNotFoundError("Pengguna tidak ditemukan"))
	}

	return service.mapUserResponse(res)
}

func (service *UserServiceImpl) LoginUser(user request.LoginUserPayload) string {

	res, err := service.UserRepository.GetByEmail(service.db, user.Email)

	if err != nil {
		panic(exception.NewLoginError("Email atau password anda salah"))
	}

	//ferify password
	err = helpers.VerifyPassword(user.Password, *res.Password)

	if err != nil {
		panic(exception.NewLoginError("Email atau password anda salah"))
	}

	token_id := uuid.New()

	token, err := helpers.GenerateToken(fmt.Sprintf("%d", res.Id), token_id.String())

	userToken := entity.UserToken{
		Token:     token_id,
		UserId:    res.Id,
		IpAddress: "",
		CreatedAt: time.Now().UTC(),
	}

	service.UserRepository.SaveToken(service.db, userToken)
	exception.PanicIfNeeded(err)

	return token
}

func (service *UserServiceImpl) HasToken(userId uint, userToken uuid.UUID) bool {
	return service.UserRepository.HasToken(service.db, userId, userToken)
}

func (service *UserServiceImpl) mapUserResponse(res entity.User) response.User {
	return response.User{
		Id:       res.Id,
		FullName: res.FullName,
		UserName: res.UserName,
		Email:    res.Email,
	}
}
