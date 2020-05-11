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
	"github.com/microsoft/wssd-sdk-for-go/pkg/lbagentclient"
)

func get(lbagentip string) error {
	authorizer, err := auth.NewAuthorizerFromEnvironment(lbagentip)
	if err != nil {
		return err
	}

	c, err := lbagentclient.NewLoadBalancerAgentClient(lbagentip, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	lbs, err := c.Get(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully got: %v\n", lbs)

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

	err := get(lbagentip)
	if err != nil {
		fmt.Printf("failed to get all lbs: %v\n", err)
	}

}
