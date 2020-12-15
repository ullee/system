module watchdog

go 1.14

require (
	constants v0.0.0
	custom-pkg/aws/cw v0.0.0
	custom-pkg/logger v0.0.0
	github.com/aws/aws-sdk-go v1.33.10
	github.com/guillermo/go.procmeminfo v0.0.0-20131127224636-be4355a9fb0e
)

replace (
	constants v0.0.0 => ../constants
	custom-pkg/aws/cw v0.0.0 => ../custom-pkg/aws/cw
	custom-pkg/logger v0.0.0 => ../custom-pkg/logger
)
