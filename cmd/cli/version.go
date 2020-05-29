/*

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
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// gitCommit is a constant representing the source version that
	// generated this build. It should be set during build via -ldflags.
	gitCommit string
	// buildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	//It should be set during build via -ldflags.
	buildDate string
	// version is the aws-auth package version
	pkgVersion string = "0.2.0"
)

// Info holds the information related to descheduler app version.
type Info struct {
	PackageVersion string `json:"pkgVersion"`
	GitCommit      string `json:"gitCommit"`
	BuildDate      string `json:"buildDate"`
	GoVersion      string `json:"goVersion"`
	Compiler       string `json:"compiler"`
	Platform       string `json:"platform"`
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	return Info{
		GitCommit:      gitCommit,
		BuildDate:      buildDate,
		GoVersion:      runtime.Version(),
		Compiler:       runtime.Compiler,
		Platform:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		PackageVersion: pkgVersion,
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of aws-auth",
	Long:  `Prints the version of aws-auth.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("aws-auth version %+v\n", Get())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
