package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"

	// "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

func insertRows(ctx context.Context, tx pgx.Tx, items []PItem) error {
	// Insert four rows into the "jixa_master" table.
	log.Println("Creating new rows...")
	for _, item := range items {
		if _, err := tx.Exec(ctx,
			"INSERT INTO jixa_master (item_name, cat_code, cost, txn_type, comment) VALUES ($1, $2, $3, $4, $5)", item.Item_name, item.Cat_code, item.Cost, item.Txn_type, item.Comment); err != nil {
			return err
		}
	}

	return nil
}

func insertCodeRows(ctx context.Context, tx pgx.Tx, catCodes []CatCode) error {
	// Insert four rows into the "jixa_master" table.
	log.Println("Creating new rows...")
	for _, catCode := range catCodes {
		if _, err := tx.Exec(ctx,
			"INSERT INTO jixa_cat_codes (code, category) VALUES ($1, $2)", catCode.Code, catCode.Category); err != nil {
			return err
		}
	}

	return nil
}

func printSummary(conn *pgx.Conn) error {
	rows, err := conn.Query(context.Background(), "SELECT id, item_name, item_head, cost, txn_type FROM jixa_master")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, cost int
		var item_name, item_head, txn_type string
		if err := rows.Scan(&id, &item_name, &item_head, &cost, &txn_type); err != nil {
			log.Fatal(err)
		}
		log.Printf("%d: %s %s %d %s\n", id, item_name, item_head, cost, txn_type)
	}
	return nil
}

func readItemsFromCSV(filename string) ([]PItem, error) {
	// 2. Read CSV file using csv.Reader
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CSV file")
	}
	// fmt.Println(data)
	var items []PItem
	for _, row := range data {
		item := PItem{Comment: "nil"}
		item.Item_name = row[0]
		cat_code, _ := strconv.Atoi(row[1])
		item.Cat_code = cat_code
		cost, _ := strconv.ParseFloat(row[2], 64)
		item.Cost = cost
		item.Txn_type = row[3]
		if len(row) > 4 && row[4] != "" {
			item.Comment = row[4]
			fmt.Println(item.Comment)
		}
		items = append(items, item)
	}
	return items, nil
}

func readCodesFromCSV(filename string) ([]CatCode, error) {
	// 2. Read CSV file using csv.Reader
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CSV file")
	}
	fmt.Println(data)
	var catCodes []CatCode
	for _, row := range data {
		var catCode CatCode
		catCode.Category = row[0]
		cat_code, _ := strconv.Atoi(row[1])
		catCode.Code = cat_code

		catCodes = append(catCodes, catCode)
	}
	return catCodes, nil
}
