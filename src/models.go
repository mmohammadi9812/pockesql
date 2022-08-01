// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import 	"gorm.io/gorm"

type Author struct {
	ID int `gorm:"column:author_id" json:"author_id"`
	Name string `gorm:"column:name" json:"name"`
	Url string `gorm:"column:url" json:"url"`

	PocketItems []PocketItem `gorm:"many2many:items_authors"`
}

type TopImage struct {
	Height       int    `mapstructure:"height"`
	ID           int    `mapstructure:"item_id"`
	Src          string `mapstructure:"src"`
	Width        int    `mapstructure:"width"`
	PocketItemID uint
}

type Image struct {
	Height       int    `mapstructure:"height"`
	ID           int    `mapstructure:"item_id"`
	Src          string `mapstructure:"src"`
	Width        int    `mapstructure:"width"`
	Caption      string `mapstructure:"caption"`
	Credit       string `mapstructure:"credit"`
	PocketItemID uint
}

type Tag struct {
	ID           int    `gorm:"primaryKey;autoIncrement:false" mapstructure:"item_id"`
	Tag          string `mapstructure:"tag"`
	PocketItemID uint
}

type DomainMetadata struct {
	gorm.Model
	Logo         string `mapstructure:"logo"`
	Name         string `mapstructure:"name"`
	PocketItemID uint
}

type Video struct {
	Height       int    `mapstructure:"height"`
	Width        int    `mapstructure:"width"`
	Length       int    `mapstructure:"length"`
	Src          string `mapstructure:"src"`
	PocketItemID uint   `mapstructure:"item_id"`
}
