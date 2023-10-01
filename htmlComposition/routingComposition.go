package htmlComposition

import "net/http"

type HttpUrl string

const (
	Start HttpUrl = "start"
)

type HttpHandling struct {
	TitleText          string
	SidebarButtonText  *string
	SidebarSubMenuText *string
	GetFunction        func(http.ResponseWriter, *http.Request)
	PostFunction       func(http.ResponseWriter, *http.Request)
	PatchFunction      func(http.ResponseWriter, *http.Request)
	DeleteFunction     func(http.ResponseWriter, *http.Request)
}

var LoadingList = []HttpUrl{Start}
var HandlerList = make(map[HttpUrl]*HttpHandling)
