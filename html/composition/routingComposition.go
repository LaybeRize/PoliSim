package composition

import (
	"PoliSim/data/database"
	"PoliSim/html/builder"
	"net/http"
)

const (
	MainBodyID = "mainBody"
	SidebarID  = "Sidebar"
	MessageID  = "messageDiv"

	Start              builder.HttpUrl = "start"
	Login              builder.HttpUrl = "login"
	Logout             builder.HttpUrl = "logout"
	CreateVote         builder.HttpUrl = "vote/create"
	RequestVotePartial builder.HttpUrl = "vote/request-partial"
	CreateUser         builder.HttpUrl = "account/create"
	EditUser           builder.HttpUrl = "account/edit"
	SearchUser         builder.HttpUrl = "account/search"
	ViewUser           builder.HttpUrl = "account/view"
	EditTitle          builder.HttpUrl = "title/edit"
	SearchTitle        builder.HttpUrl = "title/search"
	DeleteTitle        builder.HttpUrl = "title/delete"
	CreateTitle        builder.HttpUrl = "title/create"
	ViewTitles         builder.HttpUrl = "title/view"
	getTitleSubGroup   builder.HttpUrl = "title/get-sub-group/"
	GetTitleSubGroup   builder.HttpUrl = "title/get-sub-group/{mainGroup}/{subGroup}"
	CreateOrganisation builder.HttpUrl = "organisation/create"
	EditOrganisation   builder.HttpUrl = "organisation/edit"
	SearchOrganisation builder.HttpUrl = "organisation/search"
	ErrorPage          builder.HttpUrl = "errorPage"

	// NotFound is only used as a way to keep the PageTitleMap in order
	NotFound builder.HttpUrl = "notFound"

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

var PageTitleMap = make(map[builder.HttpUrl]string)
var SidebarTitleMap = make(map[builder.HttpUrl]string)
var GetHTMXFunctions = make(map[builder.HttpUrl]http.HandlerFunc)
var PostHTMXFunctions = make(map[builder.HttpUrl]http.HandlerFunc)
var PatchHTMXFunctions = make(map[builder.HttpUrl]http.HandlerFunc)
