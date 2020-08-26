package model

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

var _ TestCase = &AccuracyHttpTestCase{}

type AccuracyHttpTestCase HttpTestCase

func (tc AccuracyHttpTestCase) GetName() string {
	return tc.Name
}

func (tc AccuracyHttpTestCase) Run() TestResult {
	ctx, cancel := context.WithTimeout(context.Background(), tc.Timeout)
	defer cancel()
	url := fmt.Sprintf("http://localhost%s", tc.Path)
	req, err := http.NewRequestWithContext(ctx, tc.Method, url, bytes.NewReader(tc.Input))
	if err != nil {
		return NewTestResult(false, tc.GetName() + " request creation failed: "+err.Error())
	}
	res, err := tc.Client.Do(req)
	if err != nil {
		return NewTestResult(false, tc.GetName() + " request failed: "+err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != tc.Status {
		msg := fmt.Sprintf("status code expected %d but got %d", tc.Status, res.StatusCode)
		return NewTestResult(false, tc.GetName() + " " + msg)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return NewTestResult(false, tc.GetName() + " body read failed: "+err.Error())
	}
	if tc.Output != nil {
		if resBody == nil {
			return NewTestResult(false, tc.GetName() + " response expected, but nothing")
		}
		// TODO: Check Equality Through All Elements
		// for now, it should be encoded without additional spaces
		if string(tc.Output) != string(resBody) {
			msg := fmt.Sprintf("response expected: %s but got %s", tc.Output, resBody)
			return NewTestResult(false, tc.GetName() + " " + msg)
		}
	}
	return NewTestResult(true, tc.GetName())
}
