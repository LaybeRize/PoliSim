package database

import (
	loc "PoliSim/localisation"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"html/template"
	"strings"
	"time"
)

type DocumentType int

const (
	DocTypePost DocumentType = iota
	DocTypeDiscussion
	DocTypeVote
)

type (
	Document struct {
		ID                  string            `json:"-"`
		Type                DocumentType      `json:"-"`
		Organisation        string            `json:"-"`
		Title               string            `json:"-"`
		Author              string            `json:"-"`
		Flair               string            `json:"-"`
		Written             time.Time         `json:"-"`
		Body                template.HTML     `json:"-"`
		Public              bool              `json:"-"`
		Removed             bool              `json:"-"`
		MemberParticipation bool              `json:"-"`
		AdminParticipation  bool              `json:"-"`
		AllowedToAddTags    bool              `json:"-"`
		End                 time.Time         `json:"-"`
		Reader              []string          `json:"reader"`
		Participants        []string          `json:"participants"`
		Tags                []DocumentTag     `json:"tags"`
		Links               []VoteInfo        `json:"links"`
		VoteIDs             []string          `json:"-"`
		Comments            []DocumentComment `json:"-"`
		Result              []AccountVotes    `json:"result"`
	}
	SmallDocument struct {
		ID           string
		Type         DocumentType
		Organisation string
		Title        string
		Author       string
		Written      time.Time
		Removed      bool
	}
	DocumentTag struct {
		ID              string    `json:"id"`
		Outgoing        bool      `json:"outgoing"`
		Text            string    `json:"text"`
		Written         time.Time `json:"written"`
		BackgroundColor string    `json:"background_color"`
		TextColor       string    `json:"text_color"`
		LinkColor       string    `json:"link_color"`
		Links           []string  `json:"links"`
	}
	DocumentComment struct {
		ID      string
		Author  string
		Flair   string
		Written time.Time
		Body    template.HTML
		Removed bool
	}
)

// Todo move this to migration
const documentTableCreateStatement = `
CREATE TABLE document (
    id TEXT PRIMARY KEY,
    type INT NOT NULL,
    organisation TEXT NOT NULL,
    organisation_name TEXT NOT NULL,
    title TEXT NOT NULL,
    author  TEXT NOT NULL,
    flair  TEXT NOT NULL,
    body  TEXT NOT NULL,
    written TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    public  BOOLEAN NOT NULL,
    removed BOOLEAN NOT NULL,
    member_participation BOOLEAN NOT NULL,
    admin_participation BOOLEAN NOT NULL,
    extra_info jsonb NOT NULL,
    CONSTRAINT fk_organisation_name
        FOREIGN KEY (organisation_name) REFERENCES organisation(name) ON UPDATE CASCADE
);
CREATE TABLE document_to_account (
    document_id TEXT NOT NULL,
    account_name TEXT NOT NULL,
    participant BOOLEAN NOT NULL,
    CONSTRAINT fk_document_id
        FOREIGN KEY (document_id) REFERENCES document(id),
    CONSTRAINT fk_account_name
        FOREIGN KEY (account_name) REFERENCES account(name)
);
CREATE TABLE comment_to_document (
	comment_id TEXT PRIMARY KEY,
	document_id TEXT NOT NULL,
    author  TEXT NOT NULL,
    flair  TEXT NOT NULL,
    body  TEXT NOT NULL,
    written TIMESTAMP NOT NULL,
    removed BOOLEAN NOT NULL,
    CONSTRAINT fk_document_id
        FOREIGN KEY (document_id) REFERENCES document(id)	
);
CREATE VIEW document_linked AS
SELECT id, type, organisation, doc.organisation_name, title, author, flair, body, written, 
       end_time, public, removed, member_participation, admin_participation, extra_info, 
       dta.account_name as doc_account, ota.account_name as organisation_account, is_admin, 
       dta.participant as participant, owner_name FROM document doc
	LEFT JOIN document_to_account dta ON doc.id = dta.document_id
	LEFT JOIN organisation_to_account ota ON doc.organisation_name = ota.organisation_name
	LEFT JOIN ownership own ON ota.account_name = own.account_name OR dta.account_name = own.account_name;
`

func (s *SmallDocument) IsPost() bool { return s.Type == DocTypePost }

func (s *SmallDocument) IsDiscussion() bool { return s.Type == DocTypeDiscussion }

func (s *SmallDocument) IsVote() bool { return s.Type == DocTypeVote }

