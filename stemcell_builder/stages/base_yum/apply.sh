#!/usr/bin/env bash
#
# Copyright (c) 2009-2012 VMware, Inc.

set -e

base_dir=$(readlink -nf $(dirname $0)/../..)
source $base_dir/lib/prelude_apply.bash

# Bootstrap the base system
echo "Running yum"

yum -y --installroot=$chroot --releasever=/ groupinstall Base
yum -y --installroot=$chroot --releasever=/ groupinstall 'Development Tools'

# Configuring Repositories
cp /etc/yum.repos.d/{CentOS-Base.repo,epel.repo} $chroot/etc/yum.repos.d
cp /etc/pki/rpm-gpg/{RPM-GPG-KEY-CentOS-6,RPM-GPG-KEY-EPEL-6} $chroot/etc/pki/rpm-gpg

# add epel rpm package (not needed as we do the above hack)
# yum -y --installroot=$chroot --releasever=/ install .../epel-release-6-8.noarch.rpm

yum -y --installroot=$chroot --releasever=/ \
  install libyaml-devel libxslt-devel zlib-devel openssl-devel sudo readline-devel libffi-devel openssh-server

# these are the deb packages that are added
debs="build-essential libssl-dev lsof \
strace bind9-host dnsutils tcpdump iputils-arping \
curl wget libcurl3 libcurl3-dev bison libreadline6-dev \
zip unzip \
nfs-common flex psmisc apparmor-utils iptables sysstat \
rsync openssh-server traceroute libncurses5-dev quota \
libaio1 gdb tripwire libcap2-bin "

# add /etc/resolv.conf temporarily

# setup runit/daemontools so that it can be used to run the bosh agent (and restart it if it crashes)

touch $chroot/etc/sysconfig/network
