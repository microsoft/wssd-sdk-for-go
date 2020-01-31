// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package reimage

import (
	"context"
	"fmt"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachinescaleset"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "reimage",
		Short: "Reimage Virtual machines within the specified Virtual Machine Scale Set",
		Long:  "Reimage Virtual machines within the specified Virtual Machine Scale Set",
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
		return fmt.Errorf("Unable to find Virtual Machine Scale Set [%s]", flags.Name)
	}

	panic("vmss reimage not implemented")

	return nil

}
