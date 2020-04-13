// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/microsoft/moc/pkg/auth"
	pbcom "github.com/microsoft/moc/rpc/common"
	pblbagent "github.com/microsoft/moc/rpc/lbagent"
	"github.com/microsoft/wssd-sdk-for-go/pkg/lbagentclient"
)

func testAdd(lbagentip string) error {

	authorizer, err := auth.NewAuthorizerFromEnvironment(lbagentip)
	if err != nil {
		return err
	}

	c, err := lbagentclient.NewLoadBalancerAgentClient(lbagentip, authorizer)
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}
	fmt.Printf("Successfully added: %v\n", lbr)

	err = c.Delete(ctx, lbr)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully deleted: %v\n", lbr)

	return nil
}

func testUpdate(lbagentip string) error {
	authorizer, err := auth.NewAuthorizerFromEnvironment(lbagentip)
	if err != nil {
		return err
	}

	c, err := lbagentclient.NewLoadBalancerAgentClient(lbagentip, authorizer)
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}
	fmt.Printf("Successfully added: %v\n", lbr)

	lbr[0].Backendips = append(lbr[0].Backendips, "10.0.1.3")
	_, err = c.CreateOrUpdate(ctx, lbr)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully updated: %v\n", lbr)

	err = c.Delete(ctx, lbr)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully deleted: %v\n", lbr)

	return nil
}

func testGet(lbagentip string) error {
	authorizer, err := auth.NewAuthorizerFromEnvironment(lbagentip)
	if err != nil {
		return err
	}

	c, err := lbagentclient.NewLoadBalancerAgentClient(lbagentip, authorizer)
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}
	fmt.Printf("Successfully added: %v\n", lbr)

	config, err := c.GetConfig(ctx, pblbagent.LoadBalancerType_Haproxy)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully got config: %s\n", config)

	config, err = c.GetConfig(ctx, pblbagent.LoadBalancerType_Keepalived)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully got config: %s\n", config)

	return nil
}

func main() {
	var lbagentip string
	flag.StringVar(&lbagentip, "lbagentip", "", "ip address of lbagent")
	flag.Parse()

	if lbagentip == "" {
		fmt.Printf("enter a lbagent ip\n")
		os.Exit(1)
	}

	err := testAdd(lbagentip)
	if err != nil {
		fmt.Printf("failed testAdd: %x\n", err)
	}

	err = testUpdate(lbagentip)
	if err != nil {
		fmt.Printf("failed testUpdate: %x\n", err)
	}

	err = testGet(lbagentip)
	if err != nil {
		fmt.Printf("failed testGet: %x\n", err)
	}

	//testDelete(lbagentip)
}
