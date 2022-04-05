package develop

import (
	"github.com/neurafuse/neuracli/api/app"
	appIDTool "github.com/neurafuse/neuracli/api/app/id"
	"github.com/neurafuse/neuracli/api/base"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool) {
	var appID string
	appID, cliArgs = app.F.GetID(app.F{}, cliArgs, routeAssistant)
	base.F.Router(base.F{}, context, appID, f.getAction(cliArgs, routeAssistant, appID))
}

func (f F) getAction(cliArgs []string, routeAssistant bool, appID string) string {
	var action string
	if routeAssistant || len(cliArgs) < 2 {
		var opts []string
		if appIDTool.F.IsKindManaged(appIDTool.F{}, appID) {
			opts = []string{"create", "sync", "recreate", "update", "delete"}
		} else if appIDTool.F.IsKindCustom(appIDTool.F{}, appID) {
			opts = []string{"sync"}
		}
		var appKind string = appIDTool.F.GetKind(appIDTool.F{}, appID)
		action = terminal.GetUserSelection("Which "+context+" action do you intend to start for the "+appKind+" app "+appID+"?", opts, false, false)
	} else {
		action = cliArgs[1]
	}
	return action
}
