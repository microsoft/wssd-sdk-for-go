// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package show

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "show",
		Short: "Get all Virtual Network(s)",
		Long:  "Get all Virtual Network(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network resource")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {

	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	networks, err := vnetclient.Get(ctx, group, flags.Name)
	if err != nil {
		return err
	}
	if networks == nil || len(*networks) == 0 {
		return fmt.Errorf("Unable to find Virtual Network [%s]", flags.Name)
	}

	config.PrintYAML(networks)

	return nil
}
