package api

import (
	"../../tools-go/api/client"
	"../../tools-go/config"
	"../../tools-go/config/project"
	"../../tools-go/objects"
	"../../tools-go/objects/strings"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	f.checkSettings()
	client.F.ResetCaches(client.F{}) // TODO: Refactor + only execute if recreate/delete
	objects.CallStructInterfaceFuncByName(Packages{}, strings.Title("ciapi"), "Router", cliArgs, routeAssistant)
}

func (f F) checkSettings() {
	if !config.ValidSettings("dev", "Spec.", true) {
		project.F.SetSpec(project.F{})
	}
}
