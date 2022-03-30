package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TargetForSMS struct {
	TargetName        string `json:"nameofphone"`
	TargetPhoneNumber string `json:"phone_num"`
	TextForTarget     string `json:"phone_smstext"`
}

func router_post(ctx *gin.Context) {

	incoming_request := new(TargetForSMS)

	if err := ctx.BindJSON(&incoming_request); err != nil {

		// Error Message for user "HTTP = 200"
		bad_request := make(map[string]string)

		bad_request["HTTP_CODE"] = "400"
		bad_request["HTTP_TEXT"] = "Something went wrong checking your JSON BODY"

		ctx.IndentedJSON(http.StatusBadRequest, bad_request)

	} else {

		// Response Message for user "HTTP = 200"
		ctx.IndentedJSON(http.StatusCreated, incoming_request)

	}
}

func Validate(ctx *gin.Context, targetForSMS *TargetForSMS) {
	panic("unimplemented")
}

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")

	v1.POST("/post", router_post)

	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Run("127.0.0.1:5000")
}
