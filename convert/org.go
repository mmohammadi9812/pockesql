// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package convert

import (
	"log"
	"os"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/niklasfasching/go-org/org"
)

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

func getPriority(item src.PocketItem) org.Text {
	var prt string

	if item.TimeRead == 0 {
		prt = " [#E] "
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
	}
	return org.Text{Content: prt}
}

func itemToHeadline(item src.PocketItem) (org.Headline, error) {
	tags, err := src.ReadItemTags(item)
	if err != nil {
		return org.Headline{}, err
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
		headline.Title = append([]org.Node{org.Text{Content: "TODO "}, getPriority(item)}, headline.Title...)
	}

	return headline, nil
}


func ToOrg(filename string) {
	items, err := src.ReadPocketItems()
	if err != nil {
		log.Fatal(err)
	}
	writer := org.NewOrgWriter()

	for _, item := range items {
		headline, err := itemToHeadline(item)
		if err != nil {
			// TODO: decide whether it should be skipped or paniced
			continue
		}
		writer.WriteHeadline(headline)
	}

	err = os.WriteFile(filename, []byte(writer.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
