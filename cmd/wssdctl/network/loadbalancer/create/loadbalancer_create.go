// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssd-sdk-for-go/services/network/loadbalancer"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Create a LoadBalancer ",
		Long:  "Create a LoadBalancer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	lbclient, err := loadbalancer.NewLoadBalancerClient(server, authorizer)
	if err != nil {
		return err
	}

	lbconfig := network.LoadBalancer{}
	err = config.LoadYAMLFile(flags.FilePath, &lbconfig)
	if err != nil {
		return err
	}

	// Wait up to one minute for network creation
	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	_, err = lbclient.CreateOrUpdate(ctx, group, *(lbconfig.Name), &lbconfig)
	if err != nil {
		return err
	}

	return nil
}
