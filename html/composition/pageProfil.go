package composition

import (
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	. "PoliSim/html/builder"
)

func GetPersonalProfil(acc *extraction.AccountAuth) Node {
	return getBasePageWrapper(

		GetLoginThing(false),
		GetMessage(validation.Message{}),
	)
}

func GetLoginThing(swap bool) Node {
	return DIV(ID("password-div-id"), If(swap, HXSWAPOOB("true")),
		getFormStandardForm("form", PATCH, "/"+APIPreRoute+string(ChangePassword),
			getInput("ordPassword", "ordPassword", "", Translation["ordPassword"], "password", "", ""),
			getInput("newPassword", "newPassword", "", Translation["newPassword"], "password", "", ""),
			getInput("newPasswordAgain", "newPasswordAgain", "", Translation["newPasswordAgain"], "password", "", ""),
			getSubmitButton("loginButton", Translation["changePasswordButton"])),
	)
}
