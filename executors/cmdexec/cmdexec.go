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
	"os"
	"strconv"
	"sync"

	"github.com/heketi/heketi/pkg/logging"
	rex "github.com/heketi/heketi/pkg/remoteexec"
)

var (
	logger = logging.NewLogger("[cmdexec]", logging.LEVEL_DEBUG)
)

type RemoteCommandTransport interface {
	ExecCommands(host string, commands []string, timeoutMinutes int) (rex.Results, error)
	RebalanceOnExpansion() bool
	SnapShotLimit() int
}

type CmdExecutor struct {
	config      *CmdConfig
	Throttlemap map[string]chan bool
	Lock        sync.Mutex

	RemoteExecutor RemoteCommandTransport
	Fstab          string
	BackupLVM      bool
}

func (c *CmdExecutor) glusterCommand() string {
	return "gluster"
}

func SetWithEnvVariables(config *CmdConfig) {
	var env string

	env = os.Getenv("HEKETI_FSTAB")
	if "" != env {
		config.Fstab = env
	}

	env = os.Getenv("HEKETI_SNAPSHOT_LIMIT")
	if "" != env {
		i, err := strconv.Atoi(env)
		if err == nil {
			config.SnapShotLimit = i
		}
	}
}

func (c *CmdExecutor) InitFromConfig(config *CmdConfig) {
	c.Throttlemap = make(map[string]chan bool)

	if config.Fstab == "" {
		c.Fstab = "/etc/fstab"
	} else {
		c.Fstab = config.Fstab
	}

	c.BackupLVM = config.BackupLVM

	c.config = config
}

func (s *CmdExecutor) AccessConnection(host string) {
	var (
		c  chan bool
		ok bool
	)

	s.Lock.Lock()
	if c, ok = s.Throttlemap[host]; !ok {
		c = make(chan bool, 1)
		s.Throttlemap[host] = c
	}
	s.Lock.Unlock()

	c <- true
}

func (s *CmdExecutor) FreeConnection(host string) {
	s.Lock.Lock()
	c := s.Throttlemap[host]
	s.Lock.Unlock()

	<-c
}

func (s *CmdExecutor) SetLogLevel(level string) {
	switch level {
	case "none":
		logger.SetLevel(logging.LEVEL_NOLOG)
	case "critical":
		logger.SetLevel(logging.LEVEL_CRITICAL)
	case "error":
		logger.SetLevel(logging.LEVEL_ERROR)
	case "warning":
		logger.SetLevel(logging.LEVEL_WARNING)
	case "info":
		logger.SetLevel(logging.LEVEL_INFO)
	case "debug":
		logger.SetLevel(logging.LEVEL_DEBUG)
	}
}

func (s *CmdExecutor) Logger() *logging.Logger {
	return logger
}

func (c *CmdExecutor) RebalanceOnExpansion() bool {
	return c.config.RebalanceOnExpansion
}

func (c *CmdExecutor) SnapShotLimit() int {
	return c.config.SnapShotLimit
}
