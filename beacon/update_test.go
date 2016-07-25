package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func testUpdate(key0, value0, key1, value1, key2, value2 string) int {
	// setup environment
	readConf()
	initLogger()
	client, _ := initRedisConnection()

	postData := url.Values{}
	postData.Set(key0, value0)
	postData.Add(key1, value1)
	postData.Add(key2, value2)

	r, _ := http.NewRequest("POST", "", strings.NewReader(postData.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	update(w, r, client)
	return w.Code
}

func TestCasesUpdate(t *testing.T) {
	// update on id that does not exist
	rCode := testUpdate("id", "test.domain", "sec", "reallivetoken", "address", "129.0.0.1")
	if rCode != 400 {
		t.Error("expected 400, got", rCode, "on update id that does not exist")
	}

	// provision id to use in later tests
	tdToken, rCode := testProvision("id", "test.domain")
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "on '/provision' to use with update tests")
	}

	// correct
	rCode = testUpdate("id", "test.domain", "sec", tdToken, "address", "129.0.0.1")
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "on correct update")
	}

	// wrong auth
	rCode = testUpdate("id", "test.domain", "sec", "nottherealtoken", "address", "129.0.0.1")
	if rCode != 403 {
		t.Error("expected 403, got", rCode, "on update. wrong auth test")
	}

	// clean up the id, we are done with it
	rCode = testRemove("id", "test.domain", "sec", tdToken)
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "while removing the id from '/update' tests")
	}
}
