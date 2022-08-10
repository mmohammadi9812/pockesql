// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/niklasfasching/go-org/org"
	"github.com/xuri/excelize/v2"
)

func tagToString(tags []Tag) []string {
	if len(tags) == 0 {
		return []string{}
	}

	out := []string{}
	for _, tag := range tags {
		out = append(out, tag.Tag)
	}

	return out
}

func getPriority(item PocketItem) org.Text {
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

func itemToHeadline(item PocketItem) (org.Headline, error) {
	tags, err := ReadItemTags(item)
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
	items, err := ReadPocketItems()
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

	err = ioutil.WriteFile(filename, []byte(writer.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

type XlslCell struct {
	file *excelize.File
	item PocketItem
	err error
}

func (c *XlslCell) Set(sheet string, row int) *XlslCell {
	// TODO: make this function more flexible
	if c.err != nil {
		return c
	}
	if c.err = c.file.SetCellValue(sheet, fmt.Sprintf("A%d", row), c.item.ID); c.err != nil {
		return c
	}
	if c.err = c.file.SetCellValue(sheet, fmt.Sprintf("B%d", row), c.item.GivenTitle); c.err != nil {
		return c
	}
	if c.err = c.file.SetCellValue(sheet, fmt.Sprintf("C%d", row), c.item.GivenUrl); c.err != nil {
		return c
	}
	if c.err = c.file.SetCellValue(sheet, fmt.Sprintf("D%d", row), c.item.Excerpt); c.err != nil {
		return c
	}
	if c.err = c.file.SetCellValue(sheet, fmt.Sprintf("E%d", row), c.item.WordCount); c.err != nil {
		return c
	}
	if c.err = c.file.SetCellValue(sheet, fmt.Sprintf("F%d", row), c.item.TimeRead); c.err != nil {
		return c
	}


	return c
}

func ToXlsl(filename string) {
	// TODO: set style
	items, err := ReadPocketItems()
	if err != nil {
		log.Fatal(err)
	}
	f := excelize.NewFile()

	defer f.Close()

	sheet := "Pocket"
	indx := f.NewSheet(sheet)
	f.SetActiveSheet(indx)

	if err = f.SetColWidth(sheet, "B", "D", 50); err != nil {
		log.Fatal(err)
	}

	fields := []string{
		"ID", "Title", "Url",
		"Excerpt", "WordCount", "TimeRead"}

	columns := []string{"A", "B", "C", "D", "E", "F"}

	for i, header := range fields {
		if err := f.SetCellValue(sheet, fmt.Sprintf("%s1", columns[i]), header); err != nil {
			log.Fatal(err)
		}
	}

	for i, item := range items {
		cell := &XlslCell{
			file: f,
			item: item,
		}
		if cell = cell.Set(sheet, i+2); cell.err != nil {
			log.Fatal(err)
		}
	}

	if err = f.SaveAs(filename); err != nil {
		log.Fatal(err)
	}
}