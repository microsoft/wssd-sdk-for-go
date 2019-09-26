// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package show

import (
	"context"
	//"time"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault/secret"
	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name     string
	FilePath string
	VaultName  string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "show",
		Short: "show a secret",
		Long:  "show a secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the secret, comma separated")
	cmd.MarkFlagRequired("name")
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
	emptyValue := ""

	secrets, err := secretClient.Get(ctx, group, secretName, &keyvault.Secret{
		BaseProperties: security.BaseProperties{
			Name: &flags.Name,
		},
		Value : &emptyValue,
		VaultName: &flags.VaultName,
	})
	if err != nil {
		return err
	}
	
	if secrets == nil || len(*secrets) == 0 {
		fmt.Println("No Secrets")
		// Not an error
		return nil           
	}

	secret.PrintList(secrets)

	return nil
}
