package main

import (
	"api.mijkomp.com/config"
	"api.mijkomp.com/exception"
)

// @title					Try Out Api
// @version				1.0.0
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description		Api documentation Try Out App
// @description		Description for what is this security definition being used
// @schemes http https
func main() {
	config.New(".env")
	server := InitializedServer()

	// err := server.Listen(":5000")
	err := server.Listen("127.0.0.1:5000")
	exception.PanicIfNeeded(err)
}
