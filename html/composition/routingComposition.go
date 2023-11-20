package composition

import (
	"PoliSim/data/database"
	"PoliSim/html/builder"
	"net/http"
)

const (
	MainBodyID         = "main-body-div"
	DocumentAdminPanel = "text-doc-admin-div"
	UserSelectionID    = "user-organisation-select-div"
	SidebarID          = "sidebar-element"
	MessageID          = "message-div"
	DisplayID          = "preview-display-div"

	Start                        builder.HttpUrl = "start"
	Login                        builder.HttpUrl = "login"
	Logout                       builder.HttpUrl = "logout"
	ViewDocument                 builder.HttpUrl = "document/view"
	BlockDocumentLink            builder.HttpUrl = "document/block/"
	BlockDocument                                = BlockDocumentLink + "{uuid}"
	updateUserSelectionLink      builder.HttpUrl = "document/update/account/selection/"
	UpdateUserSelection                          = updateUserSelectionLink + "{isAdmin}"
	CreateVoteDocument           builder.HttpUrl = "document/vote/create"
	requestVotePartialLink       builder.HttpUrl = "document/vote/request-partial/"
	RequestVotePartial                           = requestVotePartialLink + "{number}"
	ViewVoteDocumentLink         builder.HttpUrl = "document/vote/view/"
	ViewVoteDocument                             = ViewVoteDocumentLink + "{uuid}"
	VoteUpdateDocumentLink       builder.HttpUrl = "document/vote/update/"
	VoteUpdateDocument                           = VoteUpdateDocumentLink + "{uuid}"
	MakeVoteLink                 builder.HttpUrl = "document/vote/send/"
	MakeVote                                     = MakeVoteLink + "{doc}/{vote}/{type}"
	CreateDiscussionDocument     builder.HttpUrl = "document/discussion/create"
	DiscussionUpdateDocumentLink builder.HttpUrl = "document/discussion/update/"
	DiscussionUpdateDocument                     = DiscussionUpdateDocumentLink + "{uuid}"
	CommentDiscussionLink        builder.HttpUrl = "document/discussion/comment/"
	CommentDiscussion                            = CommentDiscussionLink + "{uuid}"
	ViewDiscussionDocumentLink   builder.HttpUrl = "document/discussion/view/"
	ViewDiscussionDocument                       = ViewDiscussionDocumentLink + "{uuid}"
	ChangeCommentDocumentLink    builder.HttpUrl = "document/discussion/change/comment/"
	ChangeCommentDocument                        = ChangeCommentDocumentLink + "{doc}/{comment}"
	CreateTextDocument           builder.HttpUrl = "document/text/create"
	ViewTextDocumentLink         builder.HttpUrl = "document/text/view/"
	ViewTextDocument                             = ViewTextDocumentLink + "{uuid}"
	AddTagDocumentLink           builder.HttpUrl = "document/text/add/tag/"
	AddTagDocument                               = AddTagDocumentLink + "{uuid}"
	ChangeTagDocumentLink        builder.HttpUrl = "document/text/change/tag/"
	ChangeTagDocument                            = ChangeTagDocumentLink + "{doc}/{tag}"
	CreateUser                   builder.HttpUrl = "account/create"
	EditUser                     builder.HttpUrl = "account/edit"
	SearchUser                   builder.HttpUrl = "account/search"
	ViewUser                     builder.HttpUrl = "account/view"
	ViewSelf                     builder.HttpUrl = "self/view"
	ChangePassword               builder.HttpUrl = "self/password/change"
	EditTitle                    builder.HttpUrl = "title/edit"
	SearchTitle                  builder.HttpUrl = "title/search"
	DeleteTitle                  builder.HttpUrl = "title/delete"
	CreateTitle                  builder.HttpUrl = "title/create"
	ViewTitles                   builder.HttpUrl = "title/view"
	getTitleSubGroup             builder.HttpUrl = "title/get-sub-group/"
	GetTitleSubGroup                             = getTitleSubGroup + "{mainGroup}/{subGroup}"
	CreateOrganisation           builder.HttpUrl = "organisation/create"
	EditOrganisation             builder.HttpUrl = "organisation/edit"
	SearchOrganisation           builder.HttpUrl = "organisation/search"
	ViewOrganisations            builder.HttpUrl = "organisation/view"
	ViewHiddenOrganisations      builder.HttpUrl = "organisation/hidden/view"
	getOrganisationSubGroup      builder.HttpUrl = "organisation/get-sub-group/"
	GetOrganisationSubGroup                      = getOrganisationSubGroup + "{mainGroup}/{subGroup}"
	CreatePressRelease           builder.HttpUrl = "press/release/create"
	ViewHiddenNewspaperList      builder.HttpUrl = "press/view/newspaper/hidden"
	ViewHiddenNewspaper                          = ViewHiddenNewspaperList + "/{uuid}"
	publishNewspaperLink         builder.HttpUrl = "press/publish/newspaper/"
	PublishNewspaper                             = publishNewspaperLink + "{uuid}"
	ViewNewspaperList            builder.HttpUrl = "press/newspaper"
	ViewNewspaper                                = ViewNewspaperList + "/{uuid}"
	rejectArticleLink            builder.HttpUrl = "press/article/reject/"
	RejectArticle                                = rejectArticleLink + "{uuid}"
	CreateLetter                 builder.HttpUrl = "letter/personal/create"
	CreateModmail                builder.HttpUrl = "letter/modmail/create"
	ViewLetterLink               builder.HttpUrl = "letter/personal/view/"
	ViewLetter                                   = ViewLetterLink + "{account}"
	ViewSingleLetter                             = ViewLetterLink + "{account}/{uuid}"
	updateLetterLink             builder.HttpUrl = "letter/personal/"
	UpdateLetter                                 = updateLetterLink + "{account}/{uuid}/{action}"
	ChangeViewLetterAccount      builder.HttpUrl = "letter/change/account"
	ViewModMails                 builder.HttpUrl = "letter/modmails/view"
	ErrorPage                    builder.HttpUrl = "errorPage"
	MarkdownFormPage             builder.HttpUrl = "render/form/markdown"
	MarkdownJsonPage             builder.HttpUrl = "render/json/markdown"

	// NotFound is only used as a way to keep the PageTitleMap in order
	NotFound builder.HttpUrl = "notFound"

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

var PageTitleMap = make(map[builder.HttpUrl]string)
var SidebarTitleMap = make(map[builder.HttpUrl]string)
var GetHTMXFunctions = make(map[builder.HttpUrl]http.HandlerFunc)
var PostHTMXFunctions = make(map[builder.HttpUrl]http.HandlerFunc)
var PatchHTMXFunctions = make(map[builder.HttpUrl]http.HandlerFunc)
