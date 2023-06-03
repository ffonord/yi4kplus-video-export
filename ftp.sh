#!/bin/sh
#create file to test if ftp.sh is execute
date > /tmp/fuse_d/ftp_started_at
#start ftp
tcpsvd -u root -vE 0.0.0.0 21 ftpd -w /tmp/fuse_d/ >/dev/null 2>&1