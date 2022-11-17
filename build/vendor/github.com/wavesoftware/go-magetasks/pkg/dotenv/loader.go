package dotenv

import (
	"os"

	"github.com/joho/godotenv"
)

// Load the .env file.
func Load() error {
	err := godotenv.Load()
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
