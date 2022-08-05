/*
* Copyright 2022 Mohammad Mohamamdi. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"git.sr.ht/~mmohammadi9812/pockesql/src"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const (
	PAGE_SIZE   int = 500
	RETRY_SLEEP int = 3
)

var since string

var fetchCmd = &cobra.Command{
	Use:     "fetch",
	Aliases: []string{"get"},
	Short:   "Fetch items from pocket api (given that it has been authenticated)",
	Run: func(_ *cobra.Command, _ []string) {
		FetchCmd()
	},
}

type FetchOptions struct {
	Auth   src.AuthInfo
	Offset int
	Since  string
}

type Error struct {
	Messaage error
	StatusCode int
}

func createFetchRequest(args FetchOptions) (*http.Request, error) {
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

func fetchPocket(args FetchOptions) (src.RawItems, error) {
	var (
		req *http.Request
		err error
	)

	if req, err = createFetchRequest(args); err != nil {
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

func createFetchOptions(auth src.AuthInfo) (FetchOptions, error) {
	args := FetchOptions{
		Auth:   auth,
		Offset: 0,
	}

	if since != "" {
		layout := "2000-01-01 15:04"
		t, err := time.Parse(layout, since)
		if err != nil {
			return FetchOptions{}, err
		}

		args.Since = strconv.FormatInt(t.Unix(), 10)
	}

	return args, nil
}

func Save(args FetchOptions, totalItems int64) Error {
	bar := progressbar.Default(int64(totalItems))

	for {
		raw, err := fetchPocket(args)
		if err != nil {
			return Error{StatusCode: -9, Messaage: fmt.Errorf("An error occurred while trying to fetch items")}
		}

		if raw == nil && err == nil {
			break
		}

		n, err := src.CreateItems(raw)
		if err != nil || n < 0 {
			return Error{StatusCode: n, Messaage: err}
		}
		args.Offset += PAGE_SIZE

		// TODO: progress-bar micro increment instead of macro increments (n)
		bar.Add(n)
	}

	return Error{}
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

	args, err := createFetchOptions(auth)
	if err != nil {
		log.Fatal(err)
	}

	e := Save(args, totalItems)
	if e.Messaage != nil {
		log.Fatalf("\n\n[%d] An error occured:\n%v\n\n", e.StatusCode, e.Messaage)
	}
}
