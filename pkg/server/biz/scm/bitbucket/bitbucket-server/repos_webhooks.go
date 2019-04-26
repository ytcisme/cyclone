/*
Copyright 2017 caicloud authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"fmt"
	"net/http"
)

// Ref to: https://confluence.atlassian.com/bitbucketserver/event-payload-938025882.html
const (
	// RefsChanged is the event key about refs changed
	RefsChanged = "repo:refs_changed"
	// PrOpened is the event key about pull request opened.
	PrOpened = "pr:opened"
	// PrModified is the event key about pull request modified.
	PrModified = "pr:modified"
	// PrCommentAdded is the event key about a comment added on the pull request.
	PrCommentAdded = "pr:comment:added"
)

// Webhook contains webhook information.
type Webhook struct {
	ID     string   `json:"id,omitempty"`
	Name   string   `json:"name"`
	Events []string `json:"events,omitempty"`
	URL    string   `json:"url"`
	Active bool     `json:"active"`
}

// Webhooks is a set of webhooks.
type Webhooks struct {
	Pagination
	Values []Webhook `json:"values"`
}

// CreateWebhook create a new webhook.
func (server *RepositoriesService) CreateWebhook(ctx context.Context, project string, repo string, webhookReq Webhook) (*Webhook, *http.Response, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks", project, repo)
	req, err := server.v1Client.NewRequest("POST", u, webhookReq, nil)
	if err != nil {
		return nil, nil, err
	}
	var webhook Webhook
	resp, err := server.v1Client.Do(req, &webhook)
	return &webhook, resp, err
}

// DeleteWebhook delete a webhook.
func (server *RepositoriesService) DeleteWebhook(ctx context.Context, project string, repo string, webhookID string) (*http.Response, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks/%s", project, repo, webhookID)
	req, err := server.v1Client.NewRequest("DELETE", u, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := server.v1Client.Do(req, nil)
	return resp, err
}

// ListWebhook list webhooks on a repository.
func (server *RepositoriesService) ListWebhook(ctx context.Context, project string, repo string) (*Webhooks, *http.Response, error) {
	u := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/webhooks", project, repo)
	req, err := server.v1Client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var webhooks Webhooks
	resp, err := server.v1Client.Do(req, &webhooks)
	return &webhooks, resp, err
}
