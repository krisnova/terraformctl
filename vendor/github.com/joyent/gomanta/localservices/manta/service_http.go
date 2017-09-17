//
// gomanta - Go library to interact with Joyent Manta
//
// Manta double testing service - HTTP API implementation
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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/joyent/gomanta/manta"
	"github.com/julienschmidt/httprouter"
)

// ErrorResponse defines a single HTTP error response.
type ErrorResponse struct {
	Code        int
	Body        string
	contentType string
	errorText   string
	headers     map[string]string
	manta       *Manta
}

var (
	ErrNotAllowed = &ErrorResponse{
		http.StatusMethodNotAllowed,
		"Method is not allowed",
		"text/plain; charset=UTF-8",
		"MethodNotAllowedError",
		nil,
		nil,
	}
	ErrNotFound = &ErrorResponse{
		http.StatusNotFound,
		"Resource Not Found",
		"text/plain; charset=UTF-8",
		"NotFoundError",
		nil,
		nil,
	}
	ErrBadRequest = &ErrorResponse{
		http.StatusBadRequest,
		"Malformed request url",
		"text/plain; charset=UTF-8",
		"BadRequestError",
		nil,
		nil,
	}
)

func (e *ErrorResponse) Error() string {
	return e.errorText
}

func (e *ErrorResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.contentType != "" {
		w.Header().Set("Content-Type", e.contentType)
	}
	body := e.Body
	if e.headers != nil {
		for h, v := range e.headers {
			w.Header().Set(h, v)
		}
	}
	// workaround for https://code.google.com/p/go/issues/detail?id=4454
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	if e.Code != 0 {
		w.WriteHeader(e.Code)
	}
	if len(body) > 0 {
		w.Write([]byte(body))
	}
}

type mantaHandler struct {
	manta  *Manta
	method func(m *Manta, w http.ResponseWriter, r *http.Request, p httprouter.Params) error
}

func (h *mantaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	path := r.URL.Path
	// handle trailing slash in the path
	if strings.HasSuffix(path, "/") && path != "/" {
		ErrNotFound.ServeHTTP(w, r)
		return
	}
	err := h.method(h.manta, w, r, p)
	if err == nil {
		return
	}
	var resp http.Handler
	resp, _ = err.(http.Handler)
	if resp == nil {
		resp = &ErrorResponse{
			http.StatusInternalServerError,
			`{"internalServerError":{"message":"Unkown Error",code:500}}`,
			"application/json",
			err.Error(),
			nil,
			h.manta,
		}
	}
	resp.ServeHTTP(w, r)
}

func writeResponse(w http.ResponseWriter, code int, body []byte) {
	// workaround for https://code.google.com/p/go/issues/detail?id=4454
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write(body)
}

// sendJSON sends the specified response serialized as JSON.
func sendJSON(code int, resp interface{}, w http.ResponseWriter, r *http.Request) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	writeResponse(w, code, data)
	return nil
}

func getJobId(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[3]
}

func (manta *Manta) handler(method func(m *Manta, w http.ResponseWriter, r *http.Request, p httprouter.Params) error) httprouter.Handle {
	handler := &mantaHandler{manta, method}
	return handler.ServeHTTP
}

