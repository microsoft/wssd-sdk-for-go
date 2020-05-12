// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package resize

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/storage/virtualharddisk"

	wssdcommon "github.com/microsoft/moc/common"
)

type flags struct {
	// Name of the Virtual Harddisk to get
	Name      string
	SizeBytes int
	Output    string
	Query     string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "resize",
		Aliases: []string{"get"},
		Short:   "resize to a specific Virtual Machine",
		Long:    "detach from a specific Virtual Machine ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine resource")
	cmd.MarkFlagRequired("name")
	cmd.Flags().IntVar(&flags.SizeBytes, "size-bytes", 0, "name of the virtual machine resource")
	cmd.MarkFlagRequired("size-bytes")

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

	vhd := (*vhds)[0]
	size := int64(flags.SizeBytes)
	vhd.DiskSizeBytes = &size
	_, err = vhdclient.CreateOrUpdate(ctx, container, flags.Name, &vhd)
	if err != nil {
		return err
	}

	return nil
}
