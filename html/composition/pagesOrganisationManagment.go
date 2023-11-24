package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func GetCreateOrganisationPage(org *validation.OrganisationModification, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getDataListFromMap("mainGroupNames", extraction.OrganisationMainGroups),
		getDataListFromMap("subGroupNames", extraction.OrganisationSubGroups),
		getPageHeader(CreateOrganisation),
		getFormStandardForm("form", POST, "/"+HTMXPreRouter+string(CreateOrganisation), CLASS("w-[800px]"),
			getSimpleTextInput("name", "name", org.Name, Translation["organisationName"]),
			getSimpleTextInput("flair", "flair", org.Flair, Translation["flair"]),
			getInput("mainGroup", "mainGroup", org.MainGroup, Translation["mainGroup"], "text", "mainGroupNames", ""),
			getInput("subGroup", "subGroup", org.SubGroup, Translation["subGroup"], "text", "subGroupNames", ""),
			getDropDown("status", "status", Translation["organisationStatus"], false,
				database.Stati, database.StatusTranslation, database.StatusString(org.Status)),
			DIV(CLASS("flex flex-row"),
				getEditableList(org.User, "user", "displayNames",
					Translation["addOrganisationUserButton"], "w-[400px]"),
				getEditableList(org.Admins, "admins", "displayNames",
					Translation["addOrganisationAdminButton"], "w-[400px] ml-2"),
			),
			getSubmitButton("createOrganisationButton", Translation["createOrganisationButton"]),
		),
		GetMessage(val),
	)
}

func GetModifyOrganisationPage(org *validation.OrganisationModification, val validation.Message) Node {
	display, err := extraction.ReturnListOfDisplayNames()
	if err != nil {
		val.Message = Translation["errorQueryingNames"] + "\n" + val.Message
	}
	return getBasePageWrapper(
		getDataList("displayNames", display),
		getDataList("organisationNames", extraction.OrganisationNamesList),
		getDataListFromMap("mainGroupNames", extraction.OrganisationMainGroups),
		getDataListFromMap("subGroupNames", extraction.OrganisationSubGroups),
		getPageHeader(EditOrganisation),
		getFormStandardForm("form", POST, "/"+HTMXPreRouter+string(EditOrganisation), CLASS("w-[800px]"),
			getInput("name", "name", org.Name, Translation["organisationName"], "text", "organisationNames", ""),
			getSubmitButtonOverwriteURL(Translation["searchOrganisationButton"], PATCH, "/"+HTMXPreRouter+string(SearchOrganisation)),
			getSimpleTextInput("flair", "flair", org.Flair, Translation["flair"]),
			getInput("mainGroup", "mainGroup", org.MainGroup, Translation["mainGroup"], "text", "mainGroupNames", ""),
			getInput("subGroup", "subGroup", org.SubGroup, Translation["subGroup"], "text", "subGroupNames", ""),
			getDropDown("status", "status", Translation["organisationStatus"], false,
				database.Stati, database.StatusTranslation, database.StatusString(org.Status)),
			DIV(CLASS("flex flex-row"),
				getEditableList(org.User, "user", "displayNames",
					Translation["addOrganisationUserButton"], "w-[400px]"),
				getEditableList(org.Admins, "admins", "displayNames",
					Translation["addOrganisationAdminButton"], "w-[400px] ml-2"),
			),
			getSubmitButton("modifyOrganisationButton", Translation["changeOrganisationButton"]),
		),
		GetMessage(val),
	)
}
