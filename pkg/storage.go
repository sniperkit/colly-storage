// Copyright 2018 Adam Tauber
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

// Storage is an interface which handles...
type Storage interface {
	// Init initializes the storage instance
	Init() error
	// Get returns the []byte representation of an entry and a bool set to true if the value isn't empty
	Get(key string) (responseBytes []byte, ok bool)
	// Set stores the []byte representation of  an entry against a key
	Set(key string, responseBytes []byte) error
	// Delete removes the entry value associated with the key
	Delete(key string) error
	// Debug
	Debug(action string) error
	// Clear
	Clear() error
	//  Action
	Action(name string, args ...interface{}) (map[string]*interface{}, error)
}

// Visited receives and stores a request ID that is visited by the Collector
// Visited(requestID uint64) error
// IsVisited returns true if the request was visited before IsVisited
// is called
// IsVisited(requestID uint64) (bool, error)
// Cookies retrieves stored cookies for a given host
// Cookies(u *url.URL) string
// SetCookies stores cookies for a given host
// SetCookies(u *url.URL, cookies string)

// Get returns the stored entry if present, and nil otherwise.
func Get(s Storage, query string) (output []byte, err error) {
	outputStr, ok := s.Get(storageQuery(query))
	if !ok {
		return []byte{}, nil
	}
	return output, err
	// b := bytes.NewBuffer(outputStr)
	// return http.ReadResponse(bufio.NewReader(b), req)
}

// storageQuery returns the entry for the query.
func storageQuery(query string) string {
	return ""
}
