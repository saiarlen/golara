package helpers

type ResponseCode struct {
	Msg string `json:"msg"`
}

var responseCodes = map[string]ResponseCode{
	"SA1001": {"Otp sent to mobile number"},
	"SA1002": {"User found"},
}

func GetResCode(rCode string) ResponseCode {
	resCode := responseCodes[rCode]
	return resCode
}
