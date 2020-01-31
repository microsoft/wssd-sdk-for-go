// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package open_port

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name     string
	FilePath string
	Port     int
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "open-port",
		Short: "opens a port on the Virtual Machine",
		Long:  "opens a port on the Virtual Machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine")
	cmd.MarkFlagRequired("name")
	cmd.Flags().IntVar(&flags.Port, "port", 0, "port to open")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	_, err = virtualmachine.NewVirtualMachineClient(server, authorizer)
	if err != nil {
		return err
	}
	_, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	panic("port open not implemented")

	return nil
}
