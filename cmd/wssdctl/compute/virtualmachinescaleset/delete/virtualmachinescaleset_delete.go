// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
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
		Args:  cobra.NoArgs,
		Use:   "delete",
		Short: "Delete a specific VirtualMachine Scale Set",
		Long:  "Delete a specific VirtualMachine Scale Set",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name(s) of the virtual machine scale set")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")
	server := viper.GetString("server")
	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}
	vmclient, err := virtualmachinescaleset.NewVirtualMachineScaleSetClient(server, authorizer)
	if err != nil {
		return err
	}

	vmName := flags.Name
	if len(vmName) == 0 {
		configPath := viper.GetString("config")
		vmconfig := compute.VirtualMachineScaleSet{}
		err = config.LoadYAMLFile(configPath, &vmconfig)
		if err != nil {
			return err
		}
		vmName = *(vmconfig.Name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	err = vmclient.Delete(ctx, group, vmName)
	if err != nil {
		return err
	}

	return nil
}
