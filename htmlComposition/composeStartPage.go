package htmlComposition

import (
	. "PoliSim/componentHelper"
	"PoliSim/dataExtraction"
	"PoliSim/dataValidation"
	"PoliSim/database"
)

func GetStartPage(acc *dataExtraction.AccountAuth, val dataValidation.ValidationMessage) Node {
	return getBasePageWrapper(
		getCustomPageHeader(Translation["welcomMessage"]),
		IfElse(acc.Role == database.NotLoggedIn,
			// if the user is not logged in give him the possibility to log in
			getFormStandardForm("form", POST, "/"+APIPreRoute+string(Login),
				getSimpleTextInput("username", "username", "", Translation["username"]),
				getInput("password", "password", "", Translation["password"], "password", "", ""),
				getSubmitButton(Translation["loginButton"])),
			// otherwise display his name and a button to log out
			El(DIV, Attr(CLASS, "flex flex-col"),
				El(P, Attr(CLASS, "mt-4"), Text(Translation["loggedInAccountMessage"], acc.DisplayName)),
				El(BUTTON, Attr(TYPE, "submit"), Attr(HXPOST, "/"+APIPreRoute+string(Logout)),
					Attr(HXTARGET, "#"+MainBodyID), Attr(HXSWAP, "outerHTML"), Attr(HXINCLUDE, "#"+InformationID),
					Attr(CLASS, buttonClassAttribute+" mt-2"), Text(Translation["logoutButton"])),
			)),
		GetMessage(val),
		Raw(RawStartPageContent),
	)
}
