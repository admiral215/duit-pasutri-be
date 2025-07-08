package utils

import "github.com/google/uuid"

func ToUUID(value string) uuid.UUID {
	return uuid.MustParse(value)
}
