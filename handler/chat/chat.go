package chat

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var HubList map[string]*Hub
var HubMutex sync.Mutex

func init() {
	HubList = make(map[string]*Hub)
}

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
		slog.Debug("Missing Permission", "name", user)
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
		ConnectURL: template.URL("/connect/chat/" + url.PathEscape(id) + "/" + url.PathEscape(user)),
	})
}

func GetOlderMessages(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	user := request.PathValue("user")
	allowed, err := database.IsAccountAllowedToPostWith(acc, user)
	if err != nil || !allowed {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	id := request.PathValue("id")
	err = database.QueryForRoomIdAndUser(id, user, acc.Name)
	if err != nil {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		slog.Debug(err.Error())
		return
	}
	timeStamp, err := time.ParseInLocation(helper.ISOTimeFormat, request.PathValue("time"), time.UTC)
	if err != nil {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomTimeWasInvalid})
		return
	}

	pagePart := &handler.ChatLoadNextMessages{
		CanLoadMoreMessages: false,
		Account:             acc,
		Recipient:           user,
		Button: handler.ChatButtonObject{
			Room:      id,
			Recipient: user,
		},
	}

	pagePart.Messages, err = database.LoadLastMessages(loadMessageAmount+1, timeStamp, id, user)
	if err != nil {
		pagePart.Messages = make([]database.Message, 0)
	}

	if len(pagePart.Messages) > loadMessageAmount {
		pagePart.Messages = pagePart.Messages[:loadMessageAmount]
		pagePart.CanLoadMoreMessages = true
		pagePart.Button.NextTime = pagePart.Messages[loadMessageAmount-1].SendDate
	}

	handler.MakeSpecialPagePart(writer, pagePart)
}
