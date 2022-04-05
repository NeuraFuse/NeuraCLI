package inference

import (
	"github.com/neurafuse/neuracli/api/app/inference/client"
	"github.com/neurafuse/neuracli/api/base"
	baseCI "github.com/neurafuse/tools-go/ci/base"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool, appID string) {
	var action string = f.getAction(cliArgs, appID)
	switch action {
	case "request":
		client.F.Router(client.F{}, cliArgs, routeAssistant, appID)
	default:
		f.server(cliArgs, appID, action)
	}
}

func (f F) getAction(cliArgs []string, appID string) string {
	var action string
	if len(cliArgs) < 2 {
		action = terminal.GetUserSelection("Which "+context+" action do you intend to start for the app "+appID+"?", []string{"client", "create", "recreate", "update", "delete"}, false, false)
	} else {
		action = cliArgs[1]
	}
	return action
}

func (f F) server(cliArgs []string, appID, action string) {
	if action == "create" || action == "cr" || action == "update" || action == "up" {
		f.Create(cliArgs, appID)
	} else if action == "recreate" || action == "re" {
		f.Recreate(cliArgs, appID)
	} else if action == "delete" || action == "del" {
		f.Delete()
	}
}

func (f F) Create(cliArgs []string, appID string) {
	base.F.Create(base.F{}, context, f.getAction(cliArgs, appID), appID, baseCI.F.GetResType(baseCI.F{}, context))
}

func (f F) Recreate(cliArgs []string, appID string) {
	f.Delete()
	f.Create(cliArgs, appID)
}

func (f F) Delete() {
	base.F.Delete(base.F{}, context, baseCI.F.GetResType(baseCI.F{}, context))
}
