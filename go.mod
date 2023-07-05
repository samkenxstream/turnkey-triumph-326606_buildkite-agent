module github.com/buildkite/agent/v3

go 1.14

require (
	cloud.google.com/go/compute/metadata v0.2.3
	github.com/DataDog/datadog-go v3.7.2+incompatible
	github.com/aws/aws-sdk-go v1.32.10
	github.com/buildkite/bintest/v3 v3.1.0
	github.com/buildkite/interpolate v0.0.0-20200526001904-07f35b4ae251
	github.com/buildkite/shellwords v0.0.0-20180315084142-c3f497d1e000
	github.com/buildkite/yaml v0.0.0-20210326113714-4a3f40911396
	github.com/creack/pty v1.1.12
	github.com/denisbrodbeck/machineid v1.0.0
	github.com/google/go-querystring v0.0.0-20170111101155-53e6ce116135
	github.com/mattn/go-zglob v0.0.0-20180803001819-2ea3427bfa53
	github.com/mitchellh/go-homedir v1.0.0
	github.com/nightlyone/lockfile v0.0.0-20180618180623-0ad87eef1443
	github.com/oleiade/reflections v0.0.0-20160817071559-0e86b3c98b2f
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pborman/uuid v0.0.0-20170112150404-1b00554d8222
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/qri-io/jsonpointer v0.0.0-20180309164927-168dd9e45cf2 // indirect
	github.com/qri-io/jsonschema v0.0.0-20180607150648-d0d3b10ec792
	github.com/rjeczalik/interfaces v0.1.1
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/stretchr/testify v1.8.1
	github.com/urfave/cli v1.22.4
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
	golang.org/x/oauth2 v0.4.0
	golang.org/x/sys v0.4.0
	google.golang.org/api v0.103.0
	google.golang.org/grpc v1.53.0 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.28.0
)
