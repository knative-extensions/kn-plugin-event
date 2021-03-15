package retcode

import "hash/crc32"

// Calc will calculate an POSIX retcode from an error.
func Calc(err error) int {
	if err == nil {
		return 0
	}
	return int(crc32.ChecksumIEEE([]byte(err.Error())))%254 + 1
}
