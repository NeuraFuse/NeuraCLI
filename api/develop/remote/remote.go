package remote

import (
	"../../../../tools-go/env"
	"../../../../tools-go/runtime"
	"../../../../tools-go/terminal"
	"../../base"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool, appID string) {
	base.F.Router(base.F{}, context, appID, f.getAction(cliArgs, routeAssistant, appID))
}

func (f F) getAction(cliArgs []string, routeAssistant bool, appID string) string {
	var action string
	if appID != "" {
		if routeAssistant || len(cliArgs) < 3 {
			action = terminal.GetUserSelection("Which "+context+" action do you want to start for the app "+appID+"?", []string{"create", "recreate", "update", "delete"}, false, false)
		} else {
			action = cliArgs[2]
		}
	} else {
		action = "sync"
	}
	return action
}