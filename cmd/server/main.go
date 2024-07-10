package main

import (
	"HomeWork5/server/db"
	"fmt"
)

func main() {
	_, err := db.NewDB()
	if err != nil {
		fmt.Println(err)
	}
}
