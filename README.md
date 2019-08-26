# aws-auth
> Makes the management of the aws-auth config map for EKS Kubernetes clusters easier

## Install

```
$ go get github.com/eytan-avisror/aws-auth
```

## Usage from command line

```
$ aws-auth
aws-auth modifies the aws-auth configmap on eks clusters

Usage:
  aws-auth [command]

Available Commands:
  help        Help about any command
  remove      remove removes an auth-map from mapRoles or mapUsers
  upsert      upsert updates or inserts an auth-map to mapRoles
  version     Version of aws-auth

Flags:
  -h, --help   help for aws-auth
```

Given a config map with the following data:
```
$ kubectl get configmap aws-auth -n kube-system -o yaml
apiVersion: v1
kind: ConfigMap
metadata:

apiVersion: v1
data:
  mapRoles: |
    - rolearn: arn:aws:iam::555555555555:role/devel-worker-nodes-NodeInstanceRole-74RF4UBDUKL6
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
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

```
$ aws-auth remove --mapusers --userarn arn:aws:iam::555555555555:user/a-user
```

Remove by full match (only mapUsers[0] will be removed)
```
$ aws-auth remove --mapusers --userarn arn:aws:iam::555555555555:user/a-user --username admin --groups system:masters
```

Bootstrap a new node group role
```
$ aws-auth uspert --maproles --userarn arn:aws:iam::555555555555:role/my-new-node-group-NodeInstanceRole-74RF4UBDUKL6 --username system:node:{{EC2PrivateDNSName}} --groups system:bootstrappers system:nodes
```

## Usage as a library
