// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import "gorm.io/gorm"

type Author struct {
	ID   int    `mapstructure:"author_id"`
	Name string `mapstructure:"name"`
	Url  string `mapstructure:"url"`

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
	ID           int
	Tag          string `mapstructure:"tag"`

	PocketItems []PocketItem `gorm:"many2many:items_tags"`
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
