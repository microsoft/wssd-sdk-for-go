module github.com/microsoft/wssd-sdk-for-go

go 1.12

require (
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/microsoft/moc v0.8.0-alpha.31
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	go.opencensus.io v0.22.3
	google.golang.org/grpc v1.27.1
	k8s.io/klog v1.0.0
)

replace github.com/microsoft/moc => X:\go\src\github.com\microsoft\moc

