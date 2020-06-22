#!/bin/bash

# Exit when any command fails
set -e

source lbsource.sh

confirm_root

# Create unprivileged user 'lambow' in the chroot and disable login
machinectl shell $LB_machinename \
    /usr/sbin/useradd lambow -m -s /usr/sbin/nologin
