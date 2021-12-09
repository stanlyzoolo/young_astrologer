package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	"os"
	"young_astrologer/nasaservice"
	"young_astrologer/storage"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("FAILED to load env variables from .env", err)
	}
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Set URL for GET - method request
	var ns nasaservice.NASA
	url := prepareURL(&ns)

	// Prepare database and table
	var apod storage.APOD

	connNASA, err := prepareDatabase()
	if err != nil {
		logger.Error("failed database preparing process",
			zap.Error(err))
	}

	// Send request to the NASA APOD source
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("400 Bad Request",
			zap.String("package", "main"),
			zap.String("func", "main"),
			zap.Error(err))
	}

	err = apod.Metadata(resp)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Sugar().Info("To see Astronomy Picture of the Day (APOD), please, use the link below!")
	logger.Sugar().Info(fmt.Sprintf("Press Ctrl+click on link to view the picture in browser: %s", apod.URL))

	// Speak with CLI
	if interactive() {
		connNASA.Exec(context.Background(),
			"INSERT INTO nasa_pictures (unique_id, copyright, date, explanation, hdurl, media_type,	service_version, title,	url, image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",

			apod.ItemID,
			apod.Copyright,
			apod.Date,
			apod.Explanation,
			apod.HDURL,
			apod.MediaType,
			apod.ServiceVersion,
			apod.Title,
			apod.URL,
			apod.Image)
	}

	defer connNASA.Close(context.Background())
}

func prepareDatabase() (*pgx.Conn, error) {
	logger, _ := zap.NewDevelopment()

	// Create connection string for PostgreSQL driver 'pgx' with 'postgres' default user
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/",
		os.Getenv("PSQLDEFAULTUSER"),
		os.Getenv("PSQLDEFAULTPASS"),
		os.Getenv("PSQLHOST"),
		os.Getenv("PSQLPORT"))

	// Establish connection
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		logger.Sugar().Info("---------------The connection has been successfully established-----------------")
	}
	defer conn.Close(context.Background())

	// New role credentials
	role := storage.Role{
		Name:     os.Getenv("PSQLNASAUSER"),
		Password: os.Getenv("PSQLNASAPASSWORD"),
		Database: os.Getenv("PSQLNASADB"),
	}

	// Check new role existence
	logger.Sugar().Info("Checking existence new role name from your .env file...")
	ok := role.IsExist(context.Background(), conn, logger)

	if !ok {
		logger.Sugar().Info("Sorry, new role name has already exist. Set new role name in .env file.\nService stopped. Bye!")
	} else {
		// Required steps to set up new role and database
		conn.Exec(context.Background(), role.New())
		logger.Sugar().Info("---------------The role has been successfully created---------------------------")
		conn.Exec(context.Background(), role.Alter())

		conn.Exec(context.Background(), role.CreateDB())
		logger.Sugar().Info("---------------Database has been successfully created-----------------")
		conn.Exec(context.Background(), role.Grant())
		logger.Sugar().Info("---------------All privileges has been successfully granted---------------------")

		conn.Close(context.Background())

		// Create connection string for PostgreSQL driver 'pgx' with 'postgres' default user
		dbNASAURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			os.Getenv("PSQLNASAUSER"),
			os.Getenv("PSQLNASAPASSWORD"),
			os.Getenv("PSQLHOST"),
			os.Getenv("PSQLPORT"),
			os.Getenv("PSQLNASADB"))

		// Connection to new apod database
		connNASA, _ := pgx.Connect(context.Background(), dbNASAURL)
		logger.Sugar().Info("---------------Required table has been successfully created-----------------")
		connNASA.Exec(context.Background(), role.CreateTable())

		return connNASA, nil
	}

	return nil, err
}

func prepareURL(ns *nasaservice.NASA) string {
	ns.SetDomain()
	ns.SetAPIKey()
	ns.SetService()

	ns.Domain.Path = ns.Service
	query := ns.Domain.Query()
	query.Set("api_key", ns.APIKey)
	ns.Domain.RawQuery = query.Encode()
	URL := ns.Domain.String()

	return URL
}

func interactive() bool {
	fmt.Println("If you like this picture, we can save it for you. Just type yes/no:")
	
	sc := bufio.NewScanner(os.Stdin)
	
	for sc.Scan() {
		if sc.Text() == "yes" {
			return true
		} else {
			fmt.Println("We hope the better picture will appear tomorrow")			
			return false
		}
	}

	return false
}
