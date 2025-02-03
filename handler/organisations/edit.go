package organisations

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func GetEditOrganisationPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.EditOrganisationPage{Organisation: nil}
	var err error

	query := helper.GetAdvancedURLValues(request)
	if query.Has("name") {
		page.Organisation, page.User, page.Admin, err = database.GetFullOrganisationInfo(query.GetTrimmedString("name"))

		page.IsError = true
		if err != nil {
			handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
				Message: loc.OrganisationNoOrganisationWithThatName})
			return
		}

		page.User = append(page.User, "")
		page.Admin = append(page.Admin, "")
		page.IsError = false
		page.Message = loc.OrganisationFoundOrganisation

		page.AccountNames, err = database.GetNonBlockedNames()
		if err != nil {
			page.Message += "\n" + loc.ErrorSearchingForAccountNames
		}

		handler.MakeFullPage(writer, acc, page)
		return
	}

	page.Organisations, err = database.GetOrganisationNameList()
	if err != nil {
		page.IsError = true
		page.Message = loc.OrganisationErrorSearchingForOrganisationList
	}

	handler.MakeFullPage(writer, acc, page)
}

func PatchEditOrganisationPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	oldOrganisationName := values.GetTrimmedString("oldName")
	organisationUpdate := &database.Organisation{
		Name:       values.GetTrimmedString("name"),
		Visibility: database.OrganisationVisibility(values.GetInt("visiblity")),
		MainType:   values.GetTrimmedString("main-group"),
		SubType:    values.GetTrimmedString("sub-group"),
		Flair:      values.GetTrimmedString("flair"),
	}

	userNames := values.GetTrimmedArray("[]user")
	adminNames := values.GetTrimmedArray("[]admin")

	if organisationUpdate.Name == "" || organisationUpdate.MainType == "" || organisationUpdate.SubType == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.OrganisationGeneralInformationEmpty})
		return
	}

	if len([]rune(organisationUpdate.Name)) > maxNameLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.OrganisationGeneralNameTooLong, maxNameLength)})
		return
	}

	if len([]rune(organisationUpdate.MainType)) > maxMainTypeLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.OrganisationGeneralMainGroupTooLong, maxMainTypeLength)})
		return
	}

	if len([]rune(organisationUpdate.SubType)) > maxSubTypeLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.OrganisationGeneralSubGroupTooLong, maxSubTypeLength)})
		return
	}

	if strings.Contains(organisationUpdate.Flair, ",") ||
		strings.Contains(organisationUpdate.Flair, ";") ||
		len([]rune(organisationUpdate.Flair)) > maxFlairLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.OrganisationGeneralFlairContainsInvalidCharactersOrIsTooLong, maxFlairLength)})
		return
	}

	if !organisationUpdate.VisibilityIsValid() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.OrganisationGeneralInvalidVisibility})
		return
	}

	organisationUpdate.ClearInvalidFlair()

	err = database.UpdateOrganisation(oldOrganisationName, organisationUpdate)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.OrganisationErrorUpdatingOrganisation})
		return
	}

	err = database.AddOrganisationMember(organisationUpdate, userNames, adminNames)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.OrganisationErrorUpdatingOrganisationMember})
		return
	}

	page := &handler.EditOrganisationPage{Organisation: organisationUpdate, User: userNames, Admin: adminNames}
	if _, actualUser, actualAdmins, err := database.GetFullOrganisationInfo(organisationUpdate.Name); err == nil {
		page.User = append(actualUser, "")
		page.Admin = append(actualAdmins, "")
	}
	page.IsError = false
	page.Message = loc.OrganisationSuccessfullyUpdated
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message += "\n" + loc.ErrorSearchingForAccountNames
	}
	handler.MakePage(writer, acc, page)
}

func PutOrganisationSearchPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.Admin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	name := values.GetTrimmedString("name")
	_, err = database.GetOrganisationByName(name)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.OrganisationNotFoundByName})
		return
	}

	writer.Header().Add("HX-Redirect", "/edit/organisation?name="+url.QueryEscape(name))
	writer.WriteHeader(http.StatusFound)
}