func (s *SmallDocument) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return s.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return s.Written.Format(loc.TimeFormatString)
}

func (d *Document) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Document) Scan(src interface{}) error {
	switch src.(type) {
	case []byte:
		return json.Unmarshal(src.([]byte), d)
	case string:
		return json.Unmarshal([]byte(src.(string)), d)
	default:
		return errors.New("value can not be unmarshalled into document")
	}
}

func (d *Document) ShowRemovedMessage(acc *Account) bool {
	if acc.IsAtLeastAdmin() || !d.Removed {
		return false
	}
	return true
}

func (d *Document) HasResults() bool {
	return d.Result != nil
}

func (d *Document) HasComments() bool { return len(d.Comments) != 0 }

func (d *Document) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return d.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return d.Written.Format(loc.TimeFormatString)
}

func (d *Document) GetTimeEnd(a *Account) string {
	if a.Exists() {
		return d.End.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return d.End.Format(loc.TimeFormatString)
}

func (d *Document) GetAuthor() string {
	if d.Flair == "" {
		return d.Author
	}
	return d.Author + "; " + d.Flair
}

func (d *Document) GetReader() string {
	if d.Public {
		return loc.DocumentIsPublic
	} else if len(d.Reader) == 0 {
		return loc.DocumentOnlyForMember
	}
	return fmt.Sprintf(loc.DocumentFormatStringForReader, strings.Join(d.Reader, ", "))
}

func (d *Document) GetParticipants() string {
	if d.Type == DocTypePost {
		return loc.DocumentTagAddInfo
	} else if d.MemberParticipation && len(d.Participants) == 0 {
		return loc.DocumentParticipationEveryMember
	} else if d.MemberParticipation {
		return fmt.Sprintf(loc.DocumentParticipationEveryMemberPlus,
			strings.Join(d.Reader, ", "))
	} else if d.AdminParticipation && len(d.Participants) == 0 {
		return loc.DocumentParticipationOnlyAdmins
	} else if d.AdminParticipation {
		return fmt.Sprintf(loc.DocumentParticipationOnlyAdminsPlus, strings.Join(d.Reader, ", "))
	}
	return fmt.Sprintf(loc.DocumentParticipationFormatString, strings.Join(d.Reader, ", "))
}

func (d *Document) IsPost() bool { return d.Type == DocTypePost }

func (d *Document) IsDiscussion() bool { return d.Type == DocTypeDiscussion }

func (d *Document) IsVote() bool { return d.Type == DocTypeVote }

func (d *Document) Ended() bool { return time.Now().After(d.End) }

func (t *DocumentTag) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *DocumentTag) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return t.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return t.Written.Format(loc.TimeFormatString)
}

func (t *DocumentTag) HasLinks() bool {
	return len(t.Links) != 0
}

func (c *DocumentComment) GetBody(acc *Account) template.HTML {
	if !c.Removed || acc.IsAtLeastAdmin() {
		return c.Body
	}
	return loc.DocumentCommentContentRemovedHTML
}

func (c *DocumentComment) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return c.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return c.Written.Format(loc.TimeFormatString)
}

func (c *DocumentComment) GetAuthor() string {
	if c.Flair == "" {
		return c.Author
	}
	return c.Author + "; " + c.Flair
}

