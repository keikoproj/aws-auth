package main

import (
	"fmt"
	"os"

	awsauth "github.com/eytan-avisror/aws-auth/pkg/mapper"
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
	myUpsertRole := &awsauth.UpsertArguments{
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

	myDeleteRole := &awsauth.RemoveArguments{
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
