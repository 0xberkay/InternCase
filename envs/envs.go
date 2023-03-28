package envs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBhost string
	DBport string
	DBuser string
	DBpass string
	DBname string
	DBssl  string
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBhost = os.Getenv("dbhost")
	DBport = os.Getenv("dbport")
	DBuser = os.Getenv("dbuser")
	DBpass = os.Getenv("dbpass")
	DBname = os.Getenv("dbname")
	DBssl = os.Getenv("dbsslmode")

}
