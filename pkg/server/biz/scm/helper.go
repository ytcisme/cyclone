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

package scm

import "strings"

const (
	// ListPerPageOpt represents the value for PerPage in list options.
	ListPerPageOpt = 30
)

// ParseServerURL is a helper func to parse the server url, such as https://github.com/ to return server(github.com).
func ParseServerURL(url string) string {
	strs := strings.SplitN(strings.TrimSuffix(url, "/"), "://", -1)
	return strs[1]
}

// ParseRepo parses owner and name from full repo name.
// For example, parse caicloud/cyclone will return owner(caicloud) and name(cyclone).
func ParseRepo(url string) (string, string) {
	strs := strings.SplitN(url, "/", -1)
	if len(strs) < 2 {
		return strs[0], ""
	}

	return strs[0], strs[1]
}

// Repository represents the information of a repository.
type Repository struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
