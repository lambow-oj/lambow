# judge/aux

Auxillary scripts for setting up the judge host. All scripts must be executed
under this directory (`judge/aux`)

- [`make_chroot.sh`](./make_chroot.sh)
  Create a chroot environment where all judging processes will run in.
- [`clean.sh`](./clean.sh)
  Remove the chroot environment from system entirely.
- [`setup.sh`](./setup.sh)
  Setup utilities in the chroot envrionment.
