Snapshots are read-only copies of volumes (once the snapshot is activated) that
can be used for cloning new volumes from. Its is possible to create a snapshot
through the [Snapshot Create API](../api/api.md#create-a-snapshot) or the
commandline client.

From the command line client, you can type the following to create a snapshot
from an existing volume:

```
$ heketi-cli volume snapshot-create <vol_uuid> [-name=<snap_name>]
```

The new snapshot can be used to create a new volume:

```
$ heketi-cli volume snapshot-clone <snap_uuid> [-name=<clone_name>]
```

The clones of snapshots are new volumes with the same properties as the
original. The cloned volumes can be deleted with the `heketi-cli volume delete`
command. In a similar fashion, snapshots can be removed with the `heketi-cli
volume snapshot-delete` command.

# Proposed CLI

```
$ heketi-cli volume snapshot-create <vol_uuid> [--name=<snap_uuid>] [--description=<string>]
$ heketi-cli volume snapshot-clone <snap_uuid> [--name=<vol_uuid>]
$ heketi-cli volume snapshot-delete <snap_uuid>
$ heketi-cli volume snapshot-list <vol_uuid>
$ heketi-cli volume snapshot-info <snap_uuid>
```

Convenience call for direct volume cloning:

```
$ heketi-cli volume clone <vol_uuid>
```


# API Proposal

The API is layed out here for file volume types only. The same API will be
added for block-volumes at a later time.

### Create a Snapshot
* **Method:** _POST_
* **Endpoint**:`/volumes/{volume_uuid}/snapshots`
* **Content-Type**: `application/json`
* **Response HTTP Status Code**: 202, See [Asynchronous Operations](#asynchronous-operations)
* **Temporary Resource Response HTTP Status Code**: 303, `Location` header will contain `/volumes/{volume_uuid}/snapshots/{snapshot_uuid}`. See [Snapshot Info](#snapshot_info) for JSON response.
* **JSON Request**:
    * name: _string_, _optional_, Name of snapshot. If not provided, the name of the snapshot will be `snap_{id}`, for example `snap_728faa5522838746abce2980`
    * description: _string_, _optional_, Description of the snapshot. If not provided, the description will be empty.
    * Example:

```json
{
    "name": "midnight",
    "description": "nightly snapshot"
}
```

### Clone a Volume from a Snapshot
* **Method:** _POST_
* **Endpoint**:`/volumes/{volume_uuid}/snapshots/{snapshot_uuid}/clone`
* **Content-Type**: `application/json`
* **Response HTTP Status Code**: 202, See [Asynchronous Operations](#asynchronous-operations)
* **Temporary Resource Response HTTP Status Code**: 303, `Location` header will contain `/volumes/{id}`. See [Volume Info](#volume_info) for JSON response.
* **JSON Request**:
    * name: _string_, _optional_, Name of volume. If not provided, the name of the volume will be `snap_{id}`, for example `snap_728faa5522838746abce2980`
    * Example:

```json
{
    "name": "new-vol-from-snap"
}
```

### Delete a Snapshot
* **Method:** _DELETE_
* **Endpoint**:`/volumes/{volume_uuid}/snapshots/{snapshot_uuid}`
* **Response HTTP Status Code**: 202, See [Asynchronous Operations](#async)
* **Temporary Resource Response HTTP Status Code**: 204

### List Snapshots
* **Method:** _GET_
* **Endpoint**:`/volumes/{volume_uuid}/snapshots`
* **Response HTTP Status Code**: 200
* **JSON Response**:
    * snapshots: _array strings_, List of snapshot UUIDs.
    * Example:

```json
{
    "snapshots": [
        "aa927734601288237463aa",
        "70927734601288237463aa"
    ]
}
```

### Get Snapshot Information
* **Method:** _GET_
* **Endpoint**:`/volumes/{volume_uuid}/snapshots/{snapshot_uuid}`
* **Response HTTP Status Code**: 200
* **JSON Request**: None
* **JSON Response**:
    * id: _string_, Snapshot UUID
    * name: _string_, Name of the snapshot
    * description: _string_, Description of the snapshot
    * created: _int_, Seconds after the Epoch when the snapshot was taken
    * volume: _string_, UUID of the volume that this snapshot belongs to
    * cluster: _string_, UUID of cluster which contains this snapshot
    * Example:

```json
{
    "id": "70927734601288237463aa",
    "name": "midnight",
    "description": "nightly snapshot",
    "created": "1518712323",
    "volume": "aa927734601288237463aa",
    "cluster": "67e267ea403dfcdf80731165b300d1ca"
}
```

### Clone a Volume directly

This is a flavor that clones a volume directly.
Implicitly, it creates a snapshot, activates
and clones it, and then deletes the snapshot again.
It is a convenience method for cloning for users
only interested in the clone and not in the snapshots.

* **Method:** _POST_
* **Endpoint**:`/volumes/{volume_uuid}/clone`
* **Content-Type**: `application/json`
* **Response HTTP Status Code**: 202, See [Asynchronous Operations](#asynchronous-operations)
* **Temporary Resource Response HTTP Status Code**: 303, `Location` header will contain `/volumes/{volume_uuid}/snapshots/{snapshot_uuid}`. See [Snapshot Info](#snapshot_info) for JSON response.
* **JSON Request**:
    * name: _string_, _optional_, Name of the clone. If not provided, the name of the snapshot will be `snap_{id}`, for example `snap_728faa5522838746abce2980`
    * description: _string_, _optional_, Description of the snapshot. If not provided, the description will be empty.
    * Example:

```json
{
    "name": "my_clone",
    "description": "my own clone"
}
```


# Future API Extensions


### Activate a Snapshot
* **Method:** _POST_
* **Endpoint**:`/volumes/{volume_uuid}/snapshots/{snapshot_uuid}/activate`

### Deactivate a Snapshot
* **Method:** _POST_
* **Endpoint**:`/volumes/{volume_uuid}/snapshots/{snapshot_uuid}/deactivate`

# Kubernetes Snapshotting Proposal

[Volume
Snapshotting](https://github.com/kubernetes-incubator/external-storage/blob/master/snapshot/doc/volume-snapshotting-proposal.md)
in the Kubernetes external-storage provisioner.

# Gluster Snapshot CLI Reference
```
$ gluster --log-file=/dev/null snapshot help

gluster snapshot commands
=========================

snapshot activate <snapname> [force] - Activate snapshot volume.
snapshot clone <clonename> <snapname> - Snapshot Clone.
snapshot config [volname] ([snap-max-hard-limit <count>] [snap-max-soft-limit <percent>]) | ([auto-delete <enable|disable>])| ([activate-on-create <enable|disable>]) - Snapshot Config.
snapshot create <snapname> <volname> [no-timestamp] [description <description>] [force] - Snapshot Create.
snapshot deactivate <snapname> - Deactivate snapshot volume.
snapshot delete (all | snapname | volume <volname>) - Snapshot Delete.
snapshot help - display help for snapshot commands
snapshot info [(snapname | volume <volname>)] - Snapshot Info.
snapshot list [volname] - Snapshot List.
snapshot restore <snapname> - Snapshot Restore.
snapshot status [(snapname | volume <volname>)] - Snapshot Status.
```

- [Snapshot](https://github.com/gluster/glusterfs-specs/blob/master/done/GlusterFS%203.6/Gluster%20Volume%20Snapshot.md)
- [Cloning](https://github.com/gluster/glusterfs-specs/blob/master/done/GlusterFS%203.7/Clone%20of%20Snapshot.md)
