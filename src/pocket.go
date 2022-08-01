// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import "github.com/mitchellh/mapstructure"

type PocketItem struct {
	ID                     uint           `gorm:"primaryKey;autoIncrement:false" mapstructure:"item_id"`
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

func DecodeStruct(item map[string]interface{}) (pocketItem PocketItem, err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &pocketItem,
		IgnoreUntaggedFields: true,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return
	}
	err = decoder.Decode(item)
	return
}

