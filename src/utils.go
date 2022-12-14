// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/maps"
)

func TotalItems(auth AuthInfo) (int64, error) {
	statUrl := "https://getpocket.com/v3/stats"
	data := map[string]string{
		"consumer_key": auth.ConsumerKey,
		"access_token": auth.AccessToken,
	}

	bodyBytes, err := rawSendJson(statUrl, data)
	if err != nil {
		return -1, err
	}

	var (
		objMap map[string]interface{}
		jd = json.NewDecoder(strings.NewReader(string(bodyBytes)))
	)
	jd.UseNumber()
	err = jd.Decode(&objMap)
	if err != nil {
		return -3, err
	}

	return objMap["count_list"].(json.Number).Int64()
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

func Transform(item map[string]interface{}) map[string]interface{} {
	for _, key := range []string{"tags", "authors", "videos", "images"} {
		if _, ok := item[key].(map[string]interface{}); ok && reflect.TypeOf(item[key]).Kind() == reflect.Map {
			item[key] = maps.Values(item[key].(map[string]interface{}))
		}
	}

	return item
}
