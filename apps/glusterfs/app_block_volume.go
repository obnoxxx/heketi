//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package glusterfs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/heketi/heketi/pkg/glusterfs/api"
	"github.com/heketi/heketi/pkg/utils"
)

func (a *App) BlockVolumeCreate(w http.ResponseWriter, r *http.Request) {

	var msg api.BlockVolumeCreateRequest
	err := utils.GetJsonFromRequest(r, &msg)
	if err != nil {
		http.Error(w, "request unable to be parsed", 422)
		return
	}

	if msg.Size < 1 {
		http.Error(w, "Invalid volume size", http.StatusBadRequest)
		logger.LogError("Invalid volume size")
		return
	}

	// TODO: factor this into a function (it's also in VolumeCreate)
	// Check that the clusters requested are available
	err = a.db.View(func(tx *bolt.Tx) error {

		// :TODO: All we need to do is check for one instead of gathering all keys
		clusters, err := ClusterList(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		if len(clusters) == 0 {
			http.Error(w, fmt.Sprintf("No clusters configured"), http.StatusBadRequest)
			logger.LogError("No clusters configured")
			return ErrNotFound
		}

		// Check the clusters requested are correct
		for _, clusterid := range msg.Clusters {
			_, err := NewClusterEntryFromId(tx, clusterid)
			if err != nil {
				http.Error(w, fmt.Sprintf("Cluster id %v not found", clusterid), http.StatusBadRequest)
				logger.LogError(fmt.Sprintf("Cluster id %v not found", clusterid))
				return err
			}
		}

		return nil
	})
	if err != nil {
		return
	}

	blockvol := NewBlockVolumeEntryFromRequest(&msg)

	// Add device in an asynchronous function
	a.asyncManager.AsyncHttpRedirectFunc(w, r, func() (string, error) {

		logger.Info("Creating block volume %v", blockvol.Info.Id)
		err := blockvol.Create(a.db, a.executor, a.allocator)
		if err != nil {
			logger.LogError("Failed to create block volume: %v", err)
			return "", err
		}

		logger.Info("Created block volume %v", blockvol.Info.Id)

		// Done
		return "/blockvolumes/" + blockvol.Info.Id, nil
	})
}

func (a *App) BlockVolumeList(w http.ResponseWriter, r *http.Request) {

	var list api.BlockVolumeListResponse

	// Get all the cluster ids from the DB
	err := a.db.View(func(tx *bolt.Tx) error {
		var err error

		list.BlockVolumes, err = BlockVolumeList(tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send list back
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		panic(err)
	}
}

func (a *App) BlockVolumeInfo(w http.ResponseWriter, r *http.Request) {

	// Get volume id from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get volume information
	var info *api.BlockVolumeInfoResponse
	err := a.db.View(func(tx *bolt.Tx) error {
		entry, err := NewBlockVolumeEntryFromId(tx, id)
		if err == ErrNotFound {
			http.Error(w, "Id not found", http.StatusNotFound)
			return err
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		info, err = entry.NewInfoResponse(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		return nil
	})
	if err != nil {
		return
	}

	// Write msg
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		panic(err)
	}
}

func (a *App) BlockVolumeDelete(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get volume entry
	var volume *BlockVolumeEntry
	err := a.db.View(func(tx *bolt.Tx) error {

		// Access volume entry
		var err error
		volume, err = NewBlockVolumeEntryFromId(tx, id)
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

	a.asyncManager.AsyncHttpRedirectFunc(w, r, func() (string, error) {

		// Actually destroy the Volume here
		err := volume.Destroy(a.db, a.executor)

		// If it fails for some reason, we will need to add to the DB again
		// or hold state on the entry "DELETING"

		// Show that the key has been deleted
		if err != nil {
			logger.LogError("Failed to delete volume %v: %v", volume.Info.Id, err)
			return "", err
		}

		logger.Info("Deleted volume [%s]", id)
		return "", nil

	})
}

func (a *App) BlockVolumeExpand(w http.ResponseWriter, r *http.Request) {
	logger.Debug("In VolumeExpand")

	// Get the id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	var msg api.BlockVolumeExpandRequest
	err := utils.GetJsonFromRequest(r, &msg)
	if err != nil {
		http.Error(w, "request unable to be parsed", 422)
		return
	}
	logger.Debug("Msg: %v", msg)

	// Check the message
	if msg.Size < 1 {
		http.Error(w, "Invalid volume size", http.StatusBadRequest)
		return
	}
	logger.Debug("Size: %v", msg.Size)

	// Get volume entry
	var volume *BlockVolumeEntry
	err = a.db.View(func(tx *bolt.Tx) error {

		// Access volume entry
		var err error
		volume, err = NewBlockVolumeEntryFromId(tx, id)
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

	// Expand volume in an asynchronous function
	a.asyncManager.AsyncHttpRedirectFunc(w, r, func() (string, error) {

		logger.Info("Expanding block volume %v", volume.Info.Id)
		err := volume.Expand(a.db, a.executor, a.allocator, msg.Size)
		if err != nil {
			logger.LogError("Failed to expand block volume %v", volume.Info.Id)
			return "", err
		}

		logger.Info("Expanded block volume %v", volume.Info.Id)

		// Done
		return "/blockvolumes/" + volume.Info.Id, nil
	})
}