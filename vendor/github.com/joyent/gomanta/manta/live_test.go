//
// gomanta - Go library to interact with Joyent Manta
//
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
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joyent/gocommon/client"
	"github.com/joyent/gomanta/manta"
	"github.com/joyent/gosign/auth"
	gc "launchpad.net/gocheck"
	"strings"
)

func registerMantaTests(creds *auth.Credentials) {
	gc.Suite(&LiveTests{creds: creds})
}

type LiveTests struct {
	creds      *auth.Credentials
	testClient *manta.Client
}

func (s *LiveTests) SetUpTest(c *gc.C) {
	client := client.NewClient(s.creds.MantaEndpoint.URL, "", s.creds, log.New(os.Stderr, "", log.LstdFlags))
	c.Assert(client, gc.NotNil)
	s.testClient = manta.New(client)
	c.Assert(s.testClient, gc.NotNil)
}

// Helper method to create a test directory
func (s *LiveTests) createDirectory(c *gc.C, path string) {
	err := s.testClient.PutDirectory(path)
	c.Assert(err, gc.IsNil)
}

// Helper method to create a test object
func (s *LiveTests) createObject(c *gc.C, path, objName string) {
	err := s.testClient.PutObject(path, objName, []byte("Test Manta API"))
	c.Assert(err, gc.IsNil)
}

// Helper method to delete a test directory
func (s *LiveTests) deleteDirectory(c *gc.C, path string) {
	err := s.testClient.DeleteDirectory(path)
	c.Assert(err, gc.IsNil)
}

// Helper method to delete a test object
func (s *LiveTests) deleteObject(c *gc.C, path, objName string) {
	err := s.testClient.DeleteObject(path, objName)
	c.Assert(err, gc.IsNil)
}

// Helper method to create a test job
func (s *LiveTests) createJob(c *gc.C, jobName string) string {
	phases := []manta.Phase{
		{Type: "map", Exec: "wc", Init: ""},
		{Type: "reduce", Exec: "awk '{ l += $1; w += $2; c += $3 } END { print l, w, c }'", Init: ""},
	}
	jobUri, err := s.testClient.CreateJob(manta.CreateJobOpts{Name: jobName, Phases: phases})
	c.Assert(err, gc.IsNil)
	c.Assert(jobUri, gc.NotNil)
	return strings.Split(jobUri, "/")[3]
}

// Storage API
func (s *LiveTests) TestPutDirectory(c *gc.C) {
	s.createDirectory(c, "test")

	// cleanup
	s.deleteDirectory(c, "test")
}

func (s *LiveTests) TestListDirectory(c *gc.C) {
	s.createDirectory(c, "test")
	defer s.deleteDirectory(c, "test")
	s.createObject(c, "test", "obj")
	defer s.deleteObject(c, "test", "obj")

	opts := manta.ListDirectoryOpts{}
	dirs, err := s.testClient.ListDirectory("test", opts)
	c.Assert(err, gc.IsNil)
	c.Assert(dirs, gc.NotNil)
}

func (s *LiveTests) TestDeleteDirectory(c *gc.C) {
	s.createDirectory(c, "test")

	s.deleteDirectory(c, "test")
}

func (s *LiveTests) TestPutObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.deleteDirectory(c, "dir")
	s.createObject(c, "dir", "obj")
	defer s.deleteObject(c, "dir", "obj")
}

func (s *LiveTests) TestGetObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.deleteDirectory(c, "dir")
	s.createObject(c, "dir", "obj")
	defer s.deleteObject(c, "dir", "obj")

	obj, err := s.testClient.GetObject("dir", "obj")
	c.Assert(err, gc.IsNil)
	c.Assert(obj, gc.NotNil)
	c.Check(string(obj), gc.Equals, "Test Manta API")
}

func (s *LiveTests) TestDeleteObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.deleteDirectory(c, "dir")
	s.createObject(c, "dir", "obj")

	s.deleteObject(c, "dir", "obj")
}

