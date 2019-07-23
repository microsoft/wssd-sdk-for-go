// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachinescaleset

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachinescaleset"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualmachinescaleset",
		Aliases: []string{"vmss"},
		Short:   "Get a specific/all Virtual Machine Scale Set(s)",
		Long:    "Get a specific/all Virtual Machine Scale Set(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine scale set resource")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	client, err := virtualmachinescaleset.NewVirtualMachineScaleSetClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vmss, err := client.Get(ctx, flags.Name)
	if err != nil {
		return err
	}
	// If a single VM was requested
	if len(flags.Name) > 0 {
		if vmss == nil || len(*vmss) == 0 {
			return fmt.Errorf("Unable to find Virtual Machine Scale Set [%s]", flags.Name)
		}

	} else {
		if vmss == nil || len(*vmss) == 0 {
			fmt.Println("No VirtualMachineScaleSet Resources")
			// Not an error
			return nil
		}
	}

	virtualmachinescaleset.PrintList(vmss)
	return nil

}
