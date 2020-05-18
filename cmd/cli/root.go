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
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-auth",
	Short: "aws-auth modifies the aws-auth configmap on eks clusters",
	Long:  `aws-auth modifies the aws-auth configmap on eks clusters`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getKubernetesClient(kubePath string) (kubernetes.Interface, error) {
	var config *rest.Config

	if kubePath == "" {
		userHome, _ := os.UserHomeDir()
		kubePath = fmt.Sprintf("%v/.kube/config", userHome)
		if os.Getenv("KUBECONFIG") != "" {
			kubePath = os.Getenv("KUBECONFIG")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubePath)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
