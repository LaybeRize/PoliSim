//go:build EN

package loc

import (
	"html/template"
	"strings"
)

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

	// language=HTML
	homePageElement = ``
)

var replaceMap = map[string]string{
	"$$home-page$$": homePageElement,
	"{{/*_home*/}}Herzlich willkommen, {{.Account.Name}}": "Welcome, {{.Account.Name}}",
	"{{/*_home*/}}Abmelden":                               "Sign out",
	"{{/*_home*/}}Nutzername":                             "Username",
	"{{/*_home*/}}Passwort":                               "Password",
	"{{/*_home*/}}Einloggen":                              "Sign in",

	"{{/*_notFound*/}}Die Seite, die Sie suchen, existiert nicht": "The Page you requested does not exist",

	"{{/*accountCreate*/}}Anzeigename":          "Display Name",
	"{{/*accountCreate*/}}Nutzername":           "Username",
	"{{/*accountCreate*/}}Passwort":             "Password",
	"{{/*accountCreate*/}}Rolle":                "Role",
	"{{/*accountCreate*/}}Nutzer":               "User",
	"{{/*accountCreate*/}}Presse-Nutzer":        "Press User",
	"{{/*accountCreate*/}}Presse-Administrator": "Press Administrator",
	"{{/*accountCreate*/}}Administrator":        "Administrator",
	"{{/*accountCreate*/}}Oberadministrator":    "Head Administrator",
	"{{/*accountCreate*/}}Nutzer erstellen":     "Create User",

	"{{/*accountEdit*/}}Zurück zur Suche":     "Back to Search",
	"{{/*accountEdit*/}}Anzeigename":          "Display Name",
	"{{/*accountEdit*/}}Nutzername":           "Username",
	"{{/*accountEdit*/}}Rolle":                "Role",
	"{{/*accountEdit*/}}Nutzer":               "User",
	"{{/*accountEdit*/}}Presse-Nutzer":        "Press User",
	"{{/*accountEdit*/}}Presse-Administrator": "Press Administrator",
	"{{/*accountEdit*/}}Administrator":        "Administrator",
	"{{/*accountEdit*/}}Oberadministrator":    "Head Administrator",
	"{{/*accountEdit*/}}Blockiert":            "Blocked",
	"{{/*accountEdit*/}}Account-Besitzer":     "Account Owner",
	"{{/*accountEdit*/}}Nutzer anpassen":      "Update User",
	"{{/*accountEdit*/}}Nutzer suchen":        "Search for User",

	"{{/*documentColorEdit*/}}Farbpaletten":                   "Color Palette",
	"{{/*documentColorEdit*/}}Farbpalette auswählen":          "Select Color Palette",
	"{{/*documentColorEdit*/}}Name":                           "Name",
	"{{/*documentColorEdit*/}}Hintergrundfarbe":               "Background Color",
	"{{/*documentColorEdit*/}}Textfarbe":                      "Text Color",
	"{{/*documentColorEdit*/}}Link-Farbe":                     "Link Color",
	"{{/*documentColorEdit*/}}Farbpalette erstellen/anpassen": "Create/Update Color Palette",
	"{{/*documentColorEdit*/}}Farbpalette löschen":            "Delete Color Palette",

	"{{/*documentCreate*/}}Titel":              "Title",
	"{{/*documentCreate*/}}Autor":              "Author",
	"{{/*documentCreate*/}}Organisation":       "Organisation",
	"{{/*documentCreate*/}}Inhalt":             "Content",
	"{{/*documentCreate*/}}Dokument erstellen": "Create Document",
	"{{/*documentCreate*/}}Vorschau":           "Preview",

	"{{/*documentCreateDiscussion*/}}Titel":                                                              "Title",
	"{{/*documentCreateDiscussion*/}}Autor":                                                              "Author",
	"{{/*documentCreateDiscussion*/}}Organisation":                                                       "Organisation",
	"{{/*documentCreateDiscussion*/}}Ende der Diskussion ({{.NavInfo.Account.TimeZone.String}})":         "End Timestamp for Discussion ({{.NavInfo.Account.TimeZone.String}})",
	"{{/*documentCreateDiscussion*/}}Diskussion ist öffentlich (Pflicht in öffentlichen Organisationen)": "Discussion is Public (Must be checked when posting in a Public Organisation)",
	"{{/*documentCreateDiscussion*/}}Alle Organisationsmitglieder dürfen teilnehmen":                     "All Organisation Member are allowed to Participate",
	"{{/*documentCreateDiscussion*/}}Alle Organisationsadministratoren dürfen teilnehmen":                "All Organisation Administrators are allowed to Participate",
	"{{/*documentCreateDiscussion*/}}Inhalt":                                                             "Content",
	"{{/*documentCreateDiscussion*/}}Leser und Teilnehmer überprüfen":                                    "Check Reader and Participants",
	"{{/*documentCreateDiscussion*/}}Diskussion erstellen":                                               "Create Discussion",
	"{{/*documentCreateDiscussion*/}}Vorschau":                                                           "Preview",

	"{{/*documentCreateVote*/}}Titel":        "Title",
	"{{/*documentCreateVote*/}}Autor":        "Author",
	"{{/*documentCreateVote*/}}Organisation": "Organisation",
	"{{/*documentCreateVote*/}}Ende der Abstimmung (Endet immer um 23:50 UTC des ausgewählten Tages)": "End Timestamp for Vote (Always ends at 23:50 UTC of the selected day)",
	"{{/*documentCreateVote*/}}Abstimmung ist öffentlich (Pflicht in öffentlichen Organisationen)":    "Vote is Public (Must be checked when posting in a Public Organisation)",
	"{{/*documentCreateVote*/}}Alle Organisationsmitglieder dürfen teilnehmen":                        "All Organisation Member are allowed to Participate",
	"{{/*documentCreateVote*/}}Alle Organisationsadministratoren dürfen teilnehmen":                   "All Organisation Administrators are allowed to Participate",
	"{{/*documentCreateVote*/}}Abstimmungsliste":                                                      "Vote List",
	"{{/*documentCreateVote*/}}ID der ausgewählten Abstimmung übertragen":                             "Insert ID of selected Vote",
	"{{/*documentCreateVote*/}}Angehängte Abstimmungen":                                               "Attached Votes",
	"{{/*documentCreateVote*/}}Inhalt":                                                                "Content",
	"{{/*documentCreateVote*/}}Leser und Teilnehmer überprüfen":                                       "Check Reader and Participants",
	"{{/*documentCreateVote*/}}Abstimmungsdokument erstellen":                                         "Create Vote Document",
	"{{/*documentCreateVote*/}}Vorschau":                                                              "Preview",

	"{{/*documentCreateVoteElement*/}}Abstimmungsnummer":                                              "Vote Number",
	"{{/*documentCreateVoteElement*/}}Abstimmungs-ID":                                                 "Vote ID",
	"{{/*documentCreateVoteElement*/}}Abstimmungsart":                                                 "Vote Type",
	"{{/*documentCreateVoteElement*/}}Eine Stimme pro Nutzer":                                         "Single Choice Vote",
	"{{/*documentCreateVoteElement*/}}Mehrere Stimmen pro Nutzer":                                     "Multiple Choice Vote",
	"{{/*documentCreateVoteElement*/}}Rangwahl":                                                       "Option Ranking",
	"{{/*documentCreateVoteElement*/}}Gewichtete Wahl":                                                "Weighted Vote",
	"{{/*documentCreateVoteElement*/}}Maximale Stimmen pro Nutzer (Nur relevant für Gewichtete Wahl)": "Maximum Amount of Votes per User (only relevant for Weighted Vote)",
	"{{/*documentCreateVoteElement*/}}Zeige Teilnehmerbezogene Stimmen während der Wahl":              "Show current Ballots during Vote Period",
	"{{/*documentCreateVoteElement*/}}Geheime Wahl":                                                   "Secret Ballot",
	"{{/*documentCreateVoteElement*/}}Frage":                                                          "Question",
	"{{/*documentCreateVoteElement*/}}Antwort hinzufügen":                                             "Add New Option",
	"{{/*documentCreateVoteElement*/}}Antworten":                                                      "Options",
	"{{/*documentCreateVoteElement*/}}Abstimmung erstellen/bearbeiten":                                "Create/Update Vote",

	"{{/*documentSearch*/}}\"Die Anfrage hat zu einem Fehler auf der Serverseite geführt\"":                 "\"Requested could not be processed. Internal Server Error\"",
	"{{/*documentSearch*/}}Blockierte Dokumente anzeigen":                                                   "Show Blocked Documents",
	"{{/*documentSearch*/}}Anzahl der Ergebnisse":                                                           "Number of Entries per Page",
	"{{/*documentSearch*/}}Suchen":                                                                          "Search",
	"{{/*documentSearch*/}}Es konnten keine Einträge gefunden werden, die den Suchkriterien gerecht werden": "No Entries found",
	"{{/*documentSearch*/}}<strong>{{if .Removed}}[Entfernt]{{else}}{{.Title}}{{end}}</strong>":             "<strong>{{if .Removed}}[Removed]{{else}}{{.Title}}{{end}}</strong>",
	"{{/*documentSearch*/}}<i>Veröffentlicht am: {{.GetTimeWritten $acc}}</i>":                              "<i>Publish Date: {{.GetTimeWritten $acc}}</i>",
	"{{/*documentSearch*/}}Veröffentlicht von <i>{{.Author}}</i> im <i>{{.Organisation}}</i>":               "Written by <i>{{.Author}}</i> for <i>{{.Organisation}}</i>",
	"{{/*documentSearch*/}}&laquo; Vorherige Seite":                                                         "&laquo; Previous Page",
	"{{/*documentSearch*/}}Nächste Seite &raquo;":                                                           "Next Page &raquo;",

	"{{/*documentView*/}}Das Dokument wurde entfernt":                                                              "Document has been removed",
	"{{/*documentView*/}}Hinzugefügt am {{$tag.GetTimeWritten $acc}}":                                              "Added {{$tag.GetTimeWritten $acc}}",
	"{{/*documentView*/}}Geschrieben von: {{.Document.GetAuthor}}":                                                 "Written by: {{.Document.GetAuthor}}",
	"{{/*documentView*/}}Veröffentlicht in: {{.Document.Organisation}}":                                            "Published in: {{.Document.Organisation}}",
	"{{/*documentView*/}}Verfasst am: {{.Document.GetTimeWritten .NavInfo.Account}}":                               "Publish Date: {{.Document.GetTimeWritten .NavInfo.Account}}",
	"{{/*documentView*/}}Die {{if .Document.IsDiscussion}}Diskussion{{else}}Abstimmung{{end}} ist bereits vorbei.": "The {{if .Document.IsDiscussion}}Discussion{{else}}Vote{{end}} has ended.",
	"{{/*documentView*/}}Ende war am {{.Document.GetTimeEnd .NavInfo.Account}}":                                    "It ended {{.Document.GetTimeEnd .NavInfo.Account}}",
	"{{/*documentView*/}}Endet: {{.Document.GetTimeEnd .NavInfo.Account}}":                                         "Ends: {{.Document.GetTimeEnd .NavInfo.Account}}",
	"{{/*documentView*/}}Tag erstellen":                                                                            "Create Document-Tag",
	"{{/*documentView*/}}Willst du den Tag so hinzufügen?":                                                         "Confirm Tag Creation",
	"{{/*documentView*/}}Farbpaletten":                                                                             "Color Palette",
	"{{/*documentView*/}}Farbe aus Farbpalette kopieren":                                                           "Copy Colors from Selected Color Palette",
	"{{/*documentView*/}}Inhalt":                                                                                   "Content",
	"{{/*documentView*/}}Hintergrundfarbe":                                                                         "Background Color",
	"{{/*documentView*/}}Textfarbe":                                                                                "Text Color",
	"{{/*documentView*/}}Link-Farbe":                                                                               "Link Color",
	"{{/*documentView*/}}Referenzen zu anderen Dokumenten":                                                         "Document References",
	"{{/*documentView*/}}Tag hinzufügen":                                                                           "Add Tag to Document",
	"{{/*documentView*/}}{{if .Document.Removed}}Dokument wieder freigeben{{else}}Dokument blockieren{{end}}":      "{{if .Document.Removed}}Restore Document{{else}}Remove Document{{end}}",
	"{{/*documentView*/}}Kommentare":                                                                               "Comments",
	"{{/*documentView*/}}Geschrieben von: {{$comment.GetAuthor}}":                                                  "Written by: {{$comment.GetAuthor}}",
	"{{/*documentView*/}}Verfasst am: {{$comment.GetTimeWritten $acc}}":                                            "Publish Date: {{$comment.GetTimeWritten $acc}}",
	"{{/*documentView*/}}{{if $comment.Removed}}Kommentar wieder freigeben{{else}}Kommentar blockieren{{end}}":     "{{if $comment.Removed}}Restore Comment{{else}}Remove Comment{{end}}",
	"{{/*documentView*/}}Autor":                                                                                    "Author",
	"{{/*documentView*/}}Kommentar schreiben":                                                                      "Write Comment",
	"{{/*documentView*/}}Vorschau":                                                                                 "Preview",
	"{{/*documentView*/}}Abstimmungen":                                                                             "Votes",

	"{{/*documentViewVote*/}}Frage: {{.VoteInstance.Question}}":                                                                  "Question: {{.VoteInstance.Question}}",
	"{{/*documentViewVote*/}}Abstimmender Account":                                                                               "Voter",
	"{{/*documentViewVote*/}}{{$pos}}. Antwort: {{$answer}}":                                                                     "Option {{$pos}}: {{$answer}}",
	"{{/*documentViewVote*/}}{{$pos}}. Antwort: {{$answer}}":                                                                     "Option {{$pos}}: {{$answer}}",
	"{{/*documentViewVote*/}}Es dürfen maximal {{.VoteInstance.MaxVotes}} Stimmen vergeben werden":                               "You can cast a maximum of {{.VoteInstance.MaxVotes}} Votes",
	"{{/*documentViewVote*/}}Stimmen für die {{$pos}}. Antwort: {{$answer}}":                                                     "Votes for Option {{$pos}}: {{$answer}}",
	"{{/*documentViewVote*/}}Position 0 und kleiner bedeutet, dass die Antwort keinen Rang erhält. Der 1. Rang ist der höchste.": "Position 0 and smaller mean the Option is assigned no rank. The 1. rank is the highest.",
	"{{/*documentViewVote*/}}Rang der {{$pos}}. Antwort: {{$answer}}":                                                            "Rank of Option {{$pos}}: {{$answer}}",
	"{{/*documentViewVote*/}}Ungültige Stimme abgeben":                                                                           "Cast Invalid Vote",
	"{{/*documentViewVote*/}}Stimme abgeben":                                                                                     "Cast Vote",
	"{{/*documentViewVote*/}}Teilnahme setzte einen Account voraus":                                                              "Participation in the Vote requires an Account",
	"{{/*documentViewVote*/}}Das Ergebnis ist erst nach Ende der Abstimmung einsehbar":                                           "Results are available after the Vote has ended",

	"{{/*letterAdminSearch*/}}Brief ID":     "Letter ID",
	"{{/*letterAdminSearch*/}}Brief öffnen": "Open Letter",

	"{{/*letterCreate*/}}Titel":                "Title",
	"{{/*letterCreate*/}}Autor":                "Author",
	"{{/*letterCreate*/}}Empfänger hinzufügen": "Add Recipient",
	"{{/*letterCreate*/}}Empfänger":            "Recipients",
	"{{/*letterCreate*/}}Mit Unterschrift":     "With Signatures",
	"{{/*letterCreate*/}}Inhalt":               "Content",
	"{{/*letterCreate*/}}Brief überprüfen":     "Check Letter",
	"{{/*letterCreate*/}}Brief erstellen":      "Create Letter",
	"{{/*letterCreate*/}}Vorschau":             "Preview",

	"{{/*letterPersonalSearch*/}}\"Die Anfrage hat zu einem Fehler auf der Serverseite geführt\"": "\"Requested could not be processed. Internal Server Error\"",
	"{{/*letterPersonalSearch*/}}Account":               "Account",
	"{{/*letterPersonalSearch*/}}-- Alle Accounts --":   "-- All Accounts --",
	"{{/*letterPersonalSearch*/}}Anzahl der Ergebnisse": "Number of Entries per Page",
	"{{/*letterPersonalSearch*/}}Suchen":                "Search",
	"{{/*letterPersonalSearch*/}}Es konnten keine Einträge gefunden werden, die den Suchkriterien gerecht werden": "No Entries found",
	"{{/*letterPersonalSearch*/}}<strong>{{.Title}}</strong> von {{.Author}}":                                     "<strong>{{.Title}}</strong> by {{.Author}}",
	"{{/*letterPersonalSearch*/}}Empfänger: {{.Recipient}}":                                                       "Recipient: {{.Recipient}}",
	"{{/*letterPersonalSearch*/}}<i>Versendet am: {{.GetTimeWritten $acc}}":                                       "<i>Publish Date: {{.GetTimeWritten $acc}}",
	"{{/*letterPersonalSearch*/}}&laquo; Vorherige Seite":                                                         "&laquo; Previous Page",
	"{{/*letterPersonalSearch*/}}Nächste Seite &raquo;":                                                           "Next Page &raquo;",

	"{{/*letterView*/}}Geschrieben von: {{.Letter.GetAuthor}}":                          "Written by: {{.Letter.GetAuthor}}",
	"{{/*letterView*/}}<i>Verfasst am: {{.Letter.GetTimeWritten .NavInfo.Account}}</i>": "<i>Publish Date: {{.Letter.GetTimeWritten .NavInfo.Account}}</i>",
	"{{/*letterView*/}}Als {{.Letter.Recipient}} zustimmen":                             "Accept as {{.Letter.Recipient}}",
	"{{/*letterView*/}}Als {{.Letter.Recipient}} ablehnen":                              "Decline as {{.Letter.Recipient}}",

	"{{/*templates*/}}Home":                                                          "Home",
	"{{/*templates*/}}Notizen":                                                       "Notes",
	"{{/*templates*/}}Notiz erstellen":                                               "Create Note",
	"{{/*templates*/}}Zeitungen":                                                     "Newspapers",
	"{{/*templates*/}}Zeitungsartikel erstellen":                                     "Write Article",
	"{{/*templates*/}}Übersichten":                                                   "Overviews",
	"{{/*templates*/}}Titelübersicht":                                                "Title Overview",
	"{{/*templates*/}}Organisationsübersicht":                                        "Organisation Overview",
	"{{/*templates*/}}Dokumente":                                                     "Documents",
	"{{/*templates*/}}Dokument erstellen":                                            "Create Document",
	"{{/*templates*/}}Diskussion erstellen":                                          "Create Discussion",
	"{{/*templates*/}}Abstimmung erstellen":                                          "Create Vote",
	"{{/*templates*/}}Abstimmungen verwalten":                                        "Manage Vote Preparation",
	"{{/*templates*/}}Tag-Farben verwalten":                                          "Manage Tag Colors",
	"{{/*templates*/}}Profil":                                                        "Profil",
	"{{/*templates*/}}Meine Briefe":                                                  "My Letters",
	"{{/*templates*/}}Brief schreiben":                                               "Write Letter",
	"{{/*templates*/}}Meine Dokumente":                                               "My Documents",
	"{{/*templates*/}}Administration":                                                "Administration",
	"{{/*templates*/}}Zeitung verwalten":                                             "Manage Newspapers",
	"{{/*templates*/}}Brief untersuchen":                                             "Search Letter",
	"{{/*templates*/}}Nutzer verwalten":                                              "Manage Accounts",
	"{{/*templates*/}}Organisation verwalten":                                        "Manage Organisations",
	"{{/*templates*/}}Titel verwalten":                                               "Manage Titles",
	"{{/*templates*/}}Nutzer erstellen":                                              "Create Account",
	"{{/*templates*/}}Organisation erstellen":                                        "Create Organisation",
	"{{/*templates*/}}Titel erstellen":                                               "Create Title",
	"{{/*templates*/}}-- Organisation auswählen --":                                  "-- Select Organisation --",
	"{{/*templates*/}}Leser hinzufügen":                                              "Add Reader",
	"{{/*templates*/}}Teilnehmer hinzufügen":                                         "Add Participant",
	"{{/*templates*/}}Es wurden keine gültigen Stimmen abgegeben":                    "No valid Votes have been casted",
	"{{/*templates*/}}Frage: {{.Question}}":                                          "Question: {{.Question}}",
	"{{/*templates*/}}Abstimmende Person":                                            "Voter",
	"{{/*templates*/}}{{if .Type.IsRankedVoting}}Rang{{else}}Stimme(n){{end}}":       "{{if .Type.IsRankedVoting}}Rank{{else}}Vote(s){{end}}",
	"{{/*templates*/}}{{if $anonym}}{{$voter}}. Wahlzettel{{else}}{{$voter}}{{end}}": "{{if $anonym}}Ballot {{$voter}}{{else}}{{$voter}}{{end}}",
	"{{/*templates*/}}Ungültige Stimmen: {{.GetIllegalVotes}}":                       "Invalid Votes: {{.GetIllegalVotes}}",
}

func LocaliseTemplateString(input []byte) string {
	result := string(input)
	for key, value := range replaceMap {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}
