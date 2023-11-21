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
	uuidStr := chi.URLParam(r, "uuid")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")

	commentChannel := make(chan database.Votes)

	go addVoteToChannel(r.Context(), uuidStr, commentChannel)

	for info := range commentChannel {
		var event string
		var err error

		event, err = formatServerSentEvent(func() ([]string, error) {
			newErr := error(nil)
			buff := bytes.NewBuffer([]byte{})
			if info.Info.VoteMethod == database.RankedVotes {
				newErr = composition.GetInfoRankedView(&info, false).Render(buff)
			} else {
				newErr = composition.GetInfoStandardView(&info, false).Render(buff)
			}
			return []string{buff.String(), fmt.Sprintf(composition.VoteInfoDiv, info.UUID), ""}, newErr
		})

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

func addVoteToChannel(ctx context.Context, id string, channel chan database.Votes) {
	uniqueID := uuid.New().String()
	logic.AddVoteChannel(id, uniqueID, channel)
outerloop:
	for {
		select {
		case <-ctx.Done():
			break outerloop
		}
	}

	logic.RemoveVoteChannel(id, uniqueID)
	close(channel)
}

func GetSSEReaderForDiscussionService(w http.ResponseWriter, r *http.Request) {
	_, isAdmin := CheckUserPrivileges(r, database.HeadAdmin, database.Admin)
	uuidStr := chi.URLParam(r, "uuid")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")

	commentChannel := make(chan logic.CommentUpdate)

	go addCommentToChannel(r.Context(), uuidStr, commentChannel)

	for info := range commentChannel {
		event, err := formatServerSentEvent(func() ([]string, error) {
			id := composition.AdditionDiv
			buff := bytes.NewBuffer([]byte{})
			err := composition.GetCommentRendered(uuidStr, &info.Discussion, isAdmin).Render(buff)
			if err != nil {
				return []string{}, err
			}
			replacer := fmt.Sprintf(fmt.Sprintf(composition.CommentSingleDivID, info.Discussion.UUID))
			if info.Change {
				id = replacer
			} else {
				err = composition.GetNewAdditionSSEDiv().Render(buff)
			}
			if err != nil {
				return []string{}, err
			}
			return []string{buff.String(), id, replacer}, err
		})

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

func addCommentToChannel(ctx context.Context, id string, channel chan logic.CommentUpdate) {
	uniqueID := uuid.New().String()
	logic.AddCommentChannel(id, uniqueID, channel)
outerloop:
	for {
		select {
		case <-ctx.Done():
			break outerloop
		}
	}

	logic.RemoveCommentChannel(id, uniqueID)
	close(channel)
}

func formatServerSentEvent(f func() ([]string, error)) (string, error) {
	array, err := f()
	if err != nil || len(array) != 3 {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: change\n"))
	sb.WriteString(fmt.Sprintf("data: {\"data\": \"%s\", \"id\":\"%s\", \"replace\":\"%s\"}\n\n",
		strings.ReplaceAll(array[0], "\"", "\\\""),
		array[1], array[2]))

	return sb.String(), nil
}
