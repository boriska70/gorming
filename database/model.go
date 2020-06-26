package database

import (
	"database/sql/driver"
	"encoding/json"
)

//Things describes the table for gorming experiments
type Things struct {
	ID           int64   //`gorm:"type:bigint,primary_key"`
	Name         string  //`gorm:"type:text"`
	Useful       bool    //`gorm:"type:bool"`
	ThingDetails Details `json:"thing_details" sql:"type:jsonb"` //json:"details"
}

// func (t Things) String() string {
// 	return "things"
// }

//Details keeps data about things quality and ability to sell it
type Details struct {
	Quality   Quality `json:"quality"`
	CanBeSold int     `json:"can_be_sold"`
}

//Quality keeps data about real quality of the thing and does it look good
type Quality struct {
	Look int `json:"look"`
	Real int `json:"real"`
}

//Value converts struct to text for storing to DB
func (d Details) Value() (driver.Value, error) {
	valueBytes, err := json.Marshal(d)
	return string(valueBytes), err
}

//Scan converts stored text to struct
func (d *Details) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &d); err != nil {
		return err
	}
	return nil
}
