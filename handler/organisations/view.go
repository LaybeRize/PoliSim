package organisations

import (
	"PoliSim/database"
	"PoliSim/handler"
	loc "PoliSim/localisation"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func GetOrganisationView(writer http.ResponseWriter, request *http.Request) {
	acc, _ := database.RefreshSession(writer, request)
	page := &handler.ViewOrganisationPage{}
	var err error
	page.Hierarchy, err = database.GetOrganisationMapForUser(acc)
	if err != nil {
		slog.Debug(err.Error())
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
	if acc.Exists() && acc.Role <= database.Admin {
		part.Organisation, user, admin, err = database.GetFullOrganisationInfo(part.Name)
	} else {
		part.Organisation, user, admin, err = database.GetFullOrganisationInfoForUserView(acc, part.Name)
	}
	if err == nil {
		if len(user) == 0 {
			part.User = loc.OrganisationHasNoMember
		} else {
			part.User = fmt.Sprintf(loc.OrganisationMemberList, strings.Join(user, ", "))
		}
		if len(admin) == 0 {
			part.Admin = loc.OrganisationHasNoAdministrator
		} else {
			part.Admin = fmt.Sprintf(loc.OrganisationAdministratorList, strings.Join(admin, ", "))
		}
	} else {
		part.Organisation = nil
	}

	handler.MakeSpecialPagePart(writer, part)
}
