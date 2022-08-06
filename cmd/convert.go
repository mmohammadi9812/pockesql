// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"io/ioutil"
	"log"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/niklasfasching/go-org/org"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use: "convert format [filename]",
	Aliases: []string{"to"},
	Short: "convert sqlite data to other supported formats",
	Long: `
This command converts your saved pocket data (from sqlite) to other shapes
As of now, only org mode is supported
`,
	Run: func(command *cobra.Command, args []string) {
		if len(args) < 1 {
			command.Usage()
			return
		}
		if args[0] == "org" {
			filename := "pocket.org"
			if len(args) == 2 {
				filename = args[1]
			}
			ToOrg(filename)
		} else {
			log.Fatalf("Converting to %s is not supported yet\n", args[0])
		}
	},
}

func tagToString(tags []src.Tag) []string {
	if len(tags) == 0 {
		return []string{}
	}

	out := []string{}
	for _, tag := range tags {
		out = append(out, tag.Tag)
	}

	return out
}

func getPriorities(items []src.PocketItem) map[uint]org.Text {
	var prts map[uint]org.Text = make(map[uint]org.Text)
	for _, item := range items {
		if item.TimeRead == 0 {
			prt := " [#E] "
			si := item.SortId / 100
			switch si {
			case 0:
				prt = " [#A] "
			case 1:
				prt = " [#B] "
			case 2:
				prt = " [#C] "
			case 3:
				prt = " [#D] "
			}
			prts[item.ID] = org.Text{Content: prt}
		}
	}
	return prts
}

func ToOrg(filename string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	items, err := src.ReadPocketItems()
	if err != nil {
		log.Fatal(err)
	}
	prts := getPriorities(items)

	writer := org.NewOrgWriter()

	for _, item := range items {
		tags, err := src.ReadItemTags(item)
		if err != nil {
			// TODO: decide whether it should be skipped or paniced
			continue
		}
		orgTags := tagToString(tags)
		headline := org.Headline{
			Lvl: 1,
			Title: []org.Node{org.RegularLink{
				URL: item.GivenUrl,
				Description: []org.Node{org.Text{Content: item.GivenTitle, IsRaw: true}},
			}},
			Tags: orgTags,
			Children: []org.Node{org.Text{Content: item.Excerpt, IsRaw: true}},
		}
		if item.TimeRead == 0 {
			headline.Title = append([]org.Node{org.Text{Content: "TODO "}, prts[item.ID]}, headline.Title...)
		}
		writer.WriteHeadline(headline)
	}

	err = ioutil.WriteFile(filename, []byte(writer.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}
}