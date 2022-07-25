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

	"github.com/mitchellh/mapstructure"
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

func InsertItem(conn *sqlite.Conn, val interface{}) (err error) {
	tableFields := getTableFields(val)
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
	values := getStructValues(&val)

	err = sqlitex.Execute(conn, stmt, &sqlitex.ExecOptions{
		Args: values,
	})
	return
}

func getAuthors(itemMap map[string]interface{}) ([]Author, error) {
	var (
		result []Author
		rawAuthors map[string]map[string]string
		ok bool
	)
	if rawAuthors, ok = itemMap["authors"].(map[string]map[string]string); !ok {
		return nil, fmt.Errorf("authors not found")
	}
	for _, ra := range rawAuthors {
		var author Author
		authorDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName: "json",
			Result: &author,
		})
		if err != nil {
			return nil, err
		}
		authorDecoder.Decode(ra)
		result = append(result, author)
	}

	return result, nil
}

func insertAllAuthors(conn *sqlite.Conn, authors []Author) (err error) {
	err = CreateTable(conn, Author{})
	if err != nil {
		return
	}
	for _, author := range authors {
		if err = InsertItem(conn, author); err != nil {
			return
		}
	}
	return nil
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

		pocketItem, err := DecodeStruct(itemMap)
		if err := InsertItem(conn, pocketItem); err != nil {
			return i-1, err
		}

		authors, err := getAuthors(itemMap)
		if err != nil {
			return i-1, err
		}

		if err = insertAllAuthors(conn, authors); err != nil {
			return i-1, err
		}
		// TODO: insert into items_authors table
	}

	return len(items), nil
}
