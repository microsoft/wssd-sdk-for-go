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
	"github.com/microsoft/wssd-sdk-for-go/services/storage/virtualharddisk"

	wssdcommon "github.com/microsoft/moc/common"
)

type flags struct {
	// Name of the Virtual Harddisk to get
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
		Short:   "Get a specific Virtual Harddisk(s)",
		Long:    "Get a specific Virtual Harddisk(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine resource")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&flags.Output, "output", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVar(&flags.Query, "query", "", "Output Format")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	container := viper.GetString("container")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vhdclient, err := virtualharddisk.NewVirtualHardDiskClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vhds, err := vhdclient.Get(ctx, container, flags.Name)
	if err != nil {
		return err
	}
	if vhds == nil || len(*vhds) == 0 {
		return fmt.Errorf("Unable to find Virtual Harddisk [%s]", flags.Name)
	}

	config.PrintFormatList(*vhds, flags.Query, flags.Output)

	return nil
}
