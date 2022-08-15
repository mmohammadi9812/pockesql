// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~mmohammadi9812/pockesql/convert"
	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:     "convert format [filename]",
	Aliases: []string{"to"},
	Short:   "convert sqlite data to other supported formats",
	Long: `
This command converts your saved pocket data (from sqlite) to other shapes
As of now, only org mode is supported
`,
	Run: func(command *cobra.Command, args []string) {
		if len(args) < 1 {
			command.Usage()
			return
		}
		switch args[0] {
		case "org":
			filename := "pocket.org"
			if len(args) == 2 {
				filename = args[1]
			}
			convert.ToOrg(filename)
		case "excel":
			filename := "pocket.xlsx"
			if len(args) == 2 {
				filename = args[1]
				if !isValid(filename) {
					log.Fatalf("expected valid file name, got %s\n", filename)
				}
			}
			convert.ToXlsl(filename)
		default:
			log.Fatalf("Converting to %s is not supported yet\n", args[0])
		}
	},
	PreRun: func(_ *cobra.Command, _ []string) {
		if _, err := os.Stat(src.GetDatabaseName()); errors.Is(err, os.ErrNotExist) {
			FetchCmd()
		}
	},
}

func isValid(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case "xlam", "xlsm", "xlsx", "xltm", "xltx":
		return true

	default:
		return false
	}
}
