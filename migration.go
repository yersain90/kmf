package main

import (
	"database/sql"
	"log"
)

const currency = `
    create table R_CURRENCY (
	    id INT PRIMARY KEY IDENTITY(1,1),	
	    title VARCHAR(60) NOT NULL,
	    code VARCHAR(3) NOT NULL,
	    value NUMERIC(18,2) NOT NULL,
	    adate DATE NOT NULL,
    )
`

func migrate(dbDriver *sql.DB) {
	statement, err := dbDriver.Prepare(currency)
	if err == nil {
		_, creationError := statement.Exec()
		if creationError == nil {
			log.Println("Table created successfully")
		} else {
			log.Println(creationError.Error())
		}
	} else {
		log.Println(err.Error())
	}
}
