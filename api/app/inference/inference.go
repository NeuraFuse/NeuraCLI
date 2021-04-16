package inference

import (
	baseCI "../../../../neurakube/infrastructure/ci/base"
	"../../../../tools-go/env"
	"../../../../tools-go/runtime"
	"../../../../tools-go/terminal"
	"../../base"
	"./client"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool, appID string) {
	action := f.getAction(cliArgs, routeAssistant, appID)
	switch action {
	case "request":
		client.F.Router(client.F{}, cliArgs, routeAssistant, appID)
	default:
		f.server(appID, action)
	}
}

func (f F) getAction(cliArgs []string, routeAssistant bool, appID string) string {
	var action string
	if routeAssistant || len(cliArgs) < 2 {
		action = terminal.GetUserSelection("Which "+context+" action do you want to start for the app "+appID+"?", []string{"client", "create", "recreate", "update", "delete"}, false, false)
	} else {
		action = cliArgs[1]
	}
	return action
}

func (f F) server(appID, action string) {
	base.F.Prepare(base.F{}, context, action)
	if action == "create" || action == "cr" || action == "update" || action == "up" {
		f.Create(appID)
	} else if action == "recreate" || action == "re" {
		f.Recreate(appID)
	} else if action == "delete" || action == "del" {
		f.Delete()
	}
}

func (f F) Create(appID string) {
	base.F.Create(base.F{}, context, appID, baseCI.F.GetResType(baseCI.F{}, context))
}

func (f F) Recreate(appID string) {
	f.Delete()
	f.Create(appID)
}

func (f F) Delete() {
	base.F.Delete(base.F{}, context, baseCI.F.GetResType(baseCI.F{}, context))
}