// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package update

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachinescaleset"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "update",
		Short: "Update the specified Virtual Machine Scale Set",
		Long:  "Update the specified Virtual Machine Scale Set",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine scale set resource")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}
	
	client, err := virtualmachinescaleset.NewVirtualMachineScaleSetClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	_, err = client.Get(ctx, group, flags.Name)
	if err != nil {
		return err
	}

	_, err = client.CreateOrUpdate(ctx, group, flags.Name, nil)
	if err != nil {
		return err
	}

	return nil

}
