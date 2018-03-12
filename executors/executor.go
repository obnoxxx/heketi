//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package executors

import "encoding/xml"

type Executor interface {
	GlusterdCheck(host string) error
	PeerProbe(exec_host, newnode string) error
	PeerDetach(exec_host, detachnode string) error
	DeviceSetup(host, device, vgid string) (*DeviceInfo, error)
	GetDeviceInfo(host, device, vgid string) (*DeviceInfo, error)
	DeviceTeardown(host, device, vgid string) error
	BrickCreate(host string, brick *BrickRequest) (*BrickInfo, error)
	BrickDestroy(host string, brick *BrickRequest) error
	BrickDestroyCheck(host string, brick *BrickRequest) error
	VolumeCreate(host string, volume *VolumeRequest) (*Volume, error)
	VolumeDestroy(host string, volume string) error
	VolumeDestroyCheck(host, volume string) error
	VolumeExpand(host string, volume *VolumeRequest) (*Volume, error)
	VolumeReplaceBrick(host string, volume string, oldBrick *BrickInfo, newBrick *BrickInfo) error
	VolumeInfo(host string, volume string) (*Volume, error)
	VolumeSnapshotCreate(host string, volume *VolumeSnapshotRequest) (*VolumeSnapshot, error)
	VolumeSnapshotClone(host string, volume *VolumeSnapshotRequest) (*Volume, error)
	VolumeSnapshotDestroy(host string, snapshot string) error
	// TODO(removeme): VolumeSnapshotList is a DB operation, not remote
	VolumeSnapshotInfo(host string, snapshot string) (*VolumeSnapshot, error)
	HealInfo(host string, volume string) (*HealInfo, error)
	SetLogLevel(level string)
	BlockVolumeCreate(host string, blockVolume *BlockVolumeRequest) (*BlockVolumeInfo, error)
	BlockVolumeDestroy(host string, blockHostingVolumeName string, blockVolumeName string) error
}

// Enumerate durability types
type DurabilityType int

const (
	DurabilityNone DurabilityType = iota
	DurabilityReplica
	DurabilityDispersion
)

// Returns the size of the device
type DeviceInfo struct {
	// Size in KB
	Size       uint64
	ExtentSize uint64
}

// Brick description
type BrickRequest struct {
	VgId             string
	Name             string
	TpSize           uint64
	Size             uint64
	PoolMetadataSize uint64
	Gid              int64
	// Path is the brick mountpoint (named Path for symmetry with BrickInfo)
	Path string
}

// Returns information about the location of the brick
type BrickInfo struct {
	Path string
	Host string
}

type VolumeRequest struct {
	Bricks               []BrickInfo
	Name                 string
	Type                 DurabilityType
	GlusterVolumeOptions []string

	// Dispersion
	Data       int
	Redundancy int

	// Replica
	Replica int
}

type VolumeSnapshotRequest struct {
	// name of the volume to clone
	Volume      string
	// new, cloned volume name
	Name        string
	Description string
}

// TODO: automagically parse the XML output into types/structs
//
// # gluster --mode=script --xml snapshot info mysnap
// <?xml version="1.0" encoding="UTF-8" standalone="yes"?>
// <cliOutput>
//   <opRet>0</opRet>
//   <opErrno>0</opErrno>
//   <opErrstr/>
//   <snapInfo>
//     <count>1</count>
//     <snapshots>
//       <snapshot>
//         <name>mysnap</name>
//         <uuid>b0a12f9e-192b-4691-82e9-1bdb3c33e9f5</uuid>
//         <description/>
//         <createTime>2018-03-12 14:35:16</createTime>
//         <volCount>1</volCount>
//         <snapVolume>
//           <name>4516d565579c47cf82081e84f8049ae9</name>
//           <status>Stopped</status>
//           <originVolume>
//             <name>vol_10dca02524ed01e4a6cded5eacc04b96</name>
//             <snapCount>2</snapCount>
//             <snapRemaining>254</snapRemaining>
//           </originVolume>
//         </snapVolume>
//       </snapshot>
//     </snapshots>
//   </snapInfo>
// </cliOutput>
type VolumeSnapshot struct {
	XMLName     xml.Name `xml:"snapshot"`
	Name        string   `xml:"name"`
	UUID        string   `xml:"uuid"`
	Description string   `xml:"description"`
	CreateTime  string   `xml:"createTime"`
	// TODO: do we care about, or need other fields? originVolume/name was passed by the caller
}

type Brick struct {
	UUID      string `xml:"uuid,attr"`
	Name      string `xml:"name"`
	HostUUID  string `xml:"hostUuid"`
	IsArbiter int    `xml:"isArbiter"`
}

type Bricks struct {
	XMLName   xml.Name `xml:"bricks"`
	BrickList []Brick  `xml:"brick"`
}

type BrickHealStatus struct {
	HostUUID        string `xml:"hostUuid,attr"`
	Name            string `xml:"name"`
	Status          string `xml:"status"`
	NumberOfEntries string `xml:"numberOfEntries"`
}

type Option struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type Options struct {
	XMLName    xml.Name `xml:"options"`
	OptionList []Option `xml:"option"`
}

type Volume struct {
	XMLName         xml.Name `xml:"volume"`
	VolumeName      string   `xml:"name"`
	ID              string   `xml:"id"`
	Status          int      `xml:"status"`
	StatusStr       string   `xml:"statusStr"`
	BrickCount      int      `xml:"brickCount"`
	DistCount       int      `xml:"distCount"`
	StripeCount     int      `xml:"stripeCount"`
	ReplicaCount    int      `xml:"replicaCount"`
	ArbiterCount    int      `xml:"arbiterCount"`
	DisperseCount   int      `xml:"disperseCount"`
	RedundancyCount int      `xml:"redundancyCount"`
	Type            int      `xml:"type"`
	TypeStr         string   `xml:"typeStr"`
	Transport       int      `xml:"transport"`
	Bricks          Bricks
	OptCount        int `xml:"optCount"`
	Options         Options
}

type Volumes struct {
	XMLName    xml.Name `xml:"volumes"`
	Count      int      `xml:"count"`
	VolumeList []Volume `xml:"volume"`
}

type VolInfo struct {
	XMLName xml.Name `xml:"volInfo"`
	Volumes Volumes  `xml:"volumes"`
}

type HealInfoBricks struct {
	BrickList []BrickHealStatus `xml:"brick"`
}

type HealInfo struct {
	XMLName xml.Name       `xml:"healInfo"`
	Bricks  HealInfoBricks `xml:"bricks"`
}

type BlockVolumeRequest struct {
	Name              string
	Size              int
	GlusterVolumeName string
	GlusterNode       string
	Hacount           int
	BlockHosts        []string
	Auth              bool
}

type BlockVolumeInfo struct {
	Name              string
	Size              int
	GlusterVolumeName string
	GlusterNode       string
	Hacount           int
	BlockHosts        []string
	Iqn               string
	Username          string
	Password          string
}
