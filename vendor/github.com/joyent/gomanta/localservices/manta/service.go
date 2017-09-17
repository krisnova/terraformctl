//
// gomanta - Go library to interact with Joyent Manta
//
// Manta double testing service - internal direct API implementation
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2016 Joyent Inc.
//
// Written by Daniele Stroppa <daniele.stroppa@joyent.com>
//

package manta

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/joyent/gomanta/localservices"
	"github.com/joyent/gomanta/manta"
)

const (
	storagePrefix = "/%s/stor/%s"
	jobsPrefix    = "/%s/jobs/%s"
	separator     = "/"
	typeDirectory = "directory"
	typeObject    = "object"
)

type Manta struct {
	localservices.ServiceInstance
	mu          sync.Mutex // protects access to the following fields
	objects     map[string]manta.Entry
	objectsData map[string][]byte
	jobs        map[string]*manta.Job
}

func New(serviceURL, userAccount string) *Manta {
	URL, err := url.Parse(serviceURL)
	if err != nil {
		panic(err)
	}
	hostname := URL.Host
	if !strings.HasSuffix(hostname, separator) {
		hostname += separator
	}

	mantaDirectories := make(map[string]manta.Entry)

	path := fmt.Sprintf("/%s", userAccount)
	mantaDirectories[path] = createDirectory(userAccount)
	path = fmt.Sprintf("/%s/stor", userAccount)
	mantaDirectories[path] = createDirectory("stor")
	path = fmt.Sprintf("/%s/jobs", userAccount)
	mantaDirectories[path] = createDirectory("jobs")

	mantaService := &Manta{
		objects:     mantaDirectories,
		objectsData: make(map[string][]byte),
		jobs:        make(map[string]*manta.Job),
		ServiceInstance: localservices.ServiceInstance{
			Scheme:      URL.Scheme,
			Hostname:    hostname,
			UserAccount: userAccount,
		},
	}

	return mantaService
}

func createDirectory(directoryName string) manta.Entry {
	return manta.Entry{
		Name:  directoryName,
		Type:  typeDirectory,
		Mtime: time.Now().Format(time.RFC3339),
	}
}

func createJobObject(objName string, objData []byte) (manta.Entry, error) {
	etag, err := localservices.NewUUID()
	return manta.Entry{
		Name:  objName,
		Type:  typeObject,
		Mtime: time.Now().Format(time.RFC3339),
		Etag:  etag,
		Size:  len(objData),
	}, err
}

func (m *Manta) IsObject(name string) bool {
	m.mu.Lock()
	_, exist := m.objectsData[name]
	m.mu.Unlock()
	return exist
}

func (m *Manta) IsDirectory(name string) bool {
	m.mu.Lock()
	_, exist := m.objects[name]
	m.mu.Unlock()
	return !m.IsObject(name) && exist
}

