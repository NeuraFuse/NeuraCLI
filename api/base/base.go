package base

import (
	"../../../neurakube/infrastructure/ci"
	baseCI "../../../neurakube/infrastructure/ci/base"
	"../../../tools-go/api/client"
	"../../../tools-go/apps/ide/vscode"
	"../../../tools-go/apps/kubernetes/devspace"
	"../../../tools-go/apps/tensorflow/tensorboard"
	"../../../tools-go/config"
	infraconfig "../../../tools-go/config/infrastructure"
	"../../../tools-go/container"
	"../../../tools-go/env"
	"../../../tools-go/errors"
	"../../../tools-go/kubernetes/deployments"
	"../../../tools-go/kubernetes/namespaces"
	"../../../tools-go/kubernetes/pods"
	"../../../tools-go/logging"
	"../../../tools-go/objects/strings"
	"../../../tools-go/runtime"
	"../../../tools-go/terminal"
	"../../../tools-go/vars"
	"../../../tools-go/users"
)

type F struct{}

func (f F) Router(context, appID, action string) { // TODO: Implement id
	if appID == "" {
		f.Prepare(context, action)
	}
	if action == "create" || action == "cr" || action == "update" || action == "up" {
		f.Create(context, appID, baseCI.F.GetResType(baseCI.F{}, context))
	} else if action == "recreate" || action == "re" {
		f.Recreate(context, appID)
	} else if action == "delete" || action == "del" {
		f.Delete(context, baseCI.F.GetResType(baseCI.F{}, context))
	}
}

func (f F) Prepare(context, action string) {
	logging.Log([]string{"", vars.EmojiRemote, vars.EmojiProcess}, "Preparing "+baseCI.F.GetResType(baseCI.F{}, context)+"..", 0)
	f.checkUserConfigs(context)
	client.F.Sync(client.F{}) // Update neurakube
}

func (f F) Create(context, appID, resType string) {
	logging.Log([]string{"\n", vars.EmojiRemote, vars.EmojiProcess}, "Creating "+baseCI.F.GetResType(baseCI.F{}, context)+"..", 0)
	if context == "develop/remote" { // TODO: Refactor
		context = "remote"
	}
	var namespace string
	var imageAddrs string
	var contextID string
	var remoteURL string
	var waitForStatusInLog string
	var initWaitDuration int
	if appID != "" {
		namespace, imageAddrs, contextID, remoteURL, waitForStatusInLog, initWaitDuration = f.createNeuraKubeManaged(context)
	} else {
		namespace, contextID, imageAddrs = f.createSelfManaged(context)
	}
	if context == "remote" {
		f.developApp(context, remoteURL, namespace, imageAddrs)
	}
	tensorboard.F.Info(tensorboard.F{}, remoteURL)
	logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSuccess}, "Created "+resType+".\n", 0)
	pods.F.Logs(pods.F{}, namespace, contextID, waitForStatusInLog, false, initWaitDuration)
}

func (f F) createSelfManaged(context string) (string, string, string) {
	var namespace string = terminal.GetUserSelection("In which cluster namespace is your app deployed?", namespaces.F.Get(namespaces.F{}, false), false, false)
	var contextID string = terminal.GetUserSelection("How is the app named (kubernetes deployment)?", deployments.F.GetList(deployments.F{}, namespace, false), false, false)
	var containerID int = 0
	if len(pods.F.GetContainers(pods.F{}, namespace, contextID)) > 1 {
		containerName := terminal.GetUserSelection("To which pod container do you want to connect?", pods.F.GetContainerNamesList(pods.F{}, namespace, contextID), false, false)
		containerID = pods.F.GetContainerIDByName(pods.F{}, namespace, contextID, containerName)
	}
	var imageAddrs string = pods.F.GetContainerImgAddrs(pods.F{}, namespace, contextID, containerID)
	return namespace, contextID, imageAddrs
}

func (f F) createNeuraKubeManaged(context string) (string, string, string, string, string, int) {
	var namespace string = baseCI.F.GetNamespace(baseCI.F{})
	var imageAddrs string = container.F.GetImgAddrs(container.F{}, context, false, false)
	var contextID string = ci.F.GetContextID(ci.F{}, context)
	var remoteURL string = f.sentRequest(context, "create")
	var waitForStatusInLog string = f.getWaitForStatusInLog(context)
	var initWaitDuration int = ci.F.GetInitWaitDuration(ci.F{}, context)
	f.sentRequest(context, "prepare")
	if context == "remote" { // TODO: Refactor
		context = "develop/remote"
	}
	container.F.CheckUpdates(container.F{}, context+"-base", true, false)
	//portforward.Connect("neurakube", 1000, 80)
	f.selfHostedInfra(context)
	return namespace, imageAddrs, contextID, remoteURL, waitForStatusInLog, initWaitDuration
}

