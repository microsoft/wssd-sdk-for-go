package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microsoft/moc/pkg/certs"
	"github.com/microsoft/moc/rpc/common"
	wssdcommonproto "github.com/microsoft/moc/rpc/common"
	wssdcompute "github.com/microsoft/moc/rpc/nodeagent/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func Test_getVirtualMachine(t *testing.T) {
	proxy := NewProxy()
	defer proxy.Target.Close()
	caCert, _, err := certs.GenerateClientCertificate("ValidCertificate")
	if err != nil {
		t.Fatalf(err.Error())
	}
	certBytes := certs.EncodeCertPEM(caCert)
	caCertString := string(certBytes)
	var (
		vmName            = "VM-Name"
		port       uint16 = 1234
		disableRDP        = true
		httpProxy         = proxy.Target.URL
		httpsProxy        = proxy.Target.URL
		noProxy           = []string{"localhost", "127.0.0.1", ".svc", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.0.0.0/8", ".corp.microsoft.com", ".masd.stbtest.microsoft.com"}
		trustedCa         = caCertString
	)

	type args struct {
		vm *wssdcompute.VirtualMachine
	}

	type want struct {
		vm *compute.VirtualMachine
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "when input vm has nil linux configuration then output vm has nil linux configuration",
			args: args{
				vm: &wssdcompute.VirtualMachine{
					Name:                  vmName,
					PowerState:            common.PowerState_Running,
					HighAvailabilityState: common.HighAvailabilityState_STABLE,
					Network:               &wssdcompute.NetworkConfiguration{},
					Storage: &wssdcompute.StorageConfiguration{
						Osdisk: &wssdcompute.Disk{},
					},
					Status: &common.Status{
						ProvisioningStatus: &common.ProvisionStatus{
							CurrentState: common.ProvisionState_CREATED,
						},
						Health: &common.Health{
							CurrentState: common.HealthState_OK,
						},
						LastError:      &common.Error{},
						Version:        &common.Version{},
						DownloadStatus: &common.DownloadStatus{},
					},
					Os: &wssdcompute.OperatingSystemConfiguration{
						LinuxConfiguration: nil,
						WindowsConfiguration: &wssdcompute.WindowsConfiguration{
							RDPConfiguration: &wssdcompute.RDPConfiguration{
								DisableRDP: disableRDP,
								Port:       uint32(port),
							},
						},
						ProxyConfiguration: &wssdcommonproto.ProxyConfiguration{
							HttpProxy:  httpProxy,
							HttpsProxy: httpsProxy,
							NoProxy:    noProxy,
							TrustedCa:  trustedCa,
						},
					},
				},
			},
			want: want{
				vm: &compute.VirtualMachine{
					ID:   proto.String(""),
					Tags: map[string]*string{},
					Name: &vmName,
					VirtualMachineProperties: &compute.VirtualMachineProperties{
						SecurityProfile: &compute.SecurityProfile{
							EnableTPM:    proto.Bool(false),
							UefiSettings: nil,
						},
						HardwareProfile: &compute.HardwareProfile{
							VMSize: compute.VirtualMachineSizeTypesDefault,
						},
						StorageProfile: &compute.StorageProfile{
							DataDisks: &[]compute.DataDisk{},
							OsDisk: &compute.OSDisk{
								VhdName: proto.String(""),
							},
							VmConfigContainerName: proto.String(""),
						},
						OsProfile: &compute.OSProfile{
							ComputerName:       proto.String(""),
							LinuxConfiguration: nil,
							WindowsConfiguration: &compute.WindowsConfiguration{
								EnableAutomaticUpdates: proto.Bool(false),
								TimeZone:               proto.String(""),
								RDP: &compute.RDPConfiguration{
									DisableRDP: &disableRDP,
									Port:       &port,
								},
							},
							OsBootstrapEngine: compute.CloudInit,
							ProxyConfiguration: &compute.ProxyConfiguration{
								HttpProxy:  &httpProxy,
								HttpsProxy: &httpsProxy,
								NoProxy:    &noProxy,
								TrustedCa:  &trustedCa,
							},
						},
						NetworkProfile: &compute.NetworkProfile{
							NetworkInterfaces: &[]compute.NetworkInterfaceReference{},
						},
						Statuses: map[string]*string{
							"DownloadStatus": proto.String(""),
							"UploadStatus":   proto.String("<nil>"),
							"Error":          proto.String(""),
							"HealthState":    proto.String("currentState:OK "),
							"ProvisionState": proto.String("currentState:CREATED "),
							"Version":        proto.String(""),
							"PowerState":     proto.String(common.PowerState_Running.String()),
						},
						ProvisioningState:       proto.String(common.ProvisionState_CREATED.String()),
						DisableHighAvailability: proto.Bool(false),
						IsPlaceholder:           proto.Bool(false),
						HighAvailabilityState:   proto.String(common.HighAvailabilityState_STABLE.String()),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wssdClient := client{}
			vm := wssdClient.getVirtualMachine(tt.args.vm)
			assert.Equal(t, tt.want.vm, vm)
		})
	}
}

// Proxy is a simple proxy server for unit tests.
type Proxy struct {
	Target *httptest.Server
}

// NewProxy creates a new proxy server for unit tests.
func NewProxy() *Proxy {
	target := httptest.NewServer(http.DefaultServeMux)
	return &Proxy{Target: target}
}
