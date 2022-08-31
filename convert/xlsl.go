// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package convert

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/xuri/excelize/v2"
)

var columnOrders = []string{
		"ID", "Title", "Url",
		"Excerpt", "WordCount",
		"TimeRead",
}

type XlslCell struct {
	file *excelize.File
	item src.PocketItem
	err error
}

func SetHeaderColumn(file *excelize.File, sheet string) error {
	for i, header := range columnOrders {
		cn := fmt.Sprintf("%s1", strconv.QuoteRune(rune('A' + i)))
		if err := file.SetCellValue(sheet, cn, header); err != nil {
			return err
		}
	}

	return nil
}

func (c *XlslCell) Set(sheet string, row int) *XlslCell {
	// TODO: make this function more flexible
	if c.err != nil {
		return c
	}

	for i, column := range columnOrders {
		cn := fmt.Sprintf("%s%v", strconv.QuoteRune(rune('A' + i)), row)
		v := reflect.ValueOf(c.item).FieldByName(column)
		c.err = c.file.SetCellValue(sheet, cn, v)
		if c.err != nil {
			return c
		}
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

	if err := SetHeaderColumn(f, sheet); err != nil {
		log.Fatal(err)
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