package utils

import "github.com/google/uuid"

func Uuid() string {
	return uuid.Must(uuid.NewV7()).String()
}
