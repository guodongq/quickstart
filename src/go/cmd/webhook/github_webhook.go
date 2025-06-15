package main

import (
	"fmt"
	"github.com/google/go-github/v72/github"
	"log"
	"net/http"
	"strings"
)

var secret = "123"

const (
	Label   = "/label"
	UnLabel = "/un-label"
	LGTM    = "/lgtm" // rebase
	Merge   = "/merge"
	Close   = "/close"
	Reopen  = "/reopen"
	ReOpen  = "/re-open"
	Approve = "/approve"
	Update  = "/update"
)

func main() {
	http.HandleFunc("/webhook/notify", func(w http.ResponseWriter, r *http.Request) {
		payload, validateErr := github.ValidatePayload(r, []byte(secret))
		if validateErr != nil {
			http.Error(w, "The GitHub signature header is invalid.", http.StatusUnauthorized)
			log.Printf("invalid signature: %s\n", validateErr.Error())
			return
		}

		event, parseErr := github.ParseWebHook(github.WebHookType(r), payload)
		if parseErr != nil {
			http.Error(w, "The payload parsed failed", http.StatusInternalServerError)
			log.Printf("could not parse webhook: %s\n", parseErr)
			return
		}

		switch e := event.(type) {
		case *github.PushEvent:
		case *github.PullRequestEvent:
			pullRequestEvent := *e
			_ = pullRequestEvent
		case *github.PullRequestReviewEvent:
		case *github.IssueEvent:
			issueEvent := *e
			_ = issueEvent
		case *github.IssueCommentEvent:
			action := e.GetAction()
			fmt.Printf("IssueCommentEvent: %s\n", action)
			commentBody := e.GetComment().GetBody()
			if action == "edited" || action == "created" {
				issueCommentEvent := *e
				_ = issueCommentEvent
				// avoid recursion comment by bot
				//if issueCommentEvent.GetSender().GetLogin() == botName {
				//	_, _ = fmt.Fprintf(w, "ok")
				//	return
				//}
				if strings.Contains(commentBody, Label) {
					//addLabelsByComment(commentBody, githubClient, issueCommentEvent)
				}
				if strings.Contains(commentBody, UnLabel) {
					//removeLabelFromIssue(commentBody, githubClient, issueCommentEvent)
				}
				if strings.Contains(commentBody, Approve) {
					//approvePullRequest(githubClient, issueCommentEvent)
				}
				if strings.Contains(commentBody, LGTM) {
					//mergePullRequest(githubClient, issueCommentEvent, "rebase")
				}
				if strings.Contains(commentBody, Merge) {
					//mergePullRequest(githubClient, issueCommentEvent, "merge")
				}
				if strings.Contains(commentBody, Close) {
					//closeOrOpenIssue(githubClient, issueCommentEvent, false)
				}
				if strings.Contains(commentBody, Reopen) || strings.Contains(commentBody, ReOpen) {
					//closeOrOpenIssue(githubClient, issueCommentEvent, true)
				}
				if strings.Contains(commentBody, Update) {
					//updatePullRequest(githubClient, issueCommentEvent)
				}
			}
		default:
			log.Printf("unknown event type %s\n", github.WebHookType(r))
			_, _ = fmt.Fprintf(w, "ok")
			return
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
