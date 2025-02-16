//go:build DE

package loc

import (
	"html/template"
	"log"
	"strings"
)

func init() {
	log.Println("Using the German Language Configuration")
}

const (
	LanguageTag = "de"

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
	ErrorLoadingFlairInfoForAccount        = "Fehler beim Laden der Flairs für den Autor"
	ErrorTitleTooLong                      = "Titel überschreitet die maximal erlaubte Länge von %d Zeichen"
	ErrorSearchingForAccountNames          = "Es ist ein Fehler bei der Suche nach der Accountnamensliste aufgetreten"
	MissingPermissions                     = "Fehlende Berechtigung"
	MissingPermissionForAccountInfo        = "Fehlende Berechtigung um die Informationen für diesen Account anzufordern"
	ErrorLoadingAccountNames               = "Konnte Accountnamen nicht laden"

	// Database Documents

	DocumentCommentContentRemovedHTML    = template.HTML("<code>[Inhalt wurde entfernt]</code>")
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
	AccountPasswordTooShort                    = "Das Passwort hat weniger als %d Zeichen"
	AccountSelectedInvalidRole                 = "Die ausgewählte Rolle für den Nutzer ist nicht erlaubt"
	AccountNotAllowedToCreateAccountOfThatRank = "Sie sind nicht berechtigt, einen Account mit den selben oder höheren Berechtigungen zu erstellen"
	AccountPasswordHashingError                = "Es ist ein Fehler beim Verschlüsseln des Passworts aufgetreten"
	AccountCreationError                       = "Der Nutzer konnte nicht erstellt werden\nBitte überprüfe, ob Anzeigename oder Nutzername einzigartig sind"
	AccountSuccessfullyCreated                 = "Account erfolgreich erstellt\nDer Nutzername ist: %s\nDas Passwort ist: %s"

	AccountSearchedNameDoesNotCorrespond  = "Der gesuchte Name ist mit keinem Account verbunden"
	AccountErrorFindingNamesForOwner      = "Konnte Namen für mögliche Accountbesitzer nicht laden"
	AccountFoundSearchedName              = "Gesuchten Account gefunden"
	AccountErrorSearchingNameList         = "Es ist ein Fehler bei der Suche nach den Namenslisten aufgetreten"
	AccountErrorNoAccountToModify         = "Konnte keinen Account zum Bearbeiten finden"
	AccountNoPermissionToEdit             = "Sie besitzen nicht die Berechtigung, diesen Account anzupassen"
	AccountRoleIsNotAllowed               = "Die ausgewählte Rolle ist nicht erlaubt"
	AccountErrorWhileUpdating             = "Es ist ein Fehler beim Speichern der Accountupdates aufgetreten"
	AccountPressUserOwnerIsPressUser      = "Ein Presse-Nutzer kann kein Besitzer eines anderen Presse-Nutzers sein"
	AccountPressUserOwnerRemovingError    = "Es ist ein Fehler beim Entfernen des bisherigen Besitzers aufgetreten"
	AccountPressUserOwnerAddError         = "Es ist ein Fehler beim Hinzufügen des neuen Besitzers aufgetreten"
	AccountSuccessfullyUpdated            = "Account erfolgreich angepasst"
	AccountSearchedNamesDoesNotCorrespond = "Die gesuchten Namen sind mit keinem Account verbunden"

	AccountFontSizeMustBeBiggerThen          = "Die Seitenskalierung darf nicht kleiner als %d%% sein"
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
	DocumentInvalidColor                           = "Eine der übergebenen Farben ist kein valider, 6-stelliger Hexadezimal-Code"
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

	DocumentCreatePostHasInvalidVisibility = "Die angegeben Sichtbarkeit für das Dokument ist nicht vereinbar mit der ausgewählten Organisation"
	DocumentCreatePostNotAllowedError      = "Die Kombination aus Account und Organisation ist nicht erlaubt"
	DocumentCreatePostError                = "Fehler beim Erstellen des Dokuments"

	DocumentTimeNotInAreaDiscussion              = "Der angegebene Zeitstempel ist entweder in weniger als 24 Stunden oder in mehr als 15 Tagen"
	DocumentCreateDiscussionHasInvalidVisibility = "Die angegeben Sichtbarkeit für die Diskussion ist nicht vereinbar mit der ausgewählten Organisation"
	DocumentCreateDiscussionNotAllowedError      = "Die Kombination aus Account und Organisation ist nicht erlaubt"
	DocumentCreateDiscussionError                = "Fehler beim Erstellen der Diskussion"
	DocumentMissingPermissionForComment          = "Fehlende Berechtigung, um mit diesem Account ein Kommentar zu erstellen"
	DocumentErrorWhileSavingComment              = "Fehler beim Speichern des Kommentars"

	DocumentSearchErrorVotes               = "Es ist ein Fehler bei der Suche nach den Abstimmungen des Accounts aufgetreten"
	DocumentTimeNotInAreaVote              = "Das angegebene Datum ist nach der Zeitanpassung entweder in weniger als 24 Stunden oder in mehr als 15 Tagen"
	DocumentCreateVoteHasInvalidVisibility = "Die angegeben Sichtbarkeit für die Abstimmung ist nicht vereinbar mit der ausgewählten Organisation"
	DocumentCreateVoteNotAllowedError      = "Die Kombination aus Account und Organisation ist nicht erlaubt"
	DocumentCreateVoteHasNoAttachedVotes   = "Das Abstimmungsdokument hat keine angehängten Abstimmungen"
	DocumentCreateVoteError                = "Fehler beim Erstellen der Abstimmung"

	DocumentCouldNotFilterReaders      = "Konnte Lesernamensliste nicht filtern"
	DocumentCouldNotFilterParticipants = "Konnte Teilnehmernamensliste nicht filtern"

	DocumentCouldNotLoadPersonalVote  = "Konnte die ausgewählte Abstimmung nicht laden"
	DocumentInvalidVoteNumber         = "Die ausgewählte Nummer für die Abstimmung ist nicht zulässig"
	DocumentInvalidVoteType           = "Der ausgewählte Abstimmungstyp für die Abstimmung ist nicht zulässig"
	DocumentInvalidNumberMaxVotes     = "Die maximale Stimmenzahl pro Nutzer darf nicht kleiner als 1 sein für den ausgewählten Abstimmungstypen"
	DocumentAmountAnswersTooSmall     = "Es muss mindestens eine Antwort zur Abstimmung stehen"
	DocumentVoteMustHaveAQuestion     = "Die Abstimmung muss eine Frage haben, über die abgestimmt wird"
	DocumentVoteQuestionTooLong       = "Die Abstimmungsfrage darf nicht länger als %d Zeichen sein"
	DocumentErrorSavingUserVote       = "Es ist ein Fehler beim Speichern der Abstimmung aufgetreten"
	DocumentSuccessfullySavedUserVote = "Abstimmung erfolgreich gespeichert"

	DocumentNotAllowedToVoteWithThatAccount = "Fehlende Berechtigung um mit diesem Account abzustimmen"
	DocumentNotAllowedToVoteOnThis          = "Für diese Abstimmung kann keine Stimme abgegeben werden"
	DocumentVoteIsInvalid                   = "Die abgegebene Stimme ist invalide"
	DocumentVotePositionInvalid             = "Die ausgewählte Position der Antwort ist nicht gültig"
	DocumentVoteShareNotSmallerZero         = "Die Anzahl an Stimmen pro Antwort darf nicht kleiner als 0 sein"
	DocumentVoteSumTooBig                   = "Die Summe aller abgegebenen Stimmen überschreitet das festgelegte Maximum"
	DocumentVoteRankTooBig                  = "Einer der Ränge ist größer als maximal erlaubt"
	DocumentVoteInvalidDoubleRank           = "Der selbe Rang darf nicht doppelt vergeben werden"
	DocumentAlreadyVotedWithThatAccount     = "Mit dem Account wurde bereits abgestimmt"
	DocumentErrorWhileVoting                = "Fehler beim Versuch, die Stimme abzugeben\nÜberprüfe, ob der Account stimmberechtigt ist"
	DocumentSuccessfullyVoted               = "Stimme erfolgreich abgegeben"

	// Handler Letter

	LetterErrorLoadingRecipients          = "Konnte mögliche Empfängernamen nicht laden"
	LetterNotAllowedToPostWithThatAccount = "Der Brief darf nicht mit dem angegebenen Account als Autor verschickt werden"
	LetterRecipientListUnvalidated        = "Konnte Empfängerliste nicht validieren"
	LetterNeedAtLeastOneRecipient         = "Die Anzahl an Empfängern für den Brief darf nicht 0 sein"
	LetterAllowedToBeSent                 = "Der Brief darf so versendet werden"
	LetterErrorWhileSending               = "Es ist ein Fehler beim Erstellen des Briefs aufgetreten"
	LetterSuccessfullySendLetter          = "Brief erfolgreich erstellt"

	// Handler Newspaper

	NewspaperCouldNotLoadAllNewspaperForAccount = "Konnte nicht alle möglichen Zeitungen für ausgewählten Account finden"
	NewspaperSubtitleTooLong                    = "Untertitel überschreitet die maximal erlaubte Länge von %d Zeichen"
	NewspaperMissingPermissionForNewspaper      = "Fehlende Berechtigung, um mit diesem Account in dieser Zeitung zu posten"
	NewspaperErrorWhileCreatingArticle          = "Fehler beim Erstellen des Artikels"
	NewspaperSuccessfullyCreatedArticle         = "Artikel erfolgreich erstellt"

	NewspaperErrorLoadingNewspaperNames   = "Konnte Zeitungsnamen nicht laden"
	NewspaperErrorWhileCreatingNewspaper  = "Fehler beim Erstellen der Zeitung (überprüfe ob die Zeitung bereits existiert)"
	NewspaperSuccessfullyCreatedNewspaper = "Zeitung erfolgreich erstellt"
	NewspaperErrorWhileSearchingNewspaper = "Fehler bei der Suche der Zeitung"
	NewspaperSuccessfullyFoundNewspaper   = "Zeitung gefunden"
	NewspaperErrorWhileChangingNewspaper  = "Fehler beim Anpassen der Zeitung"
	NewspaperErrorWhileAddingReporters    = "Fehler beim Hinzufügen der neuen Autoren zur Zeitung"
	NewspaperSuccessfullyChangedNewspaper = "Zeitung angepasst"

	NewspaperErrorDuringPublication       = "Es ist ein Fehler beim Publizieren aufgetreten"
	NewspaperRejectionMessageEmpty        = "Der Zurückweisungsgrund darf nicht leer sein"
	NewspaperErrorFindingArticleToReject  = "Konnte keinen Artikel mit der angegeben ID finden, welcher noch nicht publiziert wurde"
	NewspaperErrorDeletingArticle         = "Konnte den Artikel nicht löschen"
	NewspaperFormatTitleForRejection      = "Zurückweisung des Artikels '%s' geschrieben für %s"
	NewspaperFormatContentForRejection    = "# Zurückweisungsgrund\n\n%s\n\n# Artikelinhalt\n\n```%s```"
	NewspaperErrorCreatingRejectionLetter = "Fehler beim Erstellen des Briefs an den Autor des Artikels"

	// Handler Notes

	NoteAuthorIsInvalid         = "Mit dem ausgewählten Autor ist es nicht möglich eine Notiz zu verfassen"
	NoteErrorWhileCreatingNote  = "Es ist ein Fehler beim Erstellen der Notiz aufgetreten"
	NoteSuccessfullyCreatedNote = "Notiz erfolgreich erstellt"

	// Handler Organisations

	OrganisationGeneralInformationEmpty                          = "Organisationsname, Hauptgruppe oder Untergruppe ist leer"
	OrganisationGeneralNameTooLong                               = "Organisationsname überschreitet die maximal erlaubte Länge von %d Zeichen"
	OrganisationGeneralMainGroupTooLong                          = "Hauptgruppe überschreitet die maximal erlaubte Länge von %d Zeichen"
	OrganisationGeneralSubGroupTooLong                           = "Untergruppe überschreitet die maximal erlaubte Länge von %d Zeichen"
	OrganisationGeneralFlairContainsInvalidCharactersOrIsTooLong = "Flair enthält ein Komma, Semikolon oder ist länger als %d Zeichen"
	OrganisationGeneralInvalidVisibility                         = "Die ausgewählte Sichtbarkeit ist nicht valide"

	OrganisationErrorWhileCreating  = "Es ist ein Fehler beim Erstellen der Organisation aufgetreten (Überprüfe, ob der Name der Organisation einzigartig ist)"
	OrganisationSuccessfullyCreated = "Organisation erfolgreich erstellt"

	OrganisationNoOrganisationWithThatName        = "Der gesuchte Name ist mit keiner Organisation verbunden"
	OrganisationFoundOrganisation                 = "Gesuchte Organisation gefunden"
	OrganisationErrorSearchingForOrganisationList = "Es ist ein Fehler bei der Suche nach der Organisationsnamensliste aufgetreten"
	OrganisationErrorUpdatingOrganisation         = "Es ist ein Fehler beim Überarbeiten der Organisation aufgetreten"
	OrganisationErrorUpdatingOrganisationMember   = "Konnte Organisationsmitglieder nicht erfolgreich updaten"
	OrganisationSuccessfullyUpdated               = "Organisation erfolgreich angepasst"
	OrganisationNotFoundByName                    = "Konnte keine Organisation finden, die den Namen trägt"

	OrganisationHasNoMember        = "Diese Organisation hat keine Mitglieder"
	OrganisationMemberList         = "Mitglieder: %s"
	OrganisationHasNoAdministrator = "Diese Organisation hat keine Administratoren"
	OrganisationAdministratorList  = "Administratoren: %s"

	// Handler Titles

	TitleGeneralInformationEmpty                          = "Titelname, Hauptgruppe oder Untergruppe ist leer"
	TitleGeneralNameTooLong                               = "Titelname überschreitet die maximal erlaubte Länge von %d Zeichen"
	TitleGeneralMainGroupTooLong                          = "Hauptgruppe überschreitet die maximal erlaubte Länge von %d Zeichen"
	TitleGeneralSubGroupTooLong                           = "Untergruppe überschreitet die maximal erlaubte Länge von %d Zeichen"
	TitleGeneralFlairContainsInvalidCharactersOrIsTooLong = "Flair enthält ein Komma, Semikolon oder ist länger als %d Zeichen"

	TitleErrorWhileCreating  = "Es ist ein Fehler beim Erstellen des Titels aufgetreten (Überprüfe, ob der Name des Titel einzigartig ist)"
	TitleSuccessfullyCreated = "Titel erfolgreich erstellt"

	TitleNoTitleWithThatName           = "Der gesuchte Name ist mit keinem Titel verbunden"
	TitleFoundTitle                    = "Gesuchter Titel gefunden"
	TitleErrorSearchingForTitleList    = "Es ist ein Fehler bei der Suche nach der Titelnamensliste aufgetreten"
	TitleErrorWhileUpdatingTitle       = "Es ist ein Fehler beim Überarbeiten des Titels aufgetreten"
	TitleErrorWhileUpdatingTitleHolder = "Konnte Titel-Halter nicht erfolgreich updaten"
	TitleSuccessfullyUpdated           = "Titel erfolgreich angepasst"
	TitleNotFoundByName                = "Konnte keinen Titel finden, der den Namen trägt"

	TitleHasNoHolder        = "Dieser Titel wird von niemandem gehalten"
	TitleHolderFormatString = "Titel-Halter: %s"

	// Handler Markdown Go

	MarkdownParseError = "`Anfrage konnte nicht verarbeitet werden`"

	// Handler Pages Go

	PagesHomePage               = "Home"
	PagesNotFoundPage           = "Seite nicht gefunden"
	PagesCreateAccountPage      = "Nutzer erstellen"
	PagesMyProfilePage          = "Mein Profil"
	PagesEditAccountPage        = "Accounts anpassen"
	PagesNotesPage              = "Notizen anschauen"
	PagesCreateNotesPage        = "Notiz erstellen"
	PagesSearchNotesPage        = "Notizen durchsuchen"
	PagesCreateTitlePage        = "Titel erstellen"
	PagesEditTitlePage          = "Titel bearbeiten"
	PagesCreateOrganisationPage = "Organisation erstellen"
	PagesEditOrganisationPage   = "Organisation bearbeiten"
	PagesViewTitlePage          = "Titelübersicht"
	PagesViewOrganisationPage   = "Organisationsübersicht"
	PagesManageNewspaperPage    = "Zeitungen verwalten"
	PagesCreateArticlePage      = "Artikel erstellen"
	PagesViewPublicationPage    = "Zeitung"
	PagesSearchPublicationsPage = "Zeitungen durchsuchen"
	PagesSearchLetterPage       = "Briefe durchsuchen"
	PagesCreateLetterPage       = "Brief erstellen"
	PagesAdminSearchLetterPage  = "Briefsuche mit ID"
	PagesViewLetterPage         = "Briefansicht"
	PagesDocumentViewPage       = "Dokumentansicht"
	PagesCreateDocumentPage     = "Dokument erstellen"
	PagesCreateDiscussionPage   = "Diskussion erstellen"
	PagesCreateVoteElementPage  = "Abstimmungen verwalten"
	PagesCreateVotePage         = "Abstimmungsdokument erstellen"
	PagesSearchDocumentsPage    = "Dokumente durchsuchen"
	PagesViewVotePage           = "Abstimmungsansicht"
	PagesEditColorPage          = "Farbpaletten anpassen"
	PagesPersonDocumentPage     = "Persönliche Dokumente"
)

func SetHomePage(text []byte) {
	replaceMap["_home"]["$$home-page$$"] = string(text)
}

var replaceMap = map[string]map[string]string{
	"_home": {
		"$$home-page$$": "",
	},

	"base": {
		"{{/*base-1-language*/}}": LanguageTag,
	},
}

func LocaliseTemplateString(input []byte, name string) string {
	result := string(input)
	for key, value := range replaceMap {
		if name == key {
			for left, right := range value {
				result = strings.ReplaceAll(result, left, right)
			}
		}
	}
	return result
}
