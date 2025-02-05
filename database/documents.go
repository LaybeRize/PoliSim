package database

import (
	loc "PoliSim/localisation"
	"fmt"
	"html/template"
	"log/slog"
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
		ID                  string
		Type                DocumentType
		Organisation        string
		Title               string
		Author              string
		Flair               string
		Written             time.Time
		Body                template.HTML
		Public              bool
		Removed             bool
		MemberParticipation bool
		AdminParticipation  bool
		AllowedToAddTags    bool
		End                 time.Time
		Reader              []string
		Participants        []string
		Tags                []DocumentTag
		Links               []VoteInfo
		VoteIDs             []string
		Comments            []DocumentComment
		Result              []AccountVotes
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
)

func (s *SmallDocument) IsPost() bool { return s.Type == DocTypePost }

func (s *SmallDocument) IsDiscussion() bool { return s.Type == DocTypeDiscussion }

func (s *SmallDocument) IsVote() bool { return s.Type == DocTypeVote }

func (s *SmallDocument) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return s.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return s.Written.Format(loc.TimeFormatString)
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

type DocumentTag struct {
	ID              string
	Outgoing        bool
	Text            string
	Written         time.Time
	BackgroundColor string
	TextColor       string
	LinkColor       string
	Links           []string
	QueriedLinks    []any
}

func (t *DocumentTag) GetTimeWritten(a *Account) string {
	if a.Exists() {
		return t.Written.In(a.TimeZone).Format(loc.TimeFormatString)
	}
	return t.Written.Format(loc.TimeFormatString)
}

func (t *DocumentTag) HasLinks() bool {
	return len(t.QueriedLinks) != 0
}

