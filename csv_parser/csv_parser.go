// proj project main.go
package csv_parser

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"proj/src/fields_mapper"

	_ "github.com/lib/pq"
)

func Parse(fileName string, db *sql.DB) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(file)

	header, err := r.Read()
	if err != nil {
		panic(err)
	}

	hotelFields := fields_mapper.HotelFields()
	hotelReviewsFields := fields_mapper.HotelReviewsFields()
	hotelAddressesFields := fields_mapper.HotelAddressesFields()
	hotelImagesFields := fields_mapper.HotelImagesFields()

	var record []string

	for {
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		var fields strings.Builder
		var values strings.Builder

		var reviewFields strings.Builder
		var reviewValues strings.Builder

		var addressFields strings.Builder
		var addressValues strings.Builder

		// IMAGES
		var imageArray []string

		// REFACTOR
		for i, header_name := range header {
			if _, ok := hotelImagesFields[header_name]; ok {
				imageArray = append(imageArray, record[i])
			} else if val, ok := hotelAddressesFields[header_name]; ok {
				fmt.Fprintf(&addressFields, "%s, ", val)
				if header_name == "addressline1" || header_name == "addressline2" || header_name == "city" || header_name == "state" || header_name == "country" || header_name == "countryisocode" || header_name == "region" || header_name == "continent_name" {
					fmt.Fprintf(&addressValues, `'%s', `, record[i])
				} else {
					fmt.Fprintf(&addressValues, `%s, `, record[i])
				}
			} else if val, ok := hotelReviewsFields[header_name]; ok {
				fmt.Fprintf(&reviewFields, "%s, ", val)
				fmt.Fprintf(&reviewValues, `%s, `, record[i])
			} else if val, ok := hotelFields[header_name]; ok {
				if len(record[i]) > 0 {
					fmt.Fprintf(&fields, "%s, ", val)
					if header_name == "overview" {
						fmt.Fprintf(&values, `'%s', `, strings.Replace(record[i], "'", "''", -1))
					} else if header_name == "hotel_name" || header_name == "hotel_formerly_name" || header_name == "url" || header_name == "rates_currency" {
						fmt.Fprintf(&values, `'%s', `, record[i])
					} else {
						fmt.Fprintf(&values, `%s, `, record[i])
					}

				}
			}
		}

		// INSERT NEW HOTEL
		var query string = fmt.Sprintf("INSERT INTO hotels(%s) VALUES(%s) RETURNING id", strings.TrimSuffix(fields.String(), ", "), strings.TrimSuffix(values.String(), ", "))
		fmt.Println(query)
		var id int
		err = db.QueryRow(query).Scan(&id)
		if err != nil {
			panic(err)
		}
		// END INSERT NEW HOTEL

		// INSERT NEW REVIEW ROW
		fmt.Fprintf(&reviewFields, "%s", "hotel_id")
		fmt.Fprintf(&reviewValues, `%s`, fmt.Sprintf("%d", id))
		InsertDb(reviewFields.String(), reviewValues.String(), "hotel_reviews", db)
		// END INSERT NEW REVIEW ROW

		// INSERT NEW ADDRESS ROW
		fmt.Fprintf(&addressFields, "%s", "hotel_id")
		fmt.Fprintf(&addressValues, `%s`, fmt.Sprintf("%d", id))
		InsertDb(addressFields.String(), addressValues.String(), "hotel_addresses", db)
		// END INSERT NEW ADDRESS ROW

		// INSERT NEW IMAGES ROW
		for _, image := range imageArray {
			imageValues := fmt.Sprintf("'%s', %d", image, id)
			InsertDb("img_url, hotel_id", imageValues, "hotel_images", db)
		}
		// END INSERT NEW IMAGES ROW
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
