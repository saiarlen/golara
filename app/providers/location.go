package providers

import (
	"ekycapp/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetLocation(lattitude, longitude string) (string, error) {

	agent := fiber.Get("https://api.opencagedata.com/geocode/v1/json?q=" + lattitude + "%2C+" + longitude + "&key=" + config.Denv("GOOGLE_LOCATION_KEY") + "&pretty=1")

	if err := agent.Parse(); err != nil {
		return "", err
	}

	resCode, responseBody, errs := agent.String()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if resCode != fiber.StatusOK {
		return "", fmt.Errorf(responseBody)
	}

	return responseBody, nil

}
