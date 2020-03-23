// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachinescaleset"

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
		Use:   "list",
		Short: "Get all Virtual Machine Scale Set(s)",
		Long:  "Get all Virtual Machine Scale Set(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.Output, "output", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVar(&flags.Query, "query", "", "Output Format")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}
	client, err := virtualmachinescaleset.NewVirtualMachineScaleSetClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vmss, err := client.Get(ctx, group, flags.Name)
	if err != nil {
		return err
	}
	if vmss == nil || len(*vmss) == 0 {
		fmt.Println("No VirtualMachineScaleSet Resources")
		// Not an error
		return nil
	}

	config.PrintFormatList(*vmss, flags.Query, flags.Output)
	return nil

}
