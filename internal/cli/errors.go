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
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/client/pkg/output/term"
)

func errorHandler(_ error, cmd *cobra.Command) bool {
	ctx := cmd.Context()
	if logfile := outlogging.LogFileFrom(ctx); logfile != nil {
		logpath := logfile.Name()
		if term.IsFancy(cmd.ErrOrStderr()) {
			logpath = color.CyanString(logpath)
		}
		cmd.PrintErrln()
		cmd.PrintErrln("The logs could help to debug the failure reason.")
		cmd.PrintErrln("Take a look at the log file:", logpath)
	}
	return false
}
