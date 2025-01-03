package titles

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
	"strings"
)

func GetTitleView(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	handler.MakeFullPage(writer, acc, &handler.ViewTitlePage{TitleHierarchy: database.TitleMap})
}

func GetSingleViewTitle(writer http.ResponseWriter, request *http.Request) {
	part := &handler.SingleTitelUpdate{
		Found: false,
		Title: request.URL.Query().Get("name"),
	}

	title, holder, err := database.GetTitleAndHolder(part.Title)
	if err == nil {
		part.Found = true
		part.Flair = title.Flair
		if len(holder) == 0 {
			part.Holder = "Dieser Titel wird von niemandem gehalten"
		} else {
			part.Holder = "Titel-Halter: " + strings.Join(holder, ", ")
		}
	}

	handler.MakeSpecialPagePart(writer, part)
}
