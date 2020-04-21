// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"github.com/spf13/cobra"
)

type flags struct {
	Name     string
	FilePath string
	Output string
	Query  string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "list",
		Aliases: []string{"vm"},
		Short:   "list Containers",
		Long:    "list Containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")
	cmd.Flags().StringVarP(&flags.Output, "output", "o", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVarP(&flags.Query, "query", "q", "", "Output Format")

	return cmd
}

func runE(flags *flags) error {
	panic("container list not implemented")

	return nil
}
