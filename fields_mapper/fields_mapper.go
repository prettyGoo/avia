package fields_mapper

func HotelFields() map[string]string {
	m := map[string]string{
		"hotel_id":            "partner_hotel_id",
		"numberrooms":         "number_rooms",
		"star_rating":         "star_rating",
		"numberfloors":        "number_floors",
		"hotel_name":          "hotel_name",
		"hotel_formerly_name": "hotel_formerly_name",
		"yearopened":          "year_opened",
		"yearrenovatd":        "year_renovated",
		"url":                 "partner_url",
		"hotelpage":           "partner_url",
		"rates_from":          "rates_from",
		"rates_currency":      "rates_currency",
		"overview":            "description",
		"description":         "description",
		"description_short":   "description_short",
	}

	return m
}

func HotelReviewsFields() map[string]string {
	m := map[string]string{
		"count":             "rating_count",
		"number_of_reviews": "reviews_count",
		"rating_average":    "rating_average",
	}
	return m
}

func HotelAddressesFields() map[string]string {
	m := map[string]string{
		"longitude":      "longitude",
		"latitude":       "latitude",
		"zip":            "zip",
		"email":          "email",
		"phone":          "phone",
		"address":        "address_line_1",
		"addressline1":   "address_line_1",
		"addressline2":   "address_line_2",
		"city":           "city",
		"state":          "state",
		"country":        "country",
		"countryisocode": "iso_country_code",
		"region":         "region",
		"continent_name": "continent",
	}

	return m
}

func HotelImagesFields() map[string]string {
	m := map[string]string{
		"photo1":   "img_url",
		"photo2":   "img_url",
		"photo3":   "img_url",
		"photo4":   "img_url",
		"photo5":   "img_url",
		"orig_url": "img_url",
	}

	return m
}
