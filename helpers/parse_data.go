package helpers

import (
	"fmt"
	"strconv"
	"time"

	"api.mijkomp.com/exception"
	"github.com/google/uuid"
)

func ParseUserId(userId interface{}) uint {

	strId := fmt.Sprintf("%v", userId)

	id, _ := strconv.ParseUint(strId, 10, 64)

	return uint(id)
}

func ParseFloat64(str string) float64 {
	id, _ := strconv.ParseFloat(str, 64)
	return id
}

func ParseInt(str string) int {
	val, _ := strconv.Atoi(str)
	return int(val)
}

func ParseNullableInt(str string) *int {
	val, err := strconv.ParseInt(str, 10, 32)

	if err != nil {
		return nil
	}
	res := int(val)
	return &res
}

func ParseNullableInt64(str string) *int64 {
	val, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		return nil
	}

	return &val
}

func ParseInt64(str string) int64 {
	val, _ := strconv.Atoi(str)
	return int64(val)
}

func ParseUint(strId string) uint {
	id, _ := strconv.ParseUint(strId, 10, 64)
	return uint(id)
}

func ParseNullableUint(str string) *uint {
	val, err := strconv.ParseUint(str, 10, 64)

	if err != nil {
		return nil
	}

	res := uint(val)
	return &res
}

func ParseNullableString(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func ParseUUID(str string) uuid.UUID {

	uuid, err := uuid.Parse(str)
	exception.PanicIfNeeded(err)

	return uuid
}

func ParseNullableUUID(str string) *uuid.UUID {
	if str == "" {
		return nil
	}
	uuid, err := uuid.Parse(str)
	exception.PanicIfNeeded(err)
	return &uuid
}

func ParseNullableBool(str string) *bool {
	if str == "" {
		return nil
	}

	boolean, err := strconv.ParseBool(str)
	exception.PanicIfNeeded(err)

	return &boolean
}

func ParseBool(str string) bool {

	boolean, err := strconv.ParseBool(str)
	exception.PanicIfNeeded(err)

	return boolean
}

func ParseNullableTime(str string) *time.Time {

	if str == "" {
		return nil
	}

	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, str)

	exception.PanicIfNeeded(err)

	return &parsedTime
}
