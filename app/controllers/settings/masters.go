package settings

import (
	"ekycapp/app/models"
	"ekycapp/config"
	"encoding/json"
	"fmt"
)

type Masters struct {
	Option            string
	Email             string
	Sms               string
	NotifyAdminMobile []string
	NotifyAdminEmail  []string
}

// Global variable to cache masters data
var mastersCache *Masters

func GetMasters() (*Masters, error) {

	ResetMastersCache()
	if mastersCache != nil {
		return mastersCache, nil // Return cached data
	}

	// Fetch from DB if cache is empty
	var mastersList []models.ApxMasters
	err := config.DB.Find(&mastersList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch masters: %v", err)
	}

	// Map database values to struct
	mastersCache = &Masters{}
	for _, master := range mastersList {
		switch master.OptionKey {
		case "option":
			mastersCache.Option = master.OptionValue

		case "email":
			mastersCache.Email = master.OptionValue

		case "sms":
			mastersCache.Sms = master.OptionValue

		case "notify_admin_mobile":
			var result []string
			json.Unmarshal([]byte(master.OptionValue), &result)
			mastersCache.NotifyAdminMobile = result

		case "notify_admin_email":
			var result []string
			json.Unmarshal([]byte(master.OptionValue), &result)
			mastersCache.NotifyAdminEmail = result

		default:
			return nil, fmt.Errorf("masters.go: Unhandled OptionKey: %s error: %v", master.OptionKey, err)

		}
	}

	return mastersCache, nil
}

func ResetMastersCache() {
	mastersCache = nil
}
