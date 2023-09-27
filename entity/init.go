package entity

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"gitlab.com/mefit/mefit-server/utils"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	//Needed for gorm working with postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/qor/qor"
	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/initializer"
	"gitlab.com/mefit/mefit-server/utils/log"
)

type manager struct {
}

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		log.Logger().WithFields(logrus.Fields{"module": "gorm", "type": "sql"}).Debug(v[3])
	}
	if v[0] == "log" {
		log.Logger().WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Debug(v[2])
	}
}

var (
	db *gorm.DB
)

//Deprecated TODO must be replaced with Query
func getDB() *gorm.DB {
	return db
}

func GetDB() *gorm.DB {
	return db
}

//WithTransaction exec sql with transaction
func WithTransaction(f func(*gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// Note the use of tx as the database handle once you are within a transaction

	if err := f(tx); err != nil {
		tx.Rollback()
		log.Logger().Errorf("transaction rollback duo %v ", err)
		return err
	}

	tx.Commit()
	if tx.Error != nil {
		log.Logger().Errorf("while commit transaction duo %v ", tx.Error)
		return utils.ErrInternal
	}
	return nil
}

func getDBCreds() (string, string, string, string, string, string) {
	log.Logger().Print("Parsing db credentials")
	var (
		dbHost string
		dbPort string
		dbUser string
		dbName string
		dbPass string
		dbSsl  string
	)
	if _, ok := os.LookupEnv("VCAP_SERVICES"); ok {
		vcapServices := make(map[string]interface{})
		json.Unmarshal([]byte(os.Getenv("VCAP_SERVICES")), &vcapServices)
		creds := vcapServices["postgresql"].([]interface{})[0].(map[string]interface{})["credentials"].(map[string]interface{})
		dbName = creds["database"].(string)
		dbUser = creds["username"].(string)
		dbHost = creds["host"].(string)
		dbPort = strconv.Itoa(int(creds["port"].(float64)))
		dbPass = creds["password"].(string)
	} else {
		dbHost = config.Config().GetString("postgres_host")
		dbPort = config.Config().GetString("postgres_port")
		dbUser = config.Config().GetString("postgres_user")
		dbName = config.Config().GetString("postgres_db")
		dbPass = config.Config().GetString("postgres_password")
	}

	dbSsl = config.Config().GetDefaultString("postgres_ssl", "disable")
	return dbHost, dbPort, dbUser, dbName, dbPass, dbSsl
}

func (m manager) Initialize() func() {
	debugMode := config.Config().GetBool("DEBUG")
	dbHost, dbPort, dbUser, dbName, dbPass, dbSsl := getDBCreds()

	var once sync.Once
	once.Do(func() {
		var err error
		if dbPass != "" {
			dbPass = "password=" + dbPass
		}
		db, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s %s sslmode=%s", dbHost, dbPort, dbUser, dbName, dbPass, dbSsl))
		if err != nil {
			log.Logger().Panic(err.Error())
		}
		db.Exec("CREATE EXTENSION IF NOT EXISTS hstore")
		db.Set("gorm:table_options", "charset=utf8")
		db.SingularTable(true)
		if !debugMode {
			db.SetLogger(&GormLogger{})
		}
		db.LogMode(true)
		err = db.DB().Ping()
		if err != nil {
			log.Logger().Fatal("Database ping error! Error:", err.Error())
		}
		//FIXME: maybe we have golang related issue here
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		db.DB().SetMaxIdleConns(5)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		db.DB().SetMaxOpenConns(50)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		db.DB().SetConnMaxLifetime(time.Minute * 15)
		log.Logger().Print("Database ping success!")
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return "mefit_" + defaultTableName
		}
		log.Logger().Print("DB initialized")
	})

	log.Logger().Print("init entity models")
	// Migrate the schema

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Profile{})
	db.AutoMigrate(&BazaarPayment{})
	// db.AutoMigrate(&Program{})
	// db.AutoMigrate(&Level{})
	// db.AutoMigrate(&SubLevel{})
	db.AutoMigrate(&Plan{})
	db.AutoMigrate(&Workout{})
	db.AutoMigrate(&Exercise{})
	db.AutoMigrate(&ExerciseSection{})
	db.AutoMigrate(&Class{})
	db.AutoMigrate(&Movement{})
	// db.AutoMigrate(&Article{})
	db.AutoMigrate(&Keyword{})
	db.AutoMigrate(&MuscleGroup{})
	db.AutoMigrate(&WorkoutHistory{})
	db.AutoMigrate(&WorkoutType{})
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Payment{})
	db.AutoMigrate(&PurchasedProduct{})

	return func() {
		db.Close()
	}
}

func init() {
	initializer.Register(manager{}, initializer.VeryHighPriority)
}
