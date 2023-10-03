package htmlComposition

import (
	"PoliSim/database"
	"net/http"
)

type HttpUrl string

const (
	MainBodyID    = "mainBody"
	SidebarID     = "Sidebar"
	InformationID = "informationDiv"

	Start  HttpUrl = "start"
	Login  HttpUrl = "login"
	Logout HttpUrl = "logout"
	Test   HttpUrl = "test"

	APIPreRoute = "htmx/"
)

type HttpHandling struct {
	TitleText          string
	SidebarButtonText  string
	HasSidebarButton   bool
	SidebarSubMenuText string
	HasSidebarSubMenu  bool
	RoleLevel          database.RoleLevel
}

var LoadingList = []HttpUrl{Start}
var PageTitleMap = make(map[HttpUrl]string)
var SidebarTitleMap = make(map[HttpUrl]string)
var GetHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var PostHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var PatchHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var DeleteHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
