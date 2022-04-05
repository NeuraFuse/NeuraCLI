package id

import (
	"sync"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/runtime"
)

type F struct{}

var ids map[string]string = make(map[string]string)
var idsSync = sync.RWMutex{}

func (f F) GetKinds() []string {
	var kinds []string
	kinds = []string{"managed", "custom"}
	return kinds
}

func (f F) IsKindManaged(appID string) bool {
	var managed bool
	var kindRecent string = f.GetKind(appID)
	if kindRecent == f.GetKinds()[0] {
		managed = true
	}
	return managed
}

func (f F) IsKindCustom(appID string) bool {
	return !f.IsKindManaged(appID)
}

func (f F) GetKind(appID string) string {
	var kind string
	idsSync.Lock()
	kind = ids[appID]
	idsSync.Unlock()
	if kind == "" {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable get kind for appID: "+appID+"!", true, true, true)
	}
	return kind
}

func (f F) SetKind(appID, kind string) {
	idsSync.Lock()
	ids[appID] = kind
	idsSync.Unlock()
}
