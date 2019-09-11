// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"fmt"
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/storage/virtualharddisk"
)

type flags struct {
	Name string
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
	
	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")

	vhdClient, err := virtualharddisk.NewVirtualHardDiskClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	
	vhds, err := vhdClient.Get(ctx, group, flags.Name)
	if err != nil {
		return err
	}
	if vhds == nil || len(*vhds) == 0 {
		fmt.Println("No Virtual Hard Disk Resources")
		// Not an error
		return nil
	}

	virtualharddisk.PrintList(vhds)

	return nil
}
