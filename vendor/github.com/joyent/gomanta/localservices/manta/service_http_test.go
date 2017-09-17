//
// gomanta - Go library to interact with Joyent Manta
//
// Manta double testing service - HTTP API tests
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2016 Joyent Inc.
//
// Written by Daniele Stroppa <daniele.stroppa@joyent.com>
//

package manta_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	gc "launchpad.net/gocheck"

	"github.com/joyent/gocommon/testing"
	lm "github.com/joyent/gomanta/localservices/manta"
	"github.com/joyent/gomanta/manta"
)

type MantaHTTPSuite struct {
	testing.HTTPSuite
	service *lm.Manta
}

var _ = gc.Suite(&MantaHTTPSuite{})

type MantaHTTPSSuite struct {
	testing.HTTPSuite
	service *lm.Manta
}

var _ = gc.Suite(&MantaHTTPSSuite{HTTPSuite: testing.HTTPSuite{UseTLS: true}})

const (
	fakeStorPrefix = "/fakeuser/stor"
	fakeJobsPrefix = "/fakeuser/jobs"
)

func (s *MantaHTTPSuite) SetUpSuite(c *gc.C) {
	s.HTTPSuite.SetUpSuite(c)
	c.Assert(s.Server.URL[:7], gc.Equals, "http://")
	s.service = lm.New(s.Server.URL, "fakeuser")
}

func (s *MantaHTTPSuite) TearDownSuite(c *gc.C) {
	s.HTTPSuite.TearDownSuite(c)
}

func (s *MantaHTTPSuite) SetUpTest(c *gc.C) {
	s.HTTPSuite.SetUpTest(c)
	s.service.SetupHTTP(s.Mux)
}

func (s *MantaHTTPSuite) TearDownTest(c *gc.C) {
	s.HTTPSuite.TearDownTest(c)
}

// assertJSON asserts the passed http.Response's body can be
// unmarshalled into the given expected object, populating it with the
// successfully parsed data.
func assertJSON(c *gc.C, resp *http.Response, expected interface{}) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	c.Assert(err, gc.IsNil)
	err = json.Unmarshal(body, &expected)
	c.Assert(err, gc.IsNil)
}

// assertBody asserts the passed http.Response's body matches the
// expected response, replacing any variables in the expected body.
func assertBody(c *gc.C, resp *http.Response, expected *lm.ErrorResponse) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	c.Assert(err, gc.IsNil)
	expBody := expected.Body
	// cast to string for easier asserts debugging
	c.Assert(string(body), gc.Equals, string(expBody))
}

// sendRequest constructs an HTTP request from the parameters and
// sends it, returning the response or an error.
func (s *MantaHTTPSuite) sendRequest(method, path string, body []byte, headers http.Header) (*http.Response, error) {
	if headers == nil {
		headers = make(http.Header)
	}
	requestURL := "http://" + s.service.Hostname + strings.TrimLeft(path, "/")
	req, err := http.NewRequest(method, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Close = true
	for header, values := range headers {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}
	// workaround for https://code.google.com/p/go/issues/detail?id=4454
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	return http.DefaultClient.Do(req)
}

// jsonRequest serializes the passed body object to JSON and sends a
// the request with authRequest().
func (s *MantaHTTPSuite) jsonRequest(method, path string, body interface{}, headers http.Header) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return s.sendRequest(method, path, jsonBody, headers)
}

// Helpers
func (s *MantaHTTPSuite) createDirectory(c *gc.C, path string) {
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json; type=directory")
	resp, err := s.sendRequest("PUT", fmt.Sprintf("%s/%s", fakeStorPrefix, path), nil, reqHeaders)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusNoContent)
}

func (s *MantaHTTPSuite) createObject(c *gc.C, objPath string) {
	resp, err := s.sendRequest("PUT", fmt.Sprintf("%s/%s", fakeStorPrefix, objPath), []byte(object), nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusNoContent)
}

func (s *MantaHTTPSuite) delete(c *gc.C, path string) {
	resp, err := s.sendRequest("DELETE", fmt.Sprintf("%s/%s", fakeStorPrefix, path), nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusNoContent)
}

func getTestJob(c *gc.C, jobName string) []byte {
	phases := []manta.Phase{{Type: "map", Exec: "wc", Init: ""}, {Type: "reduce", Exec: "awk '{ l += $1; w += $2; c += $3 } END { print l, w, c }'", Init: ""}}
	opts := manta.CreateJobOpts{Name: jobName, Phases: phases}
	optsByte, err := json.Marshal(opts)
	c.Assert(err, gc.IsNil)
	return optsByte
}

func (s *MantaHTTPSuite) createJob(c *gc.C, jobName string) string {
	resp, err := s.sendRequest("POST", fakeJobsPrefix, getTestJob(c, jobName), nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusCreated)
	return strings.Split(resp.Header.Get("Location"), "/")[3]
}

