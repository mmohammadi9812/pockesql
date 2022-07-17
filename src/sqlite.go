// Copyright 2022 Mohammad Mohamamdi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type TableField struct{
	Name string
	Type string
	IsPrimary bool
}

func SqliteConn(filename string) (*sqlite.Conn, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fn := filepath.Join(pwd, filename)
	conn, err := sqlite.OpenConn(fn, sqlite.OpenReadWrite)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func getStructValues(inputStruct interface{}) []interface{} {
	t := reflect.TypeOf(inputStruct)
	vals := []interface{}{}
	for i := 0; i < t.NumField(); i++ {
		r := reflect.ValueOf(inputStruct)
		f := reflect.Indirect(r).FieldByName(t.Field(i).Name)
		if k := t.Field(i).Type.Kind(); k == reflect.Int {
			vals = append(vals, int(f.Int()))
		} else {
			vals = append(vals, f.String())
		}
	}
	return vals
}

func getTableFields(inputStruct interface{}) []TableField {
	t := reflect.TypeOf(inputStruct)
	fields := []TableField{}
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get("json")
		fieldType := ""
		if k := t.Field(i).Type.Kind(); k == reflect.Int {
			fieldType = "INTEGER"
		} else {
			fieldType = "TEXT"
		}
		_, primaryField := t.Field(i).Tag.Lookup("primarykey")
		fields = append(fields, TableField{
			Name: fieldName,
			Type: fieldType,
			IsPrimary: primaryField,
		})
	}
	return fields
}

func createTableStmt(inputStruct interface{}) string {
	t := reflect.TypeOf(inputStruct)
	tableName := t.Name()
	tableFields := getTableFields(inputStruct)
	tf_len := len(tableFields)
	stmt := fmt.Sprintf(`
	DROP TABLE IF EXISTS %s;
	CREATE TABLE %s (
	`, tableName, tableName);
	for i, field := range tableFields {
		var formatString string
		if field.IsPrimary {
			formatString = `
			%s %s PRIMARY KEY,
			`
		} else if i != tf_len - 1 {
			formatString = `
			%s %s,
			`
		} else {
			formatString = `
			%s %s);
			`
		}
		stmt += fmt.Sprint(formatString, field.Name, field.Type)
	}
	return stmt
}

func CreateTable(conn *sqlite.Conn, inputStruct interface{}) error {
	stmt := createTableStmt(inputStruct)
	err := sqlitex.ExecuteScript(conn, stmt, nil)
	return err
}

func InsertItem(conn *sqlite.Conn, item PocketItem) (err error) {
	tableFields := getTableFields(item)
	numFields := len(tableFields)
	stmt := `INSERT INTO items (`
	for i, field := range tableFields {
		if i != numFields - 1 {
			stmt += fmt.Sprintf("%s,", field.Name)
		} else {
			stmt += fmt.Sprintf("%s)", field.Name)
		}
	}
	stmt += " VALUES (" + strings.Repeat("?, ", numFields - 1) + "?);"
	values := getStructValues(&item)

	err = sqlitex.Execute(conn, stmt, &sqlitex.ExecOptions{
		Args: values,
	})
	return
}

func SaveItems(items []map[string]interface{}) (int, error) {
	conn, err := SqliteConn("pocket.sqlite3")
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	err = CreateTable(conn, PocketItem{})
	if err != nil {
		return -1, err
	}

	for i, itemMap := range items {
		itemMap = TransformValues(itemMap)
		// TODO: implement insertItem
		pocketItem, err := DecodeStruct(itemMap)
		if err := InsertItem(conn, pocketItem); err != nil {
			return i-1, err
		}
		// TODO: implement getAuthors
		authors, err := getAuthors(itemMap)
		if err == nil {
			// TODO: implement insertAllAuthors
			insertAllAuthors(conn, authors)
		}
		// TODO: insert into items_authors table
	}
}
