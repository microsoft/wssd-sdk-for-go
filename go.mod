module github.com/microsoft/wssd-sdk-for-go

go 1.15

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
<<<<<<< HEAD
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/microsoft/moc v0.10.19-alpha.6
=======
	github.com/microsoft/moc v0.10.19-alpha.7
>>>>>>> master
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.5.1
	go.opencensus.io v0.22.3
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.27.1
	k8s.io/klog v1.0.0
)

replace (
	github.com/microsoft/moc => github.com/hvedati/moc v0.10.20-alpha.1
	github.com/miekg/dns => github.com/miekg/dns v1.1.25
)
