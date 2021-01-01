package config

import (
	"davidws/utils"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// Port is the listening port
	Port = 0
	// DBPort is the database listening Port
	DBPort = 3306
	// DBUser is the database user
	DBUser = ""
	// DBPassword is the database user password
	DBPassword = ""
	// DBName is the database name
	DBName = ""
	// DBHost is the database host
	DBHost = ""
	// RefreshTokenTimeout is the timeout of the refresh token in hours
	RefreshTokenTimeout int64 = 1
	// DBURI URL of database
	DBURI string

	// RedisAddr is the address to the Redis instance
	RedisAddr string = os.Getenv("REDIS_ADDR")
	// RedisDBPass is the password to the Redis database
	RedisDBPass string = os.Getenv("REDIS_DB_PASS")
	// RedisPass is the address to the Redis instance
	RedisPass string
	// SessionIDLength is length of the session ID
	SessionIDLength int = 156
	// SessionIDName is the name of the session ID
	SessionIDName string = os.Getenv("SESSION_ID_NAME")
	// CookieExp cookie expiration time in seconds
	CookieExp int64 = 900
)

// Load loads .env file and global variables
func Load() {
	// Load .env file
	err := godotenv.Load()

	if err != nil {
		DBURI = os.Getenv("DATABASE_URL")

		SessionIDLength, err = strconv.Atoi(os.Getenv("SESSION_ID_LENGTH"))

		if err != nil {
			log.Println("Error: ", err.Error())
			SessionIDLength = 156
		}
		return
	}

	utils.FailIfErr(err)

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))

	if err != nil {
		log.Println("Error: ", err.Error())
		Port = 3000
	}

	DBPort, err = strconv.Atoi(os.Getenv("DB_PORT"))

	utils.FailIfErr(err)

	DBUser = os.Getenv("DB_USER")
	DBHost = os.Getenv("DB_HOST")
	DBPassword = os.Getenv("DB_PASS")

	DBName = os.Getenv("DB_NAME")

	if DBHost == "" {
		log.Fatalln("Error:  nonexisten DB_HOST")
	}

	// SessionTimeout, err = strconv.Atoi(os.Getenv("SESSION_TIMEOUT"))

	utils.FailIfErr(err)

	RedisAddr = os.Getenv("REDIS_ADDR")
	RedisDBPass = os.Getenv("REDIS_DB_PASS")
	RedisPass = os.Getenv("REDIS_PASS")
}
