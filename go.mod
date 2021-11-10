module github.com/grpc-kit/cli

go 1.12

require (
	github.com/grpc-kit/pkg v0.1.1
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.5
	github.com/grpc-kit/cfg => /Users/coding/github.com/grpc-kit/cfg
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20191216164720-4f79533eabd1
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
