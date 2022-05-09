package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id          string  `db:"id"`
	Name        string  `db:"name"`
	Coef        float64 `db:"trap_coef"`
	CreatorCoef float64 `db:"trap_creator_coef"`
}

// SQL用のInsert,Update用の関数
func Push(sqlCommand string, args ...interface{}) {
	dbx, dberr := sqlx.Open("sqlite3", "./test.sqlite3")
	if dberr != nil {
		println("ERROR opening DB")
		panic(dberr)
	}
	defer dbx.Close()
	dbx.SetConnMaxLifetime(time.Minute * 3)
	dbx.SetMaxOpenConns(10)
	dbx.SetMaxIdleConns(10)

	tx := dbx.MustBegin()
	tx.MustExec(sqlCommand, args...)
	tx.Commit()
}

// SQL用のselect用の関数
func Select(rows interface{}, sqlCommand string, args ...interface{}) {
	dbx, dberr := sqlx.Open("sqlite3", "./test.sqlite3")
	if dberr != nil {
		fmt.Println("ERROR opening DB")
		panic(dberr)
	}
	defer dbx.Close()
	dbx.SetConnMaxLifetime(time.Minute * 3)
	dbx.SetMaxOpenConns(10)
	dbx.SetMaxIdleConns(10)

	if err := dbx.Select(rows, sqlCommand, args...); err != nil {
		fmt.Println("ERROR selecting")
		panic(err)
	}
}
