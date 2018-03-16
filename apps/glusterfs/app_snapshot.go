//
// Copyright (c) 2018 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package glusterfs

import (
	//"encoding/json"
	"fmt"
	//"math"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	//"github.com/heketi/heketi/pkg/db"
	"github.com/heketi/heketi/pkg/glusterfs/api"
	"github.com/heketi/heketi/pkg/utils"
)

func (a *App) SnapshotInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *App) SnapshotClone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	snap_id := vars["id"]

	var msg api.SnapshotCloneRequest
	err := utils.GetJsonFromRequest(r, &msg)
	if err != nil {
		http.Error(w, "request unable to be parsed", http.StatusUnprocessableEntity)
		return
	}
	err = msg.Validate()
	if err != nil {
		http.Error(w, "validation failed: "+err.Error(),
			http.StatusBadRequest)
		logger.LogError("validation failed: " + err.Error())
		return
	}

	var snap *SnapshotEntry
	err = a.db.View(func(tx *bolt.Tx) error {
		var err error
		snap, err = NewSnapshotEntryFromId(tx, snap_id)
		if err == ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return err
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	op := NewSnapshotCloneOperation(snap, a.db)
	if err := AsyncHttpOperation(a, w, r, op); err != nil {
		http.Error(w,
			fmt.Sprintf("Failed to create clone of snapshot "+
				"%v: %v", snap_id, err),
			http.StatusInternalServerError)
		return
	}
}

func (a *App) SnapshotDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (a *App) SnapshotList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
