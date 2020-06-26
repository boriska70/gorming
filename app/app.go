package app

import (
	db "github.com/boriska70/gorming/database"
	log "github.com/sirupsen/logrus"
)

//App describes the application structure
type App struct {
	store db.Store
}

//NewApp creates a new instance of the application
func NewApp() App {
	return App{}
}

//Initialize prepares the environment
func (app *App) Initialize() error {
	app.store = db.StoreInstance

	if err := app.store.CreateTable(db.Things{}); err != nil {
		return err
	}
	log.Debugf("table %s has been created successfully", db.Things{})

	if count, err := app.store.ThingsCount(); err == nil && count == 0 {
		app.store.ThingsBulkInsert(generateData())
	}

	log.Info("application is ready to run")
	return nil
}

//Close removes tables (optionally) and releases DB connection
func (app *App) Close(hardClean bool) {
	if err := app.store.CleanDB(hardClean, db.Things{}); err != nil {
		log.Errorf("failed to release resources: %s", err.Error())
	}
	log.Info("finished successfully")
}

//Run executes all the functionality
func (app *App) Run() {
	things, err := app.store.ThingsGetAll()
	if err != nil {
		log.WithError(err).Errorf("failed to get all things from database: %s", err.Error())
	}
	log.Infof("find all things in database: %+v", things)
}

func generateData() []db.Things {
	result := make([]db.Things, 0)

	result = append(result, db.Things{
		Name:   "table",
		Useful: true,
		ThingDetails: db.Details{
			Quality: db.Quality{
				Look: 2,
				Real: 2,
			},
			CanBeSold: 2,
		},
	})
	result = append(result, db.Things{
		Name:   "garbage",
		Useful: false,
		ThingDetails: db.Details{
			Quality: db.Quality{
				Look: 0,
				Real: 0,
			},
			CanBeSold: 0,
		},
	})
	result = append(result, db.Things{
		Name:   "laptop",
		Useful: true,
		ThingDetails: db.Details{
			Quality: db.Quality{
				Look: 2,
				Real: 0,
			},
			CanBeSold: 2,
		},
	})
	return result
}