// Directories APIs
func (m *Manta) ListDirectory(path, marker string, limit int) ([]manta.Entry, error) {
	if err := m.ProcessFunctionHook(m, path, marker, limit); err != nil {
		return nil, err
	}

	realPath := fmt.Sprintf(storagePrefix, m.ServiceInstance.UserAccount, path)

	if limit == 0 {
		limit = 256
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.objects[realPath]; !ok {
		return nil, fmt.Errorf("%s was not found", realPath)
	}

	var sortedKeys []string
	for k := range m.objects {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	if !strings.HasSuffix(realPath, separator) {
		realPath = realPath + separator
	}

	var entries []manta.Entry
	var entriesKeys []string
sortedLoop:
	for _, key := range sortedKeys {
		if strings.Contains(key, realPath) {
			for _, k := range entriesKeys {
				if strings.Contains(key, k) {
					continue sortedLoop
				}
			}
			entriesKeys = append(entriesKeys, key)
		}
	}

	for _, k := range entriesKeys {
		if marker != "" && marker > k[strings.LastIndex(k, "/")+1:] {
			continue
		}
		entries = append(entries, m.objects[k])
		if len(entries) >= limit {
			break
		}
	}

	return entries, nil
}

func getParentDirs(userAccount, path string) []string {
	var parents []string

	tokens := strings.Split(path, separator)
	for index, _ := range tokens {
		parents = append(parents, fmt.Sprintf(storagePrefix, userAccount, strings.Join(tokens[:(index+1)], separator)))
	}

	return parents
}

func (m *Manta) PutDirectory(path string) error {
	if err := m.ProcessFunctionHook(m, path); err != nil {
		return err
	}

	realPath := fmt.Sprintf(storagePrefix, m.ServiceInstance.UserAccount, path)

	// Check if parent dirs exist
	m.mu.Lock()
	defer m.mu.Unlock()
	if strings.Contains(path, separator) {
		ppath := path[:strings.LastIndex(path, separator)]
		parents := getParentDirs(m.ServiceInstance.UserAccount, ppath)
		for _, p := range parents {
			if _, ok := m.objects[p]; !ok {
				return fmt.Errorf("%s was not found", p)
			}
		}
	}

	dir := manta.Entry{
		Name:  path[(strings.LastIndex(path, separator) + 1):],
		Type:  typeDirectory,
		Mtime: time.Now().Format(time.RFC3339),
	}

	m.objects[realPath] = dir

	return nil
}

func (m *Manta) DeleteDirectory(path string) error {
	if err := m.ProcessFunctionHook(m, path); err != nil {
		return err
	}

	realPath := fmt.Sprintf(storagePrefix, m.ServiceInstance.UserAccount, path)

	// Check if empty
	ppath := realPath + separator
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, _ := range m.objects {
		if strings.Contains(k, ppath) {
			return ErrBadRequest
		}
	}

	delete(m.objects, realPath)

	return nil
}

// Objects APIs
func (m *Manta) PutObject(path, objName string, objData []byte) error {
	if err := m.ProcessFunctionHook(m, path, objName, objData); err != nil {
		return err
	}

	realPath := fmt.Sprintf(storagePrefix, m.ServiceInstance.UserAccount, path)

	// Check if parent dirs exist
	m.mu.Lock()
	defer m.mu.Unlock()
	parents := getParentDirs(m.ServiceInstance.UserAccount, path)
	for _, p := range parents {
		if _, ok := m.objects[p]; !ok {
			return fmt.Errorf("%s was not found", realPath)
		}
	}

	etag, err := localservices.NewUUID()
	if err != nil {
		return err
	}

	obj := manta.Entry{
		Name:  objName,
		Type:  typeObject,
		Mtime: time.Now().Format(time.RFC3339),
		Etag:  etag,
		Size:  len(objData),
	}

	objId := fmt.Sprintf("%s/%s", realPath, objName)
	m.objects[objId] = obj
	m.objectsData[objId] = objData

	return nil
}

func (m *Manta) GetObject(objPath string) ([]byte, error) {
	if err := m.ProcessFunctionHook(m, objPath); err != nil {
		return nil, err
	}

	objId := fmt.Sprintf(storagePrefix, m.ServiceInstance.UserAccount, objPath)
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.objects[objId]; ok {
		// TODO: Headers!
		return m.objectsData[objId], nil
	}

	return nil, fmt.Errorf("%s was not found", objId)
}

func (m *Manta) DeleteObject(objPath string) error {
	if err := m.ProcessFunctionHook(m, objPath); err != nil {
		return err
	}

	objId := fmt.Sprintf(storagePrefix, m.ServiceInstance.UserAccount, objPath)
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.objects[objId]; ok {
		delete(m.objects, objId)
		delete(m.objectsData, objId)

		return nil
	}
	return fmt.Errorf("%s was not found", objId)
}

// Link APIs
func (m *Manta) PutSnapLink(path, linkName, location string) error {
	if err := m.ProcessFunctionHook(m, path, linkName, location); err != nil {
		return err
	}

	realPath := fmt.Sprintf(storagePrefix, m.ServiceInstance.UserAccount, path)

	// Check if parent dirs exist
	m.mu.Lock()
	defer m.mu.Unlock()
	parents := getParentDirs(m.ServiceInstance.UserAccount, path)
	for _, p := range parents {
		if _, ok := m.objects[p]; !ok {
			return fmt.Errorf("%s was not found", realPath)
		}
	}

	// Check if location exist
	if _, ok := m.objects[location]; !ok {
		return fmt.Errorf("%s was not found", location)
	}

	etag, err := localservices.NewUUID()
	if err != nil {
		return err
	}

	obj := manta.Entry{
		Name:  linkName,
		Type:  typeObject,
		Mtime: time.Now().Format(time.RFC3339),
		Etag:  etag,
		Size:  len(m.objectsData[location]),
	}

	objId := fmt.Sprintf("%s/%s", realPath, linkName)
	m.objects[objId] = obj
	m.objectsData[objId] = m.objectsData[location]

	return nil
}

// Job APIs
func (m *Manta) ListJobs(live bool) ([]manta.Entry, error) {
	var jobs []manta.Entry

	if err := m.ProcessFunctionHook(m, live); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for _, job := range m.jobs {
		if live && (job.Cancelled || job.TimeDone != "") {
			continue
		}
		jobKey := fmt.Sprintf(jobsPrefix, m.ServiceInstance.UserAccount, job.Id)
		jobs = append(jobs, m.objects[jobKey])
	}
	return jobs, nil
}

func (m *Manta) CreateJob(job []byte) (string, error) {
	if err := m.ProcessFunctionHook(m, job); err != nil {
		return "", err
	}

	jsonJob := new(manta.Job)
	err := json.Unmarshal(job, jsonJob)
	if err != nil {
		return "", err
	}
	jobId, err := localservices.NewUUID()
	if err != nil {
		return "", err
	}
	jsonJob.Id = jobId
	jsonJob.State = "running"
	jsonJob.Cancelled = false
	jsonJob.InputDone = false
	jsonJob.TimeCreated = time.Now().Format(time.RFC3339)

	//create directories
	m.mu.Lock()
	defer m.mu.Unlock()
	realPath := fmt.Sprintf(jobsPrefix, m.ServiceInstance.UserAccount, jobId)
	m.objects[realPath] = createDirectory(jobId)
	realPath = fmt.Sprintf(jobsPrefix, m.ServiceInstance.UserAccount, fmt.Sprintf("%s/stor", jobId))
	m.objects[realPath] = createDirectory("stor")

	m.jobs[jsonJob.Id] = jsonJob
	return fmt.Sprintf("/%s/jobs/%s", m.ServiceInstance.UserAccount, jobId), nil
}

func (m *Manta) GetJob(id string) (*manta.Job, error) {
	if err := m.ProcessFunctionHook(m, id); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if job, ok := m.jobs[id]; ok {
		return job, nil
	}
	return nil, fmt.Errorf("/%s/jobs/%s/job.json was not found", m.ServiceInstance.UserAccount, id)
}

func (m *Manta) CancelJob(id string) error {
	if err := m.ProcessFunctionHook(m, id); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if job, ok := m.jobs[id]; ok {
		if !job.InputDone {
			job.Cancelled = true
			job.InputDone = true
			job.TimeDone = time.Now().Format(time.RFC3339)
		} else {
			return fmt.Errorf("/%s/jobs/%s/live/cancel does not exist", m.ServiceInstance.UserAccount, id)
		}
		return nil
	}

	return fmt.Errorf("/%s/jobs/%s/job.json was not found", m.ServiceInstance.UserAccount, id)
}

func (m *Manta) AddJobInputs(id string, jobInputs []byte) error {
	if err := m.ProcessFunctionHook(m, id); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if job, ok := m.jobs[id]; ok {
		var err error
		if !job.InputDone {
			// add inputs
			objId := fmt.Sprintf("/%s/jobs/%s/in.txt", m.ServiceInstance.UserAccount, id)
			m.objects[objId], err = createJobObject("in.txt", jobInputs)
			if err != nil {
				return err
			}
			m.objectsData[objId] = jobInputs

			return nil
		} else {
			return fmt.Errorf("/%s/jobs/%s/live/in does not exist", m.ServiceInstance.UserAccount, id)
		}
	}

	return fmt.Errorf("/%s/jobs/%s/job.json was not found", m.ServiceInstance.UserAccount, id)
}

func (m *Manta) EndJobInput(id string) error {
	if err := m.ProcessFunctionHook(m, id); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if job, ok := m.jobs[id]; ok {
		if !job.InputDone {
			job.InputDone = true
		} else {
			return fmt.Errorf("/%s/jobs/%s/live/in/end does not exist", m.ServiceInstance.UserAccount, id)
		}
		return nil
	}

	return fmt.Errorf("/%s/jobs/%s/job.json was not found", m.ServiceInstance.UserAccount, id)
}

func (m *Manta) GetJobOutput(id string) (string, error) {
	if err := m.ProcessFunctionHook(m, id); err != nil {
		return "", err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if job, ok := m.jobs[id]; ok {
		index := len(job.Phases) - 1
		phaseType := job.Phases[index].Type
		outputId, err := localservices.NewUUID()
		if err != nil {
			return "", err
		}
		jobOutput := fmt.Sprintf("/%s/jobs/%s/stor/%s.%d.%s", m.ServiceInstance.UserAccount, id, phaseType, index, outputId)

		return jobOutput, nil
	}

	return "", fmt.Errorf("/%s/jobs/%s/job.json was not found", m.ServiceInstance.UserAccount, id)
}

func (m *Manta) GetJobInput(id string) (string, error) {
	if err := m.ProcessFunctionHook(m, id); err != nil {
		return "", err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.jobs[id]; ok {

		objId := fmt.Sprintf("/%s/jobs/%s/in.txt", m.ServiceInstance.UserAccount, id)
		if _, ok := m.objects[objId]; ok {
			return string(m.objectsData[objId]), nil
		}

		return "", nil
	}

	return "", fmt.Errorf("/%s/jobs/%s/job.json was not found", m.ServiceInstance.UserAccount, id)
}

func (m *Manta) GetJobFailures(id string) (string, error) {
	if err := m.ProcessFunctionHook(m, id); err != nil {
		return "", err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.jobs[id]; ok {
		return "", nil
	}

	return "", fmt.Errorf("/%s/jobs/%s/job.json was not found", m.ServiceInstance.UserAccount, id)
}

func (m *Manta) GetJobErrors(id string) ([]manta.JobError, error) {
	if err := m.ProcessFunctionHook(m, id); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.jobs[id]; ok {
		return nil, nil
	}

	return nil, fmt.Errorf("/%s/jobs/%s/job.json was not found", m.ServiceInstance.UserAccount, id)
}
