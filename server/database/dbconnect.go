package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type DbConfig struct {
	host       string
	port       int
	dbname     string
	dbuser     string
	dbpassword string
}

// This shouldn't be use in production, only when running locally
func DbConfigFromEnvironment() *DbConfig {
	dbUser, userOk := os.LookupEnv("DB_USER")
	if !userOk {
		panic("Environment variable DB_USER must be set to connect to database")
	}

	dbPassword, passwordOk := os.LookupEnv("DB_PASSWORD")
	if !passwordOk {
		panic("Environment variable DB_PASSWORD must be set to connect to database")
	}

	dbPort, portOk := os.LookupEnv("DB_PORT")
	if !portOk {
		panic("Environment variable DB_PORT must be set to connect to database")
	}

	dbHost, hostOk := os.LookupEnv("DB_HOST")
	if !hostOk {
		panic("Environment variable DB_HOST must be set to connect to database")
	}

	dbPortInt, portIntOk := strconv.Atoi(dbPort)
	if portIntOk != nil {
		fmt.Printf("Port: %s\n", dbPort)
		panic("Failed to convert port number to int: " + dbPort)
	}

	dbName, nameOk := os.LookupEnv("DB_NAME")
	if !nameOk {
		panic("Environment variable DB_NAME must be set to connect to database")
	}

	return &DbConfig{
		host:       dbHost,
		port:       dbPortInt,
		dbname:     dbName,
		dbpassword: dbPassword,
		dbuser:     dbUser,
	}
}

func newConnectionString(config *DbConfig) string {
	return "postgres://" + config.dbuser + ":" + config.dbpassword + "@" + config.host + ":" + strconv.Itoa(config.port) + "/" + config.dbname + "?sslmode=disable"
}

// Create a database connection
func NewDbConnection(config *DbConfig) (*DbConnection, error) {
	psqlConnectInfo := newConnectionString(config)

	log.Printf("%s\n", psqlConnectInfo)
	db, dbConnectionError := sql.Open("postgres", psqlConnectInfo)
	if dbConnectionError != nil {
		return nil, dbConnectionError
	}

	if dbPingError := db.Ping(); dbPingError != nil {
		return nil, dbPingError
	}

	fmt.Printf("Connection to database on port :: %d\n", config.port)

	return &DbConnection{
		db: db,
	}, nil
}
