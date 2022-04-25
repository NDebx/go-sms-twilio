package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// Create a struct that represents the JSON Payload that will be sent to the end user
type TargetForSMS struct {
	TargetName        string `json:"nameofphone"`
	TargetPhoneNumber string `json:"phone_num"`
	TextForTarget     string `json:"phone_smstext"`
}

func SendSMS(TargetName string, TargetPhoneNumber string, TextForTarget string) {

	// Creating a viper instance
	viper_ := viper.New()
	viper_.SetConfigFile("config.yml")
	viper_.ReadInConfig()

	// Your Twilio Account SID and TOKEN
	accountSid := viper_.GetString("twilioAccountSid")
	authToken := viper_.GetString("twilioAccountAuthToken")

	// Create a client instance
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// Creating Sending Payload for sms-ing your end user.
	params := &openapi.CreateMessageParams{}
	params.SetTo(TargetPhoneNumber)
	params.SetFrom(viper_.GetString("twilioPhoneNumber"))
	params.SetBody(TextForTarget)

	// Send SMS Payload
	resp, err := client.ApiV2010.CreateMessage(params)
	if err != nil {

		// If error occurs the console will print the error
		fmt.Println(err.Error())
		err = nil

	} else {

		// Message was successfully sent
		fmt.Println("Message Sid: " + *resp.Sid)
		fmt.Printf("Mail send to: %v", TargetName)

	}

}

func RouterPost(ctx *gin.Context) {

	// The incoming JSON data
	IncomingRequest := new(TargetForSMS)

	if err := ctx.BindJSON(&IncomingRequest); err != nil {

		// Error Message for user "HTTP = 200"
		BadRequest := make(map[string]string)

		// Custom error message for the end client
		BadRequest["HTTP_CODE"] = "400"
		BadRequest["HTTP_TEXT"] = "Something went wrong checking your JSON BODY"

		// Response Message for user "HTTP = 400"
		ctx.IndentedJSON(http.StatusBadRequest, BadRequest)

	} else {

		// Send SMS to client
		// This function need the required JSON data, that you passed via the REST endpoint. In my case http://127.0.0.1:5000/api/v1/post/send/sms/client
		go SendSMS(IncomingRequest.TargetName, IncomingRequest.TargetPhoneNumber, IncomingRequest.TextForTarget)

		// Response Message for user "HTTP = 200"
		ctx.IndentedJSON(http.StatusCreated, IncomingRequest)

	}
}

func main() {

	// Create the Go Gin engine instance
	router := gin.Default()

	// Create an API group
	API_v1 := router.Group("/api")

	// Provide and REST endpoint
	API_v1.POST("/v1/post/send/sms/client", RouterPost)

	// Run the Router
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Run("127.0.0.1:8080")

}
