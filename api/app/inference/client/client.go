package client

import (
	"fmt"

	"../../../../../tools-go/api/client"
	"../../../../../tools-go/env"
	"../../../../../tools-go/errors"
	"../../../../../tools-go/logging"
	"../../../../../tools-go/runtime"
	"../../../../../tools-go/terminal"
	"../../../../../tools-go/vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool, appID string) {
	request := f.getRequest(cliArgs, routeAssistant, appID)
	response := f.sendRequest(appID, request)
	f.processResponse(response)
}

func (f F) processResponse(response string) {
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Reponse:", 0)
	fmt.Println(response)
}

func (f F) sendRequest(appID string, request interface{}) string {
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiProcess}, "Sending request to "+vars.NeuraKubeName+"..", 0)
	return client.F.Router(client.F{}, context+"/gpt", "GET", context+"/"+appID, "", "", request.(string), nil)
}

func (f F) getRequest(cliArgs []string, routeAssistant bool, appID string) interface{} {
	var request interface{}
	dataType := f.getRequestDataType(cliArgs, routeAssistant, appID)
	switch dataType {
	case "text":
		request = f.getRequestTextData(cliArgs, routeAssistant, appID)
	default:
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unsupported data type: "+dataType, true, true, true)
	}
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiProcess}, "Preparing "+context+" for app "+appID+"..", 0)
	return request
}

func (f F) getRequestDataType(cliArgs []string, routeAssistant bool, appID string) string {
	var dataType string
	if routeAssistant || len(cliArgs) < 2 {
		dataType = terminal.GetUserInput("What data type should your request for the app " + appID + " have?")
	} else {
		dataType = cliArgs[1]
	}
	return dataType
}

func (f F) getRequestTextData(cliArgs []string, routeAssistant bool, appID string) string {
	var request string
	if routeAssistant || len(cliArgs) < 3 {
		request = terminal.GetUserInput("What context do you want to send to the app " + appID + "?")
	} else {
		request = cliArgs[2]
	}
	return request
}