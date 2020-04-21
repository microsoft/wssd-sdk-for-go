// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package show

import (
	"github.com/spf13/cobra"
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
		Args:    cobra.NoArgs,
		Use:     "show",
		Aliases: []string{"vm"},
		Short:   "Create a Container",
		Long:    "Create a Container",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")
	cmd.Flags().StringVar(&flags.Output, "output", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVar(&flags.Query, "query", "", "Output Format")

	return cmd
}

func runE(flags *flags) error {
	panic("container show not implemented")

	return nil
}
