// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RawItems map[string]map[string]interface{}

func saveAssoc(db *gorm.DB, pi PocketItem) (err error) {
	for _, v := range []interface{}{&pi.Authors, &pi.Tags} {
		err = db.Clauses(clause.Insert{Modifier: "or IGNORE"}).Create(v).Error
		if err != nil {
			return
		}
		err = db.Model(&Author{}).Association("PocketItems").Append(&pi)
		if err != nil {
			return
		}
	}

	return nil
}

func CreateItems(items RawItems) (int, error) {
	db, err := gorm.Open(sqlite.Open("pocket.sqlite3"), &gorm.Config{})
	if err != nil {
		return -2, err
	}

	var pocketItem PocketItem
	err = db.AutoMigrate(
		&pocketItem.Tags,
		&pocketItem.DomainMetadata,
		&pocketItem.Authors,
		&pocketItem.TopImage,
		&pocketItem.Images,
		&pocketItem.Videos,
		&pocketItem)
	if err != nil {
		return -3, err
	}

	for iid, itemMap := range items {
		itemMap = Transform(itemMap)

		pocketItem, err := DecodeStruct(itemMap)
		if err != nil {
			return -4, err
		}
		var itemId int
		if itemId, err = strconv.Atoi(iid); err != nil {
			return -5, err
		}
		pocketItem.ID = uint(itemId)

		err = db.Clauses(clause.Insert{Modifier: "or IGNORE"}).Create(&pocketItem).Error
		if err != nil {
			return -6, err
		}

		if err = saveAssoc(db, pocketItem); err != nil {
			return -7, err
		}
	}

	return len(items), nil
}

func UpdateItems(items RawItems) (int, error) {
	panic("unimplemented")
}

func ReadPocketItems() ([]PocketItem, error) {
	db, err := gorm.Open(sqlite.Open("pocket.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var items []PocketItem
	result := db.Order("created_at desc").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}
