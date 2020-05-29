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
	"log"
	"time"

	"github.com/keikoproj/aws-auth/pkg/mapper"
	"github.com/spf13/cobra"
)

var removeArgs = &mapper.MapperArguments{
	OperationType: mapper.OperationRemove,
}

// deleteCmd represents the base command when called without any subcommands
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove removes a user or role from the aws-auth configmap",
	Long:  `remove removes a user or role from the aws-auth configmap`,
	Run: func(cmd *cobra.Command, args []string) {
		k, err := getKubernetesClient(removeArgs.KubeconfigPath)
		if err != nil {
			log.Fatal(err)
		}

		worker := mapper.New(k, true)
		if err := worker.RemoveByUsername(removeArgs); err != nil {
			log.Fatal(err)
		}
	},
}

// removeByUsernameCmd removes all map roles and map users in an auth cm based on the input username
func removeByUsernameCmd() *cobra.Command {
	var removeArgs = &mapper.MapperArguments{}
	var command = &cobra.Command{
		Use:   "remove-by-username",
		Short: "remove-by-username removes all map roles and map users from the aws-auth configmap",
		Run: func(cmd *cobra.Command, args []string) {
			k, err := getKubernetesClient(removeArgs.KubeconfigPath)
			if err != nil {
				log.Fatal(err)
			}

			worker := mapper.New(k, true)

			if err := worker.RemoveByUsername(removeArgs); err != nil {
				log.Fatal(err)
			}
		},
	}

	command.Flags().StringVar(&removeArgs.KubeconfigPath, "kubeconfig", "", "Kubeconfig path")
	command.Flags().StringVar(&removeArgs.Username, "username", "", "Username to remove")
	command.Flags().BoolVar(&removeArgs.WithRetries, "retry", false, "Retry on failure with exponential backoff")
	command.Flags().DurationVar(&removeArgs.MinRetryTime, "retry-min-time", time.Millisecond*200, "Minimum wait interval")
	command.Flags().DurationVar(&removeArgs.MaxRetryTime, "retry-max-time", time.Second*30, "Maximum wait interval")
	command.Flags().IntVar(&removeArgs.MaxRetryCount, "retry-max-count", 12, "Maximum number of retries before giving up")
	return command
}

func init() {
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(removeByUsernameCmd())
	removeCmd.Flags().StringVar(&removeArgs.KubeconfigPath, "kubeconfig", "", "Kubeconfig path")
	removeCmd.Flags().StringVar(&removeArgs.Username, "username", "", "Username to remove")
	removeCmd.Flags().StringVar(&removeArgs.RoleARN, "rolearn", "", "Role ARN to remove")
	removeCmd.Flags().StringVar(&removeArgs.UserARN, "userarn", "", "User ARN to remove")
	removeCmd.Flags().StringSliceVar(&removeArgs.Groups, "groups", []string{}, "Groups to remove")
	removeCmd.Flags().BoolVar(&removeArgs.MapRoles, "maproles", false, "Removes a role")
	removeCmd.Flags().BoolVar(&removeArgs.MapUsers, "mapusers", false, "Removes a user")
	removeCmd.Flags().BoolVar(&removeArgs.WithRetries, "retry", false, "Retry on failure with exponential backoff")
	removeCmd.Flags().DurationVar(&removeArgs.MinRetryTime, "retry-min-time", time.Millisecond*200, "Minimum wait interval")
	removeCmd.Flags().DurationVar(&removeArgs.MaxRetryTime, "retry-max-time", time.Second*30, "Maximum wait interval")
	removeCmd.Flags().IntVar(&removeArgs.MaxRetryCount, "retry-max-count", 12, "Maximum number of retries before giving up")
}
