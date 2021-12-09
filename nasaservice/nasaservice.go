package nasaservice

import (
	"fmt"
	"net/url"
	"os"
)

type NASA struct {
	Domain  *url.URL
	Service string
	APIKey  string
}

func (ns *NASA) SetDomain() error {
	domain, ok := os.LookupEnv("NASA_APIS")

	if !ok {
		fmt.Println("NASA_APIS environment variable doesn`t exist")
	}

	URL, err := url.Parse(domain)

	if err != nil {
		fmt.Println("failed to parse url from NASA source")
	}

	ns.Domain = URL

	return nil
}

func (ns *NASA) SetService() error {
	service, ok := os.LookupEnv("NASA_SERVICE")
	if !ok {
		fmt.Println("NASA_SERVICE environment variable doesn`t exist")
	}

	if len(service) == 0 {
		fmt.Println("NASA_SERVICE environment variable exists, but not set")
	}

	ns.Service = service

	return nil
}

func (ns *NASA) SetAPIKey() error {
	APIKey, ok := os.LookupEnv("SERVICE_APIKEY")

	if !ok {
		fmt.Println("SERVICE_APIKey environment variable doesn`t exist")
	}

	if len(APIKey) == 0 {
		fmt.Println("SERVICE_APIKey environment variable exists, but not set")
	}

	ns.APIKey = APIKey

	return nil
}
