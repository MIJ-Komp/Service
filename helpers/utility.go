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

func SplitImageIds(imageIds string) []uuid.UUID {
	images := strings.Split(imageIds, ";")

	imgIds := []uuid.UUID{}
	if imageIds == "" {
		return imgIds
	}
	for _, img := range images {
		imgIds = append(imgIds, ParseUUID(img))
	}

	return imgIds
}

func JoinImageIds(imageIds []uuid.UUID) string {
	imgIds := ""

	for i, img := range imageIds {
		if i != len(imageIds) {
			imgIds += img.String() + ";"
		} else {
			imgIds += img.String()
		}
	}

	return imgIds
}