// handleStorage handles the storage HTTP API.
func (m *Manta) handleStorage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	prefix := fmt.Sprintf("/%s/stor/", m.ServiceInstance.UserAccount)
	object := strings.TrimPrefix(r.URL.Path, prefix)
	switch r.Method {
	case "GET":
		if m.IsObject(r.URL.Path) {
			var resp []byte
			//GetObject
			obj, err := m.GetObject(object)
			if err != nil {
				return err
			}
			if obj == nil {
				obj = []byte{}
			}
			resp = obj
			// not using sendJson to avoid double json encoding
			writeResponse(w, http.StatusOK, resp)
			return nil
		} else if m.IsDirectory(r.URL.Path) {
			//ListDirectory
			var (
				marker string
				limit  int
			)
			opts := &manta.ListDirectoryOpts{}
			body, errB := ioutil.ReadAll(r.Body)
			if errB != nil {
				return errB
			}
			if len(body) > 0 {
				if errJ := json.Unmarshal(body, opts); errJ != nil {
					return errJ
				}
				marker = opts.Marker
				limit = opts.Limit
			}
			entries, err := m.ListDirectory(object, marker, limit)
			if err != nil {
				return err
			}
			if entries == nil {
				entries = []manta.Entry{}
			}
			resp := entries
			return sendJSON(http.StatusOK, resp, w, r)
		} else {
			return ErrNotFound
		}
	case "POST":
		return ErrNotAllowed
	case "PUT":
		if r.Header.Get("Content-Type") == "application/json; type=directory" {
			// PutDirectory
			err := m.PutDirectory(object)
			if err != nil {
				return err
			}
			return sendJSON(http.StatusNoContent, nil, w, r)
		} else if r.Header.Get("Location") != "" {
			// PutSnaplink
			path := object[:strings.LastIndex(object, "/")]
			objName := object[strings.LastIndex(object, "/")+1:]
			err := m.PutSnapLink(path, objName, r.Header.Get("Location"))
			if err != nil {
				return err
			}
			return sendJSON(http.StatusNoContent, nil, w, r)
		} else {
			// PutObject
			path := object[:strings.LastIndex(object, "/")]
			objName := object[strings.LastIndex(object, "/")+1:]
			defer r.Body.Close()
			objectData, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			err = m.PutObject(path, objName, objectData)
			if err != nil {
				return err
			}
			return sendJSON(http.StatusNoContent, nil, w, r)
		}
	case "DELETE":
		if m.IsObject(r.URL.Path) {
			//DeleteObject
			err := m.DeleteObject(object)
			if err != nil {
				return err
			}
			return sendJSON(http.StatusNoContent, nil, w, r)
		} else if m.IsDirectory(r.URL.Path) {
			//DeleteDirectory
			err := m.DeleteDirectory(object)
			if err != nil {
				return err
			}
			return sendJSON(http.StatusNoContent, nil, w, r)
		} else {
			return ErrNotFound
		}
	}
	return fmt.Errorf("unknown request method %q for %s", r.Method, r.URL.Path)
}

// handleJob handles the Job HTTP API.
func (m *Manta) handleJobs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var live = false
	switch r.Method {
	case "GET":
		if strings.HasSuffix(r.URL.Path, "jobs") {
			// ListJobs
			if state := r.FormValue("state"); state == "running" {
				live = true
			}
			jobs, err := m.ListJobs(live)
			if err != nil {
				return err
			}
			if len(jobs) == 0 {
				jobs = []manta.Entry{}
			}
			resp := jobs
			return sendJSON(http.StatusOK, resp, w, r)
		} else if strings.HasSuffix(r.URL.Path, "status") {
			// GetJob
			job, err := m.GetJob(getJobId(r.URL.Path))
			if err != nil {
				return err
			}
			if job == nil {
				job = &manta.Job{}
			}
			resp := job
			return sendJSON(http.StatusOK, resp, w, r)
		} else if strings.HasSuffix(r.URL.Path, "out") {
			// GetJobOutput
			out, err := m.GetJobOutput(getJobId(r.URL.Path))
			if err != nil {
				return err
			}
			resp := out
			return sendJSON(http.StatusOK, resp, w, r)
		} else if strings.HasSuffix(r.URL.Path, "in") {
			// GetJobInput
			in, err := m.GetJobInput(getJobId(r.URL.Path))
			if err != nil {
				return err
			}
			resp := in
			return sendJSON(http.StatusOK, resp, w, r)
		} else if strings.HasSuffix(r.URL.Path, "fail") {
			// GetJobFailures
			fail, err := m.GetJobFailures(getJobId(r.URL.Path))
			if err != nil {
				return err
			}
			resp := fail
			return sendJSON(http.StatusOK, resp, w, r)
		} else if strings.HasSuffix(r.URL.Path, "err") {
			// GetJobErrors
			jobErr, err := m.GetJobErrors(getJobId(r.URL.Path))
			if err != nil {
				return err
			}
			//if jobErr == nil {
			//	jobErr = []manta.JobError{}
			//}
			resp := jobErr
			return sendJSON(http.StatusOK, resp, w, r)
		} else {
			return ErrNotAllowed
		}
	case "POST":
		if strings.HasSuffix(r.URL.Path, "jobs") {
			body, errb := ioutil.ReadAll(r.Body)
			if errb != nil {
				return errb
			}
			if len(body) == 0 {
				return ErrBadRequest
			}
			jobId, err := m.CreateJob(body)
			if err != nil {
				return err
			}
			w.Header().Add("Location", jobId)
			return sendJSON(http.StatusCreated, nil, w, r)
		} else if strings.HasSuffix(r.URL.Path, "cancel") {
			// CancelJob
			err := m.CancelJob(getJobId(r.URL.Path))
			if err != nil {
				return err
			}
			return sendJSON(http.StatusAccepted, nil, w, r)
		} else if strings.HasSuffix(r.URL.Path, "end") {
			// EndJobInputs
			err := m.EndJobInput(getJobId(r.URL.Path))
			if err != nil {
				return err
			}
			return sendJSON(http.StatusAccepted, nil, w, r)
		} else if strings.HasSuffix(r.URL.Path, "in") {
			// AddJobInputs
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			if len(body) == 0 {
				return ErrBadRequest
			}
			err = m.AddJobInputs(getJobId(r.URL.Path), body)
			if err != nil {
				return err
			}
			return sendJSON(http.StatusNoContent, nil, w, r)
		} else {
			return ErrNotAllowed
		}
	case "PUT":
		return ErrNotAllowed
	case "DELETE":
		return ErrNotAllowed
	}
	return fmt.Errorf("unknown request method %q for %s", r.Method, r.URL.Path)
}

