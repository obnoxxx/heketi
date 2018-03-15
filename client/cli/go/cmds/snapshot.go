//
// Copyright (c) 2018 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package cmds

import (
	//"encoding/json"
	"errors"
	"fmt"
	//"os"
	//"strings"

	//client "github.com/heketi/heketi/client/api/go-client"
	//"github.com/heketi/heketi/pkg/glusterfs/api"
	//"github.com/heketi/heketi/pkg/kubernetes"
	"github.com/spf13/cobra"
)

var (
	clonename string
)

func init() {
	RootCmd.AddCommand(snapshotCommand)

	snapshotCommand.AddCommand(snapshotDeleteCommand)
	snapshotDeleteCommand.SilenceUsage = true

	snapshotCommand.AddCommand(snapshotCloneCommand)
	snapshotCloneCommand.Flags().StringVar(&clonename, "name", "",
		"\n\tOptional: Name of the newly cloned volume. Only set if really necessary")
	snapshotCloneCommand.SilenceUsage = true

	snapshotCommand.AddCommand(snapshotInfoCommand)
	snapshotInfoCommand.SilenceUsage = true

	snapshotCommand.AddCommand(snapshotListCommand)
	snapshotListCommand.SilenceUsage = true
}

var snapshotCommand = &cobra.Command{
	Use:   "snapshot",
	Short: "Heketi Snapshot Management",
	Long:  "Heketi Snapshot Management",
}

var snapshotDeleteCommand = &cobra.Command{
	Use:     "delete",
	Short:   "Deletes the snapshot",
	Long:    "Deletes the snapshot",
	Example: "  $ heketi-cli snapshot delete 886a86a868711bef83001",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ensure proper number of args
		s := cmd.Flags().Args()
		if len(s) < 1 {
			return errors.New("Snapshot id missing")
		}

		snapshotId := cmd.Flags().Arg(0)

		// TODO: implement heketi.VolumeSnapshotDelete()
		//heketi := client.NewClient(options.Url, options.User, options.Key)
		//err := heketi.VolumeSnapshotDelete(snapshotId)
		err := errors.New("delete is not implemented yet")
		if err == nil {
			fmt.Fprintf(stdout, "Snapshot %s deleted\n", snapshotId)
		}

		return err
	},
}

var snapshotCloneCommand = &cobra.Command{
	Use:     "clone",
	Short:   "Clones a snapshot into a new volume",
	Long:    "Clones a snapshot into a new volume",
	Example: "  $ heketi-cli snapshot clone 886a86a868711bef83001",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := cmd.Flags().Args()
		if len(s) < 1 {
			return errors.New("Snapshot id missing")
		}

		//snapshotId := cmd.Flags().Arg(0)
		//heketi := client.NewClient(options.Url, options.User, options.Key)
		//err := heketi.VolumeSnapshotClone(snapshotId)

		err := errors.New("clone is not implemented yet")

		return err
	},
}

var snapshotInfoCommand = &cobra.Command{
	Use:     "info",
	Short:   "Shows the information of a snapshot",
	Long:    "Shows the information of a snapshot",
	Example: "  $ heketi-cli snapshot info 886a86a868711bef83001",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := cmd.Flags().Args()
		if len(s) < 1 {
			return errors.New("Snapshot id missing")
		}

		//snapshotId := cmd.Flags().Arg(0)
		//heketi := client.NewClient(options.Url, options.User, options.Key)
		//snapshot, err := heketi.VolumeSnapshotInfo(snapshotId)

		err := errors.New("info is not implemented yet")

		return err
	},
}

var snapshotListCommand = &cobra.Command{
	Use:     "list",
	Short:   "Lists all snapshots",
	Long:    "Lists all snapshots",
	Example: "  $ heketi-cli snapshot list",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := cmd.Flags().Args()
		if len(s) != 1 {
			return errors.New("list command does not expect arguments")
		}

		//volumeId := cmd.Flags().Arg(0)
		//heketi := client.NewClient(options.Url, options.User, options.Key)
		//snapshots, err := heketi.VolumeSnapshotList(snapshotId)

		err := errors.New("list is not implemented yet")

		return err
	},
}
