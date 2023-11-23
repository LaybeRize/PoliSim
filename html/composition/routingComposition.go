package composition

import (
	"PoliSim/data/database"
	. "PoliSim/html/builder"
	"net/http"
)

const (
	MainBodyID         = "main-body-div"
	DocumentAdminPanel = "text-doc-admin-div"
	UserSelectionID    = "user-organisation-select-div"
	SidebarID          = "sidebar-element"
	MessageID          = "message-div"
	DisplayID          = "preview-display-div"
	LetterSidebarID    = "letter-sidebar-button-id"

	Start                        HttpUrl = "start"
	Login                        HttpUrl = "login"
	Logout                       HttpUrl = "logout"
	ViewDocument                 HttpUrl = "document/view"
	BlockDocumentLink            HttpUrl = "document/block/"
	BlockDocument                        = BlockDocumentLink + "{uuid}"
	updateUserSelectionLink      HttpUrl = "document/update/account/selection/"
	UpdateUserSelection                  = updateUserSelectionLink + "{isAdmin}"
	CreateVoteDocument           HttpUrl = "document/vote/create"
	requestVotePartialLink       HttpUrl = "document/vote/request-partial/"
	RequestVotePartial                   = requestVotePartialLink + "{number}"
	ViewVoteDocumentLink         HttpUrl = "document/vote/view/"
	ViewVoteDocument                     = ViewVoteDocumentLink + "{uuid}"
	VoteUpdateDocumentLink       HttpUrl = "document/vote/update/"
	VoteUpdateDocument                   = VoteUpdateDocumentLink + "{uuid}"
	MakeVoteLink                 HttpUrl = "document/vote/send/"
	MakeVote                             = MakeVoteLink + "{doc}/{vote}/{type}"
	sseReaderVoteLink            HttpUrl = "document/vote/read/"
	SseReaderVote                        = sseReaderVoteLink + "{uuid}"
	sseReaderDiscussionLink      HttpUrl = "document/discussion/read/"
	SseReaderDiscussion                  = sseReaderDiscussionLink + "{uuid}"
	CreateDiscussionDocument     HttpUrl = "document/discussion/create"
	DiscussionUpdateDocumentLink HttpUrl = "document/discussion/update/"
	DiscussionUpdateDocument             = DiscussionUpdateDocumentLink + "{uuid}"
	CommentDiscussionLink        HttpUrl = "document/discussion/comment/"
	CommentDiscussion                    = CommentDiscussionLink + "{uuid}"
	ViewDiscussionDocumentLink   HttpUrl = "document/discussion/view/"
	ViewDiscussionDocument               = ViewDiscussionDocumentLink + "{uuid}"
	ChangeCommentDocumentLink    HttpUrl = "document/discussion/change/comment/"
	ChangeCommentDocument                = ChangeCommentDocumentLink + "{doc}/{comment}"
	CreateTextDocument           HttpUrl = "document/text/create"
	ViewTextDocumentLink         HttpUrl = "document/text/view/"
	ViewTextDocument                     = ViewTextDocumentLink + "{uuid}"
	AddTagDocumentLink           HttpUrl = "document/text/add/tag/"
	AddTagDocument                       = AddTagDocumentLink + "{uuid}"
	ChangeTagDocumentLink        HttpUrl = "document/text/change/tag/"
	ChangeTagDocument                    = ChangeTagDocumentLink + "{doc}/{tag}"
	CreateUser                   HttpUrl = "account/create"
	EditUser                     HttpUrl = "account/edit"
	SearchUser                   HttpUrl = "account/search"
	ViewUser                     HttpUrl = "account/view"
	ViewSelf                     HttpUrl = "self/view"
	ChangePassword               HttpUrl = "self/password/change"
	EditTitle                    HttpUrl = "title/edit"
	SearchTitle                  HttpUrl = "title/search"
	DeleteTitle                  HttpUrl = "title/delete"
	CreateTitle                  HttpUrl = "title/create"
	ViewTitles                   HttpUrl = "title/view"
	getTitleSubGroup             HttpUrl = "title/get-sub-group/"
	GetTitleSubGroup                     = getTitleSubGroup + "{mainGroup}/{subGroup}"
	CreateOrganisation           HttpUrl = "organisation/create"
	EditOrganisation             HttpUrl = "organisation/edit"
	SearchOrganisation           HttpUrl = "organisation/search"
	ViewOrganisations            HttpUrl = "organisation/view"
	ViewHiddenOrganisations      HttpUrl = "organisation/hidden/view"
	getOrganisationSubGroup      HttpUrl = "organisation/get-sub-group/"
	GetOrganisationSubGroup              = getOrganisationSubGroup + "{mainGroup}/{subGroup}"
	CreatePressRelease           HttpUrl = "press/release/create"
	ViewHiddenNewspaperList      HttpUrl = "press/view/newspaper/hidden"
	ViewHiddenNewspaper                  = ViewHiddenNewspaperList + "/{uuid}"
	publishNewspaperLink         HttpUrl = "press/publish/newspaper/"
	PublishNewspaper                     = publishNewspaperLink + "{uuid}"
	ViewNewspaperList            HttpUrl = "press/newspaper"
	ViewNewspaper                        = ViewNewspaperList + "/{uuid}"
	rejectArticleLink            HttpUrl = "press/article/reject/"
	RejectArticle                        = rejectArticleLink + "{uuid}"
	CreateLetter                 HttpUrl = "letter/personal/create"
	CreateModmail                HttpUrl = "letter/modmail/create"
	ViewLetterLink               HttpUrl = "letter/personal/view/"
	ViewLetter                           = ViewLetterLink + "{account}"
	ViewSingleLetter                     = ViewLetterLink + "{account}/{uuid}"
	updateLetterLink             HttpUrl = "letter/personal/"
	UpdateLetter                         = updateLetterLink + "{account}/{uuid}/{action}"
	ChangeViewLetterAccount      HttpUrl = "letter/change/account"
	MarkAllLetterAccount         HttpUrl = "letter/mark/all"
	ViewModMails                 HttpUrl = "letter/modmails/view"
	ErrorPage                    HttpUrl = "errorPage"
	MarkdownFormPage             HttpUrl = "render/form/markdown"
	MarkdownJsonPage             HttpUrl = "render/json/markdown"

	// NotFound is only used as a way to keep the PageTitleMap in order
	NotFound HttpUrl = "notFound"

	// APIPreRoute is a subroute for the web application to prepend to any
	// backend partial replies. It never starts with a / because that is automatically prepend anyway
	APIPreRoute = "htmx/"
)

type HttpHandling struct {
	TitleText          string
	SidebarButtonText  string
	HasSidebarButton   bool
	SidebarSubMenuText string
	HasSidebarSubMenu  bool
	RoleLevel          database.RoleLevel
}

var PageTitleMap = make(map[HttpUrl]string)
var SidebarTitleMap = make(map[HttpUrl]string)
var GetHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var PostHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
var PatchHTMXFunctions = make(map[HttpUrl]http.HandlerFunc)
