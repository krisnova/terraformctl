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
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	gc "launchpad.net/gocheck"

	"github.com/joyent/gocommon/client"
	localmanta "github.com/joyent/gomanta/localservices/manta"
	"github.com/joyent/gomanta/manta"
	"github.com/joyent/gosign/auth"
)

var privateKey []byte

func registerLocalTests(keyName string) {
	var localKeyFile string
	if keyName == "" {
		localKeyFile = os.Getenv("HOME") + "/.ssh/id_rsa"
	} else {
		localKeyFile = keyName
	}
	privateKey, _ = ioutil.ReadFile(localKeyFile)

	gc.Suite(&LocalTests{})
}

type LocalTests struct {
	LiveTests
	Server     *httptest.Server
	Mux        *httprouter.Router
	oldHandler http.Handler
	manta      *localmanta.Manta
}

func (s *LocalTests) SetUpSuite(c *gc.C) {
	// Set up the HTTP server.
	s.Server = httptest.NewServer(nil)
	s.oldHandler = s.Server.Config.Handler
	s.Mux = httprouter.New()
	s.Server.Config.Handler = s.Mux

	// Set up a Joyent Manta service.
	authentication, err := auth.NewAuth("localtest", string(privateKey), "rsa-sha256")
	c.Assert(err, gc.IsNil)

	s.creds = &auth.Credentials{
		UserAuthentication: authentication,
		MantaKeyId:         "",
		MantaEndpoint:      auth.Endpoint{URL: s.Server.URL},
	}
	s.manta = localmanta.New(s.creds.MantaEndpoint.URL, s.creds.UserAuthentication.User)
	s.manta.SetupHTTP(s.Mux)
}

func (s *LocalTests) TearDownSuite(c *gc.C) {
	s.Mux = nil
	s.Server.Config.Handler = s.oldHandler
	s.Server.Close()
}

func (s *LocalTests) SetUpTest(c *gc.C) {
	client := client.NewClient(s.creds.MantaEndpoint.URL, "", s.creds, log.New(os.Stderr, "", log.LstdFlags))
	c.Assert(client, gc.NotNil)
	s.testClient = manta.New(client)
	c.Assert(s.testClient, gc.NotNil)
}

// Helper method to create a test directory
func (s *LocalTests) createDirectory(c *gc.C, path string) {
	err := s.testClient.PutDirectory(path)
	c.Assert(err, gc.IsNil)
}

// Helper method to create a test object
func (s *LocalTests) createObject(c *gc.C, path, objName string) {
	err := s.testClient.PutObject(path, objName, []byte("Test Manta API"))
	c.Assert(err, gc.IsNil)
}

// Helper method to delete a test directory
func (s *LocalTests) deleteDirectory(c *gc.C, path string) {
	err := s.testClient.DeleteDirectory(path)
	c.Assert(err, gc.IsNil)
}

// Helper method to delete a test object
func (s *LocalTests) deleteObject(c *gc.C, path, objName string) {
	err := s.testClient.DeleteObject(path, objName)
	c.Assert(err, gc.IsNil)
}

// Helper method to create a test job
func (s *LocalTests) createJob(c *gc.C, jobName string) string {
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
func (s *LocalTests) TestPutDirectory(c *gc.C) {
	s.createDirectory(c, "test")

	// cleanup
	s.deleteDirectory(c, "test")
}

func (s *LocalTests) TestListDirectory(c *gc.C) {
	s.createDirectory(c, "test")
	defer s.deleteDirectory(c, "test")
	s.createObject(c, "test", "obj")
	defer s.deleteObject(c, "test", "obj")

	opts := manta.ListDirectoryOpts{}
	dirs, err := s.testClient.ListDirectory("test", opts)
	c.Assert(err, gc.IsNil)
	c.Assert(dirs, gc.NotNil)
}

func (s *LocalTests) TestDeleteDirectory(c *gc.C) {
	s.createDirectory(c, "test")
	s.deleteDirectory(c, "test")
}

func (s *LocalTests) TestPutObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.deleteDirectory(c, "dir")
	s.createObject(c, "dir", "obj")
	defer s.deleteObject(c, "dir", "obj")
}

func (s *LocalTests) TestGetObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.deleteDirectory(c, "dir")
	s.createObject(c, "dir", "obj")
	defer s.deleteObject(c, "dir", "obj")

	obj, err := s.testClient.GetObject("dir", "obj")
	c.Assert(err, gc.IsNil)
	c.Assert(obj, gc.NotNil)
	c.Check(string(obj), gc.Equals, "Test Manta API")
}

func (s *LocalTests) TestDeleteObject(c *gc.C) {
	s.createDirectory(c, "dir")
	defer s.deleteDirectory(c, "dir")
	s.createObject(c, "dir", "obj")

	s.deleteObject(c, "dir", "obj")
}

func (s *LocalTests) TestPutSnapLink(c *gc.C) {
	s.createDirectory(c, "linkdir")
	defer s.deleteDirectory(c, "linkdir")
	s.createObject(c, "linkdir", "obj")
	defer s.deleteObject(c, "linkdir", "obj")

	location := fmt.Sprintf("/%s/%s", s.creds.UserAuthentication.User, "stor/linkdir/obj")
	err := s.testClient.PutSnapLink("linkdir", "objlnk", location)
	c.Assert(err, gc.IsNil)

	// cleanup
	s.deleteObject(c, "linkdir", "objlnk")
}

// Jobs API
func (s *LocalTests) TestCreateJob(c *gc.C) {
	s.createJob(c, "test-job")
}

func (s *LocalTests) TestListLiveJobs(c *gc.C) {
	s.createJob(c, "test-job")

	jobs, err := s.testClient.ListJobs(true)
	c.Assert(err, gc.IsNil)
	c.Assert(jobs, gc.NotNil)
	c.Assert(len(jobs) >= 1, gc.Equals, true)
}

func (s *LocalTests) TestListAllJobs(c *gc.C) {
	s.createJob(c, "test-job")

	jobs, err := s.testClient.ListJobs(false)
	c.Assert(err, gc.IsNil)
	c.Assert(jobs, gc.NotNil)
}

func (s *LocalTests) TestCancelJob(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	err := s.testClient.CancelJob(jobId)
	c.Assert(err, gc.IsNil)
}

func (s *LocalTests) TestAddJobInputs(c *gc.C) {
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

func (s *LocalTests) TestEndJobInputs(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	err := s.testClient.EndJobInputs(jobId)
	c.Assert(err, gc.IsNil)
}

func (s *LocalTests) TestGetJob(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	job, err := s.testClient.GetJob(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(job, gc.NotNil)
}

func (s *LocalTests) TestGetJobInput(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	input, err := s.testClient.GetJobInput(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(input, gc.NotNil)
}

func (s *LocalTests) TestGetJobOutput(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	output, err := s.testClient.GetJobOutput(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(output, gc.NotNil)
}

func (s *LocalTests) TestGetJobFailures(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	fail, err := s.testClient.GetJobFailures(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(fail, gc.Equals, "")
}

func (s *LocalTests) TestGetJobErrors(c *gc.C) {
	jobId := s.createJob(c, "test-job")

	errs, err := s.testClient.GetJobErrors(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(errs, gc.IsNil)
}
