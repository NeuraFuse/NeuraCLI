package packages

import (
	"github.com/neurafuse/neuracli/api"
	"github.com/neurafuse/neuracli/api/app"
	"github.com/neurafuse/neuracli/api/app/inference"
	"github.com/neurafuse/neuracli/api/develop"
	"github.com/neurafuse/neuracli/cloud"
	"github.com/neurafuse/neurakube/infrastructure"
	"github.com/neurafuse/tools-go/ci"
	resInfras "github.com/neurafuse/tools-go/infrastructures"
	"github.com/neurafuse/tools-go/kubernetes"
	"github.com/neurafuse/tools-go/projects"
	"github.com/neurafuse/tools-go/users"
)

type Packages struct {
	Infrastructure infrastructure.F
	Cluster        kubernetes.F
	Api            api.F
	Develop        develop.F
	App            app.F
	Inference      inference.F
	Cloud          cloud.F
	Ci             ci.F
}

type ResourceTypes struct {
	Users           users.F
	Projects        projects.F
	Infrastructures resInfras.F
}
