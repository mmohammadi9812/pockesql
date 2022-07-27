// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type PocketItem struct {
	ID                     uint   `mapstructure:"item_id"`
	ResolvedId             uint   `mapstructure:"resolved_id"`
	GivenUrl               string `mapstructure:"given_url"`
	GivenTitle             string `mapstructure:"given_title"`
	Favorite               int    `mapstructure:"favorite"`
	Status                 int    `mapstructure:"status"`
	TimeAdded              int    `mapstructure:"time_added"`
	TimeUpdated            int    `mapstructure:"time_updated"`
	TimeRead               int    `mapstructure:"time_read"`
	TimeFavorited          int    `mapstructure:"time_favorited"`
	SortId                 int    `mapstructure:"sort_id"`
	ResolvedTitle          string `mapstructure:"resolved_title"`
	ResolvedUrl            string `mapstructure:"resolved_url"`
	Excerpt                string `mapstructure:"excerpt"`
	IsArticle              int    `mapstructure:"is_article"`
	IsIndex                int    `mapstructure:"is_index"`
	HasVideo               int    `mapstructure:"has_video"`
	HasImage               int    `mapstructure:"has_image"`
	WordCount              int    `mapstructure:"word_count"`
	Lang                   string `mapstructure:"lang"`
	DomainMetadata         string `mapstructure:"domain_metadata"`
	ListenDurationEstimate int    `mapstructure:"listen_duration_estimate"`
	TimeToRead             int    `mapstructure:"time_to_read"`
	AmpUrl                 string `mapstructure:"amp_url"`
	TopImageUrl            string `mapstructure:"top_image_url"`
	Tags                   string `mapstructure:"tags"`
	Image                  string `mapstructure:"image"`
	Images                 string `mapstructure:"images"`
	Videos                 string `mapstructure:"videos"`

	Authors []Author `gorm:"many2many:items_authors"`
}

func DecodeStruct(item map[string]interface{}) (pocketItem PocketItem, err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &pocketItem, IgnoreUntaggedFields: true})
	if err != nil {
		return
	}
	err = decoder.Decode(item)
	return
}

func Transform(item map[string]interface{}) map[string]interface{} {
	keys := []string{
		"item_id",
		"resolved_id",
		"favorite",
		"status",
		"time_added",
		"time_updated",
		"time_read",
		"time_favorited",
		"is_article",
		"is_index",
		"has_video",
		"has_image",
		"word_count",
		"time_to_read",
		"listen_duration_estimate",
	}

	for _, key := range keys {
		if v, ok := item[key]; ok {
			nv, err := strconv.Atoi(v.(string))
			if err != nil {
				continue
			}
			item[key] = nv
		}
	}

	return item
}
