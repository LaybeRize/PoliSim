//go:build EN

package loc

import "strings"

const (
	AdministrationName            = "Administration"
	AdministrationAccountName     = "John Administrator"
	AdministrationAccountUsername = ""
	AdministrationAccountPassword = ""

	StandardColorName = "Standard Color"
	TimeFormatString  = "2006-01-02 15:04:05 MST"

	RequestParseError                      = "error during request processing"
	CouldNotFindAllAuthors                 = "couldn't find all possible authors"
	ErrorFindingAllOrganisationsForAccount = "couldn't find all possible organisation for selected account"
	ContentOrBodyAreEmpty                  = "title or content empty"
	ContentIsEmpty                         = "content is empty"
	ErrorLoadingFlairInfoForAccount        = "error while loading author flairs"
	ErrorTitleTooLong                      = "title is longer then the maximum of %d characters"
	ErrorSearchingForAccountNames          = "error while searching for account names"
	MissingPermissions                     = "missing permission"
	MissingPermissionForAccountInfo        = "missing permission to access account information"
	ErrorLoadingAccountNames               = "couldn't load account names"

	// Database Documents

	DocumentIsPublic                     = "Everyone can view this document."
	DocumentOnlyForMember                = "Reader: Every organisation member"
	DocumentFormatStringForReader        = "Reader: Every organisation member and %s"
	DocumentTagAddInfo                   = "Only administrator of the organisation are allowed to add tags"
	DocumentParticipationEveryMember     = "Participants: Every member of the organisation"
	DocumentParticipationOnlyAdmins      = "Participants: Every administrator of the organisation"
	DocumentParticipationEveryMemberPlus = "Participants: Every member of the organisation and %s"
	DocumentParticipationOnlyAdminsPlus  = "Participants: Every administrator of the organisation and %s"
	DocumentParticipationFormatString    = "Participants: %s"

	// Database Letter

	LetterRecipientsFormatString = "Recipients: %s"
	LetterAcceptedFormatString   = "Agreed: %s"
	LetterNoOneDeclined          = "No one has declined"
	LetterDeclinedFormatString   = "Declined: %s"
	LetterNoDecisionFormatString = "No decision: %s"

	// Database Notes

	NotesContentRemovedHTML = template.HTML("<code>[Removed]</code>")
	NotesRemovedTitelText   = "[Removed]"

	// Database Votes

	VoteNoIllegalVotesCasted = "None"

	// Handler Accounts

	AccountNotLoggedIn           = "you are already logged in"
	AccountNameOrPasswordWrong   = "username or password wrong"
	AccountSuccessfullyLoggedIn  = "successfully logged in"
	AccountCurrentlyNotLoggedIn  = "you are already logged out"
	AccountSuccessfullyLoggedOut = "successfully logged out"

	AccountDisplayNameTooLongOrNotAtAll        = "the display name of the account is either empty or longer then %d character"
	AccountUsernameTooLongOrNotAtAll           = "the username of the account is either empty or longer then %d character"
	AccountPasswordTooShort                    = "the password is shorter then %d characters"
	AccountSelectedInvalidRole                 = "the selected role for the account is not permitted"
	AccountNotAllowedToCreateAccountOfThatRank = "missing permission to create an account with the same or higher privileges then oneself"
	AccountPasswordHashingError                = "error occurred while hashing the password"
	AccountCreationError                       = "account could not be created\nplease check if display name and username are unique"
	AccountSuccessfullyCreated                 = "account successfully created\nusername: %s\npassword: %s"

	AccountSearchedNameDoesNotCorrespond  = "the searched name has no corresponding account"
	AccountErrorFindingNamesForOwner      = "couldn't load names for possible account owner"
	AccountFoundSearchedName              = "found account"
	AccountErrorSearchingNameList         = "couldn't load name list"
	AccountErrorNoAccountToModify         = "couldn't find account to edit"
	AccountNoPermissionToEdit             = "missing permission to edit account"
	AccountRoleIsNotAllowed               = "the selected role for the account is not permitted"
	AccountErrorWhileUpdating             = "error while trying to save account updates"
	AccountPressUserOwnerIsPressUser      = "press account can't be owner of another press account"
	AccountPressUserOwnerRemovingError    = "error while trying to remove current owner"
	AccountPressUserOwnerAddError         = "error while trying to add new owner"
	AccountSuccessfullyUpdated            = "account successfully updated"
	AccountSearchedNamesDoesNotCorrespond = "the searched names don't correspond to any account"

	AccountFontSizeMustBeBiggerThen          = "the font scaling can't be set to any number smaller then %d%%"
	AccountGivenTimezoneInvalid              = "chosen timezone is not allowed"
	AccountErrorSavingPersonalSettings       = "couldn't save personal settings"
	AccountPersonalSettingsSavedSuccessfully = "personal settings saved successfully\nreload page to see the effects"
	AccountWrongOldPassword                  = "the old password is invalid"
	AccountWrongRepeatPassword               = "the repeated password is not equal to the new password"
	AccountNewPasswordMinimumLength          = "the new password has less then %d characters"
	AccountErrorHashingNewPassword           = "error while hashing the new password"
	AccountErrorSavingNewPassword            = "error while saving the new password"
	AccountSuccessfullySavedNewPassword      = "password successfully changed"

	// Handler Documents Todo: translate

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

	// Handler Letter

	LetterCouldNotLoadAuthors             = "Konnte nicht alle möglichen Autoren laden"
	LetterErrorLoadingRecipients          = "Konnte mögliche Empfängernamen nicht laden"
	LetterNotAllowedToPostWithThatAccount = "Der Brief darf nicht mit dem angegeben Account als Autor verschickt werden"
	LetterRecipientListUnvalidated        = "Konnte Empfängerliste nicht validieren"
	LetterNeedAtLeastOneRecipient         = "Die Anzahl an Empfängern für den Brief darf nicht 0 sein"
	LetterAllowedToBeSent                 = "Der Brief darf so versendet werden"
	LetterErrorWhileSending               = "Es ist ein Fehler beim erstellen des Briefs aufgetreten"
	LetterSuccessfullySendLetter          = "Brief erfolgreich erstellt"

	// Handler Newspaper

	NewspaperCouldNotLoadAllNewspaperForAccount = "Konnte nicht alle möglichen Zeitungen für ausgewählten Account finden"
	NewspaperSubtitleTooLong                    = "Untertitel überschreitet die maximal erlaubte Länge von %d Zeichen"
	NewspaperMissingPermissionForNewspaper      = "Fehlende Berechtigung um mit diesem Account in dieser Zeitung zu posten"
	NewspaperErrorWhileCreatingArticle          = "Fehler beim erstellen des Artikels"
	NewspaperSuccessfullyCreatedArticle         = "Artikel erfolgreich erstellt"

	NewspaperErrorLoadingNewspaperNames   = "Konnte Zeitungsnamen nicht laden"
	NewspaperErrorWhileCreatingNewspaper  = "Fehler beim Erstellen der Zeitung (überprüfe ob die Zeitung bereits existiert)"
	NewspaperSuccessfullyCreatedNewspaper = "Zeitung erfolgreich erstellt"
	NewspaperErrorWhileSearchingNewspaper = "Fehler bei der Suche der Zeitung"
	NewspaperSuccessfullyFoundNewspaper   = "Zeitung gefunden"
	NewspaperErrorWhileChangingNewspaper  = "Fehler beim Anpassen der Zeitung"
	NewspaperErrorWhileAddingReporters    = "Fehler beim hinzufügen der neuen Autoren zur Zeitung"
	NewspaperSuccessfullyChangedNewspaper = "Zeitung angepasst"

	NewspaperErrorDuringPublication       = "Es ist ein Fehler beim Publizieren aufgetreten"
	NewspaperRejectionMessageEmpty        = "Der Zurückweisungsgrund darf nicht leer sein"
	NewspaperErrorFindingArticleToReject  = "Konnte keinen Artikel mit der angegeben ID finden, welcher noch nicht publiziert wurde"
	NewspaperErrorDeletingArticle         = "Konnte den Artikel nicht löschen"
	NewspaperFormatTitleForRejection      = "Zurückweisung des Artikels '%s' geschrieben für %s"
	NewspaperFormatContentForRejection    = "# Zurückweisungsgrund\n\n%s\n\n# Artikelinhalt\n\n```%s```"
	NewspaperErrorCreatingRejectionLetter = "Fehler beim erstellen des Briefs an den Autor des Artikels"

	// Handler Notes

	NoteAuthorIsInvalid         = "Mit dem ausgewählte Autor ist es nicht möglich eine Notiz zu verfassen"
	NoteErrorWhileCreatingNote  = "Es ist ein Fehler beim erstellen der Notiz aufgetreten"
	NoteSuccessfullyCreatedNote = "Notiz erfolgreich erstellt"

	// Handler Organisations

	OrganisationGeneralInformationEmpty                          = "Organisationsname, Hauptgruppe oder Untergruppe ist leer"
	OrganisationGeneralNameTooLong                               = "Organisationsname überschreitet die maximal erlaubte Länge von %d Zeichen"
	OrganisationGeneralMainGroupTooLong                          = "Hauptgruppe überschreitet die maximal erlaubte Länge von %d Zeichen"
	OrganisationGeneralSubGroupTooLong                           = "Untergruppe überschreitet die maximal erlaubte Länge von %d Zeichen"
	OrganisationGeneralFlairContainsInvalidCharactersOrIsTooLong = "Flair enthält ein Komma, Semikolon oder ist länger als %d Zeichen"
	OrganisationGeneralInvalidVisibility                         = "Die ausgewählte Sichtbarkeit ist nicht valide"

	OrganisationErrorWhileCreating  = "Es ist ein Fehler beim erstellen der Organisation aufgetreten (Überprüf ob der Name der Organisation einzigartig ist)"
	OrganisationSuccessfullyCreated = "Organisation erfolgreich erstellt"

	OrganisationNoOrganisationWithThatName        = "Der gesuchte Name ist mit keiner Organisation verbunden"
	OrganisationFoundOrganisation                 = "Gesuchte Organisation gefunden"
	OrganisationErrorSearchingForOrganisationList = "Es ist ein Fehler bei der Suche nach der Organisationsnamensliste aufgetreten"
	OrganisationErrorUpdatingOrganisation         = "Es ist ein Fehler beim überarbeiten der Organisation aufgetreten"
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

	TitleErrorWhileCreating  = "Es ist ein Fehler beim erstellen des Titels aufgetreten (Überprüf ob der Name des Titel einzigartig ist)"
	TitleSuccessfullyCreated = "Titel erfolgreich erstellt"

	TitleErrorWhileUpdatingTitle       = "Es ist ein Fehler beim überarbeiten des Titels aufgetreten"
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
)

var replaceMap = map[string]string{}

func LocaliseTemplateString(input []byte) string {
	result := string(input)
	for key, value := range replaceMap {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}
