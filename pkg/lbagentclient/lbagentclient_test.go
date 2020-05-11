//go test -v --args -lbagentip "10.137.196.121"

package lbagentclient

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/microsoft/moc/pkg/auth"
	pbcom "github.com/microsoft/moc/rpc/common"
	pblbagent "github.com/microsoft/moc/rpc/lbagent"
	admin_pb "github.com/microsoft/moc/rpc/lbagent/admin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var lbagentip = flag.String("lbagentip", "", "ip address of the lbagent")

func getClient() (*LoadBalancerAgentClient, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment(*lbagentip)
	if err != nil {
		return nil, err
	}

	return NewLoadBalancerAgentClient(*lbagentip, authorizer)
}

func getHealthClient() (admin_pb.HealthAgentClient, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment(*lbagentip)
	if err != nil {
		return nil, err
	}

	return GetHealthClient(lbagentip, authorizer)
}

func TestAdd(t *testing.T) {

	c, err := getClient()
	require.NoError(t, err)

	lbr := []*pblbagent.LoadBalancer{
		{
			Name:       "xxx",
			Frontendip: "*",
			Backendips: []string{"10.0.1.1", "10.0.1.2"},
			Loadbalancingrules: []*pblbagent.LoadBalancingRule{
				{
					FrontendPort: 80,
					BackendPort:  80,
					Protocol:     pbcom.Protocol_Tcp,
				},
			},
		},
		{
			Name:       "xxx1",
			Frontendip: "*",
			Backendips: []string{"10.0.1.1", "10.0.1.2"},
			Loadbalancingrules: []*pblbagent.LoadBalancingRule{
				{
					FrontendPort: 81,
					BackendPort:  80,
					Protocol:     pbcom.Protocol_Tcp,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	_, err = c.CreateOrUpdate(ctx, lbr)
	defer func() {
		err = c.Delete(ctx, lbr)
		require.NoError(t, err)
	}()
	require.NoError(t, err)
}

func TestUpdate(t *testing.T) {

	c, err := getClient()
	require.NoError(t, err)

	lbr := []*pblbagent.LoadBalancer{
		{
			Name:       "xxx",
			Frontendip: "*",
			Backendips: []string{"10.0.1.1", "10.0.1.2"},
			Loadbalancingrules: []*pblbagent.LoadBalancingRule{
				{
					FrontendPort: 80,
					BackendPort:  80,
					Protocol:     pbcom.Protocol_Tcp,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	_, err = c.CreateOrUpdate(ctx, lbr)
	defer func() {
		err = c.Delete(ctx, lbr)
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	lbr[0].Backendips = append(lbr[0].Backendips, "10.0.1.3")
	_, err = c.CreateOrUpdate(ctx, lbr)
	require.NoError(t, err)

	//TODO: validate backend changed.
}

func TestGetConfig(t *testing.T) {

	c, err := getClient()
	require.NoError(t, err)

	lbr := []*pblbagent.LoadBalancer{
		{
			Name:       "xxx",
			Frontendip: "*",
			Backendips: []string{"10.0.1.1", "10.0.1.2"},
			Loadbalancingrules: []*pblbagent.LoadBalancingRule{
				{
					FrontendPort: 80,
					BackendPort:  80,
					Protocol:     pbcom.Protocol_Tcp,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	_, err = c.CreateOrUpdate(ctx, lbr)
	defer func() {
		err = c.Delete(ctx, lbr)
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	config, err := c.GetConfig(ctx, pblbagent.LoadBalancerType_Haproxy)
	require.NoError(t, err)

	//TODO: validate config
	t.Logf("Successfully got config: %s\n", config)

	config, err = c.GetConfig(ctx, pblbagent.LoadBalancerType_Keepalived)
	require.NoError(t, err)

	//TODO: validate configs
	t.Logf("Successfully got config: %s\n", config)

}

func TestGet(t *testing.T) {
	c, err := getClient()
	require.NoError(t, err)

	lbr := []*pblbagent.LoadBalancer{
		{
			Name:       "xxx",
			Frontendip: "*",
			Backendips: []string{"10.0.1.1", "10.0.1.2"},
			Loadbalancingrules: []*pblbagent.LoadBalancingRule{
				{
					FrontendPort: 80,
					BackendPort:  80,
					Protocol:     pbcom.Protocol_Tcp,
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	_, err = c.CreateOrUpdate(ctx, lbr)
	defer func() {
		err = c.Delete(ctx, lbr)
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	lbs, err := c.Get(ctx, nil)
	require.NoError(t, err)
	t.Logf("Successfully got: %v\n", lbs)
}

func TestHealth(t *testing.T) {

	healthclient, err := getHealthClient()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	request := &admin_pb.HealthRequest{TimeoutSeconds: 30}
	_, err = healthclient.CheckHealth(ctx, request)
	require.NoError(t, err)
}
func TestMain(m *testing.M) {
	flag.Parse()
	if lbagentip == nil || *lbagentip == "" {
		fmt.Printf("Please specify lbagentip\n")
		os.Exit(-1)
	}
	viper.SetDefault("debug", true)
	os.Exit(m.Run())
}
