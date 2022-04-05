package base

import (
	infra "github.com/neurafuse/neurakube/infrastructure"
	"github.com/neurafuse/tools-go/api/client"
	"github.com/neurafuse/tools-go/apps/ide/vscode"
	"github.com/neurafuse/tools-go/apps/kubernetes/devspace"
	"github.com/neurafuse/tools-go/apps/tensorflow/tensorboard"
	"github.com/neurafuse/tools-go/ci"
	baseCI "github.com/neurafuse/tools-go/ci/base"
	"github.com/neurafuse/tools-go/config"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	projectConfig "github.com/neurafuse/tools-go/config/project"
	devConfig "github.com/neurafuse/tools-go/config/dev"
	"github.com/neurafuse/tools-go/container"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/apps/python/debugpy"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	kubeTools "github.com/neurafuse/tools-go/kubernetes/tools"

	appIDTools "github.com/neurafuse/neuracli/api/app/id"
	"github.com/neurafuse/tools-go/kubernetes/pods"
	"github.com/neurafuse/tools-go/kubernetes/portforward"
	"github.com/neurafuse/tools-go/kubernetes/services"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(context, appID, action string) { // TODO: Implement id
	infra.F.CheckDeploymentStatus(infra.F{})
	if action == "create" || action == "update" || action == "sync" {
		f.Create(context, action, appID, baseCI.F.GetResType(baseCI.F{}, context))
	} else if action == "recreate" {
		f.Recreate(context, action, appID)
	} else if action == "delete" {
		f.Delete(context, baseCI.F.GetResType(baseCI.F{}, context))
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to find action: "+action, true, true, true)
	}
}

func (f F) Create(context, action, appID, resType string) {
	logging.Log([]string{"\n", vars.EmojiRemote, vars.EmojiProcess}, "Creating "+baseCI.F.GetResType(baseCI.F{}, context)+"..", 0)
	var namespace string
	var imageAddrs string
	var remoteURL string
	var waitForStatusInLog string
	var initWaitDuration int
	f.checkProjectConfig()
	if appIDTools.F.IsKindManaged(appIDTools.F{}, appID) {
		namespace, imageAddrs, remoteURL, waitForStatusInLog, initWaitDuration = f.neuraKubeManaged(context, action)
	} else if appIDTools.F.IsKindCustom(appIDTools.F{}, appID) {
		namespace, imageAddrs = f.createSelfManaged(appID)
	}
	remoteURL = f.checkContainerNetwork(namespace, appID, action, remoteURL)
	if context == "develop" {
		f.developApp(context, appID, remoteURL, namespace, imageAddrs)
	}
	tensorboard.F.LogInfo(tensorboard.F{}, remoteURL)
	logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSuccess}, "Created "+resType+".\n", 0)
	pods.F.Logs(pods.F{}, namespace, appID, waitForStatusInLog, false, initWaitDuration)
}

func (f F) developApp(context, appID, remoteURL, namespace, imageAddrs string) {
	f.connectIDE(context, remoteURL)
	devspace.F.Sync(devspace.F{}, appID, namespace, imageAddrs, true)
}

func (f F) checkProjectConfig() {
	infra.F.CheckClusterID(infra.F{})
	f.checkContainerSync()
}

func (f F) checkContainerSync() {
	if !config.ValidSettings("project", "containers/sync", true) {
		var localAppRoot string = filesystem.GetWorkspaceFolderVar()
		if config.DevConfigActive() {
			localAppRoot = terminal.GetUserSelection("What is the local app root path?", []string{localAppRoot}, true, false)
		}
		config.Setting("set", "project", "Spec.Containers.Sync.PathMappings.LocalAppRoot", localAppRoot)
		var localIDERoot string = localAppRoot
		if config.DevConfigActive() {
			localIDERoot = terminal.GetUserSelection("What is the local IDE root path?", []string{localIDERoot}, true, false)
		}
		config.Setting("set", "project", "Spec.Containers.Sync.PathMappings.LocalIDERoot", localIDERoot)
		var containerAppRoot string
		containerAppRoot = terminal.GetUserInput("What is the container root path of the app (e.g. /app/src)?")
		config.Setting("set", "project", "Spec.Containers.Sync.PathMappings.ContainerAppRoot", containerAppRoot)
		var containerAppDataRoot string = "data"
		logging.Log([]string{"", vars.EmojiRocket, vars.EmojiWaiting}, "The container data root path is appended to the container root path which should direct to the persistent container data path.", 0)
		containerAppDataRoot = terminal.GetUserSelection("What is the container data root path of the app (e.g. (/app/src/) --> data)?", []string{containerAppDataRoot}, true, false)
		if !strings.HasSuffix(containerAppDataRoot, "/") {
			containerAppDataRoot = containerAppDataRoot + "/"
		}
		config.Setting("set", "project", "Spec.Containers.Sync.PathMappings.ContainerAppDataRoot", containerAppDataRoot)
	}
}

func (f F) connectIDE(context, remoteURL string) {
	var envIDE string = config.Setting("get", "infrastructure", "Spec."+env.F.GetContext(env.F{}, context, true)+".Environment.IDE", "")
	if envIDE == "vscode" {
		vscode.F.CreateConfig(vscode.F{}, remoteURL, baseCI.F.GetContainerPortsForApps(baseCI.F{}, []string{"debugpy"})[0][0])
	} else {
		logging.Log([]string{"\n", vars.EmojiDev, vars.EmojiWarning}, vars.NeuraCLIName+" does not support automatic configuration of remote debugging for the IDE "+envIDE+".", 0)
	}
}

