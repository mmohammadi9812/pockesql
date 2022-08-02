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
	"reflect"
	"strconv"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const (
	PAGE_SIZE   int = 500
	RETRY_SLEEP int = 3
)

var fetchCmd = &cobra.Command{
	Use:     "fetch",
	Aliases: []string{"get"},
	Short:   "Fetch items from pocket api (given that it has been authenticated)",
	Run: func(_ *cobra.Command, _ []string) {
		FetchCmd()
	},
}

type FUrlValues struct {
	Auth   src.AuthInfo
	Offset int
	Since  string
}

func createFetchUrl(args FUrlValues) (*http.Request, error) {
	qvals := url.Values{
		"consumer_key": []string{args.Auth.ConsumerKey},
		"access_token": []string{args.Auth.AccessToken},
		"sort":         []string{"oldest"},
		"state":        []string{"all"},
		"detailType":   []string{"complete"},
		"count":        []string{strconv.Itoa(PAGE_SIZE)},
		"offset":       []string{strconv.Itoa(args.Offset)},
	}

	if args.Since != "" {
		qvals.Set("since", args.Since)
	}

	req, err := http.NewRequest("GET", "https://getpocket.com/v3/get", nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = qvals.Encode()

	return req, nil
}

func fetch(args FUrlValues) (src.RawItems, error) {
	var (
		req *http.Request
		err error
	)

	if req, err = createFetchUrl(args); err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: better retry mechanism

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var page map[string]interface{}
	err = json.Unmarshal(bodyBytes, &page)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(page["list"]).Kind() == reflect.Slice &&
		len(page["list"].([]interface{})) == 0 {
		return nil, nil
	}

	var (
		l1 = page["list"].(map[string]interface{})
		l2 = make(map[string]map[string]interface{}, len(l1))
	)
	for k, v := range l1 {
		l2[k] = v.(map[string]interface{})
	}

	return l2, nil
}

func FetchCmd() {
	auth, err := src.ReadAuth()
	if err != nil {
		log.Fatal(err)
	}

	totalItems, err := src.TotalItems(auth)
	if err != nil {
		log.Fatal(err)
	}

	offset := 0
	bar := progressbar.Default(int64(totalItems))

	for {
		raw, err := fetch(FUrlValues{
			Auth:   auth,
			Offset: offset,
		})
		if err != nil {
			log.Fatal(err)
		}

		if raw == nil && err == nil {
			break
		}

		n, err := src.CreateItems(raw)
		if err != nil || n < 0 {
			log.Fatalf("Saving fetch items failed [%d]\n%v", n, err)
		}

		offset += PAGE_SIZE

		// TODO: progress-bar micro increment instead of macro increments (n)
		bar.Add(n)
	}
}
