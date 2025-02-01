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
)

var replaceMap = map[string]string{}

func LocaliseTemplateString(input []byte) string {
	result := string(input)
	for key, value := range replaceMap {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}
