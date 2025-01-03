package organisations

import (
	"PoliSim/database"
	"PoliSim/handler"
	"net/http"
	"strings"
)

func GetOrganisationView(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	page := &handler.ViewOrganisationPage{}
	var err error
	page.Hierarchy, err = database.GetOrganisationMapForUser(acc)
	if err != nil {
		page.HadError = true
	}
	handler.MakeFullPage(writer, acc, page)
}

func GetSingleOrganisationView(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	part := &handler.SingleOrganisationUpdate{
		Name: request.URL.Query().Get("name"),
	}

	var user []string
	var admin []string
	var err error
	part.Organisation, user, admin, err = database.GetFullOrganisationInfoForUserView(acc, part.Name)
	if err == nil {
		if len(user) == 0 {
			part.User = "Diese Organisation hat keine Mitglieder"
		} else {
			part.User = "Mitglieder: " + strings.Join(user, ", ")
		}
		if len(admin) == 0 {
			part.Admin = "Diese Organisation hat keine Administratoren"
		} else {
			part.Admin = "Administratoren: " + strings.Join(admin, ", ")
		}
	}

	handler.MakeSpecialPagePart(writer, part)
}
