// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package detach

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"

	wssdcommon "github.com/microsoft/moc/common"
)

type flags struct {
	// Name of the Virtual Machine to get
	Name   string
	VMName string
	Output string
	Query  string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "detach",
		Aliases: []string{"get"},
		Short:   "detach data disk to specific Virtual Machine",
		Long:    "detach data disk to specific Virtual Machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine resource")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&flags.VMName, "vm-name", "", "name of the virtual machine resource")
	cmd.MarkFlagRequired("vm-name")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vmclient, err := virtualmachine.NewVirtualMachineClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vms, err := vmclient.Get(ctx, group, flags.VMName)
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		return errors.Wrapf(errors.NotFound, "Unable to find Virtual Machine [%s]", flags.VMName)
	}

	vm := (*vms)[0]
	for i, element := range *vm.StorageProfile.DataDisks {
		if *element.VhdName == flags.Name {
			*vm.StorageProfile.DataDisks = append((*vm.StorageProfile.DataDisks)[:i], (*vm.StorageProfile.DataDisks)[i+1:]...)
			break
		}
	}

	_, err = vmclient.CreateOrUpdate(ctx, group, flags.VMName, &vm)
	if err != nil {
		return err
	}

	return nil
}
