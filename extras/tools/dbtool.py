#!/usr/bin/python

import json
import sys


INFO = 'Info'


def report(otype, oid, *msg):
    print ('{} {}: {}'.format(otype, oid, ': '.join(msg)))


def check_cluster(data, cid, cluster):
    if cid != cluster[INFO]['id']:
        report('Cluster', cid, 'id mismatch', cluster[INFO]['id'])
    for vid in cluster[INFO]['volumes']:
        if vid not in data['volumeentries']:
            report('Cluster', cid, 'unknown volume', vid)
    for vid in cluster[INFO]['blockvolumes']:
        if vid not in data['blockvolumeentries']:
            report('Cluster', cid, 'unknown block volume', vid)
    # check nodes
    for nid in cluster[INFO]['nodes']:
        if nid not in data['nodeentries']:
            report('Cluster', cid, 'unknown node', nid)


def check_brick(data, bid, brick):
    if bid != brick[INFO]['id']:
        report('Brick', bid, 'id mismatch', brick[INFO]['id'])
    # check brick points to real device
    if brick[INFO]['device'] not in data['deviceentries']:
        report('Brick', bid, 'device mismatch', brick[INFO]['device'])
    # check brick points to real node
    if brick[INFO]['node'] not in data['nodeentries']:
        report('Brick', bid, 'node mismatch', brick[INFO]['node'])
    # check brick path is not empty
    if brick[INFO]['path'] == '':
        report('Brick', bid, 'has empty path', '')
    # check brick entry links back to a volume
    if brick[INFO]['volume'] not in data['volumeentries']:
        report('Brick', bid, 'unknown volume', brick[INFO]['node'])

def check_volume(data, vid, volume):
    if vid != volume[INFO]['id']:
        report('Volume', vid, 'id mismatch', volume[INFO]['id'])
    for bid in volume['Bricks']:
        if bid not in data['brickentries']:
            report('Volume', vid, 'unknown brick', bid)
    if volume[INFO]['cluster'] not in data['clusterentries']:
        report('Volume', vid, 'unknown cluster', volume[INFO]['cluster'])
    # check block volumes
    if volume[INFO]['blockinfo']:
        for bvid in volume[INFO]['blockinfo']['blockvolume']:
            if bvid not in data['blockvolumes']:
                report('Volume', vid, 'unknown blockvolume', bvid)

# check_blockvolume
    # ...

# check_node
    # check devices
    # check cluster

def check_device(data, did, device):
    if did != device[INFO]['id']:
        report('Volume', did, 'id mismatch', device[INFO]['id'])
    for bid in device['Bricks']:
        if bid not in data['brickentries']:
            report('Device', did, 'unknown brick', bid)
    # check NodeId

def check_pendingop(data, pid, pendingop):
    if pid != pendingop['Id']:
        report('Pending op', pid, 'id mismatch', pendingop['Id'])

def check_db(data):
    referenced_pids = set()

    num_clusters = 0
    for cid, c in data['clusterentries'].items():
        check_cluster(data, cid, c)
        num_clusters += 1
    print(str(num_clusters) + " clusters")

    num_nodes = 0
    for nid, n in data['nodeentries'].items():
        num_nodes += 1
    print(str(num_nodes) + " nodes")

    num_v = 0
    num_v_pending = 0
    for vid, v in data['volumeentries'].items():
        check_volume(data, vid, v)
        num_v += 1
        if v['Pending']['Id'] != '':
            num_v_pending += 1
            referenced_pids.add(v['Pending']['Id'])

    print(str(num_v) + " volumes (" + str(num_v_pending) + " pending)")

    num_bv = 0
    num_bv_pending = 0
    for bvid, bv in data['blockvolumeentries'].items():
        num_bv += 1
        if bv['Pending']['Id'] != '':
            num_bv_pending += 1
            referenced_pids.add(bv['Pending']['Id'])
    print(str(num_bv) + " blockvolumes (" + str(num_bv_pending) + " pending)")

    num_devices = 0
    for did, d in data['deviceentries'].items():
        check_device(data, did, d)
        num_devices += 1
    print(str(num_devices) + " devices")

    num_bricks = 0
    num_bricks_pending = 0
    for bid, b in data['brickentries'].items():
        check_brick(data, bid, b)
        num_bricks += 1
        if b['Pending']['Id'] != '':
            num_bricks_pending += 1
            referenced_pids.add(b['Pending']['Id'])
    print(str(num_bricks) + " bricks (" + str(num_bricks_pending) + " pending)")

    num_pendingops = 0
    num_po_per_type = {}
    num_pids_unreferenced = 0
    num_po_unref_per_type = {}
    used_po_types = set()
    for i in range(0, 7):
        num_po_per_type[i] = 0
        num_po_unref_per_type[i] = 0
    for pid, po in data['pendingoperations'].items():
        check_pendingop(data, pid, po)
        num_pendingops += 1
        num_po_per_type[po['Type']] += 1
        used_po_types.add(po['Type'])
        if pid not in referenced_pids:
            num_pids_unreferenced += 1
            num_po_unref_per_type[po['Type']] += 1
    print(str(num_pendingops) + " pending operations (" +
            str(num_pids_unreferenced) + " unreferenced)")
    for t in used_po_types:
        print("type " + str(t) + ": " + str(num_po_per_type[t]) + " pending operations ("
                + str(num_po_unref_per_type[t]) + " unreferenced)")

try:
    filename = sys.argv[1]
except IndexError:
    sys.stderr.write("error: filename required")
    sys.exit(2)

with open(filename) as fh:
    data = json.load(fh)

check_db(data)