// setupHTTP attaches all the needed handlers to provide the HTTP API.
func (m *Manta) SetupHTTP(mux *httprouter.Router) {

	baseRoute := "/" + m.ServiceInstance.UserAccount
	mux.NotFound = ErrNotFound
	mux.MethodNotAllowed = ErrNotAllowed
	mux.RedirectFixedPath = true
	mux.RedirectTrailingSlash = true
	mux.HandleMethodNotAllowed = true

	// this particular route can't be handled by httprouter correctly due to its
	// handling of positional parameters but we can't just pass in ErrNotAllowed
	// either because it's the wrong type
	handleNotAllowed := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.Write([]byte("Method is not allowed"))
	}
	mux.POST(baseRoute+"/stor", handleNotAllowed)

	// storage APIs
	// PutSnapLink and PutMetaData have the same route spec as other routes
	// in this group, which isn't permitted by httprouter. We pick up the
	// correct path in the handler
	mux.PUT(baseRoute+"/stor/:dir", m.handler((*Manta).handleStorage))         // PutDirectory
	mux.GET(baseRoute+"/stor/:dir", m.handler((*Manta).handleStorage))         // ListDirectory
	mux.DELETE(baseRoute+"/stor/:dir", m.handler((*Manta).handleStorage))      // DeleteDirectory
	mux.PUT(baseRoute+"/stor/:dir/:obj", m.handler((*Manta).handleStorage))    // PutObject
	mux.GET(baseRoute+"/stor/:dir/:obj", m.handler((*Manta).handleStorage))    // GetObject
	mux.DELETE(baseRoute+"/stor/:dir/:obj", m.handler((*Manta).handleStorage)) // DeleteObject

	// job APIs
	mux.POST(baseRoute+"/jobs", m.handler((*Manta).handleJobs))                 // CreateJob
	mux.POST(baseRoute+"/jobs/:id/live/in", m.handler((*Manta).handleJobs))     // AddJobInputs
	mux.POST(baseRoute+"/jobs/:id/live/in/end", m.handler((*Manta).handleJobs)) // EndJobInput
	mux.POST(baseRoute+"/jobs/:id/live/cancel", m.handler((*Manta).handleJobs)) // CancelJob
	mux.GET(baseRoute+"/jobs", m.handler((*Manta).handleJobs))                  // ListJobs
	mux.GET(baseRoute+"/jobs/:id/live/status", m.handler((*Manta).handleJobs))  // GetJob
	mux.GET(baseRoute+"/jobs/:id/live/out", m.handler((*Manta).handleJobs))     // GetJobOutput
	mux.GET(baseRoute+"/jobs/:id/live/in", m.handler((*Manta).handleJobs))      // GetJobInput
	mux.GET(baseRoute+"/jobs/:id/live/fail", m.handler((*Manta).handleJobs))    // GetJobFailures
	mux.GET(baseRoute+"/jobs/:id/live/err", m.handler((*Manta).handleJobs))     // GetJobErrors
}
