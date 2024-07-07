package handlers_test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"os"
	"programmerjournal-backend/handlers"
	"testing"
)

const dbTestPath = "./test.db"

func ListTasksForDay(t *testing.T) {
	h := handlers.New(dbTestPath)
	defer os.Remove(dbTestPath)

	// Create tasks

	// Verify

	resp, err := callListTasksForDay(h, "2024-06-01")
	if err != nil {
		t.Fatalf("error calling GetUrlStats: %v", err)
	}
}

//func TestVotes(t *testing.T) {
//	testCases := []struct {
//		Name           string
//		VoteUpRequests []handlers.VoteRequest
//		ExpectedResult handlers.SummaryEntry
//	}{
//		{
//			Name:           "no votes",
//			VoteUpRequests: nil,
//			ExpectedResult: handlers.SummaryEntry{
//				Url: "/url1",
//				Sum: 0,
//			},
//		},
//		{
//			Name: "one upvote",
//			VoteUpRequests: []handlers.VoteRequest{
//				{
//					Url: "/url1", Comment: "Very good!",
//				},
//			},
//			ExpectedResult: handlers.SummaryEntry{
//				Url: "/url1",
//				Sum: 1,
//			},
//		},
//		{
//			Name: "different urls",
//			VoteUpRequests: []handlers.VoteRequest{
//				{
//					Url: "/url0", Comment: "Very good!",
//				},
//				{
//					Url: "/url1", Comment: "Very good indeed!",
//				},
//			},
//			ExpectedResult: handlers.SummaryEntry{
//				Url: "/url1",
//				Sum: 1,
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.Name, func(t *testing.T) {
//			h := handlers.New(dbTestPath)
//			defer os.Remove(dbTestPath)
//
//			for _, upvote := range tc.VoteUpRequests {
//				err := callVoteUp(h, upvote)
//				if err != nil {
//					t.Fatalf("error calling VoteUp: %v", err)
//				}
//			}
//
//			resp, err := callGetVotesForUrl(h, "/url1")
//			if err != nil {
//				t.Fatalf("error calling GetUrlStats: %v", err)
//			}
//
//			if resp.Url != tc.ExpectedResult.Url {
//				t.Fatalf("Wrong entry path, got: %v, want: %v", resp.Url, tc.ExpectedResult.Url)
//			}
//
//			if resp.Sum != tc.ExpectedResult.Sum {
//				t.Fatalf("Wrong vote sum, got: %v, want: %v", resp.Sum, tc.ExpectedResult.Sum)
//			}
//		})
//	}
//}
//
//func callVoteUp(h handlers.Handlers, voteUpReq handlers.VoteRequest) error {
//	return makeVoteCall(h, voteUpReq, "/api/voteUp")
//}
//
//func makeVoteCall(h handlers.Handlers, voteUpReq handlers.VoteRequest, path string) error {
//	payloadBuf := new(bytes.Buffer)
//	err := json.NewEncoder(payloadBuf).Encode(voteUpReq)
//	if err != nil {
//		return fmt.Errorf("error encoding voteup request: %v", err)
//	}
//
//	req, err := http.NewRequest("POST", path, payloadBuf)
//	if err != nil {
//		return err
//	}
//
//	rr := httptest.NewRecorder()
//
//	router := createRouter(h)
//	router.ServeHTTP(rr, req)
//
//	if rr.Code != http.StatusCreated {
//		return fmt.Errorf("HTTP error returned: %v", rr.Code)
//	}
//	return nil
//}

func callListTasksForDay(h handlers.Handlers, date string) (handlers.SummaryEntry, error) {
	path := "/api/listTasksForDay?date=" + neturl.QueryEscape(date)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return handlers.SummaryEntry{}, err
	}

	rr := httptest.NewRecorder()

	router := createRouter(h)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		return handlers.SummaryEntry{}, fmt.Errorf("HTTP error returned: %v", rr.Code)
	}

	var entry handlers.SummaryEntry
	err = json.NewDecoder(rr.Body).Decode(&entry)
	if err != nil {
		return handlers.SummaryEntry{}, fmt.Errorf("error decoding GetVotesForUrl response: %v", err)
	}
	return entry, nil
}

func createRouter(h handlers.Handlers) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.Root)
	r.HandleFunc("/api/votesCount", h.GetVotesForUrl)
	r.HandleFunc("/api/voteUp", h.VoteUp)
	return r
}
