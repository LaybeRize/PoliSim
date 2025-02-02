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

	StandardColorName = "Standard Farbe"
	TimeFormatString  = "02.01.2006 15:04:05 MST"

	RequestParseError                      = "Fehler bei der Verarbeitung der Eingangsinformationen"
	CouldNotFindAllAuthors                 = "Konnte nicht alle möglichen Autoren finden"
	ErrorFindingAllOrganisationsForAccount = "Konnte nicht alle erlaubten Organisationen für ausgewählten Account finden"
	ContentOrBodyAreEmpty                  = "Titel oder Inhalt sind leer"
	ContentIsEmpty                         = "Inhalt ist leer"
	ErrorLoadingFlairInfoForAccount        = "Fehler beim laden der Flairs für den Autor"
	ErrorTitleTooLong                      = "Titel überschreitet die maximal erlaubte Länge von %d Zeichen"
	ErrorSearchingForAccountNames          = "Es ist ein Fehler bei der Suche nach der Accountnamensliste aufgetreten"
	MissingPermissions                     = "Fehlende Berechtigung"
	MissingPermissionForAccountInfo        = "Fehlende Berechtigung um die Informationen für diesen Account anzufordern"

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

	// Handler Documents

	DocumentGeneralMissingPermissionForDocumentCreation = "Fehlende Berechtigung um mit diesem Account ein Dokument zu erstellen"
	DocumentGeneralFunctionNotAvailable                 = "Diese Funktion ist nicht verfügbar"
	DocumentGeneralTimestampInvalid                     = "Der angegebene Zeitstempel für das Ende ist nicht gültig"

	DocumentColorPaletteNameNotEmpty               = "Name der Farbpalette darf nicht leer sein"
	DocumentInvalidColor                           = "Einer der übergebene Farben ist kein valider, 6-stelliger Hexadezimal-Code"
	DocumentErrorCreatingColorPalette              = "Fehler beim Erstellen der Farbpalette"
	DocumentSuccessfullyCreatedChangedColorPalette = "Farbpalette erfolgreich erstellt/bearbeitet"
	DocumentStandardColorNotAllowedToBeDeleted     = "Die Standardfarbe darf nicht gelöscht werden"
	DocumentErrorDeletingColorPalette              = "Fehler beim Löschen der Farbpalette"
	DocumentSuccessfullyDeletedColorPalette        = "Farbe erfolgreich gelöscht"

	DocumentTagTextEmpty              = "Der Tag-Text ist leer"
	DocumentTagTextTooLong            = "Der Tag-Text ist länger als %d Zeichen"
	DocumentTagColorInvalidBackground = "Die Farbe für den Hintergrund ist nicht valide"
	DocumentTagColorInvalidText       = "Die Farbe für den Text ist nicht valide"
	DocumentTagColorInvalidLink       = "Die Farbe für die Links ist nicht valide"
	DocumentTagCreationError          = "Fehler beim Erstellen des Tags"

	DocumentCreatePostError = "Fehler beim erstellen des Dokuments"

	DocumentTimeNotInAreaDiscussion     = "Der angegebene Zeitstempel ist entweder in weniger als 24 Stunden oder in mehr als 15 Tagen"
	DocumentCreateDiscussionError       = "Fehler beim erstellen der Diskussion"
	DocumentMissingPermissionForComment = "Fehlende Berechtigung um mit diesem Account ein Kommentar zu erstellen"
	DocumentErrorWhileSavingComment     = "Fehler beim Speichern des Kommentars"

	DocumentSearchErrorVotes  = "Es ist ein Fehler bei der Suche nach den Abstimmung des Accounts aufgetreten"
	DocumentTimeNotInAreaVote = "Das angegebene Datum ist nach der Zeitanpassung entweder in weniger als 24 Stunden oder in mehr als 15 Tagen"
	DocumentCreateVoteError   = "Fehler beim erstellen des Abstimmung"

	DocumentCloudNotFilterReaders      = "Konnte Lesernamensliste nicht filtern"
	DocumentCloudNotFilterParticipants = "Konnte Teilnehmernamensliste nicht filtern"

	DocumentCouldNotLoadPersonalVote  = "Konnte die ausgewählte Abstimmung nicht laden"
	DocumentInvalidVoteNumber         = "Die ausgewählte Nummer für die Abstimmung ist nicht zulässig"
	DocumentInvalidVoteType           = "Der ausgewählte Abstimmungstyp für die Abstimmung ist nicht zulässig"
	DocumentInvalidNumberMaxVotes     = "Die maximale Stimmenzahl pro Nutzer darf nicht kleiner als 1 sein für den ausgewählten Abstimmungstypen"
	DocumentAmountAnswersTooSmall     = "Es muss mindestens eine Antwort zur Abstimmung stehen"
	DocumentVoteMustHaveAQuestion     = "Die Abstimmung muss eine Frage haben, über die abgestimmt wird"
	DocumentVoteQuestionTooLong       = "Die Abstimmungsfrage darf nicht länger als %d Zeichen sein"
	DocumentErrorSavingUserVote       = "Es ist ein Fehler beim speichern der Abstimmung aufgetreten"
	DocumentSuccessfullySavedUserVote = "Abstimmung erfolgreich gespeichert"

	DocumentNotAllowedToVoteWithThatAccount = "Fehlende Berechtigung um mit diesem Account abzustimmen"
	DocumentNotAllowedToVoteOnThis          = "Für diese Abstimmung kann keine Stimme abgegeben werden"
	DocumentVoteIsInvalid                   = "Die Abgegebene Stimme ist invalide"
	DocumentVotePositionInvalid             = "Die ausgewählte Position der Antwort ist nicht gültig"
	DocumentVoteShareNotSmallerZero         = "Die Anzahl an Stimmen pro Antwort darf nicht kleiner als 0 sein"
	DocumentVoteSumTooBig                   = "Die Summe aller abgegebenen Stimmen überschreitet das festgelegte Maximum"
	DocumentVoteRankTooBig                  = "Einer der Ränge ist größer als maximal erlaubt"
	DocumentVoteInvalidDoubleRank           = "Der selbe Rang darf nicht doppelt vergeben werden"
	DocumentAlreadyVotedWithThatAccount     = "Mit dem Account wurde bereits abgestimmt"
	DocumentErrorWhileVoting                = "Fehler beim Versuch die Stimme abzugeben\nÜberprüfe ob der Account stimmberechtigt ist"
	DocumentSuccessfullyVoted               = "Stimme erfolgreich abgegeben"
)

var replaceMap = map[string]string{}

func LocaliseTemplateString(input []byte) string {
	result := string(input)
	for key, value := range replaceMap {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}
