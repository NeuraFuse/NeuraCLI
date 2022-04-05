package app

import (
	"github.com/neurafuse/tools-go/ci"
	"github.com/neurafuse/tools-go/config"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
	"github.com/neurafuse/neuracli/api/app/id"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) GetID(cliArgs []string, routeAssistant bool) (string, []string) {
	var appID string = ci.F.GetContextID(ci.F{})
	var appKind string
	var appKindConfigKey string = "Spec.App.Kind"
	if config.ValidSettings("project", "app", true) {
		appKind = config.Setting("get", "project", appKindConfigKey, "")
	} else {
		var opts []string = []string{"A " + vars.NeuraKubeNameID + " managed app", "A custom managed"}
		appKind = terminal.GetUserSelection("What kind of app do you want to develop?", opts, false, false)
		config.Setting("set", "project", appKindConfigKey, appKind)
	}
	id.F.SetKind(id.F{}, appID, appKind)
	return appID, cliArgs
}