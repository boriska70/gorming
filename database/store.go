package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	errrr "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//Store is a generic store interface
type Store interface {
	CreateTable(model Things) error
	CleanDB(shouldDropTable bool, model Things) error
	ThingsBulkInsert(things []Things) error
	ThingsCount() (int, error)
	ThingsGetAll() ([]Things, error)
}

//StoreInstance provides the interface to DB related operations
var StoreInstance *thingsStore

type thingsStore struct {
	db *gorm.DB
}

//New creates new StoreInstance variable
func New() *thingsStore {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=gorming sslmode=disable")
	if err != nil {
		return nil
	}
	return &thingsStore{
		db: db,
	}
}

func init() {
	StoreInstance = New()
}

func (ts *thingsStore) CreateTable(model Things) error {
	if ts.db.HasTable(model) {
		log.Debugf("skip creation of table things - already exists")
		return nil
	}
	return ts.db.CreateTable(model).Error
}

func (ts *thingsStore) CleanDB(shouldDropTable bool, model Things) error {
	if shouldDropTable {
		deleteTable(ts.db, model)
	}
	return closeConnectionToDB(ts.db)
}

func (ts *thingsStore) ThingsBulkInsert(things []Things) error {
	errors := make([]error, 0)
	for _, thing := range things {
		if err := ts.db.Table("things").Create(&thing).Error; err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) == 0 {
		return nil
	}
	errMessage := ""
	for _, err := range errors {
		errMessage = fmt.Sprintf("%s;%s", errMessage, err.Error())
	}
	return errrr.New(errMessage)
}

func (ts *thingsStore) ThingsGetAll() ([]Things, error) {
	var things []Things
	err := ts.db.Table("things").Find(&things).Error
	return things, err
}

func (ts *thingsStore) ThingsCount() (int, error) {
	var count int
	err := ts.db.Table("things").Count(&count).Error
	return count, err
}

func closeConnectionToDB(db *gorm.DB) error {
	var err error
	if db != nil {
		log.Debugf("closing DB connection")
		if err = db.Close(); err != nil {
			log.Errorf("failed to close DB connection: %s", err.Error)
		}
	}
	return err
}

func deleteTable(db *gorm.DB, model interface{}) error {
	var err error
	if err = db.DropTableIfExists(model).Error; err != nil {
		log.Errorf("failed to delete table %s; reason: %s", model, err.Error)
	}
	return err
}
