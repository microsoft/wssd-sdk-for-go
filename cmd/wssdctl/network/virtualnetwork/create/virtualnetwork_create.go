// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetwork"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "create",
		Aliases: []string{"vnet"},
		Short:   "Create a Virtual Network",
		Long:    "Create a Virtual Network",
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
	group := viper.GetString("group")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server, authorizer)
	if err != nil {
		return err
	}

	vnetconfig := network.VirtualNetwork{}
	err = config.LoadYAMLFile(flags.FilePath, &vnetconfig)
	if err != nil {
		return err
	}

	// Wait up to one minute for network creation
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if vnetconfig.Name == nil {
		return errors.Wrapf(errors.InvalidInput, "The YAML is missing the 'Name' element")
	}

	_, err = vnetclient.CreateOrUpdate(ctx, group, *(vnetconfig.Name), &vnetconfig)
	if err != nil {
		return err
	}

	return nil
}
