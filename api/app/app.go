package app

import (
	"../../../tools-go/env"
	"../../../tools-go/runtime"
	"../../../tools-go/terminal"
	"../../../tools-go/vars"
	"../base"
	"./inference"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool) {
	id := f.GetID(cliArgs, routeAssistant)
	action := f.getAction(cliArgs, routeAssistant, id)
	if action == "inference" {
		inference.F.Router(inference.F{}, cliArgs, routeAssistant, id)
	} else {
		base.F.Router(base.F{}, context, id, action)
	}
}

func (f F) GetID(cliArgs []string, routeAssistant bool) string {
	var id string
	if routeAssistant || len(cliArgs) < 2 {
		sel := terminal.GetUserSelection("Do you want to develop a "+vars.NeuraKubeName+" managed app?", []string{"Yes", "No"}, false, true)
		if sel == "Yes" {
			id = terminal.GetUserSelection("Which "+context+" do you want to manage?", f.getIDs(), false, false)
		}
	} else {
		id = cliArgs[1]
	}
	return id
}

func (f F) getIDs() []string { // TODO: Dynamic
	var ids []string
	ids = []string{"gpt"}
	return ids
}

func (f F) getAction(cliArgs []string, routeAssistant bool, id string) string {
	var action string
	if routeAssistant || len(cliArgs) < 2 {
		action = terminal.GetUserSelection("Which action do you want to start for the app "+id+"?", []string{"create", "recreate", "update", "delete", "inference"}, false, false)
	} else {
		action = cliArgs[1]
	}
	return action
}