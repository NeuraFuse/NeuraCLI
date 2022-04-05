package router

import (
	"os"

	"github.com/neurafuse/neuracli/cli"
	"github.com/neurafuse/neuracli/router/packages"
	"github.com/neurafuse/tools-go/build"
	"github.com/neurafuse/tools-go/config"
	cliconfig "github.com/neurafuse/tools-go/config/cli"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/installer"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/metrics"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/projects"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/updater"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
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
	f.checkUser()
	installer.F.CheckLocalSetup(installer.F{})
	if !skip {
		f.greetUser()
		f.checkDevMode()
		updater.F.Check(updater.F{})
	}
	projects.F.CheckConfigs(projects.F{})
	logging.Log([]string{"", vars.EmojiAstronaut, vars.EmojiSuccess}, "Ready to go.\n", 0)
}

func (f F) routerAssistant() {
	var cliArgs []string = os.Args
	var sel string = terminal.GetUserSelection("What is your intention?", f.getMainMenuOpts(), false, false)
	if sel == cli.ShellDescription {
		cli.F.RouteShell(cli.F{})
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

func (f F) getMainMenuOpts() []string {
	var opts []string
	opts = []string{cli.AssistantDescription, cli.ShellDescription, cli.ResourceManagerDesc, cli.SettingsDescription, cli.ExitDescription}
	return opts
}

func (f F) resourceManager(cliArgs []string) {
	var selOpts []string = config.GetResourceTypes()
	var sel string = terminal.GetUserSelection("Which resource do you want to manage?", selOpts, false, false)
	objects.CallStructInterfaceFuncByName(packages.ResourceTypes{}, sel, "Router", cliArgs, true)
}

func (f F) routeHelp() {
	logging.Log([]string{"\n", vars.EmojiClient, vars.EmojiInfo}, vars.NeuraCLIName+" help\n", 0)
	os.Args = []string{"", "help"}
	f.Router()
}

func (f F) route(packageName string, cliArgs []string, routeAssistant bool) {
	objects.CallStructInterfaceFuncByName(packages.Packages{}, strings.Title(packageName), "Router", cliArgs, routeAssistant)
}

func (f F) checkUser() {
	if !config.ValidSettings("cli", "users", false) {
		logging.Log([]string{"", vars.EmojiUser, vars.EmojiInfo}, "In order to use "+vars.NeuraCLIName+" you have to configure a user account.", 0)
		usersID.F.CreateNew(usersID.F{})
		config.Setting("set", "cli", "Spec.Users.DefaultID", usersID.F.GetActive(usersID.F{}))
	} else {
		var userName string
		userName = config.Setting("get", "cli", "Spec.Users.DefaultID", "")
		usersID.F.SetActive(usersID.F{}, userName)
	}
}

func (f F) greetUser() {
	if env.F.CLI(env.F{}) {
		logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiWavingHand}, "Logged in "+usersID.F.GetActive(usersID.F{})+".", 0)
		logging.Log([]string{"", vars.EmojiAssistant, vars.EmojiThumbsUp}, "Wishing you a good "+timing.GetDayTime()+" hustle.\n", 0)
	}
}

func (f F) checkDevMode() {
	if config.DevConfigActive() {
		logging.Log([]string{"", vars.EmojiDev, vars.EmojiSettings}, "Developer mode is active.", 0)
		if env.F.IsFrameworkActive(env.F{}, vars.NeuraCLINameID) {
			metrics.DevStats(vars.OrganizationNameRepo)
		}
		var ciMode string = config.Setting("get", "dev", "Spec.CI.Mode", "")
		if ciMode == "auto" {
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiWarning}, "CI mode is automatic (unstable).\n", 0)
			build.F.CheckUpdates(build.F{}, env.F.GetActive(env.F{}, false), true)
		} else {
			logging.Log([]string{"", vars.EmojiDev, vars.EmojiProcess}, "CI mode is manual.\n", 0)
		}
	}
}
