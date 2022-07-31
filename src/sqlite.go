// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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

func Migrate(db *gorm.DB, items map[string]map[string]interface{}) error {
	var (
		pocketItem PocketItem
		authors    []Author
		tags       []Tag
		err        error
	)
	for iid, itemMap := range items {
		itemMap = Transform(itemMap)

		pocketItem, err = DecodeStruct(itemMap)
		if err != nil {
			return err
		}
		var itemId int
		if itemId, err = strconv.Atoi(iid); err != nil {
			return err
		}
		pocketItem.ID = uint(itemId)

		authors, err = getAuthors(itemMap)
		if err == nil {
			pocketItem.Authors = authors
		}
		for _, v := range pocketItem.Tags {
			v.PocketItemID = pocketItem.ID
			tags = append(tags, v)
		}
		break
	}
	return db.AutoMigrate(&tags, &pocketItem.DomainMetadata, &pocketItem.Authors, &pocketItem)
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

func getTags(itemMap map[string]interface{}) ([]Tag, error) {
	var (
		out []Tag
		r   map[string]map[string]string
		ok  bool
	)

	if r, ok = itemMap["tags"].(map[string]map[string]string); !ok {
		return nil, fmt.Errorf("tags not found")
	}

	for _, v := range r {
		var tag Tag
		tagDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &tag})
		if err != nil {
			return nil, err
		}
		tagDecoder.Decode(v)
		out = append(out, tag)
	}
	return out, nil
}

func SaveItems(items map[string]map[string]interface{}) (int, error) {
	db, err := SqliteConn("pocket.sqlite3")
	if err != nil {
		return -1, err
	}

	if err = Migrate(db, items); err != nil {
		return -2, err
	}

	i := 0
	for iid, itemMap := range items {
		itemMap = Transform(itemMap)

		pocketItem, err := DecodeStruct(itemMap)
		if err != nil {
			return i - 1, err
		}
		var itemId int
		if itemId, err = strconv.Atoi(iid); err != nil {
			return i - 1, err
		}
		pocketItem.ID = uint(itemId)

		authors, err := getAuthors(itemMap)
		if err == nil {
			pocketItem.Authors = authors
		}

		tags, err := getTags(itemMap)
		if err == nil {
			pocketItem.Tags = tags
		}

		for _, v := range pocketItem.Tags {
			v.PocketItemID = pocketItem.ID
		}

		// FIXME: sql: converting argument $25 type: unsupported type src.Image, a struct
		if err = db.Create(&pocketItem).Error; err != nil {
			return i - 1, err
		}
		i++
	}

	return len(items), nil
}
