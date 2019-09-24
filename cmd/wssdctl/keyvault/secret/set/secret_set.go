// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package set

import (
	"context"
	//"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/microsoft/wssd-sdk-for-go/services/keyvault"
	"github.com/microsoft/wssd-sdk-for-go/services/keyvault/secret"
	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name      string
	FilePath  string
	Value     string
	VaultName string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "set",
		Short: "set a secret",
		Long:  "set a secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the secret, comma separated")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&flags.Value, "value", "", "name of the secret, comma separated")
	cmd.MarkFlagRequired("value")

	cmd.Flags().StringVar(&flags.VaultName, "vault-name", "", "name of the secret, comma separated")
	cmd.MarkFlagRequired("vault-name")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")

	server := viper.GetString("server")
	secretClient, err := secret.NewSecretClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	secretName := flags.Name
 

	_, err = secretClient.CreateOrUpdate(ctx, group, secretName, &keyvault.Secret{
		BaseProperties: keyvault.BaseProperties{
			Name: &flags.Name,
		},
		Value : &flags.Value,
		VaultName: &flags.VaultName,
	})
	if err != nil {
		return err
	}

	return nil
}
