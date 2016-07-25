package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func testRemove(key0, value0, key1, value1 string) int {
	// setup environment
	readConf()
	initLogger()
	client, _ := initRedisConnection()

	postData := url.Values{}
	postData.Set(key0, value0)
	postData.Set(key1, value1)
	r, _ := http.NewRequest("POST", "", strings.NewReader(postData.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	remove(w, r, client)
	return w.Code
}

func TestCasesRemove(t *testing.T) {
	// id does not exist
	rCode := testRemove("id", "test.domain", "sec", "wrongsectoken")
	if rCode != 400 {
		t.Error("expected 400, user should not exist yet")
	}

	// first provision an id to test with
	tdToken, rCode := testProvision("id", "test.domain")
	if rCode != 200 {
		t.Error("Error provisioning id to use with '/remove' tests")
	}

	// wrong auth
	rCode = testRemove("id", "test.domain", "sec", "wrongsectoken")
	if rCode != 403 {
		t.Error("expected 403 to be returned when auth is wrong")
	}

	// correct
	rCode = testRemove("id", "test.domain", "sec", tdToken)
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "remove id failed")
	}
}
