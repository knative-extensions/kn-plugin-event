package output

import (
	"fmt"
	"log"
	"os"

	"github.com/wavesoftware/go-magetasks/pkg/output/color"
)

// Setup the output of tasks.
func Setup() {
	color.SetupMode()
}

// PrintPending works similar to log.Print, and displays the prefix.
func PrintPending(args ...interface{}) {
	args = appendPrefix(args...)
	printOrFail(args...)
}

// PrintEnd works similar to log.Print, and end line without adding the prefix.
func PrintEnd(args ...interface{}) {
	args = append(args, "\n")
	printOrFail(args...)
}

// Println works similar to log.Println.
func Println(args ...interface{}) {
	args = appendPrefix(args...)
	args = append(args, "\n")
	printOrFail(args...)
}

// Printlnf works similar to log.Printf.
func Printlnf(format string, args ...interface{}) {
	Println(fmt.Sprintf(format, args...))
}

// printOrFail works similar to log.Print.
func printOrFail(args ...interface{}) {
	_, err := fmt.Fprint(os.Stdout, args...)
	if err != nil {
		log.Fatal(err)
	}
}

func appendPrefix(args ...interface{}) []interface{} {
	return append([]interface{}{prefix(), " "}, args...)
}
