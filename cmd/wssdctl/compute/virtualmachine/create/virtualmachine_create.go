// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"

	log "k8s.io/klog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Create a Virtual Machine",
		Long:  "Create a Virtual Machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	vmclient, err := virtualmachine.NewVirtualMachineClient(server)
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

	_, err = vmclient.CreateOrUpdate(ctx, *(vmconfig.Name), *(vmconfig.ID), vmconfig)
	if err != nil {
		return err
	}

	return nil
}
