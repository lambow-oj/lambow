#!/bin/bash

# Exit when any command fails
set -e

source lbsource.sh

# Do cleanup before exit
trap cleanup KILL TERM EXIT

confirm_root

systemctl --system disable --now systemd-nspawn@$LB_machinename.service
systemctl --system daemon-reload
[[ -f /etc/systemd/nspawn/$LB_machinename.nspawn ]] && \
    rm -v /etc/systemd/nspawn/$LB_machinename.nspawn || true
[[ -d $LB_chrootdir ]] && rm -vr $LB_chrootdir || true
[[ -f $LB_chrootsyslock ]] && rm -v $LB_chrootsyslock || true
