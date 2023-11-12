# CLI Example

Example of using AWS Systems Parameter Store.

## Requires

* aws account
* aws cli
* `ysr` binary

## Setup

```sh
$ export AWS_PROFILE=$profile_name # if need
$ aws ssm put-parameter --name '/yashiro/example' --value '{"imageTag": "latest", "roleArn": "arn:aws:iam::012345678901:role/PodRole29A92600-1NSKGEZRCWPNU"}' --type String
{
    "Version": 1,
    "Tier": "Standard"
}
$ aws ssm put-parameter --name '/yashiro/example/secure' --value 'password' --type SecureString
{
    "Version": 1,
    "Tier": "Standard"
}
```

## Execute

```sh
$ ysr template -c ./yashiro.yaml example.yaml.tmpl
---
apiVersion: v1
kind: Secret
metadata:
  name: example
  labels:
    app: example
data:
  db_password: cGFzc3dvcmQ=
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: example
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::012345678901:role/PodRole29A92600-1NSKGEZRCWPNU
  labels:
    app: example
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
spec:
  selector:
    matchLabels:
      app: example
  template:
    metadata:
      labels:
        app: example
    spec:
      serviceAccountName: example
      containers:
      - name: hello-world
        image: hello-world:latest
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
```

## Teardown

```sh
aws ssm delete-parameter --name '/yashiro/example'
aws ssm delete-parameter --name '/yashiro/example/secure'
```
