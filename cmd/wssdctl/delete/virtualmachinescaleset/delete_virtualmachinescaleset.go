// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachinescaleset

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachinescaleset"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	// Name of the Virtual Machine to get
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualmachinescaleset",
		Aliases: []string{"vmss"},
		Short:   "Delete a specific VirtualMachine Scale Set",
		Long:    "Delete a specific VirtualMachine Scale Set",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name(s) of the virtual machine scale set")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	vmclient, err := virtualmachinescaleset.NewVirtualMachineScaleSetClient(server)
	if err != nil {
		return err
	}

	vmName := flags.Name
	vmId := ""
	if len(vmName) == 0 {
		config := viper.GetString("config")
		vmconfig, err := virtualmachinescaleset.LoadConfig(config)
		if err != nil {
			return err
		}
		vmName = *(vmconfig.Name)
		vmId = *(vmconfig.ID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	err = vmclient.Delete(ctx, vmId, vmId)
	if err != nil {
		return err
	}

	return nil
}
