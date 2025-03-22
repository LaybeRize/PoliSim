package chat

import (
	"PoliSim/database"
	"PoliSim/handler"
	"html/template"
	"log/slog"
	"net/http"
	"sync"
)

var HubList map[string]*Hub
var HubMutex sync.Mutex

func ConnectToWebsocket(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		http.Error(writer, "not allowed", http.StatusNotAcceptable)
		return
	}

	user := request.PathValue("user")
	allowed, err := database.IsAccountAllowedToPostWith(acc, user)
	if err != nil || !allowed {
		http.Error(writer, "not allowed", http.StatusForbidden)
		return
	}

	id := request.PathValue("id")

	HubMutex.Lock()
	err = database.QueryForRoomIdAndUser(id, user, acc.Name)
	if err != nil {
		HubMutex.Unlock()
		slog.Debug(err.Error())
		http.Error(writer, "not allowed", http.StatusForbidden)
		return
	}

	if _, ok := HubList[id]; !ok {
		hub := NewHub(id)
		go hub.run()
		HubList[id] = hub
		serveWs(hub, acc, user, writer, request)
	} else {
		serveWs(HubList[id], acc, user, writer, request)
	}
	HubMutex.Unlock()
}

func GetShowChat(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}
	user := request.PathValue("user")
	allowed, err := database.IsAccountAllowedToPostWith(acc, user)
	if err != nil || !allowed {
		handler.GetNotFoundPage(writer, request)
		return
	}

	id := request.PathValue("id")
	err = database.QueryForRoomIdAndUser(id, user, acc.Name)
	if err != nil {
		handler.GetNotFoundPage(writer, request)
		slog.Debug(err.Error())
		return
	}

	handler.MakeFullPage(writer, acc, &handler.DirectChatWindow{
		ConnectURL: template.URL("/connect/chat/" + template.URLQueryEscaper(id) + "/" + template.URLQueryEscaper(user)),
	})
}
