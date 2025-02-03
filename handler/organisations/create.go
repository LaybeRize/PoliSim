package organisations

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/helper"
	loc "PoliSim/localisation"
	"fmt"
	"net/http"
	"strings"
)

const (
	maxNameLength     = 600
	maxMainTypeLength = 400
	maxSubTypeLength
	maxFlairLength = 200
)

func GetCreateOrganisationPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.GetNotFoundPage(writer, request)
		return
	}

	page := &handler.CreateOrganisationPage{
		Admin: []string{""},
		User:  []string{""},
	}
	var err error
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.IsError = true
		page.Message = loc.ErrorLoadingAccountNames
	}

	handler.MakeFullPage(writer, acc, page)
}

func PostCreateOrganisationPage(writer http.ResponseWriter, request *http.Request) {
	acc, loggedIn := database.RefreshSession(writer, request)
	if !loggedIn || acc.Role > database.HeadAdmin {
		handler.PartialGetNotFoundPage(writer, request)
		return
	}

	values, err := helper.GetAdvancedFormValues(request)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.RequestParseError})
		return
	}

	newOrganisation := &database.Organisation{
		Name:       values.GetTrimmedString("name"),
		Visibility: database.OrganisationVisibility(values.GetInt("visiblity")),
		MainType:   values.GetTrimmedString("main-group"),
		SubType:    values.GetTrimmedString("sub-group"),
		Flair:      values.GetTrimmedString("flair"),
	}

	userNames := values.GetTrimmedArray("[]user")
	adminNames := values.GetTrimmedArray("[]admin")

	if newOrganisation.Name == "" || newOrganisation.MainType == "" || newOrganisation.SubType == "" {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.OrganisationGeneralInformationEmpty})
		return
	}

	if len([]rune(newOrganisation.Name)) > maxNameLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.OrganisationGeneralNameTooLong, maxNameLength)})
		return
	}

	if len([]rune(newOrganisation.MainType)) > maxMainTypeLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.OrganisationGeneralMainGroupTooLong, maxMainTypeLength)})
		return
	}

	if len([]rune(newOrganisation.SubType)) > maxSubTypeLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.OrganisationGeneralSubGroupTooLong, maxSubTypeLength)})
		return
	}

	if strings.Contains(newOrganisation.Flair, ",") ||
		strings.Contains(newOrganisation.Flair, ";") ||
		len([]rune(newOrganisation.Flair)) > maxFlairLength {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: fmt.Sprintf(loc.OrganisationGeneralFlairContainsInvalidCharactersOrIsTooLong, maxFlairLength)})
		return
	}

	if !newOrganisation.VisibilityIsValid() {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.OrganisationGeneralInvalidVisibility})
		return
	}

	newOrganisation.ClearInvalidFlair()

	err = database.CreateOrganisation(newOrganisation, userNames, adminNames)
	if err != nil {
		handler.MakeSpecialPagePartWithRedirect(writer, &handler.MessageUpdate{IsError: true,
			Message: loc.OrganisationErrorWhileCreating})
		return
	}

	page := &handler.CreateOrganisationPage{
		Admin: []string{""},
		User:  []string{""},
	}
	page.IsError = false
	page.Message = loc.OrganisationSuccessfullyCreated
	page.AccountNames, err = database.GetNonBlockedNames()
	if err != nil {
		page.Message = "\n" + loc.ErrorLoadingAccountNames
	}
	handler.MakeFullPage(writer, acc, page)
}
