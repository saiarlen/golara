package providers

import (
	"ekycapp/app/controllers/settings"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// SendSMS sends an SMS message.
//
// @param mnum - phone number of the recipient, templateId, MSG
// @param vars - optional parameters to pass to the SMS provider such as x
func SendSMS(mnum, id, msg string, vars ...interface{}) (bool, error) {

	var1 := "x"
	var2 := "y"

	if mnum == "" {
		return false, fmt.Errorf("empty mobile number: func:SendSMS")
	}

	if len(vars) > 0 {
		var1 = vars[0].(string)
	}
	if len(vars) > 1 {
		var2 = vars[1].(string)
	}
	masters, err := settings.GetMasters()
	if err != nil {
		return false, err
	}
	switch masters.Sms {
	case "SMSCOUNTRY":
		smsCountry(mnum, msg)
	case "MSG91":
		msgNOne(mnum, id, var1, var2)
	case "ALERTWINGS":
		altWings(mnum, id, msg)
	default:
		return false, fmt.Errorf("unsupported SMS provider :func SendSMS")
	}
	return true, nil
}

// smsCountry sends a message to the sms server.
//
// @param mnum - phone number of the person who sent the message
//
// @return error if there was a problem communicating with the server otherwise nil is returned.
func smsCountry(mnum, msg string) error {

	creds, err := settings.GetCreds()
	if err != nil {
		return err
	}

	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(creds.Sms.SmsUrl)

	args := fiber.AcquireArgs()
	defer fiber.ReleaseArgs(args)

	args.Set("User", creds.Sms.SmsUser)
	args.Set("passwd", creds.Sms.SmsPass)
	args.Set("mobilenumber", mnum)
	args.Set("message", msg)
	args.Set("sid", creds.Sms.SmsSid)
	args.Set("mtype", "N")
	args.Set("DR", "Y")

	agent.MultipartForm(args)

	if err := agent.Parse(); err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	resCode, body, _ := agent.String()
	if resCode == fiber.StatusOK {
		fmt.Println("Response Body:", body)
		return nil
	} else {
		return fmt.Errorf("SMSCOUNTRY: failed with status code: %d", resCode)
	}

}

// Send a message from msg91.
//
// @param mnum - the mobiles to send to
//
// @return nil if success error if failure see Fiber.
func msgNOne(mnum, id, var1, var2 string) error {
	creds, err := settings.GetCreds()
	if err != nil {
		return err
	}
	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)

	req := agent.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(creds.Sms.SmsUrl)

	req.Header.Add("authkey", creds.Sms.SmsPass)

	agent.JSON(
		fiber.Map{
			"template_id": id,
			"short_url":   "0",
			"recipients": []fiber.Map{
				{
					"mobiles": "91" + mnum,
					"var":     var1,
					"var1":    var2,
				},
			},
		})

	if err := agent.Parse(); err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	resCode, body, _ := agent.String()
	if resCode == fiber.StatusOK {
		fmt.Println("Response Body:", body)
		return nil
	} else {
		return fmt.Errorf("MSG91: failed with status code: %d", resCode)
	}
}

// Send a wings message from alertwings.
//
// @param mnum - phone number of the recipient
//
// @return error if there was a problem with the request or nil
func altWings(mnum, id, msg string) error {
	creds, err := settings.GetCreds()
	if err != nil {
		return err
	}
	agent := fiber.Get(creds.Sms.SmsUrl + "?authkey=" + creds.Sms.SmsPass + "&sender=" + creds.Sms.SmsSid + "&mobiles=" + mnum + "&message=" + msg + "&route=4&country=91&DLT_TE_ID=" + id)

	if err := agent.Parse(); err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	resCode, body, _ := agent.String()

	fmt.Println(err)

	if resCode == fiber.StatusOK {
		fmt.Println("Response Body:", body)
		return nil
	} else {
		return fmt.Errorf("ALTWINGS: failed with status code: %d", resCode)
	}
}
