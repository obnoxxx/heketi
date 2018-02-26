# Support for Arbiter volumes


## High level ideas

There are two different modes for arbiter that are neede:

* There are one or more nodes dedicated to only take arbiter bricks.
  The other nodes would host the data bricks. This is targeted for
  scenarios where you have two beefy storage nodes and add a more slim
  node to act as arbiter. Also applicable to situations with two
  data centers an an additional site with just a lightweight node for arbiter.

* Arbiter bricks should be spread throughout the cluster to achieve an
  overall reduction of storage on the nodes.

## Flags on Nodes / Devices

Introduce flags on nodes/devices to mark them as capable of
hosting arbiter bricks, or as only being allowed to host arbiter bricks.

TODO: figure out the details

## Request changes

Either define a new type of durability in heketi volume create request that
reflects "arbiter" or (like gluster does it) add a volume option for the
creation.

Because in kubernetes, heketi's supported durability types are hard-coded,
without changing kubernetes' glusterfs provisioner code significantly (which
is only possible for new releases), we would probably want to go down the
volume option way. This is also closer to what glusterfs offers already.


## Need to refactor the allocator some more

* Now that the ring structure in the simple allocator is always built anew for
  each call of GetNodes(), we can augment the signature to take additional
  aspects of the volume create request as parameters and build different ring
  structures for different requests. (In particular it could take into account
  differently flagges nodes, e.g. those which should taker arbiter bricks...)
* The GetNodes() function should probably return disk sets instead of single
  disks.
* Q: should the allocator already take into account the sizes of the devices?
* Q: I.e. should the code from volume_entry_allocate.go be put into a bigger
  allocator package?


