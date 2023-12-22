package main

import (
	"context"
	// "encoding/csv"
	"fmt"
	"log"
	"os"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

type PItem struct {
	Item_name string  `json:"item_name"`
	Cat_code  int     `json:"cat_code"`
	Cost      float64 `json:"cost"`
	Txn_type  string  `json:"txn_type"`
	Comment   string  `json:"comment"`
}

type CatCode struct {
	Code     int    `json:"code"`
	Category string `json:"category"`
}

func main() {
	// Read in connection string
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	config.RuntimeParams["application_name"] = "$jixa_pgx"
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	fmt.Println("What's up?\n1. Insert Items\n2. Display\n3. Insert Cat Codes")
	var input int
	fmt.Scan(&input)
	switch input {
	case 1:
		fmt.Println("Inserting items...")
		var inp_file string
		fmt.Print("Enter file name: ")
		fmt.Scan(&inp_file)
		items, err := readItemsFromCSV(inp_file)
		if err != nil {
			log.Fatal(err)
		}
		err = crdbpgx.ExecuteTx(context.Background(), conn,
			pgx.TxOptions{}, func(tx pgx.Tx) error {
				return insertRows(context.Background(), tx, items)
			})
		if err == nil {
			log.Println("New rows created.")
		} else {
			log.Fatal("error: ", err)
		}

	case 2:
		fmt.Println("Displaying...")
		// // Print out the Summary
		printSummary(conn)
	case 3:
		fmt.Println("Inserting Cat Codes...")
		codes, err := readCodesFromCSV("cat_codes.csv")
		if err != nil {
			log.Fatal(err)
		}
		err = crdbpgx.ExecuteTx(context.Background(), conn,
			pgx.TxOptions{}, func(tx pgx.Tx) error {
				return insertCodeRows(context.Background(), tx, codes)
			})
		if err == nil {
			log.Println("New rows created.")
		} else {
			log.Fatal("error: ", err)
		}
	case 4:
		fmt.Println("Inserting Expenses")
		items, err := readItemsFromCSV("expense.csv")
		if err != nil {
			log.Fatal(err)
		}
		err = crdbpgx.ExecuteTx(context.Background(), conn,
			pgx.TxOptions{}, func(tx pgx.Tx) error {
				return insertRows(context.Background(), tx, items)
			})
		if err == nil {
			log.Println("New rows created.")
		} else {
			log.Fatal("error: ", err)
		}
	default:
		fmt.Println("Invalid input")
		os.Exit(1)
	}

}
