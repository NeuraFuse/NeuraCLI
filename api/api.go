package api

import (
	infra "github.com/neurafuse/neurakube/infrastructure"
	"github.com/neurafuse/tools-go/api/client"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	// client.F.ResetCaches(client.F{}) // TODO: Bug Refactor + only execute if recreate/delete
	infra.F.CheckDeploymentStatus(infra.F{})
	if !routeAssistant {
		if len(cliArgs) > 1 {
			if cliArgs[1] == "inspect" { // TODO: Ref. + api inspect currently not covered by assistant
				client.F.Inspect(client.F{})
			}
		}
	}
	objects.CallStructInterfaceFuncByName(Packages{}, strings.Title("ciapi"), "Router", cliArgs, routeAssistant)
}
