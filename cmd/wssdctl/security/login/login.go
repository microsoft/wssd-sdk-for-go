// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package login

import (
	"context"
	//"time"
	"fmt"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
	//"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	//"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/authentication"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flags struct {
	Identity bool
	User string
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
	cmd.Flags().StringVar(&flags.User, "user", "", "User Name for Managed Identity Login")
	cmd.MarkFlagRequired("user")


	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
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

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	token, err := authenticationClient.Login(ctx, group, flags.User)
	if err != nil {
		return err
	}

	auth.SaveToken(*token)

	return nil
}
