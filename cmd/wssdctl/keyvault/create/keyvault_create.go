// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"
	//"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/keyvault/simplevault"
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
		Short: "Create a simplevault",
		Long:  "Create a simplevault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the keyvault resource(s), comma separated")
	cmd.MarkFlagRequired("name")

	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")

	server := viper.GetString("server")
	vaultClient, err := simplevault.NewSimpleVaultClient(server)
	if err != nil {
		return err
	}

	config := flags.FilePath
	kvConfig, err := simplevault.LoadConfig(config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vaultName := flags.Name

	_, err = vaultClient.CreateOrUpdate(ctx, group, vaultName, kvConfig)
	if err != nil {
		return err
	}

	return nil
}
