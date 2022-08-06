/*
* Copyright 2022 Mohammad Mohamamdi. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pockesql",
	Short: "Fetch all your pocket entries to a sqlite database",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(authCmd)

	fetchCmd.Flags().StringVarP(&since, "since", "s", "", "fetch since time {yyyy-mm-dd hh:mm}")
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(convertCmd)
}


