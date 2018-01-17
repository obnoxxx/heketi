//
// Copyright (c) 2017 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package glusterfs

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/heketi/heketi/pkg/glusterfs/api"
)

type Db struct {
	Clusters     []ClusterEntry     `json:"clusterentries"`
	Volumes      []VolumeEntry      `json:"volumeentries"`
	Bricks       []BrickEntry       `json:"brickentries"`
	Nodes        []NodeEntry        `json:"nodeentries"`
	Devices      []DeviceEntry      `json:"deviceentries"`
	BlockVolumes []BlockVolumeEntry `json:"blockvolumeentries"`
	DbAttributes []DbAttributeEntry `json:"dbattributeentries"`
}

func dbDumpInternal(db *bolt.DB) (Db, error) {
	var dump Db
	clusterEntryList := make([]ClusterEntry, 0)
	volEntryList := make([]VolumeEntry, 0)
	brickEntryList := make([]BrickEntry, 0)
	nodeEntryList := make([]NodeEntry, 0)
	deviceEntryList := make([]DeviceEntry, 0)
	blockvolEntryList := make([]BlockVolumeEntry, 0)
	dbattributeEntryList := make([]DbAttributeEntry, 0)

	err := db.View(func(tx *bolt.Tx) error {

		logger.Info("rtalurlogs: starting volume bucket")

		// Volume Bucket
		volumes, err := VolumeList(tx)
		if err != nil {
			return err
		}

		for _, volume := range volumes {
			logger.Info("rtalurlogs: adding volume entry %v", volume)
			volEntry, err := NewVolumeEntryFromId(tx, volume)
			if err != nil {
				return err
			}
			volEntryList = append(volEntryList, *volEntry)
		}

		// Brick Bucket
		logger.Info("rtalurlogs: starting brick bucket")
		bricks, err := BrickList(tx)
		if err != nil {
			return err
		}

		for _, brick := range bricks {
			logger.Info("rtalurlogs: adding brick entry %v", brick)
			brickEntry, err := NewBrickEntryFromId(tx, brick)
			if err != nil {
				return err
			}
			brickEntryList = append(brickEntryList, *brickEntry)
		}

		// Cluster Bucket
		logger.Info("rtalurlogs: starting cluster bucket")
		clusters, err := ClusterList(tx)
		if err != nil {
			return err
		}

		for _, cluster := range clusters {
			logger.Info("rtalurlogs: adding cluster entry %v", cluster)
			clusterEntry, err := NewClusterEntryFromId(tx, cluster)
			if err != nil {
				return err
			}
			clusterEntryList = append(clusterEntryList, *clusterEntry)
		}

		// Node Bucket
		logger.Info("rtalurlogs: starting node bucket")
		nodes, err := NodeList(tx)
		if err != nil {
			return err
		}

		for _, node := range nodes {
			logger.Info("rtalurlogs: adding node entry %v", node)
			if strings.HasPrefix(node, "MANAGE") || strings.HasPrefix(node, "STORAGE") {
				logger.Info("rtalurlogs, ignoring registry key")
			} else {
				nodeEntry, err := NewNodeEntryFromId(tx, node)
				if err != nil {
					return err
				}
				nodeEntryList = append(nodeEntryList, *nodeEntry)
			}
		}

		// Device Bucket
		logger.Info("rtalurlogs: starting device bucket")
		devices, err := DeviceList(tx)
		if err != nil {
			return err
		}

		for _, device := range devices {
			logger.Info("rtalurlogs: adding device entry %v", device)
			if strings.HasPrefix(device, "DEVICE") {
				logger.Info("rtalurlogs, ignoring registry key")
			} else {
				deviceEntry, err := NewDeviceEntryFromId(tx, device)
				if err != nil {
					return err
				}
				deviceEntryList = append(deviceEntryList, *deviceEntry)
			}
		}

		// BlockVolume Bucket
		blockvolumes, err := BlockVolumeList(tx)
		if err != nil {
			return err
		}

		for _, blockvolume := range blockvolumes {
			logger.Info("rtalurlogs: adding blockvolume entry %v", blockvolume)
			blockvolEntry, err := NewBlockVolumeEntryFromId(tx, blockvolume)
			if err != nil {
				return err
			}
			blockvolEntryList = append(blockvolEntryList, *blockvolEntry)
		}

		// DbAttributes Bucket
		dbattributes, err := DbAttributeList(tx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to list dbattributes %v", err)
			return err
		}

		for _, dbattribute := range dbattributes {
			logger.Info("rtalurlogs: adding dbattribute entry %v", dbattribute)
			dbattributeEntry, err := NewDbAttributeEntryFromKey(tx, dbattribute)
			if err != nil {
				return err
			}
			dbattributeEntryList = append(dbattributeEntryList, *dbattributeEntry)
		}

		return nil
	})
	if err != nil {
		return Db{}, fmt.Errorf("Could not construct dump from DB: %v", err.Error())
	}

	dump.Clusters = clusterEntryList
	dump.Volumes = volEntryList
	dump.Bricks = brickEntryList
	dump.Nodes = nodeEntryList
	dump.Devices = deviceEntryList
	dump.BlockVolumes = blockvolEntryList
	dump.DbAttributes = dbattributeEntryList

	return dump, nil
}

