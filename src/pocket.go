// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"log"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

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

	Tags    []Tag    `gorm:"many2many:items_authors" mapstructure:"-"`
	Authors []Author `gorm:"many2many:items_authors" mapstructure:"-"`
}

func DecodeStruct(item map[string]interface{}) (pocketItem PocketItem, err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &pocketItem, IgnoreUntaggedFields: true})
	if err != nil {
		return
	}
	err = decoder.Decode(item)
	return
}

func transformTags(item map[string]interface{}) map[string]interface{} {
	var (
		tags map[string]interface{}
		ok   bool
	)
	if tags, ok = item["tags"].(map[string]interface{}); !ok || tags == nil {
		return item
	}

	for k := range tags {
		iid, yea := item["tags"].(map[string]interface{})[k].(map[string]interface{})["item_id"]
		if !yea {
			continue
		}
		if _, iss := iid.(string); !iss {
			continue
		}
		nv, err := strconv.Atoi(iid.(string))
		if err != nil {
			continue
		}
		item["tags"].(map[string]interface{})[k].(map[string]interface{})["item_id"] = nv
	}

	return item
}

func transformImages(item map[string]interface{}) map[string]interface{} {
	keys := []string{
		"height",
		"item_id",
		"width",
	}

	for _, k := range keys {
		if v, ok := item["image"];
			!ok || v == nil ||
				reflect.TypeOf(item["image"].(map[string]interface{})[k]).Kind() == reflect.Int {
			break
		}
		nv, err := strconv.Atoi(item["image"].(map[string]interface{})[k].(string))
		if err != nil {
			continue
		}
		item["image"].(map[string]interface{})[k] = nv
	}

	if v, ok := item["images"]; !ok || v == nil {
		return item
	}

	var images []map[string]interface{}
	if reflect.TypeOf(item["images"]).Kind() == reflect.Slice {
		item["images"] = map[string]interface{}{
			"0": item["images"].([]map[string]interface{})[0],
		}
		log.Printf("DEBUG::converted slice to single item: %v\n", item["images"])
	}
	for ik := range item["images"].(map[string]interface{}) {
		for _, k := range keys {
			nv, err := strconv.Atoi(item["images"].(map[string]interface{})[ik].(map[string]interface{})[k].(string))
			if err != nil {
				continue
			}
			item["images"].(map[string]interface{})[ik].(map[string]interface{})[k] = nv
		}
		images = append(images, item["images"].(map[string]interface{})[ik].(map[string]interface{}))
	}
	item["images"] = images

	return item
}

func transformVideo(item map[string]interface{}) map[string]interface{} {
	keys := []string{
		"height",
		"item_id",
		"length",
		"width",
	}

	if v, ok := item["videos"].(map[string]interface{}); !ok || v == nil {
		return item
	}

	var videos []map[string]interface{}
	for ik := range item["videos"].(map[string]interface{}) {
		for _, k := range keys {
			nv, err := strconv.Atoi(item["videos"].(map[string]interface{})[ik].(map[string]interface{})[k].(string))
			if err != nil {
				continue
			}
			item["videos"].(map[string]interface{})[ik].(map[string]interface{})[k] = nv
		}
		videos = append(videos, item["videos"].(map[string]interface{})[ik].(map[string]interface{}))
	}
	item["videos"] = videos

	return item
}

func Transform(item map[string]interface{}) map[string]interface{} {
	skeys := []string{
		"item_id",
		"resolved_id",
		"favorite",
		"status",
		"time_added",
		"time_updated",
		"time_read",
		"time_favorited",
		"word_count",
	}

	for _, key := range skeys {
		if v, ok := item[key]; ok && reflect.TypeOf(v).Kind() == reflect.String {
			nv, err := strconv.Atoi(v.(string))
			if err != nil {
				continue
			}
			item[key] = nv
		}
	}

	fkeys := []string{
		"time_to_read",
		"listen_duration_estimate",
	}

	for _, key := range fkeys {
		if v, ok := item[key]; ok && reflect.TypeOf(v).Kind() == reflect.Float64 {
			item[key] = int(v.(float64))
		}
	}

	bkeys := []string{
		"has_video",
		"has_image",
		"is_article",
		"is_index",
	}

	for _, key := range bkeys {
		if v, ok := item[key]; ok {
			var nv int; var err error
			if reflect.TypeOf(v).Kind() == reflect.String {
				nv, err = strconv.Atoi(v.(string))
			} else {
				nv, ok = item[key].(int)
			}
			if err != nil || !ok {
				item[key] = false
			}
			item[key] = nv > 0
		}
	}

	item = transformTags(item)
	item = transformImages(item)
	item = transformVideo(item)

	return item
}
