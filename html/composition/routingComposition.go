package composition

import (
	"PoliSim/data/database"
	"net/http"
)

type HttpUrl string

const (
	MainBodyID    = "mainBody"
	SidebarID     = "Sidebar"
	InformationID = "informationDiv"
	MessageID     = "messageDiv"

	Start              HttpUrl = "start"
	Login              HttpUrl = "login"
	Logout             HttpUrl = "logout"
	CreateVote         HttpUrl = "vote/create"
	RequestVotePartial HttpUrl = "vote/request-partial"
	CreateUser         HttpUrl = "account/create"
	EditUser           HttpUrl = "account/edit"
	SearchUser         HttpUrl = "account/search"
	ViewUser           HttpUrl = "account/view"
	EditTitle          HttpUrl = "title/edit"
	SearchTitle        HttpUrl = "title/search"
	DeleteTitle        HttpUrl = "title/delete"
	CreateTitle        HttpUrl = "title/create"
	ViewTitles         HttpUrl = "title/view"
	getTitleSubGroup   HttpUrl = "title/get-sub-group/"
	GetTitleSubGroup   HttpUrl = "title/get-sub-group/{mainGroup}/{subGroup}"
	ErrorPage          HttpUrl = "errorPage"

	// NotFound is only used as a way to keep the PageTitleMap in order
	NotFound HttpUrl = "notFound"

	// APIPreRoute is a subroute for the web application to prepend to any
	// backend partial replies. It never starts with a / because that is automatically prepend anyway
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

var PageTitleMap = make(map[HttpUrl]string)
var SidebarTitleMap = make(map[HttpUrl]string)
var GetHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var PostHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var PatchHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
