/*
* Copyright 2022 Mohammad Mohamamdi. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
 */

package cmd

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/schollz/progressbar/v3"
)

const (
	PAGE_SIZE int = 500
	RETRY_SLEEP int = 3
)

func FetchCmd() {
	auth, err := src.ReadAuth()
	if err != nil {
		log.Fatal(err)
	}

	totalItems, err := src.TotalItems(auth)
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
			time.Sleep(time.Duration(retries) * time.Duration(RETRY_SLEEP))
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
		items, ok := page["list"].([]map[string]interface{})
		if !ok || len(items) == 0 {
			break
		}

		// FIXME: complete saveItems implementation
		src.SaveItems(items)

		offset += PAGE_SIZE

		bar.Add(1)
	}
}