func CreateDocument(document *Document) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	defer rollback(tx)
	var vis OrganisationVisibility
	err = tx.QueryRow(`SELECT visibility FROM organisation_linked WHERE visibility <> $1 AND name = $2 AND account_name = $3 AND is_admin = true;`,
		HIDDEN, document.Organisation, document.Author).Scan(&vis)
	if errors.Is(err, sql.ErrNoRows) {
		return NotAllowedError
	} else if err != nil {
		return err
	} else if (vis == SECRET && document.Public) || (vis == PUBLIC && !document.Public) {
		return DocumentHasInvalidVisibility
	}

	if document.Type != DocTypePost {
		err = tx.QueryRow(`SELECT ARRAY(SELECT name FROM account WHERE name = ANY($1) AND blocked = false), 
       ARRAY(SELECT name FROM account WHERE name = ANY($2) AND blocked = false);`,
			pq.Array(document.Reader), pq.Array(document.Participants)).
			Scan(pq.Array(&document.Reader), pq.Array(&document.Participants))
		if err != nil {
			return err
		}
	}

	document.Tags = make([]DocumentTag, 0)

	_, err = tx.Exec(`INSERT INTO document (id, type, organisation, organisation_name, title, author, flair, body, written, 
                      end_time, public, removed, member_participation, admin_participation, extra_info) 
VALUES ($1, $2, $3, $3, $4, $5, $6, $7, $8, $9, false, $10, $11, $12, $13);`,
		document.ID, document.Type, document.Organisation, document.Title, document.Author, document.Flair,
		document.Body, time.Now().UTC(), document.End, document.Public, document.MemberParticipation,
		document.AdminParticipation, &document)
	if err != nil {
		return err
	}

	if !document.Public {
		_, err = postgresDB.Exec(`INSERT INTO document_to_account (document_id, account_name, participant) 
SELECT $1 AS document_id, name, false AS participant FROM account
WHERE name = ANY($2);`, &document.ID, pq.Array(document.Reader))
		if err != nil {
			return err
		}
	}

	if document.Type != DocTypePost {
		_, err = postgresDB.Exec(`INSERT INTO document_to_account (document_id, account_name, participant) 
SELECT $1 AS document_id, name, true AS participant FROM account
WHERE name = ANY($2);`, &document.ID, pq.Array(document.Participants))
		if err != nil {
			return err
		}
	}

	if document.Type == DocTypeVote {
		//Todo make this after reworking the vote system
		// Idea: Two separate table and while copying return the value, check that it is at least one via Next() and if not call Close() on the value and
		// return the correct DocumentHasNoAttachedVotes error
		/*result, err = tx.Run(`
		MATCH (a:Account)-[r:MANAGES]->(v:Vote) WHERE a.name = $user AND v.id IN $ids
		MATCH (d:Document) WHERE d.id = $id
		DELETE r
		MERGE (d)-[:LINKS]->(v)
		RETURN v.id;`, map[string]any{
					"user": acc.Name,
					"ids":  document.VoteIDs,
					"id":   document.ID})
				if err != nil {
					return err
				} else if !result.Peek() {
					return DocumentHasNoAttachedVotes
				}*/
		_, err = tx.Exec(`UPDATE document SET extra_info = $2 WHERE id = $1`, document.ID, &document)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetDocumentForUser(id string, acc *Account) (*Document, []string, error) {
	doc := &Document{}
	err := postgresDB.QueryRow(`SELECT id, type, organisation, title, author, flair, 
       written, body, public, removed, member_participation, admin_participation, end_time, extra_info, is_admin FROM document_linked 
          WHERE id = $1 AND (public = true OR owner_name = $2 OR $3 = true) ORDER BY is_admin DESC NULLS LAST LIMIT 1;`,
		id, acc.Name, acc.IsAtLeastAdmin()).Scan(
		&doc.ID, &doc.Type, &doc.Organisation, &doc.Type, &doc.Author, &doc.Flair, &doc.Written, &doc.Body,
		&doc.Public, &doc.Removed, &doc.MemberParticipation, &doc.AdminParticipation, &doc.End, doc, &doc.AllowedToAddTags)
	doc.AllowedToAddTags = doc.AllowedToAddTags || acc.IsAtLeastAdmin()

	commentator := make([]string, 0)

	if doc.Type == DocTypeDiscussion {

		if !doc.Ended() && acc.Exists() {
			commentator, err = GetMyAccountNames(acc)
			if err != nil {
				return nil, nil, err
			}
		}

		var result *sql.Rows
		result, err = postgresDB.Query(`SELECT comment_id, author, flair, written, body, removed 
FROM comment_to_document WHERE document_id = $1 ORDER BY written;`, &doc.ID)
		if err != nil {
			return nil, nil, err
		}
		defer closeRows(result)

		doc.Comments = make([]DocumentComment, 0)
		comment := DocumentComment{}
		for result.Next() {
			err = result.Scan(&comment.ID, &comment.Author, &comment.Flair, &comment.Written, &comment.Body, &comment.Removed)
			if err != nil {
				return nil, nil, err
			}
			doc.Comments = append(doc.Comments, comment)
		}
	}

	return doc, commentator, err
}

func RemoveRestoreDocument(docId string) {
	_, _ = postgresDB.Exec(`UPDATE document SET removed = NOT removed WHERE id = $1`, docId)
}

func CreateTagForDocument(docID string, acc *Account, tag *DocumentTag) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}
	defer rollback(tx)

	err = tx.QueryRow(`SELECT id FROM document_linked WHERE id = $1 AND is_admin = true AND owner_name = $2 LIMIT 1;`,
		docID, acc.GetName()).Scan(&docID)
	if errors.Is(err, sql.ErrNoRows) {
		return NotAllowedError
	} else if err != nil {
		return err
	}
	var docIDs []string
	err = tx.QueryRow(`SELECT ARRAY(SELECT id FROM document WHERE id = ANY($1) AND id <> $2)`,
		pq.Array(tag.Links), docID).Scan(pq.Array(&docIDs))

	tag.Outgoing = true
	tag.Links = docIDs

	_, err = tx.Exec(`UPDATE document SET extra_info = jsonb_insert(extra_info, '{tags,0}', $2) WHERE id = $1`, docID, tag)
	if err != nil {
		return err
	}

	tag.Outgoing = false
	tag.Links = []string{docID}
	_, err = tx.Exec(`UPDATE document SET extra_info = jsonb_insert(extra_info, '{tags,0}', $2) WHERE id = ANY($1)`, pq.Array(docIDs), tag)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CreateDocumentComment(documentId string, comment *DocumentComment) error {
	err := postgresDB.QueryRow(`SELECT id FROM document_linked WHERE id = $1 AND end_time > $2 AND
                                     ((doc_account = $3 AND participant = true) OR 
                                      (organisation_account = $3 AND member_participation = true) OR 
                                      (organisation_account = $3 AND is_admin = true AND admin_participation = true)) LIMIT 1;`,
		documentId, time.Now().UTC(), comment.Author).Scan(&documentId)
	if errors.Is(err, sql.ErrNoRows) {
		return NotAllowedError
	} else if err != nil {
		return err
	}

	_, err = postgresDB.Exec(`INSERT INTO comment_to_document (comment_id, document_id, author, flair, body, written, removed) 
VALUES ($1, $2, $3, $4, $5, $6, false)`, comment.ID, documentId, comment.Author, comment.Flair, comment.Body, time.Now().UTC())
	return err
}

func RemoveRestoreComment(commentID string) {
	_, _ = postgresDB.Exec(`UPDATE comment_to_document SET removed = NOT removed WHERE comment_id = $1`, commentID)
}

func GetDocumentList(amount int, page int, acc *Account, showBlocked bool) ([]SmallDocument, error) {
	var query string
	parameter := []any{(page - 1) * amount, amount}
	if !acc.Exists() {
		query = `SELECT id, type, organisation, title, author, written, removed FROM document WHERE public = true AND removed = false`
	} else if !acc.IsAtLeastAdmin() {
		query = `SELECT DISTINCT ON (id) id, type, organisation, title, author, written, removed FROM document_linked WHERE (public = true OR owner_name = $3) AND removed = false`
		parameter = append(parameter, acc.Name)
	} else {
		if showBlocked {
			query = `SELECT id, type, organisation, title, author, written, removed FROM document`
		} else {
			query = `SELECT id, type, organisation, title, author, written, removed FROM document WHERE removed = false`
		}
	}

	result, err := postgresDB.Query(query+` ORDER BY written DESC OFFSET $1 LIMIT $2;`, parameter...)
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]SmallDocument, 0)
	trunc := SmallDocument{}
	for result.Next() {
		err = result.Scan(&trunc.ID, &trunc.Type, &trunc.Organisation, &trunc.Title, &trunc.Author, &trunc.Written, &trunc.Removed)
		if err != nil {
			return nil, err
		}
		arr = append(arr, trunc)
	}
	return arr, nil
}

func GetPersonalDocumentList(amount int, page int, acc *Account) ([]SmallDocument, error) {
	result, err := postgresDB.Query(`SELECT id, type, organisation, title, author, written, removed FROM document 
    LEFT JOIN ownership ON ownership.account_name = author 
WHERE owner_name = $3 ORDER BY written DESC OFFSET $1 LIMIT $2;`, (page-1)*amount, amount, acc.GetName())
	if err != nil {
		return nil, err
	}
	defer closeRows(result)
	arr := make([]SmallDocument, 0)
	trunc := SmallDocument{}
	for result.Next() {
		err = result.Scan(&trunc.ID, &trunc.Type, &trunc.Organisation, &trunc.Title, &trunc.Author, &trunc.Written, &trunc.Removed)
		if err != nil {
			return nil, err
		}
		arr = append(arr, trunc)
	}
	return arr, nil
}
