package whooplist

import (
	"database/sql"
)

type Place struct {
	Id int64

	Latitude  float64
	Longitude float64

	FactualId string

	Name     string
	Address  string
	Locality string
	Region   string
	Postcode string
	Country  string
	Tel      string
	Website  string
	Email    string

	Score int64
}

var getPlaceStmt, getPlaceByFactualStmt, addPlaceStmt,
	updatePlaceStmt *sql.Stmt

func preparePlace() {
	stmt(&getPlaceStmt,
		"SELECT id, latitude, longitude, factual_id, name, address, locality, "+
			"region, postcode, country, telephone, website, email "+
			"FROM wl.place WHERE id = $1;")

	stmt(&getPlaceByFactualStmt,
		"SELECT latitude, longitude, factual_id, name, address, locality, "+
			"region, postcode, country, telephone, website, email "+
			"FROM wl.place WHERE factual_id = $1;")

	stmt(&updatePlaceStmt,
		"UPDATE wl.place SET latitude=$1, longitude=$2, factual_id=$3, "+
			"name=$4, address=$5, locality=$6, region=$7, postcode=$8, "+
			"country=$9, telephone=$10, website=$11, email=$12 "+
			"WHERE factual_id=$3 RETURNING id;")

	stmt(&addPlaceStmt,
		"INSERT INTO wl.place (latitude, longitude, factual_id, name, "+
			"address, locality, region, postcode, country, telephone, "+
			"website, email) "+
			"SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12 "+
			"WHERE NOT EXISTS (SELECT 1 FROM wl.place WHERE factual_id=$3) "+
			"RETURNING id;")

}

func GetPlaceId(placeId int64) (place *Place, err error) {
	return getPlaceDb(getPlaceStmt.QueryRow(placeId))
}

func GetPlaceFactual(placeId string) (place *Place, err error) {
	return getPlaceDb(getPlaceByFactualStmt.QueryRow(placeId))
}

func getPlaceDb(res *sql.Row) (place *Place, err error) {
	err = res.Scan(&place.Id, &place.Latitude, &place.Longitude,
		&place.FactualId, &place.Name, &place.Address, &place.Locality,
		&place.Region, &place.Postcode, &place.Country,
		&place.Tel, &place.Website, &place.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		place = nil
	}

	return
}

func SearchPlace(str string, listId int64, page int32,
	lat, long, radius float64) (places []Place, err error) {

	places, err = factualPlaceSearch(str, lat, long, radius, page)

	if err != nil {
		places = nil
		return
	}

	err = addPlaces(places)

	return
}

func addPlaces(places []Place) (err error) {
	for k, place := range places {
		res := addPlaceStmt.QueryRow(place.Latitude, place.Longitude,
			place.FactualId, place.Name, place.Address, place.Locality,
			place.Region, place.Postcode, place.Country, place.Tel,
			place.Website, place.Email)

		err = res.Scan(&(places[k].Id))

		if err == sql.ErrNoRows {
			res := updatePlaceStmt.QueryRow(place.Latitude, place.Longitude,
				place.FactualId, place.Name, place.Address, place.Locality,
				place.Region, place.Postcode, place.Country, place.Tel,
				place.Website, place.Email)
			err = res.Scan(&(places[k].Id))
			if err != nil {
				return
			}
		}
	}
	return
}
