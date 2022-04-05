package main

import (
	"github.com/neurafuse/neuracli/router"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/vars"

func main() {
	env.F.SetFramework(env.F{}, vars.NeuraCLINameID)
	router.F.Router(router.F{})
}
