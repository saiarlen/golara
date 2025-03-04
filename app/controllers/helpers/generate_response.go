package helpers

import (
	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Rcode   string `json:"rcode"`
	Scode   int    `json:"scode"`
}

type SuccessResponse struct {
	BaseResponse
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	BaseResponse
	Errors interface{} `json:"errors"`
}

// Success sends a success response to the client.
// @param c - Context to send the response to.
// @param code - The R code to send to the client.
// @param data - The data to send to the client. Either data object or pass only nil
// @param httpCode - The http code to send to the client.

func Success(c *fiber.Ctx, code string, data interface{}, httpCode int) error {
	message := GetResCode(code)

	if data == nil {
		data = []interface{}{}
	}

	response := SuccessResponse{
		BaseResponse: BaseResponse{
			Status:  "success",
			Message: message.Msg,
			Rcode:   code,
			Scode:   1,
		},
		Data: data,
	}

	return c.Status(httpCode).JSON(response)
}

// Error sends an error response to the client.
//
// @param c - Context to send the response to
// @param code - R Code of the error to send.
// @param errors - Error values to send. Can be nil. Either data object or pass only nil
// @param httpCode - HTTP code to send to fiber.

func Error(c *fiber.Ctx, code string, errors interface{}, httpCode int) error {
	message := GetResCode(code)
	if errors == nil {
		errors = []interface{}{}
	}
	response := ErrorResponse{
		BaseResponse: BaseResponse{
			Status:  "fail",
			Message: message.Msg,
			Rcode:   code,
			Scode:   -1,
		},
		Errors: errors,
	}

	return c.Status(httpCode).JSON(response)
}
