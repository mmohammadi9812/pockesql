/*
* Copyright 2022 Mohammad Mohamamdi. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
 */

package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/schollz/progressbar/v3"
)

type AuthInfo struct {
	ConsumerKey string `json:"consumer_key"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

const PAGE_SIZE int = 500

func ReadAuth() (AuthInfo, error) {
	fbytes, err := ioutil.ReadFile("auth.json")
	if err != nil {
		return AuthInfo{}, err
	}
	inf := AuthInfo{}
	err = json.Unmarshal(fbytes, &inf)
	if err != nil {
		return AuthInfo{}, err
	}
	return inf, nil
}

func TotalItems(auth AuthInfo) (int, error) {
	statUrl := "https://getpocket.com/v3/stats"
	data := map[string]string{
		"consumer_key": auth.ConsumerKey,
		"access_token": auth.AccessToken,
	}
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(data)
	resp, err := http.Post(statUrl, "application/json", payloadBuffer)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var objMap map[string]interface{}
	err = json.Unmarshal(bodyBytes, &objMap)
	if err != nil {
		return -1, err
	}

	return objMap["count_list"].(int), nil
}

func FetchCmd() {
	auth, err := ReadAuth()
	if err != nil {
		log.Fatal(err)
	}

	totalItems, err := TotalItems(auth)
	if err != nil {
		log.Fatal(err)
	}

	var since string
	offset, retries := 0, 0
	bar := progressbar.Default(int64(totalItems))
	fetchUrl := "https://getpocket.com/v3/get"

	for {

		qvals := url.Values{
			"consumer_key": []string{auth.ConsumerKey},
			"access_token": []string{auth.AccessToken},
			"sort":         []string{"oldest"},
			"state":        []string{"all"},
			"detailType":   []string{"complete"},
			"count":        []string{strconv.Itoa(PAGE_SIZE)},
			"offset":       []string{strconv.Itoa(offset)},
		}
		if since != "" {
			qvals.Set("since", since)
		}

		req, err := http.NewRequest("GET", fetchUrl, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.URL.RawQuery = qvals.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 503 && retries < 5 {
			log.Println("Got a 503, retrying...")
			retries += 1
			continue
		}
		retries = 0
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var page map[string]interface{}
		err = json.Unmarshal(bodyBytes, &page)
		if err != nil {
			log.Fatal(err)
		}
		items, ok := page["list"].(map[string]interface{})
		if !ok || len(items) == 0 {
			break
		}

		// FIXME: do something with items!

		offset += PAGE_SIZE

		bar.Add(1)
	}
}
