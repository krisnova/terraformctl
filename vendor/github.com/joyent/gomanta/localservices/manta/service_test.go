//
// gomanta - Go library to interact with Joyent Manta
//
// Manta double testing service - internal direct API test
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
	"strings"
	"testing"

	gc "launchpad.net/gocheck"

	"encoding/json"

	lm "github.com/joyent/gomanta/localservices/manta"
	"github.com/joyent/gomanta/manta"
)

type MantaSuite struct {
	service *lm.Manta
}

const (
	testServiceURL  = "https://go-test.manta.joyent.com"
	testUserAccount = "gouser"
	object          = "1. Go Test -- Go Test -- GoTest\n2. Go Test -- Go Test -- GoTest"
)

var _ = gc.Suite(&MantaSuite{})

func Test(t *testing.T) {
	gc.TestingT(t)
}

func (s *MantaSuite) SetUpSuite(c *gc.C) {
	s.service = lm.New(testServiceURL, testUserAccount)
}

// Helpers
func getObject() ([]byte, error) {
	var bytes []byte

	r := strings.NewReader(object)
	if _, err := r.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func (s *MantaSuite) createDirectory(c *gc.C, path string) {
	err := s.service.PutDirectory(path)
	c.Assert(err, gc.IsNil)
}

func (s *MantaSuite) deleteDirectory(c *gc.C, path string) {
	err := s.service.DeleteDirectory(path)
	c.Assert(err, gc.IsNil)
}

func (s *MantaSuite) createObject(c *gc.C, path, objName string) {
	obj, err := getObject()
	c.Assert(err, gc.IsNil)
	err = s.service.PutObject(path, objName, obj)
	c.Assert(err, gc.IsNil)
}

func (s *MantaSuite) deleteObject(c *gc.C, path string) {
	err := s.service.DeleteObject(path)
	c.Assert(err, gc.IsNil)
}

func createTestJob(c *gc.C, jobName string) []byte {
	phases := []manta.Phase{{Type: "map", Exec: "wc", Init: ""}, {Type: "reduce", Exec: "awk '{ l += $1; w += $2; c += $3 } END { print l, w, c }'", Init: ""}}
	job, err := json.Marshal(manta.CreateJobOpts{Name: jobName, Phases: phases})
	c.Assert(err, gc.IsNil)
	return job
}

func (s *MantaSuite) createJob(c *gc.C, jobName string) string {
	jobUri, err := s.service.CreateJob(createTestJob(c, jobName))
	c.Assert(err, gc.IsNil)
	c.Assert(jobUri, gc.NotNil)
	return strings.Split(jobUri, "/")[3]
}

// Storage APIs
func (s *MantaSuite) TestPutDirectory(c *gc.C) {
	s.createDirectory(c, "test")
}

func (s *MantaSuite) TestPutDirectoryWithParent(c *gc.C) {
	s.createDirectory(c, "test")
	s.createDirectory(c, "test/innerdir")
}

func (s *MantaSuite) TestPutDirectoryNoParent(c *gc.C) {
	err := s.service.PutDirectory("nodir/test")
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/nodir was not found")
}

func (s *MantaSuite) TestListDirectoryEmpty(c *gc.C) {
	s.createDirectory(c, "empty")
	dirs, err := s.service.ListDirectory("empty", "", 0)
	c.Assert(err, gc.IsNil)
	c.Assert(dirs, gc.HasLen, 0)
}

func (s *MantaSuite) TestListDirectoryNoExists(c *gc.C) {
	_, err := s.service.ListDirectory("nodir", "", 0)
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/nodir was not found")
}

func (s *MantaSuite) TestListDirectory(c *gc.C) {
	s.createDirectory(c, "dir")
	for i := 0; i < 5; i++ {
		s.createObject(c, "dir", fmt.Sprintf("obj%d", i))
	}
	dirs, err := s.service.ListDirectory("dir", "", 0)
	c.Assert(err, gc.IsNil)
	c.Assert(dirs, gc.HasLen, 5)
}

func (s *MantaSuite) TestListDirectoryWithLimit(c *gc.C) {
	s.createDirectory(c, "limitdir")
	for i := 0; i < 500; i++ {
		s.createObject(c, "limitdir", fmt.Sprintf("obj%03d", i))
	}
	dirs, err := s.service.ListDirectory("limitdir", "", 0)
	c.Assert(err, gc.IsNil)
	c.Assert(dirs, gc.HasLen, 256)
	dirs, err = s.service.ListDirectory("limitdir", "", 10)
	c.Assert(err, gc.IsNil)
	c.Assert(dirs, gc.HasLen, 10)
	dirs, err = s.service.ListDirectory("limitdir", "", 300)
	c.Assert(err, gc.IsNil)
	c.Assert(dirs, gc.HasLen, 300)
}

func (s *MantaSuite) TestListDirectoryWithMarker(c *gc.C) {
	s.createDirectory(c, "markerdir")
	for i := 0; i < 500; i++ {
		s.createObject(c, "markerdir", fmt.Sprintf("obj%03d", i))
	}
	dirs, err := s.service.ListDirectory("markerdir", "obj400", 0)
	c.Assert(err, gc.IsNil)
	c.Assert(dirs, gc.HasLen, 100)
	c.Assert(dirs[0].Name, gc.Equals, "obj400")

}

func (s *MantaSuite) TestDeleteDirectoryNotEmpty(c *gc.C) {
	s.createDirectory(c, "notempty")
	s.createObject(c, "notempty", "obj")
	err := s.service.DeleteDirectory("notempty")
	c.Assert(err, gc.ErrorMatches, "BadRequestError")
}

func (s *MantaSuite) TestDeleteDirectory(c *gc.C) {
	s.createDirectory(c, "deletedir")
	s.deleteDirectory(c, "deletedir")
}

func (s *MantaSuite) TestDeleteDirectoryNoExists(c *gc.C) {
	s.deleteDirectory(c, "nodir")
}

func (s *MantaSuite) TestPutObject(c *gc.C) {
	s.createDirectory(c, "objdir")
	s.createObject(c, "objdir", "object")
}

func (s *MantaSuite) TestPutObjectDirectoryWithParent(c *gc.C) {
	s.createDirectory(c, "parent")
	s.createDirectory(c, "parent/objdir")
	s.createObject(c, "parent/objdir", "object")
}

func (s *MantaSuite) TestPutObjectDirectoryNoParent(c *gc.C) {
	obj, err := getObject()
	c.Assert(err, gc.IsNil)
	err = s.service.PutObject("nodir", "obj", obj)
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/nodir was not found")
}

func (s *MantaSuite) TestIsObjectTrue(c *gc.C) {
	s.createDirectory(c, "objdir")
	s.createObject(c, "objdir", "obj")
	isObj := s.service.IsObject("/gouser/stor/objdir/obj")
	c.Assert(isObj, gc.Equals, true)
}

func (s *MantaSuite) TestIsObjectFalse(c *gc.C) {
	isObj := s.service.IsObject("/gouser/stor/nodir/obj")
	c.Assert(isObj, gc.Equals, false)
}

func (s *MantaSuite) TestGetObject(c *gc.C) {
	s.createDirectory(c, "dir1")
	s.createObject(c, "dir1", "obj")
	expected, _ := getObject()
	obj, err := s.service.GetObject("dir1/obj")
	c.Assert(err, gc.IsNil)
	c.Assert(obj, gc.DeepEquals, expected)
}

func (s *MantaSuite) TestGetObjectWrongPath(c *gc.C) {
	obj, err := s.service.GetObject("nodir/obj")
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/nodir/obj was not found")
	c.Assert(obj, gc.IsNil)
}

func (s *MantaSuite) TestGetObjectWrongName(c *gc.C) {
	obj, err := s.service.GetObject("noobject")
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/noobject was not found")
	c.Assert(obj, gc.IsNil)
}

func (s *MantaSuite) TestDeleteObject(c *gc.C) {
	s.createDirectory(c, "delete")
	s.createObject(c, "delete", "obj")
	s.deleteObject(c, "delete/obj")
}

func (s *MantaSuite) TestDeleteObjectWrongPath(c *gc.C) {
	err := s.service.DeleteObject("nodir/obj")
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/nodir/obj was not found")
}

func (s *MantaSuite) TestDeleteObjectWrongName(c *gc.C) {
	err := s.service.DeleteObject("noobj")
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/noobj was not found")
}

func (s *MantaSuite) TestPutSnapLink(c *gc.C) {
	s.createDirectory(c, "linkdir")
	s.createObject(c, "linkdir", "obj")
	err := s.service.PutSnapLink("linkdir", "link", "/gouser/stor/linkdir/obj")
	c.Assert(err, gc.IsNil)
}

func (s *MantaSuite) TestPutSnapLinkNoLocation(c *gc.C) {
	s.createDirectory(c, "linkdir")
	s.createObject(c, "linkdir", "obj")
	err := s.service.PutSnapLink("linkdir", "link", "/gouser/stor/linkdir/noobj")
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/linkdir/noobj was not found")
}

func (s *MantaSuite) TestPutSnapLinkDirectoryWithParent(c *gc.C) {
	s.createDirectory(c, "link1")
	s.createDirectory(c, "link1/linkdir")
	s.createObject(c, "link1/linkdir", "obj")
	err := s.service.PutSnapLink("link1/linkdir", "link", "/gouser/stor/link1/linkdir/obj")
	c.Assert(err, gc.IsNil)
}

func (s *MantaSuite) TestPutSnapLinkDirectoryNoParent(c *gc.C) {
	s.createDirectory(c, "linkdir")
	s.createObject(c, "linkdir", "obj")
	err := s.service.PutSnapLink("nodir", "link", "/gouser/stor/linkdir/obj")
	c.Assert(err, gc.ErrorMatches, "/gouser/stor/nodir was not found")
}

// Jobs API
func (s *MantaSuite) TestCreateJob(c *gc.C) {
	s.createJob(c, "test-job")
}

func (s *MantaSuite) TestListAllJobs(c *gc.C) {
	s.createJob(c, "test-job")
	jobs, err := s.service.ListJobs(false)
	c.Assert(err, gc.IsNil)
	c.Assert(jobs, gc.NotNil)
}

func (s *MantaSuite) TestListLiveJobs(c *gc.C) {
	s.createJob(c, "test-job")
	jobs, err := s.service.ListJobs(true)
	c.Assert(err, gc.IsNil)
	c.Assert(jobs, gc.NotNil)
	//c.Assert(jobs, gc.Equals, 1)
}

func (s *MantaSuite) TestCancelJob(c *gc.C) {
	jobId := s.createJob(c, "test-job")
	err := s.service.CancelJob(jobId)
	c.Assert(err, gc.IsNil)
}

func (s *MantaSuite) TestJ05_AddJobInputs(c *gc.C) {
	var inBytes []byte
	var inputs = `/` + testUserAccount + `/stor/testjob/obj1
/` + testUserAccount + `/stor/testjob/obj2
`
	r := strings.NewReader(inputs)
	if _, err := r.Read(inBytes); err != nil {
		c.Skip("Could not read input, skipping...")
	}

	s.createDirectory(c, "testjob")
	s.createObject(c, "testjob", "obj1")
	s.createObject(c, "testjob", "obj2")

	jobId := s.createJob(c, "test-job")
	err := s.service.AddJobInputs(jobId, inBytes)
	c.Assert(err, gc.IsNil)
}

func (s *MantaSuite) TestJ06_EndJobInputs(c *gc.C) {
	jobId := s.createJob(c, "test-job")
	err := s.service.EndJobInput(jobId)
	c.Assert(err, gc.IsNil)
}

func (s *MantaSuite) TestJ07_GetJob(c *gc.C) {
	jobId := s.createJob(c, "test-job")
	job, err := s.service.GetJob(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(job, gc.NotNil)
}

func (s *MantaSuite) TestJ08_GetJobInput(c *gc.C) {
	jobId := s.createJob(c, "test-job")
	input, err := s.service.GetJobInput(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(input, gc.NotNil)
	//fmt.Println(input)
}

func (s *MantaSuite) TestJ09_GetJobOutput(c *gc.C) {
	jobId := s.createJob(c, "test-job")
	output, err := s.service.GetJobOutput(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(output, gc.NotNil)
	//fmt.Println(output)
}

func (s *MantaSuite) TestJ10_GetJobFailures(c *gc.C) {
	jobId := s.createJob(c, "test-job")
	fail, err := s.service.GetJobFailures(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(fail, gc.Equals, "")
	//fmt.Println(fail)
}

func (s *MantaSuite) TestJ11_GetJobErrors(c *gc.C) {
	jobId := s.createJob(c, "test-job")
	errs, err := s.service.GetJobErrors(jobId)
	c.Assert(err, gc.IsNil)
	c.Assert(errs, gc.IsNil)
	//fmt.Println(errs)
}
