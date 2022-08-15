// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package convert

import (
	"fmt"
	"log"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/xuri/excelize/v2"
)

type XlslCell struct {
	file *excelize.File
	item src.PocketItem
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
	items, err := src.ReadPocketItems()
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