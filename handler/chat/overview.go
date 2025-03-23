package chat

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

const messageID = "message-box"

func GetChatOverview(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	query := helper.GetAdvancedURLValues(request)
	page := &handler.ChatOverviewPage{}
	viewerName := query.GetTrimmedString("viewer")
	page.Viewer = viewerName

	errorMsg := make([]string, 0)
	var viewer []string
	var err error
	var allowed bool

	page.PossibleViewers, err = database.GetMyAccountNames(acc)
	if err != nil {
		page.IsError = true
		errorMsg = append(errorMsg, loc.CouldNotFindAllOwnedAccounts)
	}

	if page.Viewer != "" {
		allowed, err = database.IsAccountAllowedToPostWith(acc, page.Viewer)
		if !allowed || err != nil {
			viewer = []string{acc.Name}
			page.Viewer = acc.Name
		} else {
			viewer = []string{page.Viewer}
		}
	} else {
		if err != nil {
			viewer = []string{acc.Name}
			page.Viewer = acc.Name
		} else {
			viewer = page.PossibleViewers
		}
	}

	page.Chats, err = database.GetAllRoomsForUser(viewer)
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		errorMsg = append(errorMsg, loc.ChatFailedToLoadChats)
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.IsError = true
		errorMsg = append(errorMsg, loc.ErrorLoadingAccountNames)
	}

	if page.IsError {
		page.Message = strings.Join(errorMsg, "\n")
	}

	if viewerName != page.Viewer {
		writer.Header().Add("Hx-Push-Url", "/chat/overview?viewer="+template.URLQueryEscaper(page.Viewer))
	}
	page.ElementID = messageID
	handler.MakeFullPage(writer, acc, page)
}

func PostNewChat(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions, ElementID: messageID})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError, ElementID: messageID})
		return
	}

	member := values.GetTrimmedArray("[]member")
	baseAccount := values.GetTrimmedString("base-account")
	roomName := values.GetTrimmedString("chatroom-name")

	const maxRoomLength = 100
	if len(roomName) > maxRoomLength {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ChatRoomNameTooLong, maxRoomLength), ElementID: messageID})
		return
	} else if roomName == "" {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomNameIsEmpty, ElementID: messageID})
		return
	}

	allowed, err := database.IsAccountAllowedToPostWith(acc, baseAccount)
	if err != nil {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatCouldNotVerifyAccountCredentials, ElementID: messageID})
		return
	}
	if !allowed {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatNotAllowedToCreateWithThatAccount, ElementID: messageID})
		return
	}

	member = append(member, baseAccount)

	err = database.CreateChatRoom(roomName, member)
	if errors.Is(err, database.ChatRoomNameTaken) {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomNameAlreadyTaken, ElementID: messageID})
		return
	} else if errors.Is(err, database.DoubleChatRoomEntry) {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomWithMemberConstellationAlreadyExists, ElementID: messageID})
		return
	} else if err != nil {
		handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomCreationError, ElementID: messageID})
		slog.Debug(err.Error())
		return
	}

	handler.MakeSpecialPagePart(writer, &handler.MessageUpdate{IsError: false,
		Message: loc.ChatRoomSuccessfullyCreated, ElementID: messageID})
}

func PutFilterChatOverview(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions, ElementID: messageID})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError, ElementID: messageID})
		return
	}

	page := &handler.ChatOverviewPage{}
	page.Viewer = values.GetTrimmedString("viewer")

	errorMsg := make([]string, 0)
	var viewer []string
	var allowed bool

	page.PossibleViewers, err = database.GetMyAccountNames(acc)
	if err != nil {
		page.IsError = true
		errorMsg = append(errorMsg, loc.CouldNotFindAllOwnedAccounts)
	}

	if page.Viewer != "" {
		allowed, err = database.IsAccountAllowedToPostWith(acc, page.Viewer)
		if !allowed || err != nil {
			viewer = []string{acc.Name}
			page.Viewer = acc.Name
		} else {
			viewer = []string{page.Viewer}
		}
	} else {
		if err != nil {
			viewer = []string{acc.Name}
			page.Viewer = acc.Name
		} else {
			viewer = page.PossibleViewers
		}
	}

	page.Chats, err = database.GetAllRoomsForUser(viewer)
	if err != nil {
		page.IsError = true
		errorMsg = append(errorMsg, loc.ChatFailedToLoadChats)
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.IsError = true
		errorMsg = append(errorMsg, loc.ErrorLoadingAccountNames)
	}

	if page.IsError {
		page.Message = strings.Join(errorMsg, "\n")
	}

	if page.Viewer != "" {
		writer.Header().Add("Hx-Push-Url", "/chat/overview?viewer="+template.URLQueryEscaper(page.Viewer))
	} else {
		writer.Header().Add("Hx-Push-Url", "/chat/overview")
	}
	page.ElementID = messageID
	handler.MakePage(writer, acc, page)
}
