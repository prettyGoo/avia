// proj project main.go
package json_parser

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

type RawEntry struct {
	PartnerHotelId int    `json:"id"`
	PartnerUrl     string `json:hotelpage`
	Checkin        string `json:"check_in_time"`
	Checkout       string `json:"check_out_time"`
	Images         []struct {
		Url string `json:"url"`
	} `json:"images"`
	Kind           string  `json:"kind"`
	Phone          string  `json:"phone"`
	StarRating     int     `json:"star_rating"`
	Email          string  `json:"email"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	CountryIsoCode string  `json:"country_code"`
	EnInfo         struct {
		Address            string `json:"address"`
		City               string `json:"city"`
		Country            string `json:"country"`
		HotelName          string `json:"name"`
		RatingTotalVerbose string `json:"rating_total_verbose"`
		Description        string `json:"description"`
	} `json:"en"`
}

type Hotel struct {
	partner_hotel_id int
	partner_url      string
	hotel_name       string
	kind             string
	checkin          string
	checkout         string
}

type Address struct {
	address_line_1   string
	city             string
	country          string
	iso_country_code string
	latitude         float64
	longitude        float64
	phone            string
	email            string
}

func Parse(fileName string, db *sql.DB) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := RawEntry{}
		json.Unmarshal(scanner.Bytes(), &entry)

		hotel := Hotel{
			partner_hotel_id: entry.PartnerHotelId,
			partner_url:      entry.PartnerUrl,
			hotel_name:       entry.EnInfo.HotelName,
			kind:             entry.Kind,
			checkin:          entry.Checkin,
			checkout:         entry.Checkout,
		}

		fields, values := reflectHotel(hotel)

		// INSERT NEW HOTEL
		var query string = fmt.Sprintf("INSERT INTO hotels(%s) VALUES(%s) RETURNING id", strings.TrimSuffix(fields.String(), ", "), strings.TrimSuffix(values.String(), ", "))
		fmt.Println(query)
		var id int
		err = db.QueryRow(query).Scan(&id)
		if err != nil {
			panic(err)
		}
		// END INSERT NEW HOTEL

		address := Address{
			address_line_1:   entry.EnInfo.Address,
			city:             entry.EnInfo.City,
			country:          entry.EnInfo.Country,
			iso_country_code: entry.CountryIsoCode,
			latitude:         entry.Latitude,
			longitude:        entry.Longitude,
			phone:            entry.Phone,
			email:            entry.Email,
		}

		// INSERT NEW ADDRESS ROW
		addressFields, addressValues := reflectAddress(address)

		// because strings.Builder cannot copy and will cause painc, refactor
		var finalAddressFields strings.Builder
		var finalAddressValues strings.Builder

		fmt.Fprintf(&finalAddressFields, "%s", addressFields.String())
		fmt.Fprintf(&finalAddressFields, "%s", "hotel_id")

		fmt.Fprintf(&finalAddressValues, `%s`, addressValues.String())
		fmt.Fprintf(&finalAddressValues, `%s`, fmt.Sprintf("%d", id))

		fmt.Printf("INSERT INTO hotel_addresses(%s) VALUES(%s)", finalAddressFields.String(), finalAddressValues.String())
		InsertDb(finalAddressFields.String(), finalAddressValues.String(), "hotel_addresses", db)
		// END INSERT NEW ADDRESS ROW
	}
}

func InsertDb(fields string, values string, tableName string, db *sql.DB) {
	fmt.Printf("INSERT INTO %s(%s) VALUES(%s)", tableName, fields, values)
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", tableName, fields, values)
	sqlStatement, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	_, err = sqlStatement.Exec()
	if err != nil {
		panic(err)
	}
}

func reflectHotel(hotel Hotel) (strings.Builder, strings.Builder) {
	var fields strings.Builder
	var values strings.Builder

	v := reflect.ValueOf(hotel)
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		n := v.Type().Field(j).Name

		fmt.Fprintf(&fields, `%s, `, n)
		if fmt.Sprintf("%s", f.Kind()) == "string" {
			fmt.Fprintf(&values, `$$'%s'$$, `, v.Field(j))
		} else {
			fmt.Fprintf(&values, `%d, `, v.Field(j))
		}

	}

	return fields, values
}

func reflectAddress(address Address) (strings.Builder, strings.Builder) {
	var fields1 strings.Builder
	var values1 strings.Builder

	v := reflect.ValueOf(address)
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		n := v.Type().Field(j).Name

		fmt.Fprintf(&fields1, `%s, `, n)
		if fmt.Sprintf("%s", f.Kind()) == "string" {
			fmt.Fprintf(&values1, `$$'%s'$$, `, v.Field(j))
		} else if fmt.Sprintf("%s", f.Kind()) == "int" {
			fmt.Fprintf(&values1, `%d, `, v.Field(j))
		} else {
			fmt.Fprintf(&values1, `%.2f, `, v.Field(j))
		}

	}

	return fields1, values1
}
