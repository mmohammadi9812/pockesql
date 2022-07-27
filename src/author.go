// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

type Author struct {
	ID int `gorm:"column:author_id"`
	Name string `gorm:"column:name"`
	Url string `gorm:"column:url"`

	PocketItems []PocketItem `gorm:"many2many:items_authors"`
}
