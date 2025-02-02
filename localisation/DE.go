//go:build DE

package loc

import (
	"html/template"
	"strings"
)

const (
	AdministrationName            = "Administration"
	AdministrationAccountName     = "Max Musteradministrator"
	AdministrationAccountUsername = ""
	AdministrationAccountPassword = ""
	StandardColorName             = "Standard Farbe"
	TimeFormatString              = "02.01.2006 15:04:05 MST"
	RequestParseError             = "Fehler bei der Verarbeitung der Eingangsinformationen"

	// Database Documents

	DocumentIsPublic                     = "Jeder kann dieses Dokument lesen."
	DocumentOnlyForMember                = "Leser: Alle Organisationsmitglieder"
	DocumentFormatStringForReader        = "Leser: Alle Organisationsmitglieder plus %s"
	DocumentTagAddInfo                   = "Nur Administratoren der Organisation dürfen Tags hinzufügen"
	DocumentParticipationEveryMember     = "Beteiligte: Alle Mitglieder der Organisation"
	DocumentParticipationOnlyAdmins      = "Beteiligte: Alle Administratoren der Organisation"
	DocumentParticipationEveryMemberPlus = "Beteiligte: Alle Mitglieder der Organisation plus %s"
	DocumentParticipationOnlyAdminsPlus  = "Beteiligte: Alle Administratoren der Organisation plus %s"
	DocumentParticipationFormatString    = "Beteiligte: %s"

	// Database Letter

	LetterRecipientsFormatString = "Empfänger: %s"
	LetterAcceptedFormatString   = "Zugestimmt: %s"
	LetterNoOneDeclined          = "Niemand hat abgelehnt"
	LetterDeclinedFormatString   = "Abgelehnt: %s"
	LetterNoDecisionFormatString = "Keine Entscheidung: %s"

	// Database Notes

	NotesContentRemovedHTML = template.HTML("<code>[Inhalt wurde entfernt]</code>")
	NotesRemovedTitelText   = "[Entfernt]"

	// Database Votes

	VoteNoIllegalVotesCasted = "Keine"

	// Handler Accounts

	AccountNotLoggedIn           = "Sie sind bereits angemeldet"
	AccountNameOrPasswordWrong   = "Nutzername oder Passwort falsch"
	AccountSuccessfullyLoggedIn  = "Erfolgreich angemeldet"
	AccountCurrentlyNotLoggedIn  = "Sie sind nicht angemeldet"
	AccountSuccessfullyLoggedOut = "Erfolgreich abgemeldet"

	AccountDisplayNameTooLongOrNotAtAll        = "Der Anzeigename des Accounts ist entweder leer oder überschreitet das Zeichenlimit von %d"
	AccountUsernameTooLongOrNotAtAll           = "Der Nutzername des Accounts ist entweder leer oder überschreitet das Zeichenlimit von %d"
	AccountPasswordTooShort                    = "Das Password hat weniger als %d Zeichen"
	AccountSelectedInvalidRole                 = "Die ausgewählte Rolle für den Nutzer ist nicht erlaubt"
	AccountNotAllowedToCreateAccountOfThatRank = "Sie sind nicht berechtigt eine Account mit den selben oder höheren Berechtigungen zu erstellen"
	AccountPasswordHashingError                = "Es ist ein Fehler beim verschlüsseln des Passworts aufgetreten"
	AccountCreationError                       = "Der Nutzer konnte nicht erstellt werden\nBitte überprüfe ob Anzeigename oder Nutzername einzigartig sind"
	AccountSuccessfullyCreated                 = "Account erfolgreich erstellt\nDer Nutzername ist: %s\nDas Passwort ist: %s"

	AccountSearchedNameDoesNotCorrespond  = "Der gesuchte Name ist mit keinem Account verbunden"
	AccountErrorFindingNamesForOwner      = "Konnte Namen für mögliche Accountbesitzer nicht laden"
	AccountFoundSearchedName              = "Gesuchten Account gefunden"
	AccountErrorSearchingNameList         = "Es ist ein Fehler bei der Suche nach den Namenslisten aufgetreten"
	AccountErrorNoAccountToModify         = "Konnte keinen Account zum bearbeiten finden"
	AccountNoPermissionToEdit             = "Sie besitzen nicht die Berechtigung diesen Account anzupassen"
	AccountRoleIsNotAllowed               = "Die ausgewählte Rolle ist nicht erlaubt"
	AccountErrorWhileUpdating             = "Es ist ein Fehler beim Speichern der Accountsupdates aufgetreten"
	AccountPressUserOwnerIsPressUser      = "Ein Presse-Nutzer kann kein Besitzer eines anderen Presse-Nutzers sein"
	AccountPressUserOwnerRemovingError    = "Es ist ein Fehler beim Entferne des bisherigen Besitzers aufgetreten"
	AccountPressUserOwnerAddError         = "Es ist ein Fehler beim Hinzufügen des neuen Besitzers aufgetreten"
	AccountSuccessfullyUpdated            = "Account erfolgreich angepasst"
	AccountSearchedNamesDoesNotCorrespond = "Die gesuchten Namen sind mit keinem Account verbunden"

	AccountFontSizeMustBeBiggerThen          = "Die Seitenskalierung kann nicht auf eine Zahl kleiner %d gesetzt werden"
	AccountGivenTimezoneInvalid              = "Die ausgewählte Zeitzone ist nicht erlaubt"
	AccountErrorSavingPersonalSettings       = "Fehler beim speichern der persönlichen Informationen"
	AccountPersonalSettingsSavedSuccessfully = "Einstellungen erfolgreich gespeichert\nLaden Sie die Seite neu, um den Effekt zu sehen"
	AccountWrongOldPassword                  = "Das alte Passwort ist falsch"
	AccountWrongRepeatPassword               = "Die Wiederholung stimmt nicht mit dem neuen Passwort überein"
	AccountNewPasswordMinimumLength          = "Das neue Passwort ist kürzer als %d Zeichen"
	AccountErrorHashingNewPassword           = "Fehler beim Verschlüsseln des neuen Passworts"
	AccountErrorSavingNewPassword            = "Fehler beim Speichern des neuen Passworts"
	AccountSuccessfullySavedNewPassword      = "Passwort erfolgreich angepasst"
)

var replaceMap = map[string]string{}

func LocaliseTemplateString(input []byte) string {
	result := string(input)
	for key, value := range replaceMap {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}