type DocumentComment struct {
	ID      string
	Author  string
	Flair   string
	Written time.Time
	Body    template.HTML
	Removed bool
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

func CreateDocument(document *Document, acc *Account) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}

	result, err := tx.Run(`MATCH (acc:Account)-[:ADMIN]->(o:Organisation) 
WHERE acc.name = $Author AND acc.blocked = false AND o.name = $organisation AND o.visibility <> $hidden
	RETURN o.visibility;`,
		map[string]any{"Author": document.Author,
			"organisation": document.Organisation,
			"hidden":       HIDDEN})
	if err != nil {
		return err
	} else if result.Next(); result.Record() == nil {
		return notAllowedError
	} else if vis := OrganisationVisibility(result.Record().Values[0].(int64)); (vis == SECRET && document.Public) || (vis == PUBLIC && !document.Public) {
		return notAllowedError
	}

	err = tx.RunWithoutResult(`MATCH (a:Account) WHERE a.name = $author 
MATCH (o:Organisation) WHERE o.name = $organisation 
CREATE (d:Document {id: $id, title: $title, type: $type, author: $author, flair: $flair, body: $body, removed: false,
end_time: $end_time, written: $written, public: $public, member_part: $member_part, admin_part: $admin_part}) 
MERGE (a)-[:WRITTEN]->(d)
MERGE (d)-[:IN]->(o);`, map[string]any{
		"organisation": document.Organisation,
		"id":           document.ID,
		"title":        document.Title,
		"type":         document.Type,
		"author":       document.Author,
		"end_time":     document.End,
		"written":      time.Now().UTC(),
		"body":         document.Body,
		"flair":        document.Flair,
		"public":       document.Public,
		"member_part":  document.MemberParticipation,
		"admin_part":   document.AdminParticipation})
	if err != nil {
		return err
	}

	if !document.Public {
		err = tx.RunWithoutResult(`MATCH (a:Account), (d:Document) WHERE a.name IN $reader AND d.id = $id 
CREATE (a)<-[:READER]-(d);`, map[string]any{"reader": document.Reader, "id": document.ID})
		if err != nil {
			return err
		}
	}

	if document.Type == DocTypeVote {
		result, err = tx.Run(`
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
			return notAllowedError
		}
	}

	err = tx.RunWithoutResult(`MATCH (a:Account), (d:Document) WHERE a.name IN $user AND d.id = $id 
CREATE (a)<-[:PARTICIPANT]-(d);`, map[string]any{"user": document.Participants, "id": document.ID})
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func GetDocumentForUser(id string, acc *Account) (*Document, []string, error) {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return nil, nil, err
	}

	result, err := tx.Run(`MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE d.id = $id 
RETURN d, o.name;`, map[string]any{
		"id": id})
	if err != nil {
		return nil, nil, err
	} else if !result.Next() {
		return nil, nil, notAllowedError
	}
	props := GetPropsMapForRecordPosition(result.Record(), 0)
	public := props.GetBool("public")
	allowedToAddTags := false

	if !public && !acc.Exists() {
		return nil, nil, notAllowedError
	} else if acc.Exists() {
		var userCheck *dbResult
		userCheck, err = tx.Run(`
CALL { 
MATCH (a:Account)-[*..]->(o:Organisation)<-[:IN]-(d:Document) WHERE a.name = $name AND d.id = $id 
RETURN o, d 
UNION 
MATCH (o:Organisation)<-[:IN]-(d:Document)-[:READER|PARTICIPANT]->(a:Account) WHERE a.name = $name AND d.id = $id 
RETURN o,d 
} 
OPTIONAL MATCH (b:Account)-[:OWNER|ADMIN*..]->(o) WHERE b.name = $name 
RETURN d.id, b.name;`, map[string]any{
			"id":   id,
			"name": acc.Name})
		if err != nil {
			return nil, nil, err
		} else if !userCheck.Peek() && !acc.IsAtLeastAdmin() && !public {
			return nil, nil, notAllowedError
		} else if userCheck.Next() {
			slog.Debug("Document Connections", "id", id, "query output", userCheck.Record().Values)
			allowedToAddTags = userCheck.Record().Values[1] != nil
		}
	}

	doc := &Document{
		ID:                  id,
		Type:                DocumentType(props.GetInt("type")),
		Organisation:        result.Record().Values[1].(string),
		Title:               props.GetString("title"),
		Author:              props.GetString("author"),
		Flair:               props.GetString("flair"),
		Written:             props.GetTime("written"),
		Body:                template.HTML(props.GetString("body")),
		Public:              public,
		Removed:             props.GetBool("removed"),
		MemberParticipation: props.GetBool("member_part"),
		AdminParticipation:  props.GetBool("admin_part"),
		End:                 props.GetTime("end_time"),
		Tags:                make([]DocumentTag, 0),
		AllowedToAddTags:    allowedToAddTags,
		Result:              nil,
	}
	var commentator []string

	if !doc.Public {
		doc.Reader = make([]string, 0)
		result, err = tx.Run(`MATCH (d:Document)-[:READER]->(a:Account)
WHERE d.id = $id RETURN a.name;`,
			map[string]any{"id": id})
		if err != nil {
			return nil, nil, err
		}
		for result.Next() {
			doc.Reader = append(doc.Reader, result.Record().Values[0].(string))
		}
	}

	if !(doc.Type == DocTypePost) {
		doc.Participants = make([]string, 0)
		result, err = tx.Run(`MATCH (d:Document)-[:PARTICIPANT]->(a:Account)
WHERE d.id = $id RETURN a.name;`,
			map[string]any{"id": id})
		if err != nil {
			return nil, nil, err
		}
		for result.Next() {
			doc.Participants = append(doc.Participants, result.Record().Values[0].(string))
		}
	}

	if doc.Type == DocTypeDiscussion {
		doc.Comments = make([]DocumentComment, 0)
		result, err = tx.Run(`MATCH (d:Document)<-[:ON]-(c:Comment)
WHERE d.id = $id RETURN c.id, c.author, c.flair, c.written, c.body, c.removed ORDER BY c.written;`,
			map[string]any{"id": id})
		if err != nil {
			return nil, nil, err
		}
		for result.Next() {
			doc.Comments = append(doc.Comments, DocumentComment{
				ID:      result.Record().Values[0].(string),
				Author:  result.Record().Values[1].(string),
				Flair:   result.Record().Values[2].(string),
				Written: result.Record().Values[3].(time.Time),
				Body:    template.HTML(result.Record().Values[4].(string)),
				Removed: result.Record().Values[5].(bool),
			})
		}
	}

	if doc.Type == DocTypeDiscussion && !doc.Ended() && acc.Exists() {
		commentator = make([]string, 0)
		result, err = tx.Run(`
MATCH (n:Account)-[:OWNER*0..]->(a:Account) WHERE n.name = $name 
RETURN a.name;`,
			map[string]any{"id": id, "name": acc.Name})
		if err != nil {
			return nil, nil, err
		}
		for result.Next() {
			commentator = append(commentator, result.Record().Values[0].(string))
		}
	}

	if doc.Type == DocTypeVote {
		result, err = tx.Run(`MATCH (d:Document)-[:VOTED]->(r:Result)
WHERE d.id = $id RETURN r;`,
			map[string]any{"id": id})
		if err != nil {
			return nil, nil, err
		} else if result.Peek() {
			doc.Result = make([]AccountVotes, 0)
			for result.Next() {
				props = GetPropsMapForRecordPosition(result.Record(), 0)
				doc.Result = append(doc.Result, AccountVotes{
					Question:        props.GetString("question"),
					IterableAnswers: props.GetArray("answers"),
					Anonymous:       props.GetBool("anonymous"),
					Type:            VoteType(props.GetInt("type")),
					AnswerAmount:    props.GetInt("amount"),
					IllegalVotes:    nil,
					Illegal:         props.GetArray("illegal"),
					List:            nil,
					Voter:           props.GetArray("voter"),
					Votes:           props.GetArray("votes"),
				})
			}
		} else {
			doc.Links = make([]VoteInfo, 0)
			result, err = tx.Run(`MATCH (d:Document)-[:LINKS]->(v:Vote)
WHERE d.id = $id RETURN v.id, v.question;`,
				map[string]any{"id": id})
			if err != nil {
				return nil, nil, err
			}
			for result.Next() {
				doc.Links = append(doc.Links, VoteInfo{
					ID:       result.Record().Values[0].(string),
					Question: result.Record().Values[1].(string),
				})
			}
		}
	}

	result, err = tx.Run(`CALL { 
MATCH (d:Document)-[:LINKS]->(t:Tag) 
WHERE d.id = $id 
OPTIONAL MATCH (t)-[:LINKS]->(r:Document) 
RETURN t, collect(r.id) AS ids, true AS outgoing 
UNION 
MATCH (d:Document)<-[:LINKS]-(t:Tag)<-[:LINKS]-(r:Document) 
WHERE d.id = $id 
RETURN t, collect(r.id) AS ids, false AS outgoing 
} 
RETURN t, ids, outgoing ORDER BY t.written DESC;`, map[string]any{"id": id})
	if err != nil {
		return nil, nil, err
	}
	for result.Next() {
		props = GetPropsMapForRecordPosition(result.Record(), 0)
		doc.Tags = append(doc.Tags, DocumentTag{
			ID:              props.GetString("id"),
			Outgoing:        result.Record().Values[2].(bool),
			Text:            props.GetString("text"),
			Written:         props.GetTime("written"),
			BackgroundColor: props.GetString("background"),
			TextColor:       props.GetString("color"),
			LinkColor:       props.GetString("link"),
			QueriedLinks:    result.Record().Values[1].([]any),
		})
	}

	err = tx.Commit()
	return doc, commentator, err
}

func RemoveRestoreDocument(docId string) {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return
	}

	err = tx.RunWithoutResult(`MATCH (d:Document) WHERE d.id = $id  
SET d.removed = NOT d.removed;`, map[string]any{
		"id": docId,
	})
	if err != nil {
		return
	}

	_ = tx.Commit()
}

func CreateTagForDocument(docID string, acc *Account, tag *DocumentTag) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}

	result, err := tx.Run(`MATCH (a:Account)-[:ADMIN|OWNER*..]->(o:Organisation)<-[:IN]-(d:Document) 
WHERE a.name = $name AND d.id = $id 
RETURN a.name;`, map[string]any{"name": acc.Name, "id": docID})
	if err != nil {
		return err
	} else if result.Next(); result.Record() == nil {
		return notAllowedError
	}

	err = tx.RunWithoutResult(`MATCH (d:Document) WHERE d.id = $id  
CREATE (t:Tag {id: $tagId, text: $text, written: $written, background: $background, color: $color, link: $link}) 
MERGE (d)-[:LINKS]->(t);`, map[string]any{
		"id":         docID,
		"tagId":      tag.ID,
		"links":      tag.Links,
		"text":       tag.Text,
		"written":    time.Now().UTC(),
		"background": tag.BackgroundColor,
		"color":      tag.TextColor,
		"link":       tag.LinkColor,
	})
	if err != nil {
		return err
	}

	err = tx.RunWithoutResult(`MATCH (target:Document) WHERE target.id <> $id AND target.id IN $links 
MATCH (t:Tag) WHERE t.id = $tagID 
MERGE (t)-[:LINKS]->(target);`, map[string]any{"tagID": tag.ID, "id": docID, "links": tag.Links})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CreateDocumentComment(documentId string, comment *DocumentComment) error {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return err
	}

	result, err := tx.Run(`CALL {
MATCH (a:Account)<-[:PARTICIPANT]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id AND datetime($now) < datetime(d.end_time) 
RETURN a 
UNION 
MATCH (a:Account)-[:USER|ADMIN]->(:Organisation)<-[:IN]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id AND d.member_part = true AND datetime($now) < datetime(d.end_time) 
RETURN a 
UNION 
MATCH (a:Account)-[:ADMIN]->(:Organisation)<-[:IN]-(d:Document) 
WHERE a.name = $author AND a.blocked = false AND d.id = $id AND d.admin_part = true AND datetime($now) < datetime(d.end_time) 
RETURN a 
} 
RETURN a;`,
		map[string]any{"id": documentId, "author": comment.Author, "now": time.Now().UTC()})
	if err != nil {
		return err
	} else if !result.Peek() {
		return notAllowedError
	}

	err = tx.RunWithoutResult(`MATCH (a:Account) WHERE a.name = $author 
MATCH (d:Document) WHERE d.id = $doc_ID 
CREATE (c:Comment {id: $id, author: $author, flair: $flair, body: $body, removed: false, written: $written}) 
MERGE (a)-[:WRITTEN]->(c)
MERGE (c)-[:ON]->(d);`, map[string]any{
		"doc_ID":  documentId,
		"id":      comment.ID,
		"author":  comment.Author,
		"flair":   comment.Flair,
		"written": time.Now().UTC(),
		"body":    comment.Body})
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func RemoveRestoreComment(commentID string) {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		return
	}

	err = tx.RunWithoutResult(`MATCH (c:Comment) WHERE c.id = $id  
SET c.removed = NOT c.removed;`, map[string]any{
		"id": commentID,
	})
	if err != nil {
		return
	}

	_ = tx.Commit()
}

func GetDocumentList(amount int, page int, acc *Account, showBlocked bool) ([]SmallDocument, error) {
	var query string
	if !acc.Exists() {
		query = `MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE d.public = true AND d.removed = false `
	} else if !acc.IsAtLeastAdmin() {
		query = `CALL { MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE d.public = true AND d.removed = false 
RETURN d, o
UNION
MATCH (a:Account)-[*..]->(o:Organisation)<-[:IN]-(d:Document) WHERE d.public = false AND d.removed = false AND a.name = $name 
RETURN d, o 
} `
	} else {
		if showBlocked {
			query = `MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE true `
		} else {
			query = `MATCH (o:Organisation)<-[:IN]-(d:Document) WHERE d.removed = false `
		}
	}

	result, err := makeRequest(query+`RETURN d.id, d.type, o.name, d.title, d.author, d.written, d.removed 
ORDER BY d.written DESC SKIP $skip LIMIT $amount;`,
		map[string]any{
			"amount": amount,
			"skip":   (page - 1) * amount,
			"name":   acc.GetName(),
		})
	if err != nil {
		return nil, err
	}
	arr := make([]SmallDocument, 0, len(result))
	for _, record := range result {
		arr = append(arr, SmallDocument{
			ID:           record.Values[0].(string),
			Type:         DocumentType(record.Values[1].(int64)),
			Organisation: record.Values[2].(string),
			Title:        record.Values[3].(string),
			Author:       record.Values[4].(string),
			Written:      record.Values[5].(time.Time),
			Removed:      record.Values[6].(bool),
		})
	}
	return arr, nil
}

func GetPersonalDocumentList(amount int, page int, acc *Account) ([]SmallDocument, error) {
	result, err := makeRequest(`MATCH (a:Account)-[:OWNER|WRITTEN]->(d:Document) WHERE a.name = $name AND d.removed = false 
RETURN d.id, d.type, o.name, d.title, d.author, d.written, d.removed 
ORDER BY d.written DESC SKIP $skip LIMIT $amount;`,
		map[string]any{
			"amount": amount,
			"skip":   (page - 1) * amount,
			"name":   acc.Name,
		})
	if err != nil {
		return nil, err
	}
	arr := make([]SmallDocument, 0, len(result))
	for _, record := range result {
		arr = append(arr, SmallDocument{
			ID:           record.Values[0].(string),
			Type:         DocumentType(record.Values[1].(int64)),
			Organisation: record.Values[2].(string),
			Title:        record.Values[3].(string),
			Author:       record.Values[4].(string),
			Written:      record.Values[5].(time.Time),
			Removed:      record.Values[6].(bool),
		})
	}
	return arr, nil
}
