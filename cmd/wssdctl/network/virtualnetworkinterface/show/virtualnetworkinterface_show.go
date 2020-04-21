// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package show

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetworkinterface"

	wssdcommon "github.com/microsoft/moc/common"
)

type flags struct {
	Name   string
	Output string
	Query  string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "show",
		Short: "Get a all Virtual Network Interface(s)",
		Long:  "Get a all Virtual Network Interface(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network interface resource(s), comma separated")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&flags.Output, "output", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVar(&flags.Query, "query", "", "Output Format")

	return cmd
}

func runE(flags *flags) error {

	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vnetInterfaceClient, err := virtualnetworkinterface.NewVirtualNetworkInterfaceClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	networkInterfaces, err := vnetInterfaceClient.Get(ctx, group, flags.Name)
	if err != nil {
		return err
	}
	if networkInterfaces == nil || len(*networkInterfaces) == 0 {
		return fmt.Errorf("Unable to find Virtual Network Interface [%s]", flags.Name)
	}

	config.PrintFormatList(*networkInterfaces, flags.Query, flags.Output)

	return nil
}
