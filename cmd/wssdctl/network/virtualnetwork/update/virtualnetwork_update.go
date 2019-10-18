// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package update

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
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
		Use:     "update",
		Aliases: []string{"vnet"},
		Short:   "update a Virtual Network",
		Long:    "update a Virtual Network",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network resource(s), comma separated")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {

	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	vnetclient, err := virtualnetwork.NewVirtualNetworkClient(server, authorizer)
	if err != nil {
		return err
	}

	// Wait up to one minute for network creation
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	panic("vnet update not implemented")

	_, err = vnetclient.CreateOrUpdate(ctx, group, flags.Name, nil)
	if err != nil {
		return err
	}

	return nil
}