func (f F) Recreate(context, action, appID string) {
	f.Delete(context, baseCI.F.GetResType(baseCI.F{}, context))
	f.Create(context, action, appID, baseCI.F.GetResType(baseCI.F{}, context))
}

func (f F) Delete(context, resType string) {
	logging.Log([]string{"\n", vars.EmojiRemote, vars.EmojiProcess}, "Deleting "+resType+"..", 0)
	var route string = context + "/delete"
	var answer string = client.F.Router(client.F{}, vars.NeuraKubeNameID, "GET", route, "", "", "", nil)
	if answer == "success" {
		logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSuccess}, "Deleted "+resType+".\n", 0)
	} else {
		var err error = errors.New("API route: " + route + "\nAnswer: " + answer)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete "+resType+"!", false, true, true)
	}
}

func (f F) createSelfManaged(appID string) (string, string) {
	var namespace string = kubeTools.F.GetDeploymentNamespace(kubeTools.F{}, appID)
	var containerID int = kubeTools.F.GetContainerID(kubeTools.F{}, namespace, appID)
	var imageAddrs string = pods.F.GetContainerImgAddrs(pods.F{}, namespace, appID, containerID)
	return namespace, imageAddrs
}

func (f F) neuraKubeManaged(context, action string) (string, string, string, string, int) {
	if action != "sync" {
		f.sendRequest(context, "prepare")
		container.F.CheckUpdates(container.F{}, context, true, false) // Also build base image with: context+"-base"
	}
	var namespace string = baseCI.F.GetNamespace(baseCI.F{})
	var imageAddrs string = container.F.GetImgAddrs(container.F{}, context, false, false)
	var remoteURL string
	if action != "sync" {
		remoteURL = f.sendRequest(context, "create")
	}
	var waitForStatusInLog string = f.getWaitForStatusInLog(context)
	var initWaitDuration int = ci.F.GetInitWaitDuration(ci.F{}, context)
	return namespace, imageAddrs, remoteURL, waitForStatusInLog, initWaitDuration
}

func (f F) checkContainerNetwork(namespace, appID, action, remoteURL string) string {
	if projectConfig.F.NetworkMode(projectConfig.F{}, "port-forward") {
		var podID string = pods.F.GetIDFromName(pods.F{}, namespace, appID)
		var port int = strings.ToInt(debugpy.GetContainerPorts()[0])
		go portforward.Connect(namespace, podID, port, port)
		remoteURL = "localhost:"+strings.ToString(port)
	} else if projectConfig.F.NetworkMode(projectConfig.F{}, "remote-url") {
		if action == "sync" {
			var validGcloudSettings bool = config.ValidSettings("infrastructure", vars.InfraProviderGcloud, true)
			var appKindManaged bool = appIDTools.F.IsKindManaged(appIDTools.F{}, appID)
			var requestLBIP bool = validGcloudSettings && appKindManaged
			if requestLBIP {
				remoteURL = services.F.GetLoadBalancerIP(services.F{}, namespace, appID)
			}
		}
		if remoteURL == "" {
			remoteURL = projectConfig.F.GetRemoteURL(projectConfig.F{}, appID)
		}
	}
	return remoteURL
}

func (f F) getWaitForStatusInLog(context string) string {
	var waitForStatusInLog string
	if context == "develop" {
		waitForStatusInLog = "Epoch: 1"
	} else if context == "inference" {
		waitForStatusInLog = "Serving"
	}
	return waitForStatusInLog
}

func (f F) sendRequest(context, signal string) string {
	var waitDuration string = strings.ToString(ci.F.GetInitWaitDuration(ci.F{}, context))
	if signal == "prepare" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWaiting}, "Requesting "+signal+" "+context+"..", 0)
	} else if signal == "create" {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWaiting}, "Sending request to "+vars.NeuraKubeName+" and waiting for the "+context+" environment to be ready (this may take up to "+waitDuration+" minutes for a new cluster)..", 0)
	}
	var response string = client.F.Router(client.F{}, vars.NeuraKubeNameID, "GET", context, "", "", "action="+signal, nil)
	var respSuccessKey string = "success/"
	if response == "success" || strings.HasPrefix(response, respSuccessKey) {
		response = strings.TrimPrefix(response, respSuccessKey)
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Request "+signal+" "+context+" successful.\n", 0)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Received error for "+signal+" "+context+": "+response, true, true, true)
	}
	return response
}

func (f F) checkUserConfigs(context string) {
	logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSettings}, "Checking "+usersID.F.GetActive(usersID.F{})+" user configs..\n", 0)
	if !config.ValidSettings("infrastructure", context, true) {
		infraConfig.F.SetModuleSpec(infraConfig.F{}, context)
	}
	accType := baseCI.F.GetResources(baseCI.F{}, context)
	if !config.ValidSettings("infrastructure", vars.InfraProviderGcloud+"/accelerator/"+accType, true) {
		infraConfig.F.SetGcloudAccelerator(infraConfig.F{}, context, accType)
	}
	devConfig.F.RequireAPIAddrsNonLocal(devConfig.F{}, context)
}