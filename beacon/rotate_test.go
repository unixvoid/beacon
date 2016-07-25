package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func testRotate(key0, value0, key1, value1 string) (string, int) {
	// setup environment
	readConf()
	initLogger()
	client, _ := initRedisConnection()

	postData := url.Values{}
	postData.Set(key0, value0)
	postData.Add(key1, value1)

	r, _ := http.NewRequest("POST", "", strings.NewReader(postData.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	rotate(w, r, client, "tmp")
	if w.Code == 200 {
		return w.Body.String(), w.Code
	} else {
		return "", w.Code
	}
}

func TestCasesRotate(t *testing.T) {
	// rotate on id that does not exist
	_, rCode := testRotate("id", "test.domain", "sec", "reallivetoken")
	if rCode != 400 {
		t.Error("expected 400, got", rCode, "on rotate id that does not exist")
	}

	// provision id to use in later tests
	tdToken, rCode := testProvision("id", "test.domain")
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "on '/provision' to use with rotate tests")
	}

	// correct
	newTdToken, rCode := testRotate("id", "test.domain", "sec", tdToken)
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "on correct rotate")
	}

	// wrong auth
	_, rCode = testRotate("id", "test.domain", "sec", "nottherealtoken")
	if rCode != 403 {
		t.Error("expected 403, got", rCode, "on rotate. wrong auth test")
	}

	// clean up the id, we are done with it
	rCode = testRemove("id", "test.domain", "sec", newTdToken)
	if rCode != 200 {
		t.Error("expected 200, got", rCode, "while removing the id from '/rotate' tests")
	}
}
