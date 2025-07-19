package utils

import "github.com/google/uuid"

// GenerateUUID 生成36位的UUID
func GenerateUUID() string {
	random, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return random.String()
}
