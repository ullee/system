module ec2

go 1.14

require (
	constants v0.0.0
	github.com/aws/aws-sdk-go v1.33.10
)

replace (
	constants v0.0.0 => ../../../constants
)
