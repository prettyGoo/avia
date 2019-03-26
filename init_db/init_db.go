// proj project main.go
package init_db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB(db *sql.DB, noInit bool) {
	if noInit {
		return
	}

	err := createHotels(db)
	if err != nil {
		panic(err)
	}

	err = createHotelAmenities(db)
	if err != nil {
		panic(err)
	}

	err = createHotelReviews(db)
	if err != nil {
		panic(err)
	}

	err = createHotelAddresses(db)
	if err != nil {
		panic(err)
	}

	err = createHotelImages(db)
	if err != nil {
		panic(err)
	}
}

func createHotels(db *sql.DB) error {
	_, err := db.Query(
		`create table hotels
		(
			id serial not null,
			partner_hotel_id int,
			partner_source varchar(20),
			hotel_name varchar(200),
			hotel_translated_name varchar(200),
			hotel_formerly_name varchar(200),
			checkin time,
			checkout time,
			number_rooms int,
			number_floors int,
			year_opened int,
			year_renovated int,
			star_rating int,
			kind varchar(20),
			partner_url varchar(500),
			description varchar(2000),
			description_short varchar(200),
			rates_from int,
			rates_currency varchar(30)
		);
		
		create unique index hotels_id_uindex
			on hotels (id);
		
		alter table hotels
			add constraint hotels_pk
				primary key (id);
		`)

	return err
}

func createHotelAmenities(db *sql.DB) error {
	_, err := db.Query(
		`create table hotel_amenities
		(
			id serial not null,
			name varchar(100),
			hotel_id int not null
				constraint hotel_amenities_hotel_id_fk
					references hotels
						on delete cascade
		);
		
		create unique index hotel_amenities_id_uindex
			on hotel_amenities (id);
		
		alter table hotel_amenities
			add constraint hotel_amenities_pk
				primary key (id);
		`)

	return err
}

func createHotelReviews(db *sql.DB) error {
	_, err := db.Query(
		`create table hotel_reviews
		(
			id serial not null,
			reviews_count int,
			rating_count int,
			rating_average float,
			rating_total_verbose varchar(50),
			cleanness float,
            comfort float,
            location float,
            personnel float,
            price float,
            services float,
			hotel_id int not null
				constraint hotel_reviews_hotel_id_fk
					references hotels
						on delete cascade
		);
		
		create unique index hotel_reviews_id_uindex
			on hotel_reviews (id);
		
		alter table hotel_reviews
			add constraint hotel_reviews_pk
				primary key (id);
		`)

	return err
}

func createHotelAddresses(db *sql.DB) error {
	_, err := db.Query(
		`create table hotel_addresses
		(
			id serial not null,
			longitude float,
			latitude float,
			zip int,
			phone varchar(20),
			email varchar(50),
			address_line_1 varchar(100),
			address_line_2 varchar(100),
			city varchar(100),
			state varchar(100),
			country varchar(100),
			iso_country_code varchar(30),
			region varchar(100),
			continent varchar(100),

			hotel_id int not null
				constraint hotel_addresses_hotel_id_fk
					references hotels
						on delete cascade
		);
		
		create unique index hotel_addresses_id_uindex
			on hotel_addresses (id);
		
		alter table hotel_addresses
			add constraint hotel_addresses_pk
				primary key (id);
		`)

	return err
}

func createHotelImages(db *sql.DB) error {
	_, err := db.Query(
		`create table hotel_images
		(
			id serial not null,
			img_url varchar(200),

			hotel_id int not null
				constraint hotel_images_hotel_id_fk
					references hotels
						on delete cascade
		);
		
		create unique index hotel_images_id_uindex
			on hotel_images (id);
		
		alter table hotel_images
			add constraint hotel_images_pk
				primary key (id);
		`)

	return err
}

func createHotelChains(db *sql.DB) error {
	_, err := db.Query(
		`create table hotel_chains
		(
			id serial not null,
			name varchar(200),
	

			hotel_id int not null
				constraint hotel_chains_hotel_id_fk
					references hotels
						on delete cascade
		);
		
		create unique index hotel_chains_id_uindex
			on hotel_chains (id);
		
		alter table hotel_chains
			add constraint hotel_chains_pk
				primary key (id);
		`)

	return err
}
