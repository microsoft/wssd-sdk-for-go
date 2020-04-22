// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	//"time"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wssdcommon "github.com/microsoft/moc/common"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
)

type flags struct {
	Name     string
	FilePath string
	Output   string
	Query    string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "list",
		Short: "list a keyvault",
		Long:  "list a keyvault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVarP(&flags.Output, "output", "o", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVarP(&flags.Query, "query", "q", "", "Output Format")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vaultClient, err := keyvault.NewKeyVaultClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	keyvaults, err := vaultClient.Get(ctx, group, flags.Name)
	if err != nil {
		return err
	}

	if keyvaults == nil || len(*keyvaults) == 0 {
		fmt.Println("No Key Vaults")
		// Not an error
		return nil
	}

	config.PrintFormatList(*keyvaults, flags.Query, flags.Output)
	return nil
}
