module github.com/microsoft/wssd-sdk-for-go

go 1.14

require (
	github.com/microsoft/moc v0.10.9
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.5.1
	go.opencensus.io v0.22.3
	google.golang.org/grpc v1.27.1
	google.golang.org/protobuf v1.25.0 // indirect
	k8s.io/klog v1.0.0
)

replace github.com/miekg/dns => github.com/miekg/dns v1.1.25
