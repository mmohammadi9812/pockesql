// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RawItems map[string]map[string]interface{}

func getDatabaseName() string {
	err := godotenv.Load()
	defaultName := "pocket.sqlite3"
	if err != nil {
		return defaultName
	}
	f, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return defaultName
	}
	return f
}

func saveAssoc(db *gorm.DB, pi PocketItem) (err error) {
	assocs := make(map[interface{}]interface{}, 1)
	if len(pi.Authors) > 0 {
		assocs[&Author{}] = &pi.Authors
	}
	if len (pi.Tags) > 0 {
		assocs[&Tag{}] = &pi.Tags
	}
	for k, v := range assocs {
		err = db.Clauses(clause.Insert{Modifier: "or IGNORE"}).Create(v).Error
		if err != nil {
			return
		}
		err = db.Model(k).Association("PocketItems").Append(&pi)
		if err != nil {
			return
		}
	}

	return nil
}

func CreateItems(items RawItems) (int, error) {
	db, err := gorm.Open(sqlite.Open(getDatabaseName()))
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

		pocketItem.CreatedAt = time.Now()

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

func ReadPocketItems() ([]PocketItem, error) {
	db, err := gorm.Open(sqlite.Open(getDatabaseName()))
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

func ReadItemTags(pi PocketItem) ([]Tag, error) {
	db, err := gorm.Open(sqlite.Open(getDatabaseName()))
	if err != nil {
		return nil, err
	}

	var tags []Tag
	err = db.Model(&pi).Association("Tags").Find(&tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func DeleteIds(Ids []uint) error {
	log.Printf("DEBUG::Ids:%v\n\n", Ids)
	db, err := gorm.Open(sqlite.Open(getDatabaseName()))
	if err != nil {
		return err
	}

	return db.Delete(&PocketItem{}, Ids).Error
}