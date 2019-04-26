package server

import (
	"encoding/json"
	"fmt"
	"github.com/caicloud/cyclone/pkg/server/biz/scm"
	"github.com/caicloud/nirvana/log"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	tag  = "TAG"
	branch = "BRANCH"
	add = "ADD"
	update = "UPDATE"

	releaseRefPrefix = "refs/release"
	// pullRefTemplate represents reference template for pull request.
	pullRefTemplate = "refs/pull-requests/%d/from"
)


var supportEventMap = map[string]bool {
	RefsChanged: true,
	PrOpened: true,
	PrModified: true,
	PrCommentAdded: true,
}

type EventPayload struct {
	EventKey string `json:"eventKey"`
	Repository Repository `json:"repository"`
	Changes []Change `json:"changes"`
	PullRequest PullRequest `json:"pullRequest"`
	Comment Comment `json:"comment"`
}

type Change struct {
	Ref struct{
		ID string `json:"id"`
		Type string `json:"type"`
		DisplayID string `json:"displayId"`
	} `json:"ref"`
	RefID string `json:"refId"`
	FromHash string `json:"fromHash"`
	ToHash string `json:"toHash"`
	Type string `json:"type"`
}

type Comment struct {
	Text string `json:"text"`
}

// ParseEvent parses data from Bitbucket Server events.
func ParseEvent(request *http.Request) *scm.EventData {
	eventKey := request.Header.Get("X-Event-Key")
	if !supportEventMap[eventKey] {
		log.Errorf("unsupported bitbucket server event: %s", eventKey)
		return nil
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	var payload EventPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	switch eventKey {
	case RefsChanged:
		for _, change := range payload.Changes {
			if change.Ref.Type == tag && change.Type == add {
				return &scm.EventData{
					Type: scm.TagReleaseEventType,
					Repo: fmt.Sprintf("%s/%s", payload.Repository.Project.Key, payload.Repository.Slug),
					Ref:  change.RefID,
				}
			} else if change.Ref.Type == branch && strings.HasPrefix(change.RefID, releaseRefPrefix) && change.Type == add {
				return &scm.EventData{
					Type: scm.TagReleaseEventType,
					Repo: fmt.Sprintf("%s/%s", payload.Repository.Project.Key, payload.Repository.Slug),
					Ref:  change.RefID,
				}
			} else if change.Type == update {
				return &scm.EventData{
					Type:   scm.PushEventType,
					Repo:   fmt.Sprintf("%s/%s", payload.Repository.Project.Key, payload.Repository.Slug),
					Ref:    change.RefID,
					Branch: change.RefID,
				}
			}
		}
	case PrOpened:
		return &scm.EventData{
			Type:      scm.PullRequestEventType,
			Repo:      fmt.Sprintf("%s/%s", payload.Repository.Project.Key, payload.Repository.Slug),
			Ref:       fmt.Sprintf(pullRefTemplate, payload.PullRequest.ID),
			CommitSHA: payload.PullRequest.FromRef.LatestCommit,
		}
	case PrModified:
		return &scm.EventData{
			Type:      scm.PullRequestEventType,
			Repo:      fmt.Sprintf("%s/%s", payload.Repository.Project.Key, payload.Repository.Slug),
			Ref:       fmt.Sprintf(pullRefTemplate, payload.PullRequest.ID),
			CommitSHA: payload.PullRequest.FromRef.LatestCommit,
		}
	case PrCommentAdded:
		return &scm.EventData{
			Type:      scm.PullRequestCommentEventType,
			Repo:      fmt.Sprintf("%s/%s", payload.Repository.Project.Key, payload.Repository.Slug),
			Ref:       fmt.Sprintf(pullRefTemplate, payload.PullRequest.ID),
			Comment:   payload.Comment.Text,
			CommitSHA: payload.PullRequest.FromRef.LatestCommit,
		}
	default:
		log.Warningln("Skip unsupported Bitbucket Server event")
		return nil
	}
	return nil
}
