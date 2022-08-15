// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"strconv"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Aliases: []string{"u"},
	Short: "Update already present database with new contents",
	Run: func(_ *cobra.Command, _ []string) {
		UpdateCmd()
	},
}

func fetchKeys(args FetchOptions) ([]string, error) {
	items := make([]string, 1)
	for {
		raw, err := fetchPocket(args)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while trying to fetch items")
		}
		if raw == nil && err == nil {
			break
		}
		keys := maps.Keys(raw)
		items = append(items, keys...)
		args.Offset += PAGE_SIZE
	}
	return items, nil
}

func keysToIds(arr []string) map[uint]bool {
	out := make(map[uint]bool, len(arr))
	for _, s := range arr {
		u, e := strconv.ParseUint(s, 10, 64)
		if e != nil {
			continue
		}
		out[uint(u)] = true
	}
	return out
}

func UpdateCmd() {
	// TODO: long function, needs refactor
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

	since := strconv.FormatInt(currentItems[0].CreatedAt.Unix(), 10)
	args, err := createFetchOptions(auth)
	if err != nil {
		log.Fatal(err)
	}
	if currentTotal < pocketTotal {
		args.Since = since
		e := Save(args, pocketTotal - currentTotal)
		if e.Messaage != nil {
			log.Fatalf("\n\n[%d] An error occured:\n%v\n\n", e.StatusCode, e.Messaage)
		}
	} else if pocketTotal < currentTotal {
		keys, err := fetchKeys(args)
		if err != nil {
			log.Fatal(err)
		}
		remoteIds := keysToIds(keys)
		deletedIds := []uint{}
		for _, p := range currentItems {
			_, ok := remoteIds[p.ID]
			if !ok {
				deletedIds = append(deletedIds, p.ID)
			}
		}
		if err := src.DeleteIds(deletedIds); err != nil {
			log.Fatal(err)
		}
		log.Printf("Delted %d items from database\n", len(deletedIds))
	}
}
