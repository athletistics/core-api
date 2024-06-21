// Source  : https://divrhino.com/articles/rest-api-docker-go-fiber-from-scratch
// Brief   : Establish a connection to a database
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 16.06.2023

package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/internal/models"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var Mock sqlmock.Sqlmock

// ConnectDb creates a connection towards the database and keeps that connection in the public 'DB' variable.
// Once connected, the auto migration is run, updating the database structure.
func ConnectDb() {
	// prepare connection string
	dsn := fmt.Sprintf(
		"host=%s port=%v user=%s password=%s dbname=%s  sslmode=disable TimeZone=Europe/Zurich",
		os.Getenv("DB_HOST"),
		getPort(),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	// create connection and add a logger
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("joining tables")
	if err := db.SetupJoinTable(&models.User{}, "Teams", &models.Participant{}); err != nil {
		log.Fatal("Failed to join tables \n", err)
	}
	if err := db.SetupJoinTable(&models.Team{}, "Members", &models.Participant{}); err != nil {
		log.Fatal("Failed to join tables. \n", err)
	}

	log.Println("running migrations")
	if db.AutoMigrate(&models.AccessRole{}, &models.Authorisation{}, &models.Category{}, &models.Club{}, &models.Country{},
		&models.Discipline{}, &models.Email{}, &models.Event{}, &models.EventType{}, &models.Participant{},
		&models.Phone{}, &models.Result{}, &models.Sport{}, &models.Team{}, &models.UserShare{},
		&models.User{}, &models.Language{}, &models.Gender{}, &models.ResultStatic{}, &models.ResultStaticType{},
		&models.ResultTraining{},
		// External IDs
		&models.ExternalIdType{}, &models.ExternalId{},
		// NDS models
		&models.NdsActivityType{}, &models.NdsActivity{}, &models.NdsActivityInstance{}, &models.NdsAttendance{},
	) != nil {
		log.Fatal("Failed to run auto migrations. \n", err)
	}

	DB = db
}

func ConnectMockDB() {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	DB = db
	Mock = mock
}

func GetSqlDB() *sql.DB {
	db, err := DB.DB()
	if err != nil {
		log.Fatal("Unable to extract SQL DB and prepare seeder. \n", err)
	}

	return db
}

func getPort() int {
	value, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return constants.DB_PORT_DEFAULT
	}

	intVal, err := strconv.Atoi(value)
	if err != nil || intVal <= 0 {
		return constants.DB_PORT_DEFAULT
	}

	return intVal
}
