package cli

import (
	"github.com/c-bata/go-prompt"
	"github.com/neurafuse/neuracli/router/packages"
	"github.com/neurafuse/tools-go/build"
	buildconfig "github.com/neurafuse/tools-go/config/build"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
	"github.com/spf13/cobra"
)

type F struct{}

var ShellDescription string = "Open " + vars.NeuraCLIName + " shell"
var AssistantDescription string = "Start the assistant"
var ResourceManagerDesc string = "Resource Manager (Users, projects, infra.)"
var infraDescription string = "Manage your cluster setup."
var clusterDescription string = "Manage your existing clusters."
var gcloudDescription string = "Select the infrastructure provider gcloud."
var getDescription string = "Get resources of module."
var inspectDescription string = "Inspect resources of module."
var createDescription string = "Starts creation of module."
var recreateDescription string = "Starts recreation of module."
var deleteDescription string = "Starts deletion of module."
var apiDescription string = "Interact with " + vars.NeuraKubeName + " API."
var devDescription string = "Develop an application within your cluster."
var appDescription string = "Manage your cluster apps."
var cloudDescription string = "Interact with the " + vars.OrganizationName + " Cloud."
var ciDescription string = "Interact with the CI module."
var ciBuildCheckDisableDesc string = "Disable build check on startup."
var ciBuildHandoverDesc string = "Handover mode: If a build starts another build."
var exitDescription string = "Exit shell"
var SettingsDescription string = "Go to CLI settings"
var HelpDescription string = "Open Help"
var ExitDescription string = "Close " + vars.NeuraCLIName

func (f F) Router() {
	var cmdAssistant = &cobra.Command{
		Use:   "assistant",
		Short: AssistantDescription,
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			f.route("", args, true)
		},
	}
	var cmdShell = &cobra.Command{
		Use:   "shell",
		Short: ShellDescription,
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			f.RouteShell()
		},
	}
	var cmdInfraShort = &cobra.Command{
		Use:   "infra",
		Short: "Shortcut for infrastructure",
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"infrastructure"}, args...)
			f.route("infrastructure", args, false)
		},
	}
	var cmdInfra = &cobra.Command{
		Use:   "infrastructure",
		Short: infraDescription,
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"infrastructure"}, args...)
			f.route("infrastructure", args, false)
		},
	}
	var cmdCluster = &cobra.Command{
		Use:   "cluster",
		Short: clusterDescription,
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"cluster"}, args...)
			f.route("cluster", args, false)
		},
	}
	var cmdAPI = &cobra.Command{
		Use:   "api",
		Short: apiDescription,
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"api"}, args...)
			f.route("api", args, false)
		},
	}
	var cmdDev = &cobra.Command{
		Use:     "develop",
		Aliases: []string{"dev"},
		Short:   devDescription,
		Long:    ``,
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"develop"}, args...)
			f.route("develop", args, false)
		},
	}
	var cmdApp = &cobra.Command{
		Use:   "app",
		Short: appDescription,
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"app"}, args...)
			f.route("app", args, false)
		},
	}
	var cmdCloud = &cobra.Command{
		Use:   "cloud",
		Short: cloudDescription,
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"cloud"}, args...)
			f.route("cloud", args, false)
		},
	}
	// TODO: ? var depUpdate bool
	var cmdCI = &cobra.Command{
		Use:   "ci",
		Short: ciDescription,
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			args = append([]string{"ci"}, args...)
			f.route("ci", args, false)
		},
	}
	var rootCmd = &cobra.Command{Use: vars.NeuraCLINameID}
	var buildCheckDisabled bool
	rootCmd.PersistentFlags().BoolVarP(&buildCheckDisabled, build.F.GetFlags(build.F{})["build"][0], "b", false, ciBuildCheckDisableDesc) // TODO: Refactor
	var buildHandover bool
	rootCmd.PersistentFlags().BoolVarP(&buildHandover, build.F.GetFlags(build.F{})["build"][1], "c", false, ciBuildHandoverDesc)
	//cmdCI.Flags().BoolVar(&depUpdate, "dep-update", false, "Update all build dependencies.")
	buildconfig.F.Setting(buildconfig.F{}, "set", "check", !buildCheckDisabled)
	rootCmd.AddCommand(cmdAssistant, cmdShell, cmdInfra, cmdInfraShort, cmdCluster, cmdAPI, cmdDev, cmdApp, cmdCloud, cmdCI)
	rootCmd.Execute()
}

func (f F) RouteShell() {
	var cliArgs []string
	var assistant bool
	cliArgs, assistant = f.Autocomplete()
	f.route(cliArgs[0], cliArgs, assistant)
}

func (f F) Autocomplete() ([]string, bool) {
	logging.Log([]string{"\n", vars.EmojiAssistant, vars.EmojiInfo}, "You can activate the assistant at any time if you hit enter.", 0)
	var cliArgs string = prompt.Input(vars.NeuraCLINameID+" > ", f.AutocompletePrompt)
	var assistant bool
	if cliArgs == "exit" {
		terminal.Exit(0, "")
	} else if len(cliArgs) == 1 {
		assistant = true
	}
	return strings.Split(cliArgs, " "), assistant
}

func (f F) AutocompletePrompt(d prompt.Document) []prompt.Suggest {
	opts := []prompt.Suggest{
		{Text: "assistant", Description: AssistantDescription},
		{Text: "infrastructure", Description: infraDescription},
		{Text: "api", Description: apiDescription},
		{Text: "develop", Description: devDescription},
		{Text: "app", Description: appDescription},
		{Text: "cluster", Description: clusterDescription},

		{Text: "get", Description: getDescription},
		{Text: "inspect", Description: inspectDescription},
		{Text: "create", Description: createDescription},
		{Text: "recreate", Description: recreateDescription},
		{Text: "delete", Description: deleteDescription},
		{Text: "exit", Description: exitDescription},
	}
	return prompt.FilterFuzzy(opts, d.GetWordBeforeCursor(), true)
}

func (f F) route(packageName string, cliArgs []string, routeAssistant bool) {
	if packageName == "" || packageName == "assistant" {
		routeAssistant = true
	}
	if routeAssistant {
		packageName = f.GetPackageName(cliArgs)
	}
	objects.CallStructInterfaceFuncByName(packages.Packages{}, strings.Title(packageName), "Router", cliArgs, routeAssistant)
}

func (f F) GetPackageName(cliArgs []string) string {
	var packageName string
	packageName = terminal.GetUserSelection("What is your intention?", f.getBasePackagesArray(), false, false)
	return packageName
}

func (f F) getBasePackages() map[string]string {
	var basePackages map[string]string
	var basePackagesAr []string = f.getBasePackagesArray()
	for i := 0; i <= 4; i++ {
		basePackages["module-"+strings.ToString(i)] = basePackagesAr[i]
	}
	return basePackages
}

func (f F) getBasePackagesArray() []string {
	var basePackagesAr []string = []string{"infrastructure", "api", "develop", "app", "cluster"}
	return basePackagesAr
}
