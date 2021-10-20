//go:build !ignore
// +build !ignore

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/magefile/mage/mage"
)

const parseErrCode = 2

var errUnknownCommandType = errors.New("unknown command type")

func main() {
	os.Exit(ParseAndRun(os.Stdout, os.Stderr, os.Stdin, os.Args[1:]))
}

// ParseAndRun parses the command line, and then compiles and runs the mage
// files in the given directory with the given args (do not include the command
// name in the args).
// A copy of mage.ParseAndRun with overwritten magefiles dir.
func ParseAndRun(stdout, stderr io.Writer, stdin io.Reader, args []string) int {
	errlog := log.New(stderr, "", 0)
	inv, cmd, err := mage.Parse(stderr, stdout, args)
	inv.Stderr = stderr
	inv.Stdin = stdin
	inv.WorkDir = ".."
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if err != nil {
		errlog.Println("Error:", err)
		return parseErrCode
	}

	switch cmd {
	case
		mage.Version,
		mage.Init,
		mage.Clean:
		return mage.ParseAndRun(stdout, stderr, stdin, append(args, "-w", ".."))
	case
		mage.CompileStatic,
		mage.None:
		return mage.Invoke(inv)
	default:
		panic(fmt.Errorf("%w: %v", errUnknownCommandType, cmd))
	}
}