// SimpleTest defines a simple request without a body and expected response.
type SimpleTest struct {
	method  string
	url     string
	headers http.Header
	expect  *lm.ErrorResponse
}

func (s *MantaHTTPSuite) simpleTests() []SimpleTest {
	var simpleTests = []SimpleTest{
		{
			method:  "GET",
			url:     "/",
			headers: make(http.Header),
			expect:  lm.ErrNotFound,
		},
		{
			method:  "POST",
			url:     "/",
			headers: make(http.Header),
			expect:  lm.ErrNotFound,
		},
		{
			method:  "DELETE",
			url:     "/",
			headers: make(http.Header),
			expect:  lm.ErrNotFound,
		},
		{
			method:  "PUT",
			url:     "/",
			headers: make(http.Header),
			expect:  lm.ErrNotFound,
		},
		{
			method:  "GET",
			url:     "/any",
			headers: make(http.Header),
			expect:  lm.ErrNotFound,
		},
		{
			method:  "POST",
			url:     "/any",
			headers: make(http.Header),
			expect:  lm.ErrNotFound,
		},
		{
			method:  "DELETE",
			url:     "/any",
			headers: make(http.Header),
			expect:  lm.ErrNotFound,
		},
		{
			method:  "PUT",
			url:     "/any",
			headers: make(http.Header),
			expect:  lm.ErrNotFound,
		},
		{
			method:  "POST",
			url:     "/fakeuser/stor",
			headers: make(http.Header),
			expect:  lm.ErrNotAllowed,
		},
		{
			method:  "DELETE",
			url:     "/fakeuser/jobs",
			headers: make(http.Header),
			expect:  lm.ErrNotAllowed,
		},
		{
			method:  "PUT",
			url:     "/fakeuser/jobs",
			headers: make(http.Header),
			expect:  lm.ErrNotAllowed,
		},
	}
	return simpleTests
}

func (s *MantaHTTPSuite) TestSimpleRequestTests(c *gc.C) {
	simpleTests := s.simpleTests()
	for i, t := range simpleTests {
		c.Logf("#%d. %s %s -> %d", i, t.method, t.url, t.expect.Code)
		if t.headers == nil {
			t.headers = make(http.Header)
		}
		var (
			resp *http.Response
			err  error
		)
		resp, err = s.sendRequest(t.method, t.url, nil, t.headers)
		c.Assert(err, gc.IsNil)
		c.Assert(resp.StatusCode, gc.Equals, t.expect.Code)
		assertBody(c, resp, t.expect)
	}
}

// Storage API
func (s *MantaHTTPSuite) TestPutDirectory(c *gc.C) {
	s.createDirectory(c, "test")
	s.delete(c, "test")
}

func (s *MantaHTTPSuite) TestPutDirectoryWithParent(c *gc.C) {
	s.createDirectory(c, "test")
	defer s.delete(c, "test")
	s.createDirectory(c, "test/innerdir")

	s.delete(c, "test/innerdir")
}

func (s *MantaHTTPSuite) TestListDirectory(c *gc.C) {
	var expected []manta.Entry
	s.createDirectory(c, "test")
	defer s.delete(c, "test")
	s.createObject(c, "test/object")
	defer s.delete(c, "test/object")

	resp, err := s.sendRequest("GET", "/fakeuser/stor/test", nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	c.Assert(expected, gc.HasLen, 1)
}

func (s *MantaHTTPSuite) TestListDirectoryWithOpts(c *gc.C) {
	var expected []manta.Entry
	s.createDirectory(c, "dir")
	defer s.delete(c, "dir")
	s.createObject(c, "dir/object1")
	defer s.delete(c, "dir/object1")
	s.createObject(c, "dir/object2")
	defer s.delete(c, "dir/object2")
	s.createObject(c, "dir/object3")
	defer s.delete(c, "dir/object3")

	opts := manta.ListDirectoryOpts{Limit: 1, Marker: "object2"}
	optsByte, errB := json.Marshal(opts)
	c.Assert(errB, gc.IsNil)
	resp, err := s.sendRequest("GET", "/fakeuser/stor/dir", optsByte, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	c.Assert(expected, gc.HasLen, 1)
}

func (s *MantaHTTPSuite) TestDeleteDirectory(c *gc.C) {
	s.createDirectory(c, "test")
	defer s.delete(c, "test")
	s.createDirectory(c, "test/innerdir")

	s.delete(c, "test/innerdir")
}

func (s *MantaHTTPSuite) TestPutObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.delete(c, "dir")
	s.createObject(c, "dir/object")

	s.delete(c, "dir/object")
}

func (s *MantaHTTPSuite) TestGetObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.delete(c, "dir")
	s.createObject(c, "dir/object")
	defer s.delete(c, "dir/object")

	resp, err := s.sendRequest("GET", "/fakeuser/stor/dir/object", nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	// TODO: assert for headers
}

func (s *MantaHTTPSuite) TestDeleteObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.delete(c, "dir")
	s.createObject(c, "dir/object")

	s.delete(c, "dir/object")
}

