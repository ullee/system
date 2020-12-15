module watchdog-server

go 1.14

require (
	constants v0.0.0
	custom-pkg/aws/cw v0.0.0
	custom-pkg/aws/ec2 v0.0.0
	custom-pkg/aws/rds v0.0.0
	custom-pkg/logger v0.0.0
	custom-pkg/slack v0.0.0
	github.com/aws/aws-sdk-go v1.33.10
)

replace (
	constants v0.0.0 => ../constants
	custom-pkg/aws/cw v0.0.0 => ../custom-pkg/aws/cw
	custom-pkg/aws/ec2 v0.0.0 => ../custom-pkg/aws/ec2
	custom-pkg/aws/rds v0.0.0 => ../custom-pkg/aws/rds
	custom-pkg/logger v0.0.0 => ../custom-pkg/logger
	custom-pkg/slack v0.0.0 => ../custom-pkg/slack
)
