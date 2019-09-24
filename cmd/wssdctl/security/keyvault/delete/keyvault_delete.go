// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"context"
	//"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "delete",
		Short: "delete a keyvault",
		Long:  "delete a keyvault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the keyvault resource(s), comma separated")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")

	server := viper.GetString("server")
	vaultClient, err := keyvault.NewKeyVaultClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vaultName := flags.Name

	err = vaultClient.Delete(ctx, group, vaultName)
	if err != nil {
		return err
	}

	return nil
}