// DbDump ... Creates a JSON output representing the state of DB
// This is the variant to be called offline, i.e. when the server is not
// running.
func DbDump(jsonfile string, dbfile string) error {
	// Load config file
	fp, err := os.Create(jsonfile)
	if err != nil {
		return fmt.Errorf("Could not open input file: %v", err.Error())
	}
	defer fp.Close()

	db, err := bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open database: %v. exiting", err)
		os.Exit(1)
	}

	dump, err := dbDumpInternal(db)
	if err != nil {
		return fmt.Errorf("Could not construct dump from DB: %v", err.Error())
	}

	if err := json.NewEncoder(fp).Encode(dump); err != nil {
		return fmt.Errorf("Could not encode dump as JSON: %v", err.Error())
	}

	return nil
}

// DbDump ... Creates a JSON output representing the state of DB
// This is the variant to be called via the API and running in the App
func (a *App) DbDump(w http.ResponseWriter, r *http.Request) {
	dump, err := dbDumpInternal(a.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write msg
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dump); err != nil {
		panic(err)
	}
}

// DbCreate ... Creates a bolt db file based on JSON input
func DbCreate(jsonfile string, dbfile string) error {
	var dump Db
	//vars := mux.Vars(r)
	//jsonFile := vars["jsonFile"]
	// Check arguments
	//if jsonFile == "" {
	//	logger.Info("rtalurlogs, jsonFile value is %v", jsonFile)
	//	http.Error(w, ErrNotFound.Error(), http.StatusInternalServerError)
	//	return
	//}

	// Load config file
	fp, err := os.Open(jsonfile)
	if err != nil {
		return fmt.Errorf("Could not open input file: %v", err.Error())
	}
	defer fp.Close()

	dbParser := json.NewDecoder(fp)
	if err = dbParser.Decode(&dump); err != nil {
		return fmt.Errorf("Could not decode input file as JSON: %v", err.Error())
	}

	// Setup BoltDB database
	dbhandle, err := bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open database: %v. exiting", err)
		os.Exit(1)
	} else {
		err = dbhandle.Update(func(tx *bolt.Tx) error {
			// Create Cluster Bucket
			_, err := tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_CLUSTER))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create cluster bucket in DB")
				return err
			}

			// Create Node Bucket
			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_NODE))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create node bucket in DB")
				return err
			}

			// Create Volume Bucket
			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_VOLUME))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create volume bucket in DB")
				return err
			}

			// Create Device Bucket
			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_DEVICE))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create device bucket in DB")
				return err
			}

			// Create Brick Bucket
			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_BRICK))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create brick bucket in DB")
				return err
			}

			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_BLOCKVOLUME))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create blockvolume bucket in DB")
				return err
			}

			_, err = tx.CreateBucketIfNotExists([]byte(BOLTDB_BUCKET_DBATTRIBUTE))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create dbattribute bucket in DB")
				return err
			}

			return nil

		})
		if err != nil {
			logger.Err(err)
			return nil
		}
	}

	//this registration should ideally be done during initialization, but it is existing bug
	//work around it
	gob.Register(&NoneDurability{})
	gob.Register(&VolumeReplicaDurability{})
	gob.Register(&VolumeDisperseDurability{})

	err = dbhandle.Update(func(tx *bolt.Tx) error {
		for _, cluster := range dump.Clusters {
			fmt.Fprintf(os.Stderr, "adding cluster entry %v", cluster.Info.Id)
			err := cluster.Save(tx)
			if err != nil {
				return fmt.Errorf("Could not save cluster bucket: %v", err.Error())
			}
		}
		for _, volume := range dump.Volumes {
			fmt.Fprintf(os.Stderr, "adding volume entry %v", volume.Info.Id)
			// Set default durability values
			durability := volume.Info.Durability.Type
			switch {

			case durability == api.DurabilityReplicate:
				volume.Durability = NewVolumeReplicaDurability(&volume.Info.Durability.Replicate)

			case durability == api.DurabilityEC:
				volume.Durability = NewVolumeDisperseDurability(&volume.Info.Durability.Disperse)

			case durability == api.DurabilityDistributeOnly || durability == "":
				volume.Durability = NewNoneDurability()

			default:
				return fmt.Errorf("Not a known volume type: %v", err.Error())
			}

			// Set the default values accordingly
			volume.Durability.SetDurability()
			err := volume.Save(tx)
			if err != nil {
				return fmt.Errorf("Could not save volume bucket: %v", err.Error())
			}
		}
		for _, brick := range dump.Bricks {
			fmt.Fprintf(os.Stderr, "adding brick entry %v", brick.Info.Id)
			err := brick.Save(tx)
			if err != nil {
				return fmt.Errorf("Could not save brick bucket: %v", err.Error())
			}
		}
		for _, node := range dump.Nodes {
			fmt.Fprintf(os.Stderr, "adding node entry %v", node.Info.Id)
			err := node.Save(tx)
			if err != nil {
				return fmt.Errorf("Could not save node bucket: %v", err.Error())
			}
			fmt.Fprintf(os.Stderr, "registering node entry %v", node.Info.Id)
			err = node.Register(tx)
			if err != nil {
				return fmt.Errorf("Could not register node: %v", err.Error())
			}
		}
		for _, device := range dump.Devices {
			fmt.Fprintf(os.Stderr, "adding device entry %v", device.Info.Id)
			err := device.Save(tx)
			if err != nil {
				return fmt.Errorf("Could not save device bucket: %v", err.Error())
			}
			fmt.Fprintf(os.Stderr, "registering device entry %v", device.Info.Id)
			err = device.Register(tx)
			if err != nil {
				return fmt.Errorf("Could not register device: %v", err.Error())
			}
		}
		for _, blockvolume := range dump.BlockVolumes {
			fmt.Fprintf(os.Stderr, "adding blockvolume entry %v", blockvolume.Info.Id)
			err := blockvolume.Save(tx)
			if err != nil {
				return fmt.Errorf("Could not save blockvolume bucket: %v", err.Error())
			}
		}
		for _, dbattribute := range dump.DbAttributes {
			fmt.Fprintf(os.Stderr, "adding dbattribute entry %v", dbattribute.Key)
			err := dbattribute.Save(tx)
			if err != nil {
				return fmt.Errorf("Could not save dbattribute bucket: %v", err.Error())
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
