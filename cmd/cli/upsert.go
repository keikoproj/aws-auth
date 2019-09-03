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
	"github.com/keikoproj/aws-auth/pkg/mapper"
	"github.com/spf13/cobra"
)

var upsertArgs = &mapper.UpsertArguments{}

// upsertCmd represents the base command when called without any subcommands
var upsertCmd = &cobra.Command{
	Use:   "upsert",
	Short: "upsert updates or inserts a user or role to the aws-auth configmap",
	Long:  `upsert updates or inserts a user or role to the aws-auth configmap`,
	Run: func(cmd *cobra.Command, args []string) {
		k, err := getKubernetesClient(upsertArgs.KubeconfigPath)
		die(err)
		worker := mapper.New(k, true)
		err = worker.Upsert(upsertArgs)
		die(err)
	},
}

func init() {
	rootCmd.AddCommand(upsertCmd)
	upsertCmd.Flags().StringVar(&upsertArgs.KubeconfigPath, "kubeconfig", "", "Path to kubeconfig")
	upsertCmd.Flags().StringVar(&upsertArgs.Username, "username", "", "Username to upsert")
	upsertCmd.Flags().StringVar(&upsertArgs.RoleARN, "rolearn", "", "Role ARN to upsert")
	upsertCmd.Flags().StringVar(&upsertArgs.UserARN, "userarn", "", "User ARN to upsert")
	upsertCmd.Flags().StringSliceVar(&upsertArgs.Groups, "groups", []string{}, "Groups to upsert")
	upsertCmd.Flags().BoolVar(&upsertArgs.MapRoles, "maproles", false, "Upsert a role")
	upsertCmd.Flags().BoolVar(&upsertArgs.MapUsers, "mapusers", false, "Upsert a user")
}
