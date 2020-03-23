// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package restart

import (
	"context"

	log "k8s.io/klog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"

	wssdcommon "github.com/microsoft/moc/common"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "restart",
		Short: "restart a Virtual Machine",
		Long:  "restart a Virtual Machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	vmclient, err := virtualmachine.NewVirtualMachineClient(server, authorizer)
	if err != nil {
		return err
	}
	config := flags.FilePath
	vmconfig, err := virtualmachine.LoadConfig(config)
	if err != nil {
		return err
	}

	log.Infof("Loaded Configuration [%s]", vmconfig)
	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	panic("vm restart open not implemented")

	return nil
}
