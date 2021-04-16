package packages

import (
	"../../../neurakube/infrastructure"
	"../../../neurakube/infrastructure/ci"
	"../../../tools-go/kubernetes"
	"../../api"
	"../../api/app"
	"../../api/develop"
	"../../api/app/inference"
	"../../cloud"
)

type Packages struct {
	Infrastructure infrastructure.F
	Kubernetes     kubernetes.F
	Api            api.F
	Develop        develop.F
	App            app.F
	Inference      inference.F
	Cloud          cloud.F
	Ci             ci.F
}