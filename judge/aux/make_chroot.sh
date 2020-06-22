#!/bin/bash

# Exit when any command fails
set -e

source ./lbsource.sh

# Do cleanup before exit
trap cleanup KILL TERM EXIT

REQUIRED_COMMANDS=(
    debootstrap
    systemd-nspawn
)

# Check for required packages
for r in ${REQUIRED_PACKAGES[@]}; do
    command_exists $r || {
        echoe "Install 'debootstrap', 'systemd-container' first."
        exit 1
    }
done

[[ -e $LB_chrootdir ]] && {
    echoe "Chroot directory '$LB_chrootdir' already exists."
    exit 1
}

confirm_root

# Preinstall these binraies to the chroot enviroment
LB_chrootinc=(
    g++ # CC
    systemd-container # manage through systemd
)
for i in ${LB_chrootinc[@]}; do
    includearg+="$i,"
done

args="--include=$includearg"
components="--components=main,universe"

# Prompt
echo -e Press [enter] to install a \'$LB_chrootsuite\' chroot enviroment \
    into \'$LB_chrootdir\', [^C] to abort:
echo -n '>>> '
read

# Start installation =========================================================
echo "Starting, log file is written to debootstrap.log"
echo
debootstrap $args $components $LB_chrootsuite $LB_chrootdir $LB_chrootmirror \
    | tee debootstrap.log

# Start on boot
# Disable "virtual ethernet" which is enabled by default in service file
# '/lib/systemd/system/systemd-nspawn@.service'
# Reference: https://wiki.archlinux.org/index.php/Systemd-nspawn#Enable_container_on_boot
mkdir -p /etc/systemd/nspawn
cat > /etc/systemd/nspawn/lambow.nspawn <<EOF
[Network]
VirtualEthernet=no
EOF
systemctl enable machines.target
systemctl enable --now systemd-nspawn@$LB_machinename.service

# Show the chroot environment just installed
machinectl list
