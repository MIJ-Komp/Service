package helpers

import (
	"math/rand"
	"strconv"
	"time"
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
