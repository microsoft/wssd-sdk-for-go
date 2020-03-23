// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package login

import (
	"context"
	//"time"
	"fmt"

	wssdcommon "github.com/microsoft/moc/common"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/authentication"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flags struct {
	Identity      bool
	LoginFilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "login",
		Short: "Login an identity",
		Long:  "Login an identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().BoolVar(&flags.Identity, "identity", true, "Uses Managed Identity to Login User")
	cmd.MarkFlagRequired("identity")
	cmd.Flags().StringVar(&flags.LoginFilePath, "loginpath", "", "Path for the NodeAgent Certificate")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")
	server := viper.GetString("server")

	loginconfig := auth.LoginConfig{}

	err := config.LoadYAMLFile(flags.LoginFilePath, &loginconfig)
	if err != nil {
		return err
	}

	authorizer, err := auth.NewAuthorizerForAuth(loginconfig.Token, loginconfig.Certificate, server)
	if err != nil {
		return err
	}

	authenticationClient, err := authentication.NewAuthenticationClient(server, authorizer)
	if err != nil {
		return err
	}

	if !flags.Identity {
		return fmt.Errorf("Not Supported")
	}

	clientCert, accessFile, err := auth.GenerateClientKey(loginconfig)
	if err != nil {
		return err
	}

	id := security.Identity{
		Name:        &loginconfig.Name,
		Certificate: &clientCert,
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	_, err = authenticationClient.Login(ctx, group, &id)
	if err != nil {
		return err
	}

	auth.PrintAccessFile(accessFile)

	return nil
}
