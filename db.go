package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var DB *sql.DB

func initDB(cfg *Config) error {
	var err error
	connStr := fmt.Sprintf("user id=%s;password=%s;port=%s;database=%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Port, cfg.DB.Database)
	DB, err = sql.Open("mssql", connStr)
	return err
}

func SelectRates(date string, code string) ([]CurrencyRate, error) {
	list := make([]CurrencyRate, 0, 0)

	format := "02.02.2006"
	tDate, err := time.Parse(format, date)
	if err != nil {
		log.Println("error parsing date: ", err)
		return list, err
	}
	var stmt string
	var rows *sql.Rows
	if len(code) == 0 {
		stmt = `SELECT title, code, value  FROM dbo.R_CURRENCY
					where  adate = ?;`
		rows, err = DB.Query(stmt, tDate)
	} else {
		stmt = `SELECT  title, code, value  FROM dbo.R_CURRENCY
				where  adate = ? and code = ?`
		rows, err = DB.Query(stmt, tDate, code)
	}
	defer rows.Close()
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var title, code, value string
		err := rows.Scan(&title, &code, &value)
		if err != nil {
			fmt.Println("Error reading rows: " + err.Error())
		}
		list = append(list, CurrencyRate{
			Fullname:   title,
			Title:      code,
			Desciption: value,
		})
	}
	return list, nil
}

func InsertRates(date string, rates *Rates) {
	format := "02.02.2006"
	time, error := time.Parse(format, date)
	if error != nil {
		log.Println("error parsing time: ", error)
		return
	}

	queryStmt := `
			INSERT INTO R_CURRENCY( title, code, value, adate) VALUES ( ?, ?, ?, ?);`

	query, err := DB.Prepare(queryStmt)
	if err != nil {
		log.Fatal("error preparing statement: ", err)
	}
	defer query.Close()

	for _, v := range rates.Items {
		_, err = query.Exec(v.Fullname, v.Title, v.Desciption, time)
		if err != nil {
			fmt.Println("Error inserting new row: " + err.Error())
			return
		}
		//}
	}
}
