package database

import (
	"contact-api/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	procard_db_name	= "procard_dev"
	pg_db_user      = "postgres"
	pg_db_password  = ""
	pg_db_host 		= "127.0.0.1:5432"
	pg_db_port      = "5432"
)

var (
	db *gorm.DB
	user = pg_db_user // os.Getenv("pg_db_user")
	password = pg_db_password // os.Getenv("pg_db_password")
	database = procard_db_name // os.Getenv("procard_db_name")
	host     = pg_db_host // os.Getenv("pg_db_host")
	port     = pg_db_port // os.Getenv("pg_db_port")
)

func init() {

	//Load environmenatal variables
	//err := godotenv.Load()
	//
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	//Define DB connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		user, password, host, database,
	)

	//connect to database connection string
	conn, err := gorm.Open("postgres", connStr)

	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}

	if err = conn.DB().Ping(); err != nil {
		panic(err)
	}
	db = conn

	// close database when not in use
	// defer conn.Close()

	// Migrate the schema
	conn.Debug().AutoMigrate(
		&models.Account{})

	fmt.Println("Successfully connected!", db)
	// return conn
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
