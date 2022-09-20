package image

import "os"

// BaseName returning a basename of the images.
func BaseName() string {
	return env("localhost", "KO_DOCKER_REPO", "IMAGE_BASENAME")
}

// BaseNameSeparator returns a separator between name and basename of OCI image.
func BaseNameSeparator() string {
	return env("/", "IMAGE_BASENAME_SEPARATOR")
}

func env(defaultVal string, keys ...string) string {
	for _, key := range keys {
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
	}
	return defaultVal
}
