package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

// Checks if URI has correct structure
func checkURI(uri string) error {
	uriParts := strings.Split(uri, "/")
	if len(uriParts) == 6 && uriParts[1] == "projectinfo" && uriParts[2] == "v1" &&
		(uriParts[3] == "github.com" || uriParts[3] == "www.github.com") {
		return nil
	}
	return errors.New("invalid URI")
}

// Produces base URI for GitHub API
func produceBaseURI(owner, project string) string {
	return "https://api.github.com/repos/" + owner + "/" + project
}

/*
Fetches the top committer and number of commits made.
Takes a slice to be able to parse array of JSON objects.
*/
func fetchCommitter(resp http.Response) Committer {
	committers := make([]Committer, 0)
	// Decodes slice of objects
	err := json.NewDecoder(resp.Body).Decode(&committers)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}
	// Returns only the first object (top committer)
	return committers[0]
}

// Fetches languages to object
func fetchLanguages(resp http.Response, object *GitObject) {
	// Key/Value array to hold languages
	languages := map[string]int{}
	json.NewDecoder(resp.Body).Decode(&languages)
	// Ranges over all languages, and puts keys in object
	for key := range languages {
		object.Languages = append(object.Languages, key)
	}
}

func gitHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Checks if URI has correct structure
	err := checkURI(r.URL.String())
	if err != nil {
		fmt.Fprintln(w, err, http.StatusBadRequest)
	}

	// Fetch 'owner' and 'project' from URI and put in GitObject
	gitObj := GitObject{}
	gitObj.Owner = p.ByName("owner")
	gitObj.Project = p.ByName("project")

	// Produce base URI for GitHub API
	baseURI := produceBaseURI(gitObj.Owner, gitObj.Project)

	// Produce URI and fetch 'committer' and 'commits'
	committerURI := baseURI + "/contributors"
	resp, err := http.Get(committerURI)
	if err != nil {
		fmt.Fprintln(w, "Could not retrieve contributors page.", http.StatusNotFound)
	}
	committer := fetchCommitter(*resp)
	gitObj.Committer = committer.Committer
	gitObj.Commits = committer.Commits

	// Produce URI and fetch 'languages'
	languageURI := baseURI + "/languages"
	resp, err = http.Get(languageURI)
	if err != nil {
		fmt.Fprintln(w, "Could not retrieve languages page.", http.StatusNotFound)
	}
	fetchLanguages(*resp, &gitObj)

	// Marshal GitObject to JSON and respond with new JSON object
	jsonObj, err := json.Marshal(gitObj)
	if err != nil {
		fmt.Fprintln(w, "Error marshaling object to JSON.", http.StatusConflict)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonObj)
}
