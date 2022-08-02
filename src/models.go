// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

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
	ID  int
	Tag string `mapstructure:"tag"`

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

type DeletedAt sql.NullTime

type PocketItem struct {
	ID        uint `gorm:"primaryKey;autoIncrement:false" mapstructure:"item_id"`
	CreatedAt time.Time
	DeletedAt DeletedAt `gorm:"index"`

	ResolvedId             uint           `mapstructure:"resolved_id"`
	GivenUrl               string         `mapstructure:"given_url"`
	GivenTitle             string         `mapstructure:"given_title"`
	Favorite               int            `mapstructure:"favorite"`
	Status                 int            `mapstructure:"status"`
	TimeAdded              int            `mapstructure:"time_added"`
	TimeUpdated            int            `mapstructure:"time_updated"`
	TimeRead               int            `mapstructure:"time_read"`
	TimeFavorited          int            `mapstructure:"time_favorited"`
	SortId                 int            `mapstructure:"sort_id"`
	ResolvedTitle          string         `mapstructure:"resolved_title"`
	ResolvedUrl            string         `mapstructure:"resolved_url"`
	Excerpt                string         `mapstructure:"excerpt"`
	IsArticle              bool           `mapstructure:"is_article"`
	IsIndex                bool           `mapstructure:"is_index"`
	HasVideo               bool           `mapstructure:"has_video"`
	HasImage               bool           `mapstructure:"has_image"`
	WordCount              int            `mapstructure:"word_count"`
	Lang                   string         `mapstructure:"lang"`
	DomainMetadata         DomainMetadata `mapstructure:"domain_metadata"`
	ListenDurationEstimate int            `mapstructure:"listen_duration_estimate"`
	TimeToRead             int            `mapstructure:"time_to_read"`
	AmpUrl                 string         `mapstructure:"amp_url"`
	TopImageUrl            string         `mapstructure:"top_image_url"`
	TopImage               TopImage       `mapstructure:"image"`
	Images                 []Image        `mapstructure:"images"`
	Videos                 []Video        `mapstructure:"videos"`

	Tags    []Tag    `gorm:"many2many:items_tags;" mapstructure:"tags"`
	Authors []Author `gorm:"many2many:items_authors;" mapstructure:"authors"`
}
