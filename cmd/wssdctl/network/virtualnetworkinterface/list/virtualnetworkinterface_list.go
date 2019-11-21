// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetworkinterface"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
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
		Use:   "list",
		Short: "Get a all Virtual Network Interface(s)",
		Long:  "Get a all Virtual Network Interface(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Output, "output", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVar(&flags.Query, "query", "", "Output Format")
	return cmd
}

func runE(flags *flags) error {

	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	vnetInterfaceClient, err := virtualnetworkinterface.NewVirtualNetworkInterfaceClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	networkInterfaces, err := vnetInterfaceClient.Get(ctx, group, "")
	if err != nil {
		return err
	}

	if networkInterfaces == nil || len(*networkInterfaces) == 0 {
		fmt.Println("No Virtual Network Interface Resources")
		// Not an error
		return nil
	}

	config.PrintFormatList(*networkInterfaces, flags.Query, flags.Output)

	return nil
}
