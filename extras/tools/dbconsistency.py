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


def check_brick(data, bid, brick):
    if bid != brick[INFO]['id']:
        report('Brick', bid, 'id mismatch', brick[INFO]['id'])
    # check brick points to real device
    if brick[INFO]['device'] not in data['deviceentries']:
        report('Brick', bid, 'device mismatch', brick[INFO]['device'])
    # check brick points to real node
    if brick[INFO]['node'] not in data['nodeentries']:
        report('Brick', bid, 'node mismatch', brick[INFO]['node'])


def check_volume(data, vid, volume):
    if vid != volume[INFO]['id']:
        report('Volume', vid, 'id mismatch', volume[INFO]['id'])
    for bid in volume['Bricks']:
        if bid not in data['brickentries']:
            report('Volume', vid, 'unknown brick', bid)


def check_device(data, did, device):
    if did != device[INFO]['id']:
        report('Volume', did, 'id mismatch', device[INFO]['id'])
    for bid in device['Bricks']:
        if bid not in data['brickentries']:
            report('Device', did, 'unknown brick', bid)


def check_db(data):
    for cid, c in data['clusterentries'].items():
        check_cluster(data, cid, c)

    for vid, v in data['volumeentries'].items():
        check_volume(data, vid, v)

    for did, d in data['deviceentries'].items():
        check_device(data, did, d)

    for bid, b in data['brickentries'].items():
        check_brick(data, bid, b)


try:
    filename = sys.argv[1]
except IndexError:
    sys.stderr.write("error: filename required")
    sys.exit(2)

with open(filename) as fh:
    data = json.load(fh)

check_db(data)
