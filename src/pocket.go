// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type PocketItem struct {
	Id int `json:"item_id",primarykey`
	ResolvedId int `json:"resolved_id"`
	GivenUrl string `json:"given_url"`
	GivenTitle string `json:"given_title"`
	Favorite int `json:"favorite"`
	Status int `json:"status"`
	TimeAdded int `json:"time_added"`
	TimeUpdated int `json:"time_updated"`
	TimeRead int `json:"time_read"`
	TimeFavorited int `json:"time_favorited"`
	SortId int `json:"sort_id"`
	ResolvedTitle string `json:"resolved_title"`
	ResolvedUrl string `json:"resolved_url"`
	Excerpt string `json:"excerpt"`
	IsArticle int `json:"is_article"`
	IsIndex int `json:"is_index"`
	HasVideo int `json:"has_video"`
	HasImage int `json:"has_image"`
	WordCount int `json:"word_count"`
	Lang string `json:"lang"`
	DomainMetadata string `json:"domain_metadata"`
	ListenDurationEstimate int `json:"listen_duration_estimate"`
	TimeToRead int `json:"time_to_read"`
	AmpUrl string `json:"amp_url"`
	TopImageUrl string `json:"top_image_url"`
	Tags string `json:"tags"`
	Image string `json:"image"`
	Images string `json:"images"`
	Videos string `json:"videos"`
}

func DecodeStruct(item map[string]interface{}) (pocketItem PocketItem, err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &pocketItem})
	if err != nil {
		return
	}
	err = decoder.Decode(item)
	return
}

func TransformValues(item map[string]interface{}) map[string]interface{} {
	keys := []string {
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
