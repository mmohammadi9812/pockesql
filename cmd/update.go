// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"strconv"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use: "update",
	Short: "Update already present database with new contents",
	Run: func(_ *cobra.Command, _ []string) {
		UpdateCmd()
	},
}


func UpdateCmd() {
	auth, err := src.ReadAuth()
	if err != nil {
		log.Fatal(err)
	}

	pocketTotal, err := src.TotalItems(auth)
	if err != nil {
		log.Fatal(err)
	}

	currentItems, err := src.ReadPocketItems()
	if err != nil {
		log.Fatal(err)
	}
	currentTotal := int64(len(currentItems))

	if currentTotal == pocketTotal {
		return
	} else if currentTotal < pocketTotal {
		lastItem := currentItems[0]
		since := strconv.FormatInt(lastItem.CreatedAt.Unix(), 10)
		offset := 0
		bar := progressbar.Default(pocketTotal - currentTotal)

		// FIXME: this code is duplicate of code in fetch
		for {
			raw, err := fetch(FUrlValues{
				Auth:   auth,
				Offset: offset,
				Since: since,
			})
			if err != nil {
				log.Fatal(err)
			}

			if raw == nil && err == nil {
				break
			}

			n, err := src.UpdateItems(raw)
			if err != nil || n < 0 {
				log.Fatalf("Saving fetch items failed [%d]\n%v", n, err)
			}

			offset += PAGE_SIZE

			bar.Add(n)
		}
	} else {
		// FIXME: do something if currentTotal > pocketTotal
	}
}