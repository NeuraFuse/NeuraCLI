package cloud

import (
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	logging.Log([]string{"\n", vars.EmojiAPI, ""}, vars.OrganizationName+" Cloud will be available soon.", 0)
	logging.Log([]string{"", vars.EmojiRocket, ""}, "Stay tuned.\n", 0)
}