func (s *LiveTests) TestPutSnapLink(c *gc.C) {
	s.createDirectory(c, "linkdir")
	defer s.deleteDirectory(c, "linkdir")
	s.createObject(c, "linkdir", "obj")
	defer s.deleteObject(c, "linkdir", "obj")

	location := fmt.Sprintf("%s/%s", s.creds.UserAuthentication.User, "stor/linkdir/obj")
	err := s.testClient.PutSnapLink("linkdir", "objlnk", location)
	c.Assert(err, gc.IsNil)

	// cleanup
	s.deleteObject(c, "linkdir", "objlnk")
}

func (s *LiveTests) TestSignURL(c *gc.C) {
	s.createDirectory(c, "sign")
	defer s.deleteDirectory(c, "sign")
	s.createObject(c, "sign", "object")
	defer s.deleteObject(c, "sign", "object")

	location := fmt.Sprintf("/%s/%s", s.creds.UserAuthentication.User, "stor/sign/object")
	url, err := s.testClient.SignURL(location, time.Now().Add(time.Minute*5))
	c.Assert(err, gc.IsNil)
	c.Assert(url, gc.Not(gc.Equals), "")

	resp, err := http.Get(url)
	c.Assert(err, gc.IsNil)
	c.Assert(resp.StatusCode, gc.Equals, http.StatusOK)
}

// Jobs API
func (s *LiveTests) TestCreateJob(c *gc.C) {
	s.createJob(c, "test-job")
}

func (s *LiveTests) TestListLiveJobs(c *gc.C) {
	s.createJob(c, "test-job")

	jobs, err := s.testClient.ListJobs(true)
	c.Assert(err, gc.IsNil)
	c.Assert(jobs, gc.NotNil)
	c.Assert(len(jobs) >= 1, gc.Equals, true)
}

func (s *LiveTests) TestListAllJobs(c *gc.C) {
	s.createJob(c, "test-job")

	jobs, err := s.testClient.ListJobs(false)
	c.Assert(err, gc.IsNil)
	c.Assert(jobs, gc.NotNil)
}

func (s *LiveTests) TestCancelJob(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	err := s.testClient.CancelJob(jobId)
	c.Assert(err, gc.IsNil)
}

func (s *LiveTests) TestAddJobInputs(c *gc.C) {
	var inputs = `/` + s.creds.UserAuthentication.User + `/stor/testjob/obj1
/` + s.creds.UserAuthentication.User + `/stor/testjob/obj2
`
	s.createDirectory(c, "testjob")
	defer s.deleteDirectory(c, "testjob")
	s.createObject(c, "testjob", "obj1")
	defer s.deleteObject(c, "testjob", "obj1")
	s.createObject(c, "testjob", "obj2")
	defer s.deleteObject(c, "testjob", "obj2")
	jobId := s.createJob(c, "test-job")

	err := s.testClient.AddJobInputs(jobId, strings.NewReader(inputs))
	c.Assert(err, gc.IsNil)
}

func (s *LiveTests) TestEndJobInputs(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	err := s.testClient.EndJobInputs(jobId)
	c.Assert(err, gc.IsNil)
}

func (s *LiveTests) TestGetJob(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	job, err := s.testClient.GetJob(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(job, gc.NotNil)
}

func (s *LiveTests) TestGetJobInput(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	input, err := s.testClient.GetJobInput(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(input, gc.NotNil)
}

func (s *LiveTests) TestGetJobOutput(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	output, err := s.testClient.GetJobOutput(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(output, gc.NotNil)
}

func (s *LiveTests) TestGetJobFailures(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	fail, err := s.testClient.GetJobFailures(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(fail, gc.Equals, nil)
}

func (s *LiveTests) TestGetJobErrors(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	errs, err := s.testClient.GetJobErrors(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(errs, gc.IsNil)
}
