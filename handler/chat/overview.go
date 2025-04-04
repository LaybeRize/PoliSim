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

func GetChatOverview(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	query := helper.GetAdvancedURLValues(request)
	page := &handler.ChatOverviewPage{
		Query: &database.ChatSearch{
			Name:                query.GetTrimmedString("chat-name"),
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
		page.Chats, err = database.GetRoomsPageBackwards(page.Amount, page.PreviousItemTime, recName, page.Query)
	} else {
		page.Chats, err = database.GetRoomsPageForwards(page.Amount, page.NextItemTime, recName, page.Query)
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

	handler.MakeFullPage(writer, acc, page)
}

func PutFilterChatOverview(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	page := &handler.ChatOverviewPage{
		Query: &database.ChatSearch{
			Name:                values.GetTrimmedString("chat-name"),
			Viewer:              values.GetTrimmedString("viewer"),
			ShowOnlyUnreadChats: values.GetBool("only-new"),
			Owner:               acc.Name,
		},
		Amount: values.GetInt("amount"),
	}

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	var backward bool
	page.PreviousItemTime, backward = values.GetUTCTime("backward", false)
	page.NextItemTime, _ = values.GetUTCTime("forward", true)
	recName := values.GetTrimmedString("rec-name")

	if backward {
		page.Chats, err = database.GetRoomsPageBackwards(page.Amount, page.PreviousItemTime, recName, page.Query)
	} else {
		page.Chats, err = database.GetRoomsPageForwards(page.Amount, page.NextItemTime, recName, page.Query)
	}
	if err != nil {
		slog.Debug(err.Error())
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatFailedToLoadChats})
		return
	}

	page.PossibleViewers, err = database.GetMyAccountNames(acc)

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
	handler.MakePage(writer, acc, page)
}

func PostNewChat(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.MissingPermissions})
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	member := values.GetTrimmedArray("[]member")
	baseAccount := values.GetTrimmedString("base-account")
	roomName := values.GetTrimmedString("chatroom-name")

	const maxRoomLength = 100
	if len(roomName) > maxRoomLength {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.ChatRoomNameTooLong, maxRoomLength)})
		return
	} else if roomName == "" {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomNameIsEmpty})
		return
	}

	allowed, err := database.IsAccountAllowedToPostWith(acc, baseAccount)
	if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatCouldNotVerifyAccountCredentials})
		return
	}
	if !allowed {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatNotAllowedToCreateWithThatAccount})
		return
	}

	member = append(member, baseAccount)

	err = database.CreateChatRoom(roomName, member)
	if errors.Is(err, database.ChatRoomNameTaken) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomNameAlreadyTaken})
		return
	} else if errors.Is(err, database.DoubleChatRoomEntry) {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomWithMemberConstellationAlreadyExists})
		return
	} else if err != nil {
		handler.SendMessageUpdate(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.ChatRoomCreationError})
		slog.Debug(err.Error())
		return
	}

	page := &handler.ChatOverviewPage{
		Query: &database.ChatSearch{
			Name:                values.GetTrimmedString("chat-name"),
			Viewer:              values.GetTrimmedString("viewer"),
			ShowOnlyUnreadChats: values.GetBool("only-new"),
			Owner:               acc.Name,
		},
		Amount: values.GetInt("amount"),
		MessageUpdate: handler.MessageUpdate{
			Message: loc.ChatRoomSuccessfullyCreated,
			IsError: false,
		},
	}

	values.DeleteFields([]string{"[]member", "base-account", "chatroom-name"})

	if page.Amount < 10 || page.Amount > 50 {
		page.Amount = 20
	}

	page.Chats, err = database.GetRoomsPageForwards(page.Amount, time.Now().UTC(), "", page.Query)

	if len(page.Chats) > page.Amount {
		page.HasNext = true
		page.NextItemTime = page.Chats[page.Amount].Created
		page.NextItemRec = page.Chats[page.Amount].User
		page.Chats = page.Chats[:page.Amount]
	}

	page.PossibleViewers, err = database.GetMyAccountNames(acc)

	writer.Header().Add("Hx-Push-Url", "/chat/overview?"+values.Encode())
	handler.MakePage(writer, acc, page)
}
