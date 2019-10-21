// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"
	//"time"
	"fmt"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/identity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Short: "Create a identity",
		Long:  "Create a identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the identity resource(s), comma separated")
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	identityClient, err := identity.NewIdentityClient(server, authorizer)
	if err != nil {
		return err
	}

	var identityName string
	var idConfig *security.Identity

	if flags.FilePath != "" {
		idConfig = &security.Identity{}
		err = config.LoadYAMLFile(flags.FilePath, idConfig)
		if err != nil {
			return err
		}
		identityName = *idConfig.Name
	} else {
		if flags.Name == "" {
			return fmt.Errorf("Error: must specify --name or --config")
		}
		idConfig = nil
		identityName = flags.Name
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	_, err = identityClient.CreateOrUpdate(ctx, group, identityName, idConfig)
	if err != nil {
		return err
	}

	return nil
}
