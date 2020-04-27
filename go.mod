module github.com/microsoft/wssd-sdk-for-go

go 1.12

require (
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/microsoft/moc v0.8.0-alpha.12
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
	google.golang.org/grpc v1.27.1
	k8s.io/klog v1.0.0
)

replace github.com/microsoft/moc => /home/erfrimod/repo/gopath/src/github.com/microsoft/moc