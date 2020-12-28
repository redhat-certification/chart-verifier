/*
 * Copyright (C) 28/12/2020, 16:06 igors
 * This file is part of helmcertifier.
 *
 * helmcertifier is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * helmcertifier is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with helmcertifier.  If not, see <http://www.gnu.org/licenses/>.
 */

package cmd

import (
	"github.com/spf13/cobra"
	"helmcertifier/pkg/helmcertifier"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	// allChecks contains all available checks to be executed by the program.
	allChecks []string
	// chartUri contains the chart location as informed by the user; should accept anything that Helm understands as a Chart
	// URI.
	chartUri string
	// onlyChecks are the checks that should be performed, after the command initialization has happened.
	onlyChecks []string
	// exceptChecks are the checks that should not be performed.
	exceptChecks []string
)

func NewCertifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certify",
		Args:  cobra.NoArgs,
		Short: "Certifies a Helm chart by checking some of its characteristics",
		RunE: func(cmd *cobra.Command, args []string) error {

			b := helmcertifier.NewCertificationBuilder()
			result, err := b.SetUri(chartUri).Build()
			if err != nil {
				return err
			}
			cmd.Println(result)

			return nil
		},
	}

	cmd.Flags().StringVarP(&chartUri, "uri", "u", "", "uri of the Chart being certified")
	_ = cmd.MarkFlagRequired("uri")

	cmd.Flags().StringSliceVarP(&onlyChecks, "only", "o", []string{}, "only the informed checks will be performed")

	cmd.Flags().StringSliceVarP(&exceptChecks, "except", "e", []string{}, "all available checks except those informed will be performed")

	return cmd
}

// certifyCmd represents the lint command
var certifyCmd = NewCertifyCmd()

func init() {
	rootCmd.AddCommand(certifyCmd)
}
