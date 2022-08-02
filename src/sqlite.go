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

func SaveItems(items map[string]map[string]interface{}) (int, error) {
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

		if len(pocketItem.Authors) > 0 {
			err = db.Clauses(clause.Insert{Modifier: "or IGNORE"}).Create(&pocketItem.Authors).Error
			if err != nil {
				return -7, err
			}
			err = db.Model(&Author{}).Association("PocketItems").Append(&pocketItem)
			if err != nil {
				return -8, err
			}
		}

		if len(pocketItem.Tags) > 0 {
			err = db.Clauses(clause.Insert{Modifier: "or IGNORE"}).Create(&pocketItem.Tags).Error
			if err != nil {
				return -9, err
			}
			err = db.Model(&Tag{}).Association("PocketItems").Append(&pocketItem)
			if err != nil {
				return -10, err
			}
		}
	}

	return len(items), nil
}