func (f F) getWaitForStatusInLog(context string) string {
	var waitForStatusInLog string
	if context == "app" {
		waitForStatusInLog = "Epoch: 1"
	} else if context == "inference" {
		waitForStatusInLog = "Serving"
	}
	return waitForStatusInLog
}

func (f F) developApp(context, remoteURL, namespace, imageAddrs string) {
	f.connectIDE(context, remoteURL)
	devspace.F.Sync(devspace.F{}, context, namespace, imageAddrs, true)
}

func (f F) sentRequest(context, signal string) string {
	waitDuration := strings.ToString(ci.F.GetInitWaitDuration(ci.F{}, context))
	if signal == "prepare" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWaiting}, "Requesting "+signal+" "+context+"..", 0)
	} else if signal == "create" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWaiting}, "Sending request to "+vars.NeuraKubeName+" and waiting for the "+context+" environment to be ready (this may take up to "+waitDuration+" minutes for a new cluster)..", 0)
	}
	response := client.F.Router(client.F{}, vars.NeuraKubeNameRepo, "GET", context, "", "", "action="+signal, nil)
	if !strings.Contains(response, "Error") {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Request "+signal+" "+context+" successful.\n", 0)
	} else {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWaiting}, "Response: "+response, 0)
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Received "+vars.NeuraKubeName+" error message for "+signal+" "+context+"!", true, true, true)
	}
	return response
}

func (f F) Recreate(context, appID string) {
	f.Delete(context, baseCI.F.GetResType(baseCI.F{}, context))
	f.Create(context, appID, baseCI.F.GetResType(baseCI.F{}, context))
}

func (f F) Delete(context, resType string) {
	logging.Log([]string{"\n", vars.EmojiRemote, vars.EmojiProcess}, "Deleting "+resType+"..", 0)
	if context == "remote" { // TODO: Refactor
		context = "develop/remote"
	}
	route := context + "/delete"
	answer := client.F.Router(client.F{}, vars.NeuraKubeNameRepo, "GET", route, "", "", "", nil)
	if answer == "success" {
		logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSuccess}, "Deleted "+resType+".\n", 0)
	} else {
		err := errors.New("API route: " + route + "\nAnswer: " + answer)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete "+resType+"!", false, true, true)
	}
}

func (f F) checkUserConfigs(context string) {
	logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSettings}, "Checking "+users.GetIDActive()+" user configs..\n", 0)
	if !config.ValidSettings("infrastructure", context, true) {
		infraconfig.F.SetModuleSpec(infraconfig.F{}, context)
	}
	accType := baseCI.F.GetResources(baseCI.F{}, context)
	if !config.ValidSettings("infrastructure", vars.InfraProviderGcloud+"/accelerator/"+accType, true) {
		infraconfig.F.SetGcloudAccelerator(infraconfig.F{}, context, accType)
	}
	//f.checkDevconfig(context) TODO: Necessary?
}

func (f F) checkDevconfig(context string) {
	if config.Setting("get", "dev", "Spec.API.Address", "") != "cluster" {
		logging.Log([]string{"\n", vars.EmojiAPI, vars.EmojiWaiting}, "The "+context+" module is not available if you have activated the API localhost mode.", 0)
		sel := terminal.GetUserSelection("Do you want to switch to API cluster mode?", []string{}, false, true)
		if sel == "Yes" {
			config.Setting("set", "dev", "Spec.API.Address", "cluster")
		} else {
			terminal.Exit(0, "")
		}
	}
}

func (f F) connectIDE(context, remoteURL string) {
	if remoteURL == "" {
		logging.Log([]string{"", vars.EmojiRocket, vars.EmojiWaiting}, "Please provide information about the routing to your remote app.", 0)
		remoteURL = terminal.GetUserInput("What is the remote URL (or IP) on which the app is reachable via network?")
	}
	envIDE := config.Setting("get", "infrastructure", "Spec."+env.F.GetContext(env.F{}, context, true)+".Environment.IDE", "")
	if envIDE == "vscode" {
		vscode.F.CreateConfig(vscode.F{}, remoteURL, baseCI.F.GetContainerPortsForApps(baseCI.F{}, []string{"debugpy"})[0][0])
	} else {
		logging.Log([]string{"\n", vars.EmojiAPI, vars.EmojiWaiting}, "To get terminal live logs from your "+context+" environment execute:\nneuracli cluster logs "+baseCI.F.GetNamespace(baseCI.F{})+"\nin another terminal.", 0)
	}
}

func (f F) selfHostedInfra(context string) {
	if vars.InfraProviderActive == vars.InfraProviderSelfHosted {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInfo}, "To spin up your "+context+" development environment you have to connect a GPU enabled kubernetes node to your cluster.", 0)
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInfo}, "All necessary GPU drivers will be installed automatically by "+vars.OrganizationName+".", 0)
		sel := terminal.GetUserSelection("Do you have already connected a GPU node?", []string{}, false, true)
		if sel != "Yes" {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiInfo}, "Okay, so then please connect one and come back later.", 0)
			terminal.Exit(0, "")
		}
	}
}
