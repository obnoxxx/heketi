//
// Copyright (c) 2017 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package glusterfs

type Db struct {
	Clusters     []ClusterEntry     `json:"clusterentries"`
	Volumes      []VolumeEntry      `json:"volumeentries"`
	Bricks       []BrickEntry       `json:"brickentries"`
	Nodes        []NodeEntry        `json:"nodeentries"`
	Devices      []DeviceEntry      `json:"deviceentries"`
	BlockVolumes []BlockVolumeEntry `json:"blockvolumeentries"`
	DbAttributes []DbAttributeEntry `json:"dbattributeentries"`
}
