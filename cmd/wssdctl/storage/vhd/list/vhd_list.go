// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/storage/virtualharddisk"
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
		Short: "list a specific vhd",
		Long:  "list a specific vhd ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the vhd")
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

	vhdClient, err := virtualharddisk.NewVirtualHardDiskClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	vhds, err := vhdClient.Get(ctx, container, flags.Name)
	if err != nil {
		return err
	}
	if vhds == nil || len(*vhds) == 0 {
		fmt.Println("No Virtual Hard Disk Resources")
		// Not an error
		return nil
	}

	config.PrintFormatList(*vhds, flags.Query, flags.Output)

	return nil
}
