module custom-pkg/slack

go 1.14

require (
	constants v0.0.0
)

replace (
	constants v0.0.0 => ../../constants
)