#!/bin/bash

pushd(){
    command pushd $@ > /dev/null || exit 1
}
popd(){
    command popd $@ > /dev/null
}
# Returns 0 when the first command exists in current system
# install, non-zero otherwise.
command_exists(){
    command -v "$@" > /dev/null
}
# Redirect output to stderr, same usage as 'echo'.
echoe(){
    echo "$@" >&2
}
# Exit execution when current user is not root.
confirm_root(){
    if [[ $(id -u) -ne 0 ]]; then
        echoe "This script needs sudo privilege, use 'sudo $0'"
        exit 1
    fi
}
cleanup(){
    laststatus=$?
    echo
    [[ $laststatus -eq 0 ]] && echo "$0 executed successfully"
    [[ $laststatus -ne 0 ]] && echoe "$0 exited with status $laststatus"
    return $laststatus
}

LB_chrootsuite=focal # aka Ubuntu 20
LB_machinename=lb-judge
LB_chrootdir="/var/lib/machines/$LB_machinename"
LB_chrootsyslock="/var/lib/machines/.#$LB_machinename.lck"
LB_chrootmirror=https://mirrors.aliyun.com/ubuntu/
