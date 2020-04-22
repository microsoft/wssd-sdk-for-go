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
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"

	wssdcommon "github.com/microsoft/moc/common"
)

type flags struct {
	// Name of the Virtual Machine to get
	Name   string
	Output string
	Query  string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "show",
		Aliases: []string{"get"},
		Short:   "Get a specific Virtual Machine(s)",
		Long:    "Get a specific Virtual Machine(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine resource")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&flags.Output, "output", "o", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVarP(&flags.Query, "query", "q", "", "Output Format")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vmclient, err := virtualmachine.NewVirtualMachineClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vms, err := vmclient.Get(ctx, group, flags.Name)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return fmt.Errorf("Unable to find Virtual Machine [%s]", flags.Name)
	}

	config.PrintFormatList(*vms, flags.Query, flags.Output)

	return nil
}
