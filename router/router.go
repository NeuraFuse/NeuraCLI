package router

import (
	"fmt"
	"os"

	"../../tools-go/build"
	"../../tools-go/config"
	cliconfig "../../tools-go/config/cli"
	"../../tools-go/env"
	"../../tools-go/filesystem"
	"../../tools-go/installer"
	"../../tools-go/logging"
	"../../tools-go/metrics"
	"../../tools-go/objects"
	"../../tools-go/objects/strings"
	"../../tools-go/projects"
	"../../tools-go/terminal"
	"../../tools-go/timing"
	"../../tools-go/updater"
	"../../tools-go/users"
	"../../tools-go/vars"
	"../cli"
	"./packages"
)

type F struct{}

func (f F) Router() {
	f.startup()
	if len(os.Args) == 1 {
		f.routerAssistant()
	} else {
		cli.F.Router(cli.F{})
	}
	terminal.Exit(0, "")
}

func (f F) startup() {
	var exceptions []string = []string{"cluster"}
	var skip bool
	if len(os.Args) > 1 {
		if strings.ArrayContains(exceptions, os.Args[1]) {
			skip = true
			logging.Log([]string{"", vars.EmojiAstronaut, vars.EmojiSuccess}, "Skipped startup routines.\n", 0)
		}
	}
	terminal.Init(skip)
	f.checkUser(false)
	installer.F.CheckLocalSetup(installer.F{})
	f.checkUser(true)
	if !skip {
		f.greetUser()
		f.checkDevconfig()
		build.F.CheckUpdates(build.F{}, env.F.GetActive(env.F{}, false), true)
		updater.F.Check(updater.F{})
	}
	projects.F.CheckConfigs(projects.F{})
	logging.Log([]string{"", vars.EmojiAstronaut, vars.EmojiSuccess}, "Ready to go.\n", 0)
}

func (f F) routerAssistant() {
	cliArgs := os.Args
	sel := terminal.GetUserSelection("What do you want to do?", []string{cli.AssistantDescription, cli.ResourceManagerDesc, cli.ShellDescription, cli.SettingsDescription, cli.ExitDescription}, false, false)
	if sel == cli.ShellDescription {
		f.routeShellAutocomplete()
	} else if sel == cli.AssistantDescription {
		f.route(cli.F.GetPackageName(cli.F{}, cliArgs), cliArgs, true)
	} else if sel == cli.ResourceManagerDesc {
		f.resourceManager(cliArgs)
	} else if sel == cli.SettingsDescription {
		cliconfig.F.Configure(cliconfig.F{})
	} else if sel == cli.ExitDescription {
		terminal.Exit(0, "")
	}
	f.routerAssistant()
}

func (f F) resourceManager(cliArgs []string) {
	sel := terminal.GetUserSelection("What do you want to do?", []string{cli.UsersDescription, cli.InfraDescription, cli.ProjectsDescription}, false, false)
	switch sel {
	//case cli.UsersDescription:
	//case cli.InfraDescription:
	case cli.ProjectsDescription:
		projects.F.Router(projects.F{}, cliArgs, true)
	}
}

func (f F) routeShellAutocomplete() {
	cliArgs := cli.F.Autocomplete(cli.F{})
	fmt.Println(cliArgs)
	if cliArgs[0] == "exit" {
		terminal.Exit(0, "")
	}
	f.route(cliArgs[0], cliArgs, false)
}

func (f F) routeHelp() {
	logging.Log([]string{"\n", vars.EmojiClient, vars.EmojiInfo}, vars.NeuraCLIName+" help\n", 0)
	os.Args = []string{"", "help"}
	f.Router()
}

func (f F) route(packageName string, cliArgs []string, routeAssistant bool) {
	fmt.Println(packageName)
	switch packageName {
	case "infra":
		packageName = "infrastructure"
	case "dev":
		packageName = "develop"
	case "cluster":
		packageName = "kubernetes"
	}
	objects.CallStructInterfaceFuncByName(packages.Packages{}, strings.Title(packageName), "Router", cliArgs, routeAssistant)
}

func (f F) checkUser(create bool) {
	var userName string
	if !config.ValidSettings("cli", "users", false) {
		if create {
			logging.Log([]string{"", vars.EmojiUser, vars.EmojiInfo}, "In order to use "+vars.NeuraCLIName+" you have to configure a user account.", 0)
			userName = terminal.GetUserSelection("Choose an existing user or create a new one", users.GetAllIDs(), true, false)
			config.Setting("set", "cli", "Spec.Users.DefaultID", userName)
			config.Setting("reset", "cli", "Spec.Projects.DefaultID", "")
		}
	} else {
		userName = config.Setting("get", "cli", "Spec.Users.DefaultID", "")
	}
	users.SetIDActive(userName) // TODO: Refactor
	activeUserPath := users.BasePath + "/" + userName
	vars.ProjectsBasePath = activeUserPath + "/"
	if create {
		filesystem.CreateDir(activeUserPath, false)
	}
}

func (f F) greetUser() {
	if env.F.CLI(env.F{}) {
		logging.Log([]string{"", vars.EmojiUser, vars.EmojiWavingHand}, "Welcome to "+vars.NeuraKubeName+", "+users.GetIDActive()+".", 0)
		logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiThumbsUp}, "Wishing you a good "+timing.GetDayTime()+" hustle.\n", 0)
	}
}

func (f F) checkDevconfig() {
	var status string = config.Setting("get", "dev", "Spec.Status", "")
	if status == "active" {
		logging.Log([]string{"", vars.EmojiDev, vars.EmojiSuccess}, "Developer mode "+status+".", 0)
		if env.F.ActiveFramework(env.F{}, vars.NeuraCLINameRepo) {
			metrics.DevStats(vars.OrganizationNameRepo)
		}
	}
}
