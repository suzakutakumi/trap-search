package main

import (
	"crypto/rand"
	"math/big"
	"trap-search/db"
	"trap-search/search"
)

func main() {
	var count []int
	db.Select(&count, "select count(*) from user where trap_creator_coef>=?", 0.5)
	user := []db.User{}
	num, err := rand.Int(rand.Reader, big.NewInt(int64(count[0])))
	if err != nil {
		panic(err)
	}
	db.Select(&user, "select * from (select * from user where trap_creator_coef>=?) limit 1 offset ?", 0.5, num.Int64())
	//fmt.Println(user[0].Id)
	search.FromCreator(user[0].Id)

}
