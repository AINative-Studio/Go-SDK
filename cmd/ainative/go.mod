module github.com/ainative/go-sdk/cmd/ainative

go 1.21

require (
	github.com/ainative/go-sdk v1.1.0
	github.com/spf13/cobra v1.8.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/ainative/go-sdk => ../..
