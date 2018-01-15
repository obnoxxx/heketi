//
// Copyright (c) 2016 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), as published by the Free Software Foundation,
// or under the Apache License, Version 2.0 <LICENSE-APACHE2 or
// http://www.apache.org/licenses/LICENSE-2.0>.
//
// You may not use this file except in compliance with those terms.
//

package client

import (
	"io/ioutil"
	"net/http"

	"github.com/heketi/heketi/pkg/utils"
)

func (c *Client) DbDump() (string, error) {
	req, err := http.NewRequest("GET", c.host+"/db/dump", nil)
	if err != nil {
		return "", err
	}

	// Set token
	err = c.setToken(req)
	if err != nil {
		return "", err
	}

	// Send request
	r, err := c.do(req)
	if err != nil {
		return "", err
	}
	if r.StatusCode != http.StatusOK {
		return "", utils.GetErrorFromResponse(r)
	}

	// Read JSON response
	//var dbdump *glusterfs.Db
	//err = utils.GetJsonFromResponse(r, &dbdump)
	//if err != nil {
	//	return nil, err
	//}
	//return &dbdump, nil

	defer r.Body.Close()
	respBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	respJSON := string(respBytes)
	return respJSON, nil
}
