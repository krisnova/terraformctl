//
// gomanta - Go library to interact with Joyent Manta
//
// Double testing service
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2016 Joyent Inc.
//
// Written by Daniele Stroppa <daniele.stroppa@joyent.com>
//

package localservices

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/joyent/gomanta/localservices/hook"
	"github.com/julienschmidt/httprouter"
)

// An HttpService provides the HTTP API for a service double.
type HttpService interface {
	SetupHTTP(mux *httprouter.Router)
}

// A ServiceInstance is an Joyent Cloud service, one of manta or cloudapi.
type ServiceInstance struct {
	hook.TestService
	Scheme      string
	Hostname    string
	UserAccount string
}

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
