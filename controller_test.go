package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// Tests valid URI
func TestCheckUriValid(t *testing.T) {
	// Valid test cases
	testTable := []struct {
		name     string
		uri      string
		expected error
	}{
		{
			name:     "without www",
			uri:      "/projectinfo/v1/github.com/apache/kafka",
			expected: error(nil),
		},
		{
			name:     "with www",
			uri:      "/projectinfo/v1/www.github.com/apache/kafka",
			expected: error(nil),
		},
	}
	// Run all test cases
	for _, testCase := range testTable {
		err := checkURI(testCase.uri)
		if err != testCase.expected {
			t.Errorf("Case: %v - Invalid URI.", testCase.name)
		}
	}
}

// Tests invalid URI
func TestCheckUriInvalid(t *testing.T) {
	// Invalid test cases
	testTable := []struct {
		name string
		uri  string
	}{
		{
			name: "invalid projectinfo",
			uri:  "/projectinf/v1/github.com/apache/kafka",
		},
		{
			name: "invalid v1",
			uri:  "/projectinfo/v2/github.com/apache/kafka",
		},
		{
			name: "invalid github",
			uri:  "/projectinfo/v1/bitbucket.org/apache/kafka",
		},
		{
			name: "too few path variables",
			uri:  "/projectinfo/v1/github.com/apache",
		},
		{
			name: "too many path variables",
			uri:  "/projectinfo/v1/github.com/apache/kafka/kafka",
		},
	}
	// Run all test cases
	for _, testCase := range testTable {
		err := checkURI(testCase.uri)
		if err == nil {
			t.Errorf("Case: %v - Should return error. Returned nil instead.", testCase.name)
		}
	}
}

// Tests function to produce base URI for GitHub API
func TestProduceBaseUri(t *testing.T) {
	owner := "testOwner"
	project := "testProject"
	expected := "https://api.github.com/repos/testOwner/testProject"
	result := produceBaseURI(owner, project)
	if result != expected {
		t.Errorf("URI's do not match. \nExpected:\n%v\nGot:\n%v", expected, result)
	}
}

// Tests function that fetches top committer
func TestFetchCommitter(t *testing.T) {
	// Mocks an HTTP Response
	resp := http.Response{}
	resp.Body = ioutil.NopCloser(strings.NewReader(`[
	  {
		"login": "ijuma",
		"id": 24747,
		"avatar_url": "https://avatars2.githubusercontent.com/u/24747?v=4",
		"contributions": 343
	  },
	  {
		"login": "hachikuji",
		"id": 12502538,
		"avatar_url": "https://avatars3.githubusercontent.com/u/12502538?v=4",
		"contributions": 261
	  }]`))

	// Returns only first array element
	contributor := fetchCommitter(resp)
	if contributor.Committer != "ijuma" {
		t.Errorf("Committer differs. Expected 'ijua'. Got %v.", contributor.Committer)
	}
	if contributor.Commits != 343 {
		t.Errorf("Commits differs. Expected 343. Got %d.", contributor.Commits)
	}
}

// Tests function that fetches languages
func TestFetchLanguages(t *testing.T) {
	// Mocks an HTTP Response
	resp := http.Response{}
	resp.Body = ioutil.NopCloser(strings.NewReader(`{
	  "Java": 11536583,
	  "Scala": 5165565,
	  "Python": 658053,
	  "Shell": 86324,
	  "Batchfile": 27518,
	  "XSLT": 7116,
	  "HTML": 5443
	}`))
	testObject := GitObject{}
	// Fetches languages to testObject
	fetchLanguages(resp, &testObject)
	// Number of languages should be 7
	if len(testObject.Languages) != 7 {
		t.Errorf("Number of languages differs. Expected 7. Got %d.", len(testObject.Languages))
	}
}
