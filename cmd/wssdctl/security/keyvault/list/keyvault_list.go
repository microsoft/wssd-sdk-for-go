// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	//"time"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//"github.com/microsoft/wssd-sdk-for-go/services/keyvault"
	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
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
		Use:   "list",
		Short: "list a keyvault",
		Long:  "list a keyvault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

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

	keyvaults, err := vaultClient.Get(ctx, group, "")
	if err != nil {
		return err
	}

	if keyvaults == nil || len(*keyvaults) == 0 {
		fmt.Println("No Key Vaults")
		// Not an error
		return nil
	}

	config.PrintTable(keyvaults)

	return nil
}
