//go:build EN

package loc

import (
	"html/template"
	"log"
	"strings"
)

func init() {
	log.Println("Using the English Language Configuration")
}

const (
	LanguageTag = "en"

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

	DocumentCommentContentRemovedHTML    = template.HTML("<code>[Removed]</code>")
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

	// Handler Documents

	DocumentGeneralMissingPermissionForDocumentCreation = "missing permission to create document with this account"
	DocumentGeneralFunctionNotAvailable                 = "this functionality is not available"
	DocumentGeneralTimestampInvalid                     = "the given end timestamp is invalid"

	DocumentColorPaletteNameNotEmpty               = "name of color palette is not allowed to be empty"
	DocumentInvalidColor                           = "one of the given color is not a valid 6 character hex code"
	DocumentErrorCreatingColorPalette              = "couldn't create color palette"
	DocumentSuccessfullyCreatedChangedColorPalette = "successfully created/edited color palette"
	DocumentStandardColorNotAllowedToBeDeleted     = "the standard color can't be deleted"
	DocumentErrorDeletingColorPalette              = "couldn't delete color palette"
	DocumentSuccessfullyDeletedColorPalette        = "successfully deleted color palette"

	DocumentTagTextEmpty              = "tag text is empty"
	DocumentTagTextTooLong            = "tag text more then %d characters"
	DocumentTagColorInvalidBackground = "background color not a valid hex code"
	DocumentTagColorInvalidText       = "text color not a valid hex code"
	DocumentTagColorInvalidLink       = "link color not a valid hex color"
	DocumentTagCreationError          = "error while trying to create tag"

	DocumentCreatePostError = "error while trying to create the document"

	DocumentTimeNotInAreaDiscussion     = "given end timestamp is either in less then 24 hours or in more then 15 days"
	DocumentCreateDiscussionError       = "error while trying to create discussion"
	DocumentMissingPermissionForComment = "missing permission to create comment with this account"
	DocumentErrorWhileSavingComment     = "error while trying to save comment"

	DocumentSearchErrorVotes  = "error while trying to retrieve prepared votes from the account"
	DocumentTimeNotInAreaVote = "given end date is after correction in less then 24 hours or in more then 15 days"
	DocumentCreateVoteError   = "error while trying to create vote"

	DocumentCouldNotFilterReaders      = "couldn't filter reader list"
	DocumentCouldNotFilterParticipants = "couldn't filter participants list"

	DocumentCouldNotLoadPersonalVote  = "error while trying to load vote preparation"
	DocumentInvalidVoteNumber         = "given number is invalid to retrieve a vote"
	DocumentInvalidVoteType           = "given vote type is invalid"
	DocumentInvalidNumberMaxVotes     = "the maximum votes that can be casted is not allowed to be less then 1"
	DocumentAmountAnswersTooSmall     = "vote must have at least one possible option to vote for"
	DocumentVoteMustHaveAQuestion     = "vote must have a question that is voted on"
	DocumentVoteQuestionTooLong       = "vote question is not allowed to have more then %d character"
	DocumentErrorSavingUserVote       = "couldn't save vote preparation"
	DocumentSuccessfullySavedUserVote = "vote preparation successfully saved"

	DocumentNotAllowedToVoteWithThatAccount = "missing permission to vote with this account"
	DocumentNotAllowedToVoteOnThis          = "no permission to vote for this"
	DocumentVoteIsInvalid                   = "the casted vote is invalid"
	DocumentVotePositionInvalid             = "the selected position is not on the ballot"
	DocumentVoteShareNotSmallerZero         = "votes per ballot options is not allowed to be less then 0"
	DocumentVoteSumTooBig                   = "sum of all votes is more then the allowed maximum"
	DocumentVoteRankTooBig                  = "at least one rank is bigger then the amount of options that are on the ballot"
	DocumentVoteInvalidDoubleRank           = "rank can't be given to two or more options"
	DocumentAlreadyVotedWithThatAccount     = "account already voted"
	DocumentErrorWhileVoting                = "couldn't cast vote\ncheck if the account used is allowed to vote"
	DocumentSuccessfullyVoted               = "vote casted successfully"

	// Handler Letter

	LetterErrorLoadingRecipients          = "couldn't load names for possible recipients"
	LetterNotAllowedToPostWithThatAccount = "letter can't be send with given author"
	LetterRecipientListUnvalidated        = "couldn't validate recipient list"
	LetterNeedAtLeastOneRecipient         = "number of recipients is not allowed to be zero"
	LetterAllowedToBeSent                 = "letter is allowed to be send like this"
	LetterErrorWhileSending               = "couldn't send letter"
	LetterSuccessfullySendLetter          = "letter successfully send"

	// Handler Newspaper

	NewspaperCouldNotLoadAllNewspaperForAccount = "couldn't load all possible newspaper for the selected author"
	NewspaperSubtitleTooLong                    = "Subtitle is longer then maximum of %d characters"
	NewspaperMissingPermissionForNewspaper      = "missing permission to post with this author in selected newspaper"
	NewspaperErrorWhileCreatingArticle          = "couldn't create article"
	NewspaperSuccessfullyCreatedArticle         = "article successfully created"

	NewspaperErrorLoadingNewspaperNames   = "couldn't load names for newspaper"
	NewspaperErrorWhileCreatingNewspaper  = "couldn't create newspaper (check if the newspaper already exists)"
	NewspaperSuccessfullyCreatedNewspaper = "newspaper successfully created"
	NewspaperErrorWhileSearchingNewspaper = "couldn't find newspaper with given name"
	NewspaperSuccessfullyFoundNewspaper   = "successfully found newspaper"
	NewspaperErrorWhileChangingNewspaper  = "couldn't change newspaper"
	NewspaperErrorWhileAddingReporters    = "error while trying to add reporters to the newspaper"
	NewspaperSuccessfullyChangedNewspaper = "newspaper successfully changed"

	NewspaperErrorDuringPublication       = "couldn't publish publication"
	NewspaperRejectionMessageEmpty        = "rejection reason for article can't be empty"
	NewspaperErrorFindingArticleToReject  = "couldn't find article with the given ID that isn't already published"
	NewspaperErrorDeletingArticle         = "couldn't delete article"
	NewspaperFormatTitleForRejection      = "Rejection of the article '%s' written for %s"
	NewspaperFormatContentForRejection    = "# Reason for Rejection\n\n%s\n\n# Article Content\n\n```%s```"
	NewspaperErrorCreatingRejectionLetter = "an error occurred while trying to send the rejection letter"

	// Handler Notes

	NoteAuthorIsInvalid         = "missing permission to post note with selected author"
	NoteErrorWhileCreatingNote  = "couldn't create note"
	NoteSuccessfullyCreatedNote = "note successfully created"

	// Handler Organisations

	OrganisationGeneralInformationEmpty                          = "organisation name, main group name or sub group name is empty"
	OrganisationGeneralNameTooLong                               = "organisation name is longer then the allowed maximum of %d characters"
	OrganisationGeneralMainGroupTooLong                          = "main group name is longer then the allowed maximum of %d characters"
	OrganisationGeneralSubGroupTooLong                           = "sub group name is longer then the allowed maximum of %d characters"
	OrganisationGeneralFlairContainsInvalidCharactersOrIsTooLong = "flair contains comma, semicolon or is longer then %d characters"
	OrganisationGeneralInvalidVisibility                         = "chosen organisation visibility is invalid"

	OrganisationErrorWhileCreating  = "couldn't create organisation (check if organisation name is unique)"
	OrganisationSuccessfullyCreated = "organisation successfully created"

	OrganisationNoOrganisationWithThatName        = "given name is not connected to any organisation"
	OrganisationFoundOrganisation                 = "organisation found"
	OrganisationErrorSearchingForOrganisationList = "couldn't load names for organisations"
	OrganisationErrorUpdatingOrganisation         = "couldn't save organisation updates"
	OrganisationErrorUpdatingOrganisationMember   = "couldn't save changes to organisation member and administration"
	OrganisationSuccessfullyUpdated               = "successfully updated organisation"
	OrganisationNotFoundByName                    = "could not find any organisation with the given name"

	OrganisationHasNoMember        = "Organisation has no member"
	OrganisationMemberList         = "Member: %s"
	OrganisationHasNoAdministrator = "Organisation has no administrator"
	OrganisationAdministratorList  = "Administrator: %s"

	// Handler Titles

	TitleGeneralInformationEmpty                          = "title name, main group name or sub group name is empty"
	TitleGeneralNameTooLong                               = "title name is longer then the allowed maximum of %d characters"
	TitleGeneralMainGroupTooLong                          = "main group name is longer then the allowed maximum of %d characters"
	TitleGeneralSubGroupTooLong                           = "sub group name is longer then the allowed maximum of %d characters"
	TitleGeneralFlairContainsInvalidCharactersOrIsTooLong = "flair contains comma, semicolon or is longer then %d characters"

	TitleErrorWhileCreating  = "couldn't create title (check if title name is unique)"
	TitleSuccessfullyCreated = "title successfully created"

	TitleNoTitleWithThatName           = "given name is not connected to any title"
	TitleFoundTitle                    = "title found"
	TitleErrorSearchingForTitleList    = "couldn't load names for titles"
	TitleErrorWhileUpdatingTitle       = "couldn't save title updates"
	TitleErrorWhileUpdatingTitleHolder = "couldn't save changes to title owners"
	TitleSuccessfullyUpdated           = "successfully updated title"
	TitleNotFoundByName                = "could not find any title with the given name"

	TitleHasNoHolder        = "Title has no owner"
	TitleHolderFormatString = "Title owner: %s"

	// Handler Markdown Go

	MarkdownParseError = "`Request could not be processed`"

	// Handler Pages Go

	PagesHomePage               = "Home"
	PagesNotFoundPage           = "Page Not Found"
	PagesCreateAccountPage      = "Create Account"
	PagesMyProfilePage          = "My Profil"
	PagesEditAccountPage        = "Edit Accounts"
	PagesNotesPage              = "Notes"
	PagesCreateNotesPage        = "Create Note"
	PagesSearchNotesPage        = "Search Notes"
	PagesCreateTitlePage        = "Create Title"
	PagesEditTitlePage          = "Edit Titles"
	PagesCreateOrganisationPage = "Create Organisation"
	PagesEditOrganisationPage   = "Edit Organisations"
	PagesViewTitlePage          = "Title Overview"
	PagesViewOrganisationPage   = "Organisation Overview"
	PagesManageNewspaperPage    = "Manage Newspapers"
	PagesCreateArticlePage      = "Create Article"
	PagesViewPublicationPage    = "Publication"
	PagesSearchPublicationsPage = "Search Publications"
	PagesSearchLetterPage       = "Search Letters"
	PagesCreateLetterPage       = "Create Letter"
	PagesAdminSearchLetterPage  = "Letter Search with ID"
	PagesViewLetterPage         = "Letter View"
	PagesDocumentViewPage       = "Document View"
	PagesCreateDocumentPage     = "Create Document"
	PagesCreateDiscussionPage   = "Create Discussion"
	PagesCreateVoteElementPage  = "Manage Ballots"
	PagesCreateVotePage         = "Create Vote"
	PagesSearchDocumentsPage    = "Search Documents"
	PagesViewVotePage           = "Vote View"
	PagesEditColorPage          = "Manage Color Palettes"
	PagesPersonDocumentPage     = "Personal Documents"
)

