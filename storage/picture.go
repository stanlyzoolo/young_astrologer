package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// A picture of a day presents json from NASA api source.
type APOD struct {
	ItemID         int    `db:"unique_id"       json:"id"`
	Copyright      string `db:"copyright"       json:"copyright"`
	Date           string `db:"date"            json:"date"`
	Explanation    string `db:"explanation"     json:"explanation"`
	HDURL          string `db:"hdurl"           json:"hdurl"`
	MediaType      string `db:"mediaType"      json:"mediaType"`
	ServiceVersion string `db:"serviceVersion" json:"serviceVersion"`
	Title          string `db:"title"           json:"title"`
	URL            string `db:"url"             json:"url"`
	Image          []byte `db:"image"           json:"image"`
}

var logger *zap.Logger

func (p *APOD) getImage() error {
	imageURL := p.URL
	resp, err := http.Get(imageURL)

	if err != nil {
		logger.Error(err.Error())
	}

	defer resp.Body.Close()

	p.Image, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

func (p *APOD) Metadata(resp *http.Response) error {
	d, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logger.Error(err.Error())
	}

	err = json.Unmarshal(d, &p)

	if err != nil {
		logger.Error(err.Error())
	}

	err = p.getImage()

	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

func (p *APOD) Save() string {
	sql := fmt.Sprintf(`INSERT INTO %s (%v, %s, %s, %s, %s, %s, %s, %s, %s, %v)
	 VALUES ($1, $2, $3, $4. $5, $6, $7, $8, $9, $10)`,
		os.Getenv("PSQLNASATABLE"),
		p.ItemID,
		p.Copyright,
		p.Date,
		p.Explanation,
		p.HDURL,
		p.MediaType,
		p.ServiceVersion,
		p.Title,
		p.URL,
		p.Image)

	return sql
}