func (s *MantaHTTPSuite) TestPutSnaplink(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.delete(c, "dir")
	s.createObject(c, "dir/object")
	defer s.delete(c, "dir/object")
	defer s.delete(c, "dir/link")

	reqHeaders := make(http.Header)
	reqHeaders.Set("Location", "/fakeuser/stor/dir/object")
	resp, err := s.sendRequest("PUT", "/fakeuser/stor/dir/link", nil, reqHeaders)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusNoContent)
}

// Jobs API
func (s *MantaHTTPSuite) TestCreateJob(c *gc.C) {
	s.createJob(c, "test-job")
}

func (s *MantaHTTPSuite) TestListAllJobs(c *gc.C) {
	var expected []manta.Entry
	s.createJob(c, "test-job")

	resp, err := s.sendRequest("GET", fakeJobsPrefix, nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	//c.Assert(len(expected), gc.Equals, 2)
}

func (s *MantaHTTPSuite) TestListLiveJobs(c *gc.C) {
	var expected []manta.Entry
	s.createJob(c, "test-job")

	resp, err := s.sendRequest("GET", fmt.Sprintf("%s?state=running", fakeJobsPrefix), nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	//c.Assert(len(expected), gc.Equals, 1)
}

func (s *MantaHTTPSuite) TestCancelJob(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	url := fmt.Sprintf("%s/%s/live/cancel", fakeJobsPrefix, jobId)
	resp, err := s.sendRequest("POST", url, nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusAccepted)
}

func (s *MantaHTTPSuite) TestAddJobInputs(c *gc.C) {
	var inputs = `/fakeuser/stor/testjob/obj1
/fakeuser/stor/testjob/obj2
`
	s.createDirectory(c, "testjob")
	s.createObject(c, "testjob/obj1")
	s.createObject(c, "testjob/obj2")
	jobId := s.createJob(c, "test-job")

	url := fmt.Sprintf("%s/%s/live/in", fakeJobsPrefix, jobId)
	resp, err := s.sendRequest("POST", url, []byte(inputs), nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusNoContent)
}

func (s *MantaHTTPSuite) TestEndJobInputs(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	url := fmt.Sprintf("%s/%s/live/in/end", fakeJobsPrefix, jobId)
	resp, err := s.sendRequest("POST", url, nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusAccepted)
}

func (s *MantaHTTPSuite) TestGetJob(c *gc.C) {
	var expected manta.Job
	jobId := s.createJob(c, "test-job")

	url := fmt.Sprintf("%s/%s/live/status", fakeJobsPrefix, jobId)
	resp, err := s.sendRequest("GET", url, nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	c.Assert(expected.Id, gc.Equals, jobId)
	c.Assert(expected.Name, gc.Equals, "test-job")
	c.Assert(expected.Cancelled, gc.Equals, false)
	c.Assert(expected.InputDone, gc.Equals, false)
	c.Assert(expected.Phases, gc.HasLen, 2)
}

func (s *MantaHTTPSuite) TestGetJobInput(c *gc.C) {
	var expected string
	jobId := s.createJob(c, "test-job")

	url := fmt.Sprintf("%s/%s/live/in", fakeJobsPrefix, jobId)
	resp, err := s.sendRequest("GET", url, nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	fmt.Println(expected)
	c.Assert(expected, gc.HasLen, 0)
}

func (s *MantaHTTPSuite) TestGetJobOutput(c *gc.C) {
	var expected string
	jobId := s.createJob(c, "test-job")

	url := fmt.Sprintf("%s/%s/live/out", fakeJobsPrefix, jobId)
	resp, err := s.sendRequest("GET", url, nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	c.Assert(strings.Split(expected, "\n"), gc.HasLen, 1)
}

func (s *MantaHTTPSuite) TestGetJobFailures(c *gc.C) {
	var expected string
	jobId := s.createJob(c, "test-job")

	url := fmt.Sprintf("%s/%s/live/fail", fakeJobsPrefix, jobId)
	resp, err := s.sendRequest("GET", url, nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	c.Assert(expected, gc.Equals, "")
}

func (s *MantaHTTPSuite) TestGetJobErrors(c *gc.C) {
	var expected []manta.JobError
	jobId := s.createJob(c, "test-job")

	url := fmt.Sprintf("%s/%s/live/err", fakeJobsPrefix, jobId)
	resp, err := s.sendRequest("GET", url, nil, nil)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
	assertJSON(c, resp, &expected)
	c.Assert(expected, gc.HasLen, 0)
}
