/*
 Copyright 2024 The Knative Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cli

import (
	"io"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/client/pkg/output/term"
	"knative.dev/kn-plugin-event/pkg/errors"
)

const likelyErrorChainDepth = 3

var (
	numbersRE = regexp.MustCompile(`\s+\d+\s+`)
	stringRE  = regexp.MustCompile(`\s+"[^"]+"\s+`)
)

func errorHandler(err error, cmd *cobra.Command) bool {
	prettyPrintErr(err, cmd)
	ctx := cmd.Context()
	if logfile := outlogging.LogFileFrom(ctx); logfile != nil {
		logpath := logfile.Name()
		hint := "🌟 Hint:"
		if term.IsFancy(cmd.ErrOrStderr()) {
			logpath = color.CyanString(logpath)
			hint = color.YellowString(hint)
		}
		cmd.PrintErrln()
		cmd.PrintErrln(hint, "The execution logs could help debug the failure.")
		cmd.PrintErrln("         Consider, taking a look at the log file:", logpath)
	}
	return false
}

func prettyPrintErr(err error, cmd *cobra.Command) {
	prefix := "🔥 Error:"
	stderr := cmd.ErrOrStderr()
	if term.IsFancy(stderr) {
		prefix = color.RedString(prefix)
	}
	messages := make([]string, 0, likelyErrorChainDepth)
	messages = append(messages, err.Error())
	for {
		if err = errors.Cause(err); err != nil {
			messages = append(messages, err.Error())
		} else {
			break
		}
	}
	for i := 0; i < len(messages); i++ {
		j := i + 1
		if j < len(messages) {
			messages[i] = strings.Replace(messages[i], ": "+messages[j], "", 1)
		}
		if i == 0 {
			cmd.PrintErrln(prefix, colorizeMessage(messages[i], stderr))
		} else {
			padding := strings.Repeat("  ", i) + "└─ caused by:"
			if term.IsFancy(stderr) {
				padding = color.RedString(padding)
			}
			cmd.PrintErrln(padding, colorizeMessage(messages[i], stderr))
		}
	}
}

func colorizeMessage(msg string, out io.Writer) string {
	if !term.IsFancy(out) {
		return msg
	}
	msg = numbersRE.ReplaceAllStringFunc(msg, func(num string) string {
		return color.YellowString(num)
	})
	msg = stringRE.ReplaceAllStringFunc(msg, func(s string) string {
		return color.GreenString(s)
	})
	return msg
}
