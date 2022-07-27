// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SqliteConn(filename string) (*gorm.DB, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filename = filepath.Join(pwd, filename)
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getAuthors(itemMap map[string]interface{}) ([]Author, error) {
	var (
		result     []Author
		rawAuthors map[string]map[string]string
		ok         bool
	)
	if rawAuthors, ok = itemMap["authors"].(map[string]map[string]string); !ok {
		return nil, fmt.Errorf("authors not found")
	}
	for _, ra := range rawAuthors {
		var author Author
		authorDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName: "json",
			Result:  &author,
		})
		if err != nil {
			return nil, err
		}
		authorDecoder.Decode(ra)
		result = append(result, author)
	}

	return result, nil
}

func SaveItems(items []map[string]interface{}) (int, error) {
	db, err := SqliteConn("pocket.sqlite3")
	if err != nil {
		return -1, err
	}

	if err = db.AutoMigrate(&PocketItem{}, &Author{}); err != nil {
		return -2, err
	}

	for i, itemMap := range items {
		itemMap = Transform(itemMap)

		pocketItem, err := DecodeStruct(itemMap)
		if err != nil {
			return i-1, nil
		}
		authors, err := getAuthors(itemMap)
		if err != nil {
			return i-1, nil
		}
		pocketItem.Authors = authors

		if err = db.Create(&pocketItem).Error; err != nil {
			return i-1, err
		}
	}

	return len(items), nil
}
