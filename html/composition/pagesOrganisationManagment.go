package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func GetCreateOrganisationPage(org *validation.OrganisationModification, val validation.Message) Node {
	return getBasePageWrapper(
		getPageHeader(CreateOrganisation),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(CreateOrganisation), CLASS("w-[800px]"),
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
			getSubmitButton(Translation["createOrganisationButton"]),
		),
		GetMessage(val),
	)
}

func GetModifyOrganisationPage(org *validation.OrganisationModification, val validation.Message) Node {
	return getBasePageWrapper(
		getPageHeader(EditOrganisation),
		getFormStandardForm("form", POST, "/"+APIPreRoute+string(EditOrganisation), CLASS("w-[800px]"),
			getSimpleTextInput("name", "name", org.Name, Translation["organisationName"]),
			getSubmitButtonOverwriteURL(Translation["searchOrganisationButton"], PATCH, "/"+APIPreRoute+string(SearchOrganisation)),
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
			getSubmitButton(Translation["createOrganisationButton"]),
		),
		GetMessage(val),
	)
}
