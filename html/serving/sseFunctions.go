package serving

import (
	"PoliSim/data/database"
	"PoliSim/data/logic"
	"PoliSim/html/composition"
	"bytes"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func InstallSSEHandlers() {
	composition.GetHTMXFunctions[composition.SseReaderDiscussion] = GetSSEReaderForDiscussionService
	composition.GetHTMXFunctions[composition.SseReaderVote] = GetSSEReaderForVoteService
}

func GetSSEReaderForVoteService(w http.ResponseWriter, r *http.Request) {
	generlizeSSEConnection[database.Votes](w, r, addVoteToChannel, getVoteEvent)
}

func getVoteEvent(info database.Votes, isAdmin bool, uuidStr string) (*SendEventStruct, error) {
	newErr := error(nil)
	buff := bytes.NewBuffer([]byte{})
	if info.Info.VoteMethod == database.RankedVotes {
		newErr = composition.GetInfoRankedView(&info, false).Render(buff)
	} else {
		newErr = composition.GetInfoStandardView(&info, false).Render(buff)
	}
	return &SendEventStruct{
		HTML:      buff.String(),
		EventName: info.UUID,
	}, newErr
}

func addVoteToChannel(ctx context.Context, id string, channel chan database.Votes, accountID int64) {
	newID := uuid.New().String()
	logic.AddVoteChannel(id, newID, accountID, channel)
outerloop:
	for {
		select {
		case <-ctx.Done():
			break outerloop
		}
	}

	logic.RemoveVoteChannel(id, newID, accountID)
	close(channel)
}

func GetSSEReaderForDiscussionService(w http.ResponseWriter, r *http.Request) {
	generlizeSSEConnection[logic.CommentUpdate](w, r, addCommentToChannel, getCommentEvent)
}

func getCommentEvent(info logic.CommentUpdate, isAdmin bool, uuidStr string) (*SendEventStruct, error) {
	buff := bytes.NewBuffer([]byte{})
	err := composition.GetCommentRendered(uuidStr, &info.Discussion, isAdmin).Render(buff)
	if err != nil {
		return nil, err
	}
	eventName := info.Discussion.UUID
	if !info.Change {
		eventName = composition.EventAddComment
	}
	if err != nil {
		return nil, err
	}
	return &SendEventStruct{
		HTML:      buff.String(),
		EventName: eventName,
	}, err
}

func addCommentToChannel(ctx context.Context, id string, channel chan logic.CommentUpdate, accountID int64) {
	newID := uuid.New().String()
	logic.AddCommentChannel(id, newID, accountID, channel)
outerloop:
	for {
		select {
		case <-ctx.Done():
			break outerloop
		}
	}

	logic.RemoveCommentChannel(id, newID, accountID)
	close(channel)
}

type SendEventStruct struct {
	HTML      string
	EventName string
}

func generlizeSSEConnection[t any](w http.ResponseWriter, r *http.Request,
	addRoutine func(ctx context.Context, id string, channel chan t, accountID int64),
	formater func(info t, isAdmin bool, uuidStr string) (*SendEventStruct, error)) {
	acc, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	uuidStr := chi.URLParam(r, "uuid")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")

	commentChannel := make(chan t)

	go addRoutine(r.Context(), uuidStr, commentChannel, acc.ID)

	for info := range commentChannel {
		event, err := formatServerSentEvent[t](info, isAdmin, uuidStr, formater)

		if err != nil {
			fmt.Println(err)
			break
		}

		_, err = fmt.Fprint(w, event)
		if err != nil {
			fmt.Println(err)
			break
		}

		flusher.Flush()
	}
}

func formatServerSentEvent[t any](info t, isAdmin bool, uuidStr string,
	formater func(info t, isAdmin bool, uuidStr string) (*SendEventStruct, error)) (string, error) {
	convert, err := formater(info, isAdmin, uuidStr)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", convert.EventName))
	sb.WriteString(fmt.Sprintf("data: %s\n\n", convert.HTML))

	return sb.String(), nil
}
