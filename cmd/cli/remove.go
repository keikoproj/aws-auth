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
	"github.com/eytan-avisror/aws-auth/pkg/mapper"
	"github.com/spf13/cobra"
)

var removeArgs = &mapper.RemoveArguments{}

// deleteCmd represents the base command when called without any subcommands
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove removes an auth-map from mapRoles or mapUsers",
	Long:  `remove removes auth-map from mapRoles or mapUsers`,
	Run: func(cmd *cobra.Command, args []string) {
		k, err := getKubernetesClient(removeArgs.KubeconfigPath)
		die(err)
		worker := mapper.New(k)
		err = worker.Remove(removeArgs)
		die(err)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringVar(&removeArgs.KubeconfigPath, "kubeconfig", "", "Kubeconfig path")
	removeCmd.Flags().StringVar(&removeArgs.Username, "username", "", "Username to upsert")
	removeCmd.Flags().StringVar(&removeArgs.RoleARN, "rolearn", "", "Role ARN to upsert")
	removeCmd.Flags().StringSliceVar(&removeArgs.Groups, "groups", []string{}, "Groups to upsert")
	removeCmd.Flags().BoolVar(&removeArgs.MapRoles, "maproles", false, "Upsert a mapRoles")
	removeCmd.Flags().BoolVar(&removeArgs.MapUsers, "mapusers", false, "Upsert a mapUsers")
}
