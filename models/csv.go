package models

import (
	"github.com/jinzhu/gorm"
)

type CSV struct {
	gorm.Model
	SiteID                int    `json:"site_id"`
	FxiletID              int    `json:"fxilet_id"`
	Name                  string `json:"name"`
	Criticality           string `json:"criticality"`
	RelevantComputerCount int    `json:"relevant_computer_count"`
}
