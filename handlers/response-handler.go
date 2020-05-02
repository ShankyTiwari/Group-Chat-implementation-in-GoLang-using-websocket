package handlers

import (
	"encoding/json"
	"net/http"

	"../constant"
)

type apiResponseStruct struct {
	Code     int
	Status   string
	Message  string
	Response interface{}
}

// ReturnResponse will be used as Response template to send the response for API
func ReturnResponse(response http.ResponseWriter, request *http.Request, apiResponse apiResponseStruct) {
	var responseMessage, responseStatusText string

	if apiResponse.Message != "" {
		responseMessage = apiResponse.Message
	} else {
		responseMessage = constant.SuccessfulResponse
	}

	if apiResponse.Status != "" {
		responseStatusText = apiResponse.Status
	} else {
		responseStatusText = http.StatusText(http.StatusOK)
	}

	httpResponse := &apiResponseStruct{
		Code:     apiResponse.Code,
		Status:   responseStatusText,
		Message:  responseMessage,
		Response: apiResponse.Response,
	}
	jsonResponse, err := json.Marshal(httpResponse)
	if err != nil {
		panic(err)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(httpResponse.Code)
	response.Write(jsonResponse)
}
