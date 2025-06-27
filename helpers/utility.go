package helpers

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func NewRandom(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Initialize an empty string to hold the random numbers
	var result string

	// Generate {length} random numbers and append them to the result string
	for i := 0; i < length; i++ {
		randomNum := r.Intn(9)            // Random number between 0 and 9
		result += strconv.Itoa(randomNum) // Convert to string and append
	}

	return result
}

func SetDefaultPageRequest(page *int, pageSize *int) {
	if *page == 0 {
		*page = 1
	}
	if *pageSize == 0 {
		*pageSize = 10
	}
}

func SplitImageIds(imageIds *string) []uuid.UUID {
	var imgIds []uuid.UUID

	if imageIds == nil || *imageIds == "" {
		return imgIds
	}

	images := strings.Split(*imageIds, ";")
	for _, img := range images {
		id, err := uuid.Parse(strings.TrimSpace(img))
		if err == nil {
			imgIds = append(imgIds, id)
		}
		// Optional: bisa log atau handle error jika UUID tidak valid
	}

	return imgIds
}

func JoinImageIds(imageIds *[]uuid.UUID) *string {
	if imageIds == nil || len(*imageIds) == 0 {
		return nil
	}

	var builder strings.Builder
	for i, img := range *imageIds {
		builder.WriteString(img.String())
		if i != len(*imageIds)-1 {
			builder.WriteString(";")
		}
	}

	result := builder.String()
	return &result
}
