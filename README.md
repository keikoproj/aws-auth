
# aws-auth
[![Build Status](https://travis-ci.org/keikoproj/aws-auth.svg?branch=master)](https://travis-ci.org/keikoproj/aws-auth)
[![codecov](https://codecov.io/gh/keikoproj/aws-auth/branch/master/graph/badge.svg)](https://codecov.io/gh/keikoproj/aws-auth)


> Makes the management of the aws-auth config map for EKS Kubernetes clusters easier

## Use cases

- make bootstrapping a node group or removing/adding user access on EKS fast and easy

- useful for automation purposes, any workflow that needs to grant IAM access to an EKS cluster can use this library to modify the config map.

- run as part of a workflow on kubernetes using a docker image

## Install

```text
$ go get github.com/keikoproj/aws-auth
$
```

## Usage from command line

```text
$ aws-auth
aws-auth modifies the aws-auth configmap on eks clusters

Usage:
  aws-auth [command]

Available Commands:
  help        Help about any command
  remove      remove removes an auth-map from mapRoles or mapUsers
  remove-by-username remove-by-username removes all map roles and map users from the aws-auth configmap
  upsert      upsert updates or inserts an auth-map to mapRoles
  version     Version of aws-auth

Flags:
  -h, --help   help for aws-auth
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
$ aws-auth upsert --maproles --userarn arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6 --username system:node:{{EC2PrivateDNSName}} --groups system:bootstrappers system:nodes
added arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6 to aws-auth
```

You can also add retries with exponential backoff

```text
$ aws-auth upsert --maproles --userarn arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6 --username system:node:{{EC2PrivateDNSName}} --groups system:bootstrappers system:nodes --retry
```

Retries are configurable using the following flags

```text
      --retry                     Retry on failure with exponential backoff
      --retry-max-count int       Maximum number of retries before giving up (default 12)
      --retry-max-time duration   Maximum wait interval (default 30s)
      --retry-min-time duration   Minimum wait interval (default 200ms)
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
        MinRetryTime:  time.Millisecond * 100,
        MaxRetryTime:  time.Second * 30,
        MaxRetryCount: 12,
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