func SetHomePage(text []byte) {
	replaceMap["_home"]["$$home-page$$"] = string(text)
}

var replaceMap = map[string]map[string]string{
	"_home": {
		"$$home-page$$": "",
		"{{/*_home-1*/}}Herzlich willkommen, {{.Account.Name}}": "Welcome, {{.Account.Name}}",
		"{{/*_home-2*/}}Abmelden":                               "Sign out",
		"{{/*_home-3*/}}Nutzername":                             "Username",
		"{{/*_home-4*/}}Passwort":                               "Password",
		"{{/*_home-5*/}}Einloggen":                              "Sign in",
	},

	"_notFound": {
		"{{/*_notFound-1*/}}Die Seite, die Sie suchen, existiert nicht": "The Page you requested does not exist",
	},

	"accountCreate": {
		"{{/*accountCreate-1*/}}Anzeigename":          "Display Name",
		"{{/*accountCreate-2*/}}Nutzername":           "Username",
		"{{/*accountCreate-3*/}}Passwort":             "Password",
		"{{/*accountCreate-4*/}}Rolle":                "Role",
		"{{/*accountCreate-5*/}}Nutzer":               "User",
		"{{/*accountCreate-6*/}}Presse-Nutzer":        "Press User",
		"{{/*accountCreate-7*/}}Presse-Administrator": "Press Administrator",
		"{{/*accountCreate-8*/}}Administrator":        "Administrator",
		"{{/*accountCreate-9*/}}Oberadministrator":    "Head Administrator",
		"{{/*accountCreate-10*/}}Nutzer erstellen":    "Create User",
	},

	"accountEdit": {
		"{{/*accountEdit-1*/}}Zurück zur Suche":     "Back to Search",
		"{{/*accountEdit-2*/}}Anzeigename":          "Display Name",
		"{{/*accountEdit-3*/}}Nutzername":           "Username",
		"{{/*accountEdit-4*/}}Rolle":                "Role",
		"{{/*accountEdit-5*/}}Nutzer":               "User",
		"{{/*accountEdit-6*/}}Presse-Nutzer":        "Press User",
		"{{/*accountEdit-7*/}}Presse-Administrator": "Press Administrator",
		"{{/*accountEdit-8*/}}Administrator":        "Administrator",
		"{{/*accountEdit-9*/}}Oberadministrator":    "Head Administrator",
		"{{/*accountEdit-10*/}}Blockiert":           "Blocked",
		"{{/*accountEdit-11*/}}Account-Besitzer":    "Account Owner",
		"{{/*accountEdit-12*/}}Nutzer anpassen":     "Update User",
		"{{/*accountEdit-13*/}}Anzeigename":         "Display Name",
		"{{/*accountEdit-14*/}}Nutzername":          "Username",
		"{{/*accountEdit-15*/}}Nutzer suchen":       "Search for User",
	},

	"documentColorEdit": {
		"{{/*documentColorEdit-1*/}}Farbpaletten":                   "Color Palette",
		"{{/*documentColorEdit-2*/}}Farbpalette auswählen":          "Select Color Palette",
		"{{/*documentColorEdit-3*/}}Name":                           "Name",
		"{{/*documentColorEdit-4*/}}Hintergrundfarbe":               "Background Color",
		"{{/*documentColorEdit-5*/}}Textfarbe":                      "Text Color",
		"{{/*documentColorEdit-6*/}}Link-Farbe":                     "Link Color",
		"{{/*documentColorEdit-7*/}}Farbpalette erstellen/anpassen": "Create/Update Color Palette",
		"{{/*documentColorEdit-8*/}}Farbpalette löschen":            "Delete Color Palette",
	},

	"documentCreate": {
		"{{/*documentCreate-1*/}}Titel":              "Title",
		"{{/*documentCreate-2*/}}Autor":              "Author",
		"{{/*documentCreate-3*/}}Organisation":       "Organisation",
		"{{/*documentCreate-4*/}}Inhalt":             "Content",
		"{{/*documentCreate-5*/}}Dokument erstellen": "Create Document",
		"{{/*documentCreate-6*/}}Vorschau":           "Preview",
	},

	"documentCreateDiscussion": {
		"{{/*documentCreateDiscussion-1*/}}Titel":                                                              "Title",
		"{{/*documentCreateDiscussion-2*/}}Autor":                                                              "Author",
		"{{/*documentCreateDiscussion-3*/}}Organisation":                                                       "Organisation",
		"{{/*documentCreateDiscussion-4*/}}Ende der Diskussion ({{.NavInfo.Account.TimeZone.String}})":         "End Timestamp for Discussion ({{.NavInfo.Account.TimeZone.String}})",
		"{{/*documentCreateDiscussion-5*/}}Diskussion ist öffentlich (Pflicht in öffentlichen Organisationen)": "Discussion is Public (Must be checked when posting in a Public Organisation)",
		"{{/*documentCreateDiscussion-6*/}}Alle Organisationsmitglieder dürfen teilnehmen":                     "All Organisation Member are allowed to Participate",
		"{{/*documentCreateDiscussion-7*/}}Alle Organisationsadministratoren dürfen teilnehmen":                "All Organisation Administrators are allowed to Participate",
		"{{/*documentCreateDiscussion-8*/}}Inhalt":                                                             "Content",
		"{{/*documentCreateDiscussion-9*/}}Leser und Teilnehmer überprüfen":                                    "Check Reader and Participants",
		"{{/*documentCreateDiscussion-10*/}}Diskussion erstellen":                                              "Create Discussion",
		"{{/*documentCreateDiscussion-11*/}}Vorschau":                                                          "Preview",
	},

	"documentCreateVote": {
		"{{/*documentCreateVote-1*/}}Titel":        "Title",
		"{{/*documentCreateVote-2*/}}Autor":        "Author",
		"{{/*documentCreateVote-3*/}}Organisation": "Organisation",
		"{{/*documentCreateVote-4*/}}Ende der Abstimmung (Endet immer um 23:50 UTC des ausgewählten Tages)":                                                             "End Timestamp for Vote (Always ends at 23:50 UTC of the selected day)",
		"{{/*documentCreateVote-5*/}}Abstimmung ist öffentlich (Pflicht in öffentlichen Organisationen)":                                                                "Vote is Public (Must be checked when posting in a Public Organisation)",
		"{{/*documentCreateVote-6*/}}Alle Organisationsmitglieder dürfen teilnehmen":                                                                                    "All Organisation Member are allowed to Participate",
		"{{/*documentCreateVote-7*/}}Alle Organisationsadministratoren dürfen teilnehmen":                                                                               "All Organisation Administrators are allowed to Participate",
		"{{/*documentCreateVote-8*/}}Abstimmungslis<span class=\"hover-target\">te &#x1F6C8;</span>":                                                                    "Vote Li<span class=\"hover-target\">st &#x1F6C8;</span>",
		"{{/*documentCreateVote-9*/}}ID der ausgewählten Abstimmung übertragen":                                                                                         "Insert ID of selected Vote",
		"{{/*documentCreateVote-10*/}}Angehängte Abstimmungen":                                                                                                          "Attached Votes",
		"{{/*documentCreateVote-11*/}}Inhalt":                                                                                                                           "Content",
		"{{/*documentCreateVote-12*/}}Leser und Teilnehmer überprüfen":                                                                                                  "Check Reader and Participants",
		"{{/*documentCreateVote-13*/}}Abstimmungsdokument erstellen":                                                                                                    "Create Vote Document",
		"{{/*documentCreateVote-14*/}}Vorschau":                                                                                                                         "Preview",
		"{{/*documentCreateVote-15*/}}Um Abstimmungen vorzubereiten, öffne die Seite unter <strong>Dokumente</strong> &#8594; <strong>Abstimmungen verwalten</strong>.": "To prepare a vote for a document open the page under the menu point <strong>Documents</strong> &#8594; <strong>Manage Vote Preparation</strong>.",
	},

	"documentCreateVoteElement": {
		"{{/*documentCreateVoteElement-1*/}}Abstimmungsnummer":                                              "Vote Number",
		"{{/*documentCreateVoteElement-2*/}}Sicher, dass du die Abstimmung wechseln willst?":                "Are you sure you want to switch selected Vote",
		"{{/*documentCreateVoteElement-3*/}}Abstimmungs-ID":                                                 "Vote ID",
		"{{/*documentCreateVoteElement-4*/}}Abstimmungsart":                                                 "Vote Type",
		"{{/*documentCreateVoteElement-5*/}}Eine Stimme pro Nutzer":                                         "Single Choice Vote",
		"{{/*documentCreateVoteElement-6*/}}Mehrere Stimmen pro Nutzer":                                     "Multiple Choice Vote",
		"{{/*documentCreateVoteElement-7*/}}Rangwahl":                                                       "Option Ranking",
		"{{/*documentCreateVoteElement-8*/}}Gewichtete Wahl":                                                "Weighted Vote",
		"{{/*documentCreateVoteElement-9*/}}Maximale Stimmen pro Nutzer (Nur relevant für Gewichtete Wahl)": "Maximum Amount of Votes per User (only relevant for Weighted Vote)",
		"{{/*documentCreateVoteElement-10*/}}Zeige Teilnehmerbezogene Stimmen während der Wahl":             "Show already casted Ballots during Vote Period",
		"{{/*documentCreateVoteElement-11*/}}Geheime Wahl":                                                  "Secret Ballot",
		"{{/*documentCreateVoteElement-12*/}}Frage":                                                         "Question",
		"{{/*documentCreateVoteElement-13*/}}Antwort hinzufügen":                                            "Add New Option",
		"{{/*documentCreateVoteElement-14*/}}Antworten":                                                     "Options",
		"{{/*documentCreateVoteElement-15*/}}Abstimmung erstellen/bearbeiten":                               "Create/Update Vote",
		"{{/*documentCreateVoteElement-16*/}}Momentane Abstimmungsnummer":                                   "Current Vote Number",
	},

	"documentPersonalSearch": {
		"{{/*documentPersonalSearch-1*/}}\"Die Anfrage hat zu einem Fehler auf der Serverseite geführt\"":     "\"Requested could not be processed. Internal Server Error\"",
		"{{/*documentPersonalSearch-2*/}}Es konnten keine Einträge gefunden werden":                           "No Entries found",
		"{{/*documentPersonalSearch-3*/}}<strong>{{if .Removed}}[Entfernt]{{else}}{{.Title}}{{end}}</strong>": "<strong>{{if .Removed}}[Removed]{{else}}{{.Title}}{{end}}</strong>",
		"{{/*documentPersonalSearch-4*/}}<i>Veröffentlicht am: {{.GetTimeWritten $acc}}</i>":                  "<i>Publish Date: {{.GetTimeWritten $acc}}</i>",
		"{{/*documentPersonalSearch-5*/}}Veröffentlicht von <i>{{.Author}}</i> im <i>{{.Organisation}}</i>":   "Written by <i>{{.Author}}</i> for <i>{{.Organisation}}</i>",
		"{{/*documentPersonalSearch-6*/}}&laquo; Vorherige Seite":                                             "&laquo; Previous Page",
		"{{/*documentPersonalSearch-7*/}}Nächste Seite &raquo;":                                               "Next Page &raquo;",
	},

	"documentSearch": {
		"{{/*documentSearch-1*/}}\"Die Anfrage hat zu einem Fehler auf der Serverseite geführt\"":                 "\"Requested could not be processed. Internal Server Error\"",
		"{{/*documentSearch-2*/}}Blockierte Dokumente anzeigen":                                                   "Show Blocked Documents",
		"{{/*documentSearch-3*/}}Anzahl der Ergebnisse":                                                           "Number of Entries per Page",
		"{{/*documentSearch-4*/}}Suchen":                                                                          "Search",
		"{{/*documentSearch-5*/}}Es konnten keine Einträge gefunden werden, die den Suchkriterien gerecht werden": "No Entries found",
		"{{/*documentSearch-6*/}}<strong>{{if .Removed}}[Entfernt]{{else}}{{.Title}}{{end}}</strong>":             "<strong>{{if .Removed}}[Removed]{{else}}{{.Title}}{{end}}</strong>",
		"{{/*documentSearch-7*/}}<i>Veröffentlicht am: {{.GetTimeWritten $acc}}</i>":                              "<i>Publish Date: {{.GetTimeWritten $acc}}</i>",
		"{{/*documentSearch-8*/}}Veröffentlicht von <i>{{.Author}}</i> im <i>{{.Organisation}}</i>":               "Written by <i>{{.Author}}</i> for <i>{{.Organisation}}</i>",
		"{{/*documentSearch-9*/}}&laquo; Vorherige Seite":                                                         "&laquo; Previous Page",
		"{{/*documentSearch-10*/}}Nächste Seite &raquo;":                                                          "Next Page &raquo;",
	},

	"documentView": {
		"{{/*documentView-1*/}}Das Dokument wurde entfernt":                                                              "Document has been removed",
		"{{/*documentView-2*/}}Hinzugefügt am {{$tag.GetTimeWritten $acc}}":                                              "Added {{$tag.GetTimeWritten $acc}}",
		"{{/*documentView-3*/}}Geschrieben von: {{.Document.GetAuthor}}":                                                 "Written by: {{.Document.GetAuthor}}",
		"{{/*documentView-4*/}}Veröffentlicht in: {{.Document.Organisation}}":                                            "Published in: {{.Document.Organisation}}",
		"{{/*documentView-5*/}}Verfasst am: {{.Document.GetTimeWritten .NavInfo.Account}}":                               "Publish Date: {{.Document.GetTimeWritten .NavInfo.Account}}",
		"{{/*documentView-6*/}}Die {{if .Document.IsDiscussion}}Diskussion{{else}}Abstimmung{{end}} ist bereits vorbei.": "The {{if .Document.IsDiscussion}}Discussion{{else}}Vote{{end}} has ended.",
		"{{/*documentView-7*/}}Ende war am {{.Document.GetTimeEnd .NavInfo.Account}}":                                    "It ended {{.Document.GetTimeEnd .NavInfo.Account}}",
		"{{/*documentView-8*/}}Endet: {{.Document.GetTimeEnd .NavInfo.Account}}":                                         "Ends: {{.Document.GetTimeEnd .NavInfo.Account}}",
		"{{/*documentView-9*/}}Tag erstellen":                                                                            "Create Document-Tag",
		"{{/*documentView-10*/}}Willst du den Tag so hinzufügen?":                                                        "Confirm Tag Creation",
		"{{/*documentView-11*/}}Farbpaletten":                                                                            "Color Palette",
		"{{/*documentView-12*/}}Farbe aus Farbpalette kopieren":                                                          "Copy Colors from Selected Color Palette",
		"{{/*documentView-13*/}}Inhalt":                                                                                  "Content",
		"{{/*documentView-14*/}}Hintergrundfarbe":                                                                        "Background Color",
		"{{/*documentView-15*/}}Textfarbe":                                                                               "Text Color",
		"{{/*documentView-16*/}}Link-Farbe":                                                                              "Link Color",
		"{{/*documentView-17*/}}Referenzen zu anderen Dokumenten":                                                        "Document References",
		"{{/*documentView-18*/}}Tag hinzufügen":                                                                          "Add Tag to Document",
		"{{/*documentView-19*/}}{{if .Document.Removed}}Dokument wieder freigeben{{else}}Dokument blockieren{{end}}":     "{{if .Document.Removed}}Restore Document{{else}}Remove Document{{end}}",
		"{{/*documentView-20*/}}Kommentare":                                                                              "Comments",
		"{{/*documentView-21*/}}Geschrieben von: {{$comment.GetAuthor}}":                                                 "Written by: {{$comment.GetAuthor}}",
		"{{/*documentView-22*/}}Verfasst am: {{$comment.GetTimeWritten $acc}}":                                           "Publish Date: {{$comment.GetTimeWritten $acc}}",
		"{{/*documentView-23*/}}{{if $comment.Removed}}Kommentar wieder freigeben{{else}}Kommentar blockieren{{end}}":    "{{if $comment.Removed}}Restore Comment{{else}}Remove Comment{{end}}",
		"{{/*documentView-24*/}}Autor":                                                                                   "Author",
		"{{/*documentView-25*/}}Kommentar schreiben":                                                                     "Write Comment",
		"{{/*documentView-26*/}}Vorschau":                                                                                "Preview",
		"{{/*documentView-27*/}}Abstimmungen":                                                                            "Votes",
	},

	"documentViewVote": {
		"{{/*documentViewVote-1*/}}Frage: {{.VoteInstance.Question}}":                                                                  "Question: {{.VoteInstance.Question}}",
		"{{/*documentViewVote-2*/}}Abstimmender Account":                                                                               "Voter",
		"{{/*documentViewVote-3*/}}{{$pos}}. Antwort: {{$answer}}":                                                                     "Option {{$pos}}: {{$answer}}",
		"{{/*documentViewVote-4*/}}{{$pos}}. Antwort: {{$answer}}":                                                                     "Option {{$pos}}: {{$answer}}",
		"{{/*documentViewVote-5*/}}Es dürfen maximal {{.VoteInstance.MaxVotes}} Stimmen vergeben werden":                               "You can cast a maximum of {{.VoteInstance.MaxVotes}} Votes",
		"{{/*documentViewVote-6*/}}Stimmen für die {{$pos}}. Antwort: {{$answer}}":                                                     "Votes for Option {{$pos}}: {{$answer}}",
		"{{/*documentViewVote-7*/}}Position 0 und kleiner bedeutet, dass die Antwort keinen Rang erhält. Der 1. Rang ist der höchste.": "Position 0 and smaller mean the Option is assigned no rank. The 1. rank is the highest.",
		"{{/*documentViewVote-8*/}}Rang der {{$pos}}. Antwort: {{$answer}}":                                                            "Rank of Option {{$pos}}: {{$answer}}",
		"{{/*documentViewVote-9*/}}Ungültige Stimme abgeben":                                                                           "Cast Invalid Vote",
		"{{/*documentViewVote-10*/}}Stimme abgeben":                                                                                    "Cast Vote",
		"{{/*documentViewVote-11*/}}Teilnahme setzte einen Account voraus":                                                             "Participation in the Vote requires an Account",
		"{{/*documentViewVote-12*/}}Es wird über die folgende Frage abgestimmt: <strong>{{.VoteInstance.Question}}</strong>":           "The posed Vote Question: <strong>{{.VoteInstance.Question}}</strong>",
		"{{/*documentViewVote-13*/}}Die Antwortmöglichkeiten sind: <strong>{{.VoteInstance.GetAnswerAsList}}</strong>":                 "The presented Options: <strong>{{.VoteInstance.GetAnswerAsList}}</strong>",
		"{{/*documentViewVote-14*/}}Das Ergebnis ist erst nach Ende der Abstimmung einsehbar":                                          "Results are available after the Vote has ended",
	},

	"letterAdminSearch": {
		"{{/*letterAdminSearch-1*/}}Brief ID":     "Letter ID",
		"{{/*letterAdminSearch-2*/}}Brief öffnen": "Open Letter",
	},

	"letterCreate": {
		"{{/*letterCreate-1*/}}Titel":                "Title",
		"{{/*letterCreate-2*/}}Autor":                "Author",
		"{{/*letterCreate-3*/}}Empfänger hinzufügen": "Add Recipient",
		"{{/*letterCreate-4*/}}Empfänger":            "Recipients",
		"{{/*letterCreate-5*/}}Mit Unterschrift":     "With Signatures",
		"{{/*letterCreate-6*/}}Inhalt":               "Content",
		"{{/*letterCreate-7*/}}Brief überprüfen":     "Check Letter",
		"{{/*letterCreate-8*/}}Brief erstellen":      "Create Letter",
		"{{/*letterCreate-9*/}}Vorschau":             "Preview",
	},

	"letterPersonalSearch": {
		"{{/*letterPersonalSearch-1*/}}\"Die Anfrage hat zu einem Fehler auf der Serverseite geführt\"":                 "\"Requested could not be processed. Internal Server Error\"",
		"{{/*letterPersonalSearch-2*/}}Account":                                                                         "Account",
		"{{/*letterPersonalSearch-3*/}}-- Alle Accounts --":                                                             "-- All Accounts --",
		"{{/*letterPersonalSearch-4*/}}Anzahl der Ergebnisse":                                                           "Number of Entries per Page",
		"{{/*letterPersonalSearch-5*/}}Suchen":                                                                          "Search",
		"{{/*letterPersonalSearch-6*/}}Es konnten keine Einträge gefunden werden, die den Suchkriterien gerecht werden": "No Entries found",
		"{{/*letterPersonalSearch-7*/}}<strong>{{.Title}}</strong> von {{.Author}}":                                     "<strong>{{.Title}}</strong> by {{.Author}}",
		"{{/*letterPersonalSearch-8*/}}Empfänger: {{.Recipient}}":                                                       "Recipient: {{.Recipient}}",
		"{{/*letterPersonalSearch-9*/}}<i>Versendet am: {{.GetTimeWritten $acc}}":                                       "<i>Publish Date: {{.GetTimeWritten $acc}}",
		"{{/*letterPersonalSearch-10*/}}&laquo; Vorherige Seite":                                                        "&laquo; Previous Page",
		"{{/*letterPersonalSearch-11*/}}Nächste Seite &raquo;":                                                          "Next Page &raquo;",
	},

	"letterView": {
		"{{/*letterView-1*/}}Geschrieben von: {{.Letter.GetAuthor}}":                          "Written by: {{.Letter.GetAuthor}}",
		"{{/*letterView-2*/}}<i>Verfasst am: {{.Letter.GetTimeWritten .NavInfo.Account}}</i>": "<i>Publish Date: {{.Letter.GetTimeWritten .NavInfo.Account}}</i>",
		"{{/*letterView-3*/}}Als {{.Letter.Recipient}} zustimmen":                             "Accept as {{.Letter.Recipient}}",
		"{{/*letterView-4*/}}Als {{.Letter.Recipient}} ablehnen":                              "Decline as {{.Letter.Recipient}}",
		"{{/*letterView-5*/}}ID: {{.Letter.ID}}":                                              "ID: {{.Letter.ID}}",
	},

	"newspaperCreate": {
		"{{/*newspaperCreate-1*/}}Titel":                   "Title",
		"{{/*newspaperCreate-2*/}}Untertitel":              "Subtitle",
		"{{/*newspaperCreate-3*/}}Autor":                   "Author",
		"{{/*newspaperCreate-4*/}}Zeitung":                 "Newspaper",
		"{{/*newspaperCreate-5*/}}-- Zeitung auswählen --": "-- Select Newspaper --",
		"{{/*newspaperCreate-6*/}}Eilmeldung":              "Breaking News",
		"{{/*newspaperCreate-7*/}}Inhalt":                  "Content",
		"{{/*newspaperCreate-8*/}}Artikel erstellen":       "Create Article",
		"{{/*newspaperCreate-9*/}}Vorschau":                "Preview",
	},

	"newspaperManage": {
		"{{/*newspaperManage-1*/}}Zeitung erstellen":                                            "Create Newspaper",
		"{{/*newspaperManage-2*/}}Name":                                                         "Name",
		"{{/*newspaperManage-3*/}}Zeitung erstellen":                                            "Create Newspaper",
		"{{/*newspaperManage-4*/}}Zeitung verändern":                                            "Update Newspaper",
		"{{/*newspaperManage-5*/}}Name":                                                         "Name",
		"{{/*newspaperManage-6*/}}Autor hinzufügen":                                             "Add Reporter",
		"{{/*newspaperManage-7*/}}Autoren":                                                      "Reporter",
		"{{/*newspaperManage-8*/}}Zeitung suchen":                                               "Search Newspaper",
		"{{/*newspaperManage-9*/}}Zeitung anpassen":                                             "Update Newspaper",
		"{{/*newspaperManage-10*/}}Es ist ein Fehler beim Suchen der Publikationen aufgetreten": "Internal Error during Publication Retrieval",
		"{{/*newspaperManage-11*/}}Es konnten keine Publikationen gefunden werden":              "No Publications",
		"{{/*newspaperManage-12*/}}Zeitung: <strong>{{.NewspaperName}}</strong>":                "Newspaper: <strong>{{.NewspaperName}}</strong>",
		"{{/*newspaperManage-13*/}}<i>Erstellt am: {{.GetPublishedDate $acc}}</i>":              "<i>Created Date: {{.GetPublishedDate $acc}}</i>",
	},

	"newspaperPubView": {
		"{{/*newspaperPubView-1*/}}Es ist ein Fehler beim Verarbeiten der Publikation für den Nutzer aufgetreten": "Publication couldn't be loaded",
		"{{/*newspaperPubView-2*/}}Sonderausgabe vom {{.Publication.GetPublishedDate $acc}}":                      "Breaking News published {{.Publication.GetPublishedDate $acc}}",
		"{{/*newspaperPubView-3*/}}Ausgabe vom {{.Publication.GetPublishedDate $acc}}":                            "Normal Volume published {{.Publication.GetPublishedDate $acc}}",
		"{{/*newspaperPubView-4*/}}Für diese Publikation existieren noch keine Artikel":                           "No Articles yet",
		"{{/*newspaperPubView-5*/}}Publikation freigeben":                                                         "Publish Volume",
		"{{/*newspaperPubView-6*/}}Geschrieben von: {{.GetAuthor}}":                                               "Written by: {{.GetAuthor}}",
		"{{/*newspaperPubView-7*/}}<i>Verfasst am: {{.GetTimeWritten $acc}}</i>":                                  "<i>Publish Date: {{.GetTimeWritten $acc}}</i>",
		"{{/*newspaperPubView-8*/}}Artikel zurückweisen":                                                          "Reject Article",
		"{{/*newspaperPubView-9*/}}Zurückweisungsgrund":                                                           "Reason for Rejection",
		"{{/*newspaperPubView-10*/}}Artikel zurückweisen":                                                         "Reject Article",
	},

	"newspaperSearch": {
		"{{/*newspaperSearch-1*/}}\"Die Anfrage hat zu einem Fehler auf der Serverseite geführt\"":                 "\"Requested could not be processed. Internal Server Error\"",
		"{{/*newspaperSearch-2*/}}Suchanfrage":                                                                     "Search Query",
		"{{/*newspaperSearch-3*/}}Anzahl der Ergebnisse":                                                           "Number of Entries per Page",
		"{{/*newspaperSearch-4*/}}Suchen":                                                                          "Search",
		"{{/*newspaperSearch-5*/}}Es konnten keine Einträge gefunden werden, die den Suchkriterien gerecht werden": "No Entries found",
		"{{/*newspaperSearch-6*/}}<strong>{{.NewspaperName}}</strong>":                                             "<strong>{{.NewspaperName}}</strong>",
		"{{/*newspaperSearch-7*/}}<i>Veröffentlicht am: {{.GetPublishedDate $acc}}</i>":                            "<i>Publish Date: {{.GetPublishedDate $acc}}</i>",
		"{{/*newspaperSearch-8*/}}&laquo; Vorherige Seite":                                                         "&laquo; Previous Page",
		"{{/*newspaperSearch-9*/}}Nächste Seite &raquo;":                                                           "Next Page &raquo;",
	},

	"noteCreate": {
		"{{/*noteCreate-1*/}}Referenzen (Komma-seperiert)": "References (comma seperated)",
		"{{/*noteCreate-2*/}}Titel":                        "Title",
		"{{/*noteCreate-3*/}}Autor":                        "Author",
		"{{/*noteCreate-4*/}}Inhalt":                       "Content",
		"{{/*noteCreate-5*/}}Notiz erstellen":              "Create Note",
		"{{/*noteCreate-6*/}}Vorschau":                     "Preview",
	},

	"notesSearch": {
		"{{/*noteSearch-1*/}}\"Die Anfrage hat zu einem Fehler auf der Serverseite geführt\"": "\"Requested could not be processed. Internal Server Error\"",
		"{{/*noteSearch-2*/}}Suchanfrage":                 "Query",
		"{{/*noteSearch-3*/}}Blockierte Notizen anzeigen": "Show Blocked Notes",
		"{{/*noteSearch-4*/}}Anzahl der Ergebnisse":       "Number of Entries per Page",
		"{{/*noteSearch-5*/}}Suchen":                      "Search",
		"{{/*noteSearch-6*/}}Es konnten keine Einträge gefunden werden, die den Suchkriterien gerecht werden": "No Entries found",
		"{{/*noteSearch-7*/}}<strong>{{.Title}}</strong> von {{.GetAuthor}}":                                  "<strong>{{.Title}}</strong> by {{.GetAuthor}}",
		"{{/*noteSearch-8*/}}<i>Veröffentlicht am: {{.GetTimePostedAt $acc}}</i>":                             "<i>Publish Date: {{.GetTimePostedAt $acc}}</i>",
		"{{/*noteSearch-9*/}}&laquo; Vorherige Seite":                                                         "&laquo; Previous Page",
		"{{/*noteSearch-10*/}}Nächste Seite &raquo;":                                                          "Next Page &raquo;",
	},

	"notesView": {
		"{{/*noteView-1*/}}\"Die Anfrage hat zu einem Fehler auf der Serverseite geführt\"": "\"Requested could not be processed. Internal Server Error\"",
		"{{/*noteView-2*/}}Schreibe eine eigene Notiz zu allen offenen Beiträgen":           "Write Note referencing all open Notes",
		"{{/*noteView-3*/}}ID: {{.ID}}":                                                          "ID: {{.ID}}",
		"{{/*noteView-4*/}}Geschrieben von: {{.GetAuthor}}":                                      "Written by: {{.GetAuthor}}",
		"{{/*noteView-5*/}}<i>Veröffentlicht am: {{.GetTimePostedAt $acc}}</i>":                  "<i>Publish Date: {{.GetTimePostedAt $acc}}</i>",
		"{{/*noteView-6*/}}{{if .Removed}}Notiz wieder freigeben{{else}}Notiz blockieren{{end}}": "{{if .Removed}}Restore Note{{else}}Remove Note{{end}}",
		"{{/*noteView-7*/}}Schreibe eine eigene Notiz zu diesem Beitrag":                         "Write Note referencing this Note",
		"{{/*noteView-8*/}}Referenzen":                                                           "References",
		"{{/*noteView-9*/}}<strong>[Entfernt]</strong>":                                          "<strong>[Removed]</strong>",
		"{{/*noteView-10*/}}<strong>{{.Title}}</strong> von {{.Author}}":                         "<strong>{{.Title}}</strong> by {{.Author}}",
		"{{/*noteView-11*/}}Kommentare":                                                          "Comments",
		"{{/*noteView-12*/}}<strong>[Entfernt]</strong>":                                         "<strong>[Removed]</strong>",
		"{{/*noteView-13*/}}<strong>{{.Title}}</strong> von {{.Author}}":                         "<strong>{{.Title}}</strong> by {{.Author}}",
	},

	"organisationCreate": {
		"{{/*organisationCreate-1*/}}Name":                                   "Name",
		"{{/*organisationCreate-2*/}}Hauptgruppe":                            "Main Group",
		"{{/*organisationCreate-3*/}}Untergruppe":                            "Sub Group",
		"{{/*organisationCreate-4*/}}Sichtbarkeit":                           "Visibility",
		"{{/*organisationCreate-5*/}}Öffentlich":                             "Public",
		"{{/*organisationCreate-6*/}}Privat":                                 "Private",
		"{{/*organisationCreate-7*/}}Geheim":                                 "Secret",
		"{{/*organisationCreate-8*/}}Versteckt":                              "Hidden",
		"{{/*organisationCreate-9*/}}Flair":                                  "Flair",
		"{{/*organisationCreate-10*/}}Organisationsmitglied hinzufügen":      "Add Organisation Member",
		"{{/*organisationCreate-11*/}}Organisationsmitglieder":               "Organisation Member",
		"{{/*organisationCreate-12*/}}Organisationsadministrator hinzufügen": "Add Organisation Administrator",
		"{{/*organisationCreate-13*/}}Organisationsadministratoren":          "Organisation Administrators",
		"{{/*organisationCreate-14*/}}Organisation erstellen":                "Create Organisation",
	},

	"organisationEdit": {
		"{{/*organisationEdit-1*/}}Zurück zur Suche":                       "Back to Search",
		"{{/*organisationEdit-2*/}}Bisheriger Name":                        "Current Name",
		"{{/*organisationEdit-3*/}}Name":                                   "Name",
		"{{/*organisationEdit-4*/}}Hauptgruppe":                            "Main Group",
		"{{/*organisationEdit-5*/}}Untergruppe":                            "Sub Group",
		"{{/*organisationEdit-6*/}}Sichtbarkeit":                           "Visibility",
		"{{/*organisationEdit-7*/}}Öffentlich":                             "Public",
		"{{/*organisationEdit-8*/}}Privat":                                 "Private",
		"{{/*organisationEdit-9*/}}Geheim":                                 "Secret",
		"{{/*organisationEdit-10*/}}Versteckt":                             "Hidden",
		"{{/*organisationEdit-11*/}}Flair":                                 "Flair",
		"{{/*organisationEdit-12*/}}Organisationsmitglied hinzufügen":      "Add Organisation Member",
		"{{/*organisationEdit-13*/}}Organisationsmitglieder":               "Organisation Member",
		"{{/*organisationEdit-14*/}}Organisationsadministrator hinzufügen": "Add Organisation Administrator",
		"{{/*organisationEdit-15*/}}Organisationsadministratoren":          "Organisation Administrators",
		"{{/*organisationEdit-16*/}}Organisation anpassen":                 "Update Organisation",
		"{{/*organisationEdit-17*/}}Organisationsname":                     "Organisation Name",
		"{{/*organisationEdit-18*/}}Organisation suchen":                   "Search Organisation",
	},

	"organisationView": {
		"{{/*organisationView-1*/}}Fehler beim Laden der Organisationen":      "Error While Trying to Load Organisations",
		"{{/*organisationView-2*/}}Es existieren keine Organisationen":        "No Organisations yet",
		"{{/*organisationView-3*/}}Flair: {{.Organisation.Flair}}":            "Flair: {{.Organisation.Flair}}",
		"{{/*organisationView-4*/}}Kein Flair":                                "No Flair",
		"{{/*organisationView-5*/}}Organisationsinformationen nicht gefunden": "Organisation Information could not be retrieved",
	},

	"personalProfil": {
		"{{/*personalProfil-1*/}}Persönliche Einstellungen":     "Personal Settings",
		"{{/*personalProfil-2*/}}Seitenskalierung (in Prozent)": "Page Scaling in Percent",
		"{{/*personalProfil-3*/}}Persönliche Zeitzone":          "Personal Timezone",
		"{{/*personalProfil-4*/}}Einstellungen speichern":       "Save Settings",
		"{{/*personalProfil-5*/}}Passwort ändern":               "Change Password",
		"{{/*personalProfil-6*/}}Altes Passwort":                "Old Password",
		"{{/*personalProfil-7*/}}Neues Passwort":                "New Password",
		"{{/*personalProfil-8*/}}Neues Passwort wiederholen":    "Repeat New Password",
		"{{/*personalProfil-9*/}}Passwort ändern":               "Change Password",
	},

	"titleCreate": {
		"{{/*titleCreate-1*/}}Name":                    "Name",
		"{{/*titleCreate-2*/}}Hauptgruppe":             "Main Group",
		"{{/*titleCreate-3*/}}Untergruppe":             "Sub Group",
		"{{/*titleCreate-4*/}}Flair":                   "Flair",
		"{{/*titleCreate-5*/}}Titel-Halter hinzufügen": "Add Title Owner",
		"{{/*titleCreate-6*/}}Titel-Halter":            "Title Owner",
		"{{/*titleCreate-7*/}}Titel erstellen":         "Create Title",
	},

	"titleEdit": {
		"{{/*titleEdit-1*/}}Zurück zur Suche":        "Back to Search",
		"{{/*titleEdit-2*/}}Bisheriger Name":         "Current Name",
		"{{/*titleEdit-3*/}}Name":                    "Name",
		"{{/*titleEdit-4*/}}Hauptgruppe":             "Main Group",
		"{{/*titleEdit-5*/}}Untergruppe":             "Sub Group",
		"{{/*titleEdit-6*/}}Flair":                   "Flair",
		"{{/*titleEdit-7*/}}Titel-Halter hinzufügen": "Add Title Owner",
		"{{/*titleEdit-8*/}}Titel-Halter":            "Title Owner",
		"{{/*titleEdit-9*/}}Titel anpassen":          "Update Title",
		"{{/*titleEdit-10*/}}Titelname":              "Title Name",
		"{{/*titleEdit-11*/}}Titel suchen":           "Search Title",
	},

	"titleView": {
		"{{/*titleView-1*/}}Es existieren keine Titel":         "No Titles yet",
		"{{/*titleView-2*/}}Flair: {{.Flair}}":                 "Flair: {{.Flair}}",
		"{{/*titleView-3*/}}Kein Flair":                        "No Flair",
		"{{/*titleView-4*/}}Titelinformationen nicht gefunden": "Title Information could not be retrieved",
	},

	"base": {
		"{{/*base-1-language*/}}": LanguageTag,
	},

	"templates": {
		"{{/*templates-1*/}}Home":                                                           "Home",
		"{{/*templates-2*/}}Notizen":                                                        "Notes",
		"{{/*templates-3*/}}Notiz erstellen":                                                "Create Note",
		"{{/*templates-4*/}}Zeitungen":                                                      "Newspapers",
		"{{/*templates-5*/}}Zeitungsartikel erstellen":                                      "Write Article",
		"{{/*templates-6*/}}Übersichten":                                                    "Overviews",
		"{{/*templates-7*/}}Titelübersicht":                                                 "Title Overview",
		"{{/*templates-8*/}}Organisationsübersicht":                                         "Organisation Overview",
		"{{/*templates-9*/}}Dokumente":                                                      "Documents",
		"{{/*templates-10*/}}Dokument erstellen":                                            "Create Document",
		"{{/*templates-11*/}}Diskussion erstellen":                                          "Create Discussion",
		"{{/*templates-12*/}}Abstimmung erstellen":                                          "Create Vote",
		"{{/*templates-13*/}}Abstimmungen verwalten":                                        "Manage Vote Preparation",
		"{{/*templates-14*/}}Tag-Farben verwalten":                                          "Manage Tag Colors",
		"{{/*templates-15*/}}Profil":                                                        "Profil",
		"{{/*templates-16*/}}Meine Briefe":                                                  "My Letters",
		"{{/*templates-17*/}}Brief schreiben":                                               "Write Letter",
		"{{/*templates-18*/}}Meine Dokumente":                                               "My Documents",
		"{{/*templates-19*/}}Administration":                                                "Administration",
		"{{/*templates-20*/}}Zeitung verwalten":                                             "Manage Newspapers",
		"{{/*templates-21*/}}Brief untersuchen":                                             "Search Letter",
		"{{/*templates-22*/}}Nutzer verwalten":                                              "Manage Accounts",
		"{{/*templates-23*/}}Organisation verwalten":                                        "Manage Organisations",
		"{{/*templates-24*/}}Titel verwalten":                                               "Manage Titles",
		"{{/*templates-25*/}}Nutzer erstellen":                                              "Create Account",
		"{{/*templates-26*/}}Organisation erstellen":                                        "Create Organisation",
		"{{/*templates-27*/}}Titel erstellen":                                               "Create Title",
		"{{/*templates-28*/}}-- Organisation auswählen --":                                  "-- Select Organisation --",
		"{{/*templates-29*/}}Leser hinzufügen":                                              "Add Reader",
		"{{/*templates-30*/}}Leser":                                                         "Reader",
		"{{/*templates-31*/}}Teilnehmer hinzufügen":                                         "Add Participant",
		"{{/*templates-32*/}}Teilnehmer":                                                    "Participants",
		"{{/*templates-33*/}}Es wurden keine gültigen Stimmen abgegeben":                    "No valid Votes have been casted",
		"{{/*templates-34*/}}Frage: {{.Question}}":                                          "Question: {{.Question}}",
		"{{/*templates-35*/}}Abstimmende Person":                                            "Voter",
		"{{/*templates-36*/}}{{if .Type.IsRankedVoting}}Rang{{else}}Stimme(n){{end}}":       "{{if .Type.IsRankedVoting}}Rank{{else}}Vote(s){{end}}",
		"{{/*templates-37*/}}{{if $anonym}}{{$voter}}. Wahlzettel{{else}}{{$voter}}{{end}}": "{{if $anonym}}Ballot {{$voter}}{{else}}{{$voter}}{{end}}",
		"{{/*templates-38*/}}Ungültige Stimmen: {{.GetIllegalVotes}}":                       "Invalid Votes: {{.GetIllegalVotes}}",
		"{{/*templates-39*/}}CSV Herunterladen":                                             "Download CSV",
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
