//
// Copyright (c) 2018 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package cmdexec

import (
	"encoding/xml"
	"fmt"

	"github.com/heketi/heketi/executors"
	"github.com/lpabon/godbc"
)

func (s *CmdExecutor) snapshotActivate(host string, snapshot string) error {
	godbc.Require(host != "")
	godbc.Require(snapshot != "")

	type CliOutput struct {
		OpRet    int                `xml:"opRet"`
		OpErrno  int                `xml:"opErrno"`
		OpErrStr string             `xml:"opErrstr"`
		Snapshot executors.Snapshot `xml:"snapActivate"`
	}

	command := []string{
		fmt.Sprintf("gluster --mode=script --xml snapshot activate %v", snapshot),
	}

	output, err := s.RemoteExecutor.RemoteCommandExecute(host, command, 10)
	if err != nil {
		return fmt.Errorf("Unable to activate snapshot: %v", snapshot)
	}

	var snapActivate CliOutput
	err = xml.Unmarshal([]byte(output[0]), &snapActivate)
	if err != nil {
		return fmt.Errorf("Unable to activate snapshot: %v", snapshot)
	}
	logger.Debug("%+v\n", snapActivate)

	return nil
}

func (s *CmdExecutor) snapshotDeactivate(host string, snapshot string) error {
	godbc.Require(host != "")
	godbc.Require(snapshot != "")

	type CliOutput struct {
		OpRet    int                `xml:"opRet"`
		OpErrno  int                `xml:"opErrno"`
		OpErrStr string             `xml:"opErrstr"`
		Snapshot executors.Snapshot `xml:"snapDeactivate"`
	}

	command := []string{
		fmt.Sprintf("gluster --mode=script --xml snapshot deactivate %v", snapshot),
	}

	output, err := s.RemoteExecutor.RemoteCommandExecute(host, command, 10)
	if err != nil {
		return fmt.Errorf("Unable to deactivate snapshot: %v", snapshot)
	}

	var snapDeactivate CliOutput
	err = xml.Unmarshal([]byte(output[0]), &snapDeactivate)
	if err != nil {
		return fmt.Errorf("Unable to deactivate snapshot: %v", snapshot)
	}
	logger.Debug("%+v\n", snapDeactivate)

	return nil
}

func (s *CmdExecutor) SnapshotCloneVolume(host string, vcr *executors.SnapshotCloneRequest) (*executors.Volume, error) {
	godbc.Require(host != "")
	godbc.Require(vcr != nil)

	// cloning can only be done when a snapshot is acticated
	err := s.snapshotActivate(host, vcr.Snapshot)
	if err != nil {
		return nil, err
	}

	// we do not want activated snapshots sticking around
	defer s.snapshotDeactivate(host, vcr.Snapshot)

	type CliOutput struct {
		OpRet    int              `xml:"opRet"`
		OpErrno  int              `xml:"opErrno"`
		OpErrStr string           `xml:"opErrstr"`
		Volume   executors.Volume `xml:"CloneCreate"`
	}

	command := []string{
		fmt.Sprintf("gluster --mode=script --xml snapshot clone %v %v", vcr.Volume, vcr.Snapshot),
	}

	output, err := s.RemoteExecutor.RemoteCommandExecute(host, command, 10)
	if err != nil {
		return nil, fmt.Errorf("Unable to clone snapshot: %v", vcr.Snapshot)
	}

	var snapCreate CliOutput
	err = xml.Unmarshal([]byte(output[0]), &snapCreate)
	if err != nil {
		return nil, fmt.Errorf("Unable to clone snapshot: %v", vcr.Snapshot)
	}
	logger.Debug("%+v\n", snapCreate)

	return &snapCreate.Volume, nil
}

func (s *CmdExecutor) SnapshotCloneBlockVolume(host string, vcr *executors.SnapshotCloneRequest) (*executors.BlockVolumeInfo, error) {
	// TODO: cloning of block volume is not implemented yet
	return nil, fmt.Errorf("block snapshot %v can not be cloned, not implemented yet", vcr.Snapshot)
}

func (s *CmdExecutor) SnapshotDestroy(host string, snapshot string) error {
	godbc.Require(host != "")
	godbc.Require(snapshot != "")

	type CliOutput struct {
		OpRet    int                `xml:"opRet"`
		OpErrno  int                `xml:"opErrno"`
		OpErrStr string             `xml:"opErrstr"`
		Snapshot executors.Snapshot `xml:"snapDelete"`
	}

	command := []string{
		fmt.Sprintf("gluster --mode=script --xml snapshot delete %v", snapshot),
	}

	output, err := s.RemoteExecutor.RemoteCommandExecute(host, command, 10)
	if err != nil {
		return fmt.Errorf("Unable to delete snapshot: %v", snapshot)
	}

	var snapDelete CliOutput
	err = xml.Unmarshal([]byte(output[0]), &snapDelete)
	if err != nil {
		return fmt.Errorf("Unable to delete snapshot: %v", snapshot)
	}
	logger.Debug("%+v\n", snapDelete)

	return nil
}

func (s *CmdExecutor) SnapshotInfo(host string, snapshot string) (*executors.Snapshot, error) {
	godbc.Require(host != "")
	godbc.Require(snapshot != "")

	// info of a single snapshot returns a list of snapshots...
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

	type CliOutput struct {
		OpRet    int                  `xml:"opRet"`
		OpErrno  int                  `xml:"opErrno"`
		OpErrStr string               `xml:"opErrstr"`
		Snapshot []executors.Snapshot `xml:"snapshots"` // TODO: does this work without mentioning <snapInfo>?
	}

	command := []string{
		fmt.Sprintf("gluster --mode=script --xml snapshot info %v", snapshot),
	}

	output, err := s.RemoteExecutor.RemoteCommandExecute(host, command, 10)
	if err != nil {
		return nil, fmt.Errorf("Unable get information about snapshot: %v", snapshot)
	}

	var snapInfo CliOutput
	err = xml.Unmarshal([]byte(output[0]), &snapInfo)
	if err != nil {
		return nil, fmt.Errorf("Unable get information about snapshot: %v", snapshot)
	}
	logger.Debug("%+v\n", snapInfo)

	// TODO: instead of Snapshot[0], it would be better to search for Snapshot.Name == snapshot
	return &snapInfo.Snapshot[0], nil
}
