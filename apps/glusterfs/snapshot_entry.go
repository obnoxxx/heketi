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
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/heketi/heketi/pkg/glusterfs/api"
	"github.com/lpabon/godbc"
)

type SnapshotEntry struct {
	Info    api.SnapshotInfo
	Pending PendingItem
}

func NewSnapshotEntry() *SnapshotEntry {
	entry := &SnapshotEntry{}
	return entry
}

func NewSnapshotEntryFromId(tx *bolt.Tx, id string) (entry *SnapshotEntry, e error) {
	godbc.Require(tx != nil)

	entry = NewSnapshotEntry()
	e = EntryLoad(tx, entry, id)
	if e != nil {
		entry = nil
	}

	return
}

func (s *SnapshotEntry) BucketName() string {
	return BOLTDB_BUCKET_SNAPSHOT
}

func (s *SnapshotEntry) Save(tx *bolt.Tx) error {
	godbc.Require(tx != nil)
	godbc.Require(len(s.Info.Id) > 0)

	return EntrySave(tx, s, s.Info.Id)
}

func (s *SnapshotEntry) Delete(tx *bolt.Tx) error {
	return EntryDelete(tx, s, s.Info.Id)
}

func (s *SnapshotEntry) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(*s)

	return buffer.Bytes(), err
}

func (s *SnapshotEntry) Unmarshal(buffer []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(buffer))
	err := dec.Decode(s)

	return err
}
