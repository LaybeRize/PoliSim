package htmlComposition

import "net/http"

type HttpUrl string

const (
	Start  HttpUrl = "start"
	Login  HttpUrl = "login"
	Logout HttpUrl = "logout"

	APIPreRoute = "htmx/"
)

type HttpHandling struct {
	TitleText          string
	SidebarButtonText  string
	HasSidebarButton   bool
	SidebarSubMenuText string
	HasSidebarSubMenu  bool
}

var LoadingList = []HttpUrl{Start}
var HandlerList = make(map[HttpUrl]*HttpHandling)
var GetHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var PostHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var PatchHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var DeleteHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
