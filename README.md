
# aws-auth
[![unit-test](https://github.com/keikoproj/aws-auth/actions/workflows/unit-test.yaml/badge.svg?branch=master)](https://github.com/keikoproj/aws-auth/actions/workflows/unit-test.yaml)
[![codecov](https://codecov.io/gh/keikoproj/aws-auth/branch/master/graph/badge.svg)](https://codecov.io/gh/keikoproj/aws-auth)
[![Go Report Card](https://goreportcard.com/badge/github.com/keikoproj/aws-auth)](https://goreportcard.com/report/github.com/keikoproj/aws-auth)


The `aws-auth` utility and library makes the management of the `aws-auth` ConfigMap for EKS Kubernetes clusters easier and safer.

## Use cases

- make bootstrapping a node group or removing/adding user access on EKS fast and easy

- useful for automation purposes, any workflow that needs to grant IAM access to an EKS cluster can use this library to modify the config map.

- run as part of a workflow on kubernetes using a docker image

The `aws-auth` tool is referenced in the AWS EKS best practices documentation [here](https://aws.github.io/aws-eks-best-practices/security/docs/iam/#use-tools-to-make-changes-to-the-aws-auth-configmap).

## Install

`aws-auth` includes both a CLI and a [go library](#usage-as-a-library). You can install the CLI via `go get` or as a kubectl plugin via [Krew](https://krew.sigs.k8s.io/) or by downloading a binary from the [releases page](https://github.com/keikoproj/aws-auth/releases).

### go get

```text
go get github.com/keikoproj/aws-auth
aws-auth help
```

### kubectl krew

Alternatively, install aws-auth with the krew plugin manager for kubectl.

```
kubectl krew install aws-auth
kubectl aws-auth
```

### Download release artifact

The latest release artifacts can be downloaded from the [GitHub releases page](https://github.com/keikoproj/aws-auth/releases/latest).

Or you can use the following command to download the latest release artifact for your platform:

``` bash
curl -s https://api.github.com/repos/keikoproj/aws-auth/releases/latest
| grep "browser_download_url" \
| grep $(go env GOARCH) | grep $(go env GOOS) \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -qi -
```

## Usage from command line or Krew

Either download/install a released binary or add as a plugin to kubectl via Krew

```text
$ kubectl krew update
$ kubectl krew install aws-auth
Installing plugin: aws-auth
Installed plugin: aws-auth

$ kubectl krew aws-auth
aws-auth modifies the aws-auth configmap on eks clusters

Usage:
  aws-auth [command]

Available Commands:
  help               Help about any command
  remove             remove removes a user or role from the aws-auth configmap
  remove-by-username remove-by-username removes all map roles and map users from the aws-auth configmap
  upsert             upsert updates or inserts a user or role to the aws-auth configmap
  version            Version of aws-auth

Flags:
  -h, --help   help for aws-auth

Use "aws-auth [command] --help" for more information about a command.
```

Given a config map with the following data:

```text
$ kubectl get configmap aws-auth -n kube-system -o yaml
apiVersion: v1
kind: ConfigMap
metadata:
    name: aws-auth
    namespace: kube-system
data:
  mapRoles: |
    - rolearn: arn:aws:iam::555555555555:role/devel-worker-nodes-NodeInstanceRole-74RF4UBDUKL6
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
    - rolearn: arn:aws:iam::555555555555:role/abc
      username: ops-user
      groups:
        - system:masters
  mapUsers: |
    - userarn: arn:aws:iam::555555555555:user/a-user
      username: admin
      groups:
        - system:masters
    - userarn: arn:aws:iam::555555555555:user/a-user
      username: ops-user
      groups:
        - system:masters
```

Remove all access belonging to an ARN (both mapUser roles will be removed)

```text
$ aws-auth remove --mapusers --userarn arn:aws:iam::555555555555:user/a-user
removed arn:aws:iam::555555555555:user/a-user from aws-auth
```

Remove by full match (only `mapUsers[0]` will be removed)

```text
$ aws-auth remove --mapusers --userarn arn:aws:iam::555555555555:user/a-user --username admin --groups system:masters
removed arn:aws:iam::555555555555:user/a-user from aws-auth
```

Remove based on a username

This command removes all map roles and map users that have matching input username. In the above configmap, map role for roleARN *arn:aws:iam::555555555555:role/abc* and mapUser for userARN *arn:aws:iam::555555555555:user/a-user* will be removed.

```text
$ aws-auth remove-by-username --username ops-user
```


Bootstrap a new node group role

```text
$ aws-auth upsert --maproles --rolearn arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6 --username system:node:{{EC2PrivateDNSName}} --groups system:bootstrappers system:nodes
added arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6 to aws-auth
```

You can also add retries with exponential backoff

```text
$ aws-auth upsert --maproles --rolearn arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6 --username system:node:{{EC2PrivateDNSName}} --groups system:bootstrappers system:nodes --retry
```

Retries are configurable using the following flags

```text
      --retry                     Retry on failure with exponential backoff
      --retry-max-count int       Maximum number of retries before giving up (default 12)
      --retry-max-time duration   Maximum wait interval (default 30s)
      --retry-min-time duration   Minimum wait interval (default 200ms)
```

Append groups to mapping instead of overwriting by using --append

```
$ aws-auth upsert --maproles --rolearn arn:aws:iam::00000000000:role/test --username test --groups test --append
```

Avoid overwriting username by using --update-username=false

```
$ aws-auth upsert --maproles --rolearn arn:aws:iam::00000000000:role/test --username test2 --groups test --update-username=false
```

Use the `get` command to get a detailed view of mappings

```
$ aws-auth get

TYPE        	ARN                                               USERNAME                         	GROUPS
Role Mapping	arn:aws:iam::555555555555:role/my-new-node-group  system:node:{{EC2PrivateDNSName}}	system:bootstrappers, system:nodes
```

use impersonate
```
aws-auth get|update|remove --as <username> --as-group <groupname> 
```

## Usage as a library

```go


package main

import (
    awsauth "github.com/keikoproj/aws-auth/pkg/mapper"
)

func someFunc(client kubernetes.Interface) error {
    awsAuth := awsauth.New(client, false)
    myUpsertRole := &awsauth.MapperArguments{
        MapRoles: true,
        RoleARN:  "arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6",
        Username: "system:node:{{EC2PrivateDNSName}}",
        Groups: []string{
            "system:bootstrappers",
            "system:nodes",
        },
        WithRetries: true,
        MinRetryTime:   time.Millisecond * 100,
        MaxRetryTime:   time.Second * 30,
        MaxRetryCount:  12,
    }

    err = awsAuth.Upsert(myUpsertRole)
    if err != nil {
        return err
    }
}

```

## Run in a container

```shell
$ docker run \
-v ~/.kube/:/root/.kube/ \
-v ~/.aws/:/root/.aws/ \
keikoproj/aws-auth:latest \
aws-auth upsert --mapusers \
--userarn arn:aws:iam::555555555555:user/a-user \
--username admin \
--groups system:masters
```
