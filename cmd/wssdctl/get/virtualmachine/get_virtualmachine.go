// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachine

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"

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
		Use:     "virtualmachine",
		Aliases: []string{"vm"},
		Short:   "Get a specific/all Virtual Machine(s)",
		Long:    "Get a specific/all Virtual Machine(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine resource")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	vmclient, err := virtualmachine.NewVirtualMachineClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vms, err := vmclient.Get(ctx, flags.Name)
	if err != nil {
		return err
	}
	// If a single VM was requested
	if len(flags.Name) > 0 {
		if vms == nil || len(*vms) == 0 {
			return fmt.Errorf("Unable to find Virtual Machine [%s]", flags.Name)
		}

	} else {
		if vms == nil || len(*vms) == 0 {
			fmt.Println("No VirtualMachine Resources")
			// Not an error
			return nil
		}
	}

	virtualmachine.PrintList(vms)

	return nil
}
