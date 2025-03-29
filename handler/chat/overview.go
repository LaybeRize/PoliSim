package chat

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const messageID = "message-box"

func GetChatOverview(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	query := helper.GetAdvancedURLValues(request)
	page := &handler.ChatOverviewPage{
		Query: &database.ChatSearch{
			Viewer:              query.GetTrimmedString("viewer"),
			ShowOnlyUnreadChats: query.GetBool("only-new"),
			Owner:               acc.Name,
		},
		Amount: query.GetInt("amount"),
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	errorMsg := make([]string, 0)
	var err error

	page.PossibleViewers, err = database.GetMyAccountNames(acc)
	if err != nil {
		page.IsError = true
		errorMsg = append(errorMsg, loc.CouldNotFindAllOwnedAccounts)
	}

	var backward bool
	page.PreviousItemTime, backward = query.GetUTCTime("backward", false)
	page.NextItemTime, _ = query.GetUTCTime("forward", true)
	recName := query.GetTrimmedString("rec-name")

	if backward {
		page.Chats, err = database.GetRoomsPageBackwards(page.Amount, acc, page.PreviousItemTime, recName, page.Query)
	} else {
		page.Chats, err = database.GetRoomsPageForwards(page.Amount, acc, page.NextItemTime, recName, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		errorMsg = append(errorMsg, loc.ChatFailedToLoadChats)
	}

	if len(page.Chats) > 0 {
		if !backward && page.Chats[0].Created.Equal(page.NextItemTime) && page.Chats[0].User == recName {
			page.HasPrevious = true
			page.PreviousItemTime = page.NextItemTime
			page.PreviousItemRec = page.Chats[0].User
		} else if lst := len(page.Chats) - 1; backward && page.Chats[lst].Created.Equal(page.PreviousItemTime) && page.Chats[lst].User == recName {
			page.HasNext = true
			page.NextItemTime = page.PreviousItemTime
			page.NextItemRec = page.Chats[lst].User
			page.Chats = page.Chats[:lst]
		}
	}

	if !backward && len(page.Chats) > page.Amount {
		page.HasNext = true
		page.NextItemTime = page.Chats[page.Amount].Created
		page.NextItemRec = page.Chats[page.Amount].User
		page.Chats = page.Chats[:page.Amount]
	} else if backward && len(page.Chats) > page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Chats[1].Created
		page.PreviousItemRec = page.Chats[1].User
		page.Chats = page.Chats[1:]
	} else if backward && len(page.Chats) > page.Amount {
		amt := len(page.Chats) - page.Amount
		page.HasPrevious = true
		page.PreviousItemTime = page.Chats[amt].Created
		page.PreviousItemRec = page.Chats[amt].User
		page.Chats = page.Chats[amt:]
	}

	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.IsError = true
		errorMsg = append(errorMsg, loc.ErrorLoadingAccountNames)
	}

	if page.IsError {
		page.Message = strings.Join(errorMsg, "\n")
	}

	page.ElementID = messageID
	handler.MakeFullPage(writer, acc, page)
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

	page := &handler.ChatOverviewPage{
		Query: &database.ChatSearch{
			Viewer:              values.GetTrimmedString("viewer"),
			ShowOnlyUnreadChats: values.GetBool("only-new"),
			Owner:               acc.Name,
		},
		Amount:        values.GetInt("amount"),
		MessageUpdate: handler.MessageUpdate{ElementID: messageID},
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	var backward bool
	page.PreviousItemTime, backward = values.GetUTCTime("backward", false)
	page.NextItemTime, _ = values.GetUTCTime("forward", true)
	recName := values.GetTrimmedString("rec-name")

	if backward {
		page.Chats, err = database.GetRoomsPageBackwards(page.Amount, acc, page.PreviousItemTime, recName, page.Query)
	} else {
		page.Chats, err = database.GetRoomsPageForwards(page.Amount, acc, page.NextItemTime, recName, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		page.IsError = true
		page.Message = loc.ChatFailedToLoadChats
	}

	if len(page.Chats) > 0 {
		if !backward && page.Chats[0].Created.Equal(page.NextItemTime) && page.Chats[0].User == recName {
			page.HasPrevious = true
			page.PreviousItemTime = page.NextItemTime
			page.PreviousItemRec = page.Chats[0].User
		} else if lst := len(page.Chats) - 1; backward && page.Chats[lst].Created.Equal(page.PreviousItemTime) && page.Chats[lst].User == recName {
			page.HasNext = true
			page.NextItemTime = page.PreviousItemTime
			page.NextItemRec = page.Chats[lst].User
			page.Chats = page.Chats[:lst]
		}
	}

	if !backward && len(page.Chats) > page.Amount {
		page.HasNext = true
		page.NextItemTime = page.Chats[page.Amount].Created
		page.NextItemRec = page.Chats[page.Amount].User
		page.Chats = page.Chats[:page.Amount]
	} else if backward && len(page.Chats) > page.Amount && page.HasNext {
		page.HasPrevious = true
		page.PreviousItemTime = page.Chats[1].Created
		page.PreviousItemRec = page.Chats[1].User
		page.Chats = page.Chats[1:]
	} else if backward && len(page.Chats) > page.Amount {
		amt := len(page.Chats) - page.Amount
		page.HasPrevious = true
		page.PreviousItemTime = page.Chats[amt].Created
		page.PreviousItemRec = page.Chats[amt].User
		page.Chats = page.Chats[amt:]
	}

	writer.Header().Add("Hx-Push-Url", "/chat/overview?"+values.Encode())
	handler.MakeSpecialPagePart(writer, page)
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

	page := &handler.ChatOverviewPage{
		Query: &database.ChatSearch{
			Viewer:              values.GetTrimmedString("viewer"),
			ShowOnlyUnreadChats: values.GetBool("only-new"),
			Owner:               acc.Name,
		},
		Amount: values.GetInt("amount"),
		MessageUpdate: handler.MessageUpdate{
			ElementID: messageID,
			Message:   loc.ChatRoomSuccessfullyCreated,
			IsError:   false,
		},
	}

	values.DeleteFields([]string{"[]member", "base-account", "chatroom-name"})

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.Chats, err = database.GetRoomsPageForwards(page.Amount, acc, time.Now().UTC(), "", page.Query)

	if len(page.Chats) > page.Amount {
		page.HasNext = true
		page.NextItemTime = page.Chats[page.Amount].Created
		page.NextItemRec = page.Chats[page.Amount].User
		page.Chats = page.Chats[:page.Amount]
	}

	writer.Header().Add("Hx-Push-Url", "/chat/overview?"+values.Encode())
	handler.MakeSpecialPagePart(writer, page)
}
