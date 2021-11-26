module github.com/kiprotect/kodex

go 1.16

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/gin-gonic/gin v1.7.4
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/btree v1.0.0
	github.com/google/gopacket v1.1.19
	github.com/kiprotect/go-helpers v0.0.0-20211126144132-63cf08b7160f
	github.com/kr/pretty v0.2.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/streadway/amqp v1.0.0
	github.com/urfave/cli v1.22.4
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

// replace github.com/kiprotect/go-helpers => ../go-helpers
