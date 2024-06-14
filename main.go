package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// #region city
type City struct {
	ID          int    `json:"ID,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

// #endregion city
func main() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{
		User:      os.Getenv("DB_USERNAME"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_DATABASE"),
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}

	db, err := sqlx.Open("mysql", conf.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected")
	// #region get

	cityName := os.Args[1]

	var city City
	err = db.Get(&city, "SELECT * FROM city WHERE Name = ?", cityName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such city Name = '%s'\n", cityName)
		return
	}
	if err != nil {
		log.Fatalf("DB Error: %s\n", err)
	}
	// #endregion get
	log.Printf("%sの人口は%d人です\n", cityName, city.Population)

	var countryPop int
	err = db.Get(&countryPop, "SELECT Population FROM country WHERE Code = ?", city.CountryCode)
	if err != nil {
		log.Fatalf("DB Error: %s\n", err)
	}
	log.Printf("%sの人口は%sの人口の%f%%です\n", city.Name, city.CountryCode, float64(city.Population)/float64(countryPop)*100)
}
