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

package main

import (
	"fmt"
	"os"

	awsauth "github.com/keikoproj/aws-auth/pkg/mapper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	kubePath := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubePath)
	if err != nil {
		os.Exit(1)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		os.Exit(1)
	}

	awsAuth := awsauth.New(client, false)
	myUpsertRole := &awsauth.MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups: []string{
			"system:bootstrappers",
			"system:nodes",
		},
	}

	err = awsAuth.Upsert(myUpsertRole)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	myDeleteRole := &awsauth.MapperArguments{
		MapRoles: true,
		RoleARN:  "arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6",
		Username: "system:node:{{EC2PrivateDNSName}}",
		Groups: []string{
			"system:bootstrappers",
			"system:nodes",
		},
	}

	err = awsAuth.Remove(myDeleteRole)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
