# Yashiro

Yashiro is a templating engine with the external stores.

## Service

AWS

* [Systems Manager Parameter Store](https://docs.aws.amazon.com/systems-manager/)
* [Secrets Manager](https://docs.aws.amazon.com/secretsmanager/)

## Usage

See [Godoc](https://pkg.go.dev/github.com/dwango/yashiro).

```sh
go get github.com/dwango/yashiro
```

### Authorization

AWS

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ssm:GetParameter",
                "secretsmanager:GetSecretValue"
            ],
            "Resource": ["*"],
        },
    ]
}
```

## CLI Tool

### Installation

#### From release page

Download binary from [release page](https://github.com/dwango/yashiro/releases).

#### Homebrew Users

Download and install by homebrew.

```sh
brew tap dwango/yashiro
brew install ysr
```

### Example

See [example](./example/).
