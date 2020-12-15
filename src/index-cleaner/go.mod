module index-cleaner

go 1.14

require (
	constants v0.0.0
	custom-pkg/logger v0.0.0
	github.com/aws/aws-sdk-go v1.33.10 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
)

replace (
	constants v0.0.0 => ../constants
	custom-pkg/logger v0.0.0 => ../custom-pkg/logger
)
