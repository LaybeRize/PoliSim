package newspaper

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	"net/http"
	"strings"
)

func GetCreateArticlePage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.CreateArticlePage{}
	page.IsError = true
	page.Message = ""

	arr, err := database.GetOwnedAccountNames(acc)
	if err != nil {
		page.Message = "Konnte nicht alle möglichen Autoren finden"
		arr = make([]string, 0)
	}
	arr = append([]string{acc.Name}, arr...)
	page.Author = acc.Name
	page.PossibleAuthors = arr
	page.PossibleNewspaper, err = database.GetNewspaperNameListForAccount(acc.Name)
	if err != nil {
		page.Message = "\n" + "Konnte nicht alle möglichen Zeitungen für ausgewählten Account finden"
		page.Message = strings.TrimSpace(page.Message)
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateArticlePage(writer http.ResponseWriter, request *http.Request) {

}

func GetFindNewspaperForAccountPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehlende Berechtigung"})
		return
	}

	err := request.ParseForm()
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehler beim Parsen der Informationen"})
		return
	}
	baseAcc, owner, err := database.GetAccountAndOwnerByAccountName(helper.GetFormEntry(request, "author"))
	if !(baseAcc.Name == acc.Name || owner.Name == acc.Name) || err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Fehlende Berechtigung um die Informationen für diesen Account anzufordern"})
		return
	}

	page := &handler.CreateArticlePage{}
	page.PossibleNewspaper, err = database.GetNewspaperNameListForAccount(baseAcc.Name)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: "Konnte nicht alle möglichen Zeitungen für ausgewählten Account finden"})
		return
	}
	handler.MakeSpecialPagePart(writer, page)
}
