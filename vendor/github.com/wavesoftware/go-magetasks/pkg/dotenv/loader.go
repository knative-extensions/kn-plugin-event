package dotenv

import (
	"github.com/joho/godotenv"
)

// Load the .env file.
func Load() error {
	return godotenv.Load()
}
