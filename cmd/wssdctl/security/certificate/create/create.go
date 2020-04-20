// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wssdcommon "github.com/microsoft/moc/common"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/certificate"
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
		Short: "Create a certificate ",
		Long:  "Create a certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	certclient, err := certificate.NewCertificateClient(server, authorizer)
	if err != nil {
		return err
	}

	certconfig := security.Certificate{}
	err = config.LoadYAMLFile(flags.FilePath, &certconfig)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	if certconfig.Name == nil {
		return errors.Wrapf(errors.InvalidInput, "The YAML is missing the 'Name' element")
	}

	lbs, err := certclient.CreateOrUpdate(ctx, group, *(certconfig.Name), &certconfig)
	if err != nil {
		return err
	}
	config.PrintYAMLList(*lbs)

	return nil
}
