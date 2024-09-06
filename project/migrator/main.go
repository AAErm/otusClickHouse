package main

import (
	"database/sql"
	"log"
	"sync"

	usergenerator "github.com/AAErm/otusClickHouse/project/migrator/userGenerator"
)

const totalRecords = 1_500_000

func main() {
	dsn := "username:password@tcp(localhost:3306)/dbname"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	numGoroutines := 10

	recordsPerGoroutine := totalRecords / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < recordsPerGoroutine/100; j++ {
				randomUsers, err := usergenerator.GetUsers()
				if err != nil {
					log.Println("Error fetching users:", err)
					continue
				}

				for _, randomUser := range randomUsers {
					_, err := db.Exec("INSERT INTO users (FirstName, FatherName, LastName, GenderCode, Bank, YearsOld) VALUES (?, ?, ?, ?, ?, ?)",
						randomUser.FatherName, randomUser.FatherName, randomUser.LastName, randomUser.GenderCode, randomUser.Bank, randomUser.YearsOld)
					if err != nil {
						log.Println("Error inserting user:", err)
					}
				}
			}
		}()
	}

	wg.Wait()
}
