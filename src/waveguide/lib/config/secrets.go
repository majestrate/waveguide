package config

import (
	"os"
)

func GetAPISecret() []byte {
	return []byte(os.Getenv("WAVEGUIDED_API_SECRET"))
}
