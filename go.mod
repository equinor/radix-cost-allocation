module go.etcd.io/etcd/v3

go 1.16

require (
	github.com/denisenkom/go-mssqldb v0.0.0-20200428022330-06a60b6afbbc
	github.com/equinor/radix-cost-allocation v0.1.2
	github.com/kr/pretty v0.2.1 // indirect
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.20.0
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/apache/thrift v0.12.0 => github.com/apache/thrift v0.13.0
	github.com/gorilla/websocket v0.0.0-20170926233335-4201258b820c => github.com/gorilla/websocket v1.4.1
	github.com/nats-io/jwt v0.3.2 => github.com/nats-io/jwt/v2 v2.0.1
	github.com/nats-io/nats-server/v2 v2.1.2 => github.com/nats-io/nats-server/v2 v2.2.0
)
