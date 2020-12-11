# go-jailsbsd
Library to create and remove FreeBSD in Go language


This is a heavy WIP, Params name can change, structs name can change :)

Please be aware before using.

This library intends to implement both CGO and Syscall for dealing with Jails.

Now it only implements CGO, and in a future some PoC of Syscalls will also be added.

## Using

Examples can be found in examples/cgo

## TODO

* Find some better naming for the struct things
* Unit tests (always....)
* Deal with non-kernel parameters, such as:
  * exec.start
  * interface
  * etc
* Remove Jail method
* Verify a way to create an 'attach' (allow logging to the Jail)
* Verify (in this package or other) some way to deal with image layering (maybe ZFS + Snapshot of a base image, and exec.start as the way to build the image)
* Verify later if this can be used on a broader way with, let's say a FreeBSD CNI that creates interfaces, pf rules and then attach a jail to it :)