module github.com/equinor/radix-cost-allocation

go 1.16

require (
	github.com/denisenkom/go-mssqldb v0.10.0
	github.com/equinor/radix-operator v1.11.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.20.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	github.com/vrischmann/envconfig v1.3.0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.19.9
	k8s.io/apimachinery v0.19.9
	k8s.io/client-go v12.0.0+incompatible
)

replace (
	github.com/apache/thrift => github.com/apache/thrift v0.13.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.1
	github.com/nats-io/jwt/v2 => github.com/nats-io/jwt/v2 v2.0.1
	github.com/nats-io/nats-server/v2 => github.com/nats-io/nats-server/v2 v2.2.0
	k8s.io/client-go => k8s.io/client-go v0.19.9
)
