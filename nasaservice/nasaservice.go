package nasaservice

import (
	"net/url"
	"os"

	"go.uber.org/zap"
)

type NASA struct {
	Domain  *url.URL
	Service string
	APIKey  string
}

var logger *zap.Logger

func (ns *NASA) SetDomain() error {
	domain, ok := os.LookupEnv("NASA_APIS")

	if !ok {
		logger.Error("NASA_APIS environment variable doesn`t exist")
	}

	url, err := url.Parse(domain)

	if err != nil {
		logger.Error("failed to parse url from NASA source")
	}

	ns.Domain = url

	return nil
}

func (ns *NASA) SetService() error {
	service, ok := os.LookupEnv("NASA_SERVICE")
	if !ok {
		logger.Error("NASA_SERVICE environment variable doesn`t exist")
	}

	if len(service) == 0 {
		logger.Error("NASA_SERVICE environment variable exists, but not set")
	}

	ns.Service = service

	return nil
}

func (ns *NASA) SetAPIKey() error {
	APIKey, ok := os.LookupEnv("SERVICE_APIKey")

	if !ok {
		logger.Error("SERVICE_APIKey environment variable doesn`t exist")
	}

	if len(APIKey) == 0 {
		logger.Error("SERVICE_APIKey environment variable exists, but not set")
	}

	ns.APIKey = APIKey

	return nil
}
