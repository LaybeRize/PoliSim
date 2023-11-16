package composition

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func GetStartPage(acc *extraction.AccountAuth, val validation.Message) Node {
	return getBasePageWrapper(
		getCustomPageHeader(Translation["welcomeMessage"]),
		IfElse(acc.Role == database.NotLoggedIn,
			// if the user is not logged in give him the possibility to log in
			getFormStandardForm("form", POST, "/"+APIPreRoute+string(Login),
				getSimpleTextInput("username", "username", "", Translation["username"]),
				getInput("password", "password", "", Translation["password"], "password", "", ""),
				getSubmitButton("loginButton", Translation["loginButton"])),
			// otherwise display his name and a button to log out
			DIV(CLASS("flex flex-col"),
				P(CLASS("mt-4"), Text(Translation["loggedInAccountMessage"], acc.DisplayName)),
				BUTTON(TYPE("submit"), HXPOST("/"+APIPreRoute+string(Logout)),
					HXTARGET("#"+MainBodyID), HXSWAP("outerHTML"),
					CLASS(buttonClassAttribute+" mt-2"), Text(Translation["logoutButton"])),
			)),
		GetMessage(val),
		Raw(RawStartPageContent),
	)
}
