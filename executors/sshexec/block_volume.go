//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package sshexec

import (
	"encoding/xml"
	"fmt"

	"github.com/heketi/heketi/executors"
	"github.com/lpabon/godbc"
)

func (s *SshExecutor) BlockVolumeCreate(host string,
	volume *executors.BlockVolumeRequest) (*executors.BlockVolumeInfo, error) {

	godbc.Require(volume != nil)
	godbc.Require(host != "")
	godbc.Require(volume.Name != "")

	cmd := fmt.Sprintf("gluster-block --create %v --volume %v --host %v --size %v --multipath %v  --block-host %v",
		volume.Name, volume.GlusterVolumeName, volume.GlusterNode,
		volume.Size, volume.Hacount, strings.Join(volume.BlockHosts, ',')

	// Initialize the commands with the create command
	commands := []string{cmd}

	// Execute command
	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 10)
	if err != nil {
		s.BlockVolumeDestroy(host, volume.Name)
		return nil, err
	}

	// TODO: fill the Info?

	return &executors.BlockVolumeInfo{}, nil
}

func (s *SshExecutor) BlockVolumeInfo(host string, volume string, gluster_volume string)
	(*executors.VolumeInfo, error) {

	godbc.Require(volume != nil)
	godbc.Require(host != "")
	godbc.Require(volume.Name != "")

	cmd := fmt.Sprintf("gluster-block --info %v --volume %v", volume, glusterfs_volume)

	commands := []string{cmd}

	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 10)
	if err != nil {
		return nil, err
	}

	// TODO: fill info from output of gluster-block ??!!

	/* example output:
	NAME: sample-block
	VOLUME: block-test
	GBID: 6b60c53c-8ce0-4d8d-a42c-5b546bca3d09
	SIZE: 1073741824
	MULTIPATH: 3
	BLOCK CONFIG NODE(S): 192.168.1.11 192.168.1.12 192.168.1.13
	*/

	return &executors.BlockVolumeInfo{}, nil
}


/*
func (s *SshExecutor) VolumeExpand(host string,
	volume *executors.VolumeRequest) (*executors.VolumeInfo, error) {

	godbc.Require(volume != nil)
	godbc.Require(host != "")
	godbc.Require(len(volume.Bricks) > 0)
	godbc.Require(volume.Name != "")

	commands := s.createAddBrickCommands(volume,
		0, // start at the beginning of the brick list
		inSet,
		maxPerSet)

	// Rebalance if configured
	if s.RemoteExecutor.RebalanceOnExpansion() {
		commands = append(commands,
			fmt.Sprintf("gluster --mode=script volume rebalance %v start", volume.Name))
	}

	// Execute command
	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 10)
	if err != nil {
		return nil, err
	}

	return &executors.VolumeInfo{}, nil
}
*/

func (s *SshExecutor) BlockVolumeDestroy(host string, volume string) error {
	godbc.Require(host != "")
	godbc.Require(volume != "")

	commands := []string{
		fmt.Sprintf("gluster-block --delete %v --volume %v", volume, glusterfs_volume)
	}

	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 10)
	if err != nil {
		logger.LogError("Unable to delete volume %v: %v", volume, err)
	}

	return nil
}

func (s *SshExecutor) BlockVolumeDestroyCheck(host, volume string) error {
	godbc.Require(host != "")
	godbc.Require(volume != "")

	// TODO: do we need checks?

	return nil
}
