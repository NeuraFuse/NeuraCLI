package develop

import (
	"../../../tools-go/env"
	"../../../tools-go/runtime"
	"../../../tools-go/terminal"
	"../app"
	"./remote"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool) {
	appID := app.F.GetID(app.F{}, cliArgs, routeAssistant)
	module := f.getModule(cliArgs, routeAssistant)
	if module == "remote" || module == "re" {
		remote.F.Router(remote.F{}, cliArgs, routeAssistant, appID)
	}
}

func (f F) getModule(cliArgs []string, routeAssistant bool) string {
	var module string
	if routeAssistant || len(cliArgs) < 2 {
		module = terminal.GetUserSelection("Which "+context+" workflow do you want to start?", []string{"remote"}, false, false)
	} else {
		module = cliArgs[1]
	}
	return module
}