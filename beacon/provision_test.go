package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func testProvision(key, value string) (string, int) {
	// setup environment
	readConf()
	initLogger()
	client, _ := initRedisConnection()

	postData := url.Values{}
	postData.Set(key, value)
	r, _ := http.NewRequest("POST", "", strings.NewReader(postData.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	provision(w, r, client, "tmp")
	if w.Code == 200 {
		return w.Body.String(), w.Code
	} else {
		return "", w.Code
	}
}

func TestCasesProvision(t *testing.T) {
	// correct
	tdToken, rCode := testProvision("id", "test.domain")
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "on first request, id already exists")
	}
	// already exists
	_, rCode = testProvision("id", "test.domain")
	if rCode != 400 {
		t.Error("expected 400, got", rCode, ".token should exist already")
	}

	// already exists
	_, rCode = testProvision("", "test.domain")
	if rCode != 400 {
		t.Error("expected 400, got", rCode, "to be returned when the 'id' is not set")
	}

	// remove id, we are done with it
	rCode = testRemove("id", "test.domain", "sec", tdToken)
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "trying to remove token after provisioning tests")
	}
}
