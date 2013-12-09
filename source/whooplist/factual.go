package whooplist

import (
	"encoding/json"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"strconv"
)

const oauthKey = "DSGgIoFt5B0B5YjeXLOzb61p2t2H00Pc1GD1xvXe"
const oauthSecret = "qfE5c4TbLtjvgoJhph5g1MIPuxKLxjpusxRKr4Zr"
const factualV3Url = "http://api.v3.factual.com/t/"
const factualLimit = 40

var oConsumer *oauth.Consumer

type factualManyResp struct {
	Version  int    `json:"version"`
	Status   string `json:"status"`
	Response struct {
		Data         []factualPlace `json:"data"`
		IncludedRows int            `json:"included_rows"`
	} `json:"response"`
}

type factualPlace struct {
	Address         string      `json:"address"`
	AddressExtended string      `json:"address_extended"`
	AdminRegion     string      `json:"admin_region"`
	CategoryIds     []int       `json:"category_ids"`
	CategoryLabels  [][]string  `json:"category_labels"`
	ChainId         string      `json:"chain_id"`
	Country         string      `json:"country"`
	Email           string      `json:"email"`
	FactualId       string      `json:"factual_id"`
	Fax             string      `json:"fax"`
	Hours           interface{} `json:"hours"` //TODO: formalize
	Latitude        float64     `json:"latitude"`
	Locality        string      `json:"locality"`
	Longitude       float64     `json:"longitude"`
	Name            string      `json:"name"`
	Neighbourhood   string      `json:"neighbourhood"`
	Postcode        string      `json:"postcode"`
	Region          string      `json:"region"`
	Tel             string      `json:"tel"`
	Website         string      `json:"website"`
}

func initializeOauth() {
	if oConsumer != nil {
		return
	}

	oConsumer = oauth.NewConsumer(oauthKey, oauthSecret,
		oauth.ServiceProvider{})
}

func factualSearchPlace(str string, lat, long, radius float64,
	page int32) (places []Place, err error) {

	latS := strconv.FormatFloat(lat, 'f', -1, 64)
	longS := strconv.FormatFloat(long, 'f', -1, 64)
	radiusS := strconv.FormatFloat(radius, 'f', -1, 64)

	params := make(map[string]string)
	params["q"] = str
	params["geo"] = "{\"$circle\":{\"$center\":[" + latS + "," + longS +
		"],\"$meters\":" + radiusS + "}}"
	params["offset"] = strconv.Itoa((int(page) - 1) * factualLimit)
	params["limit"] = strconv.Itoa(factualLimit)

	respHttp, err := oConsumer.Get(factualV3Url+"places",
		params, &oauth.AccessToken{})
	if err != nil {
		return
	}

	resp, err := ioutil.ReadAll(respHttp.Body)
	respHttp.Body.Close()
	if err != nil {
		return
	}

	data := &factualManyResp{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return
	}

	places = make([]Place, data.Response.IncludedRows)
	for i := 0; i < data.Response.IncludedRows; i++ {
		decodeFactualPlace(&data.Response.Data[i], &places[i])
	}

	return
}

//TODO: Implement get single place, at factualV3Url+"places/"+factualId

func decodeFactualPlace(fP *factualPlace, p *Place) {
	p.Latitude = fP.Latitude
	p.Longitude = fP.Longitude
	p.FactualId = fP.FactualId
	p.Name = fP.Name
	p.Address = fP.Address
	p.Locality = fP.Locality
	p.Region = fP.Region
	p.Postcode = fP.Postcode
	p.Country = fP.Country
	p.Tel = fP.Tel
	p.Website = fP.Website
	p.Email = fP.Email
}
