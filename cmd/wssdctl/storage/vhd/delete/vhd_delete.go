// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
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
		Use:   "delete",
		Short: "delete a specific vhd",
		Long:  "delete a specific vhd ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the vhd")
	cmd.MarkFlagRequired("name")

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

	
	err = vhdClient.Delete(ctx, group, flags.Name)
	if err != nil {
		return err
	}
	
	return nil
}
