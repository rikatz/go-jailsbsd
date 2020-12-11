package main

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Params of a Jail
type Params map[string]interface{}

func main() {
	params := make(Params)
	params["name"] = "teste"
	params["host.hostname"] = "teste.bla"
	params["path"] = "/jails/katz"
	params["persist"] = ""
	// TODO: IP and int does not work
	// params["ip4.addr"] = "192.168.0.222"
	// params["securelevel"] = 3

	iov, err := params.buildIovec()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	for k, v := range iov {
		fmt.Printf("%d %v\n", k, v)
	}
	err = getSet(sysJailSet, iov, CreateFlag)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

}

func (p Params) buildIovec() ([]unix.Iovec, error) {
	iovSize := (len(p)) * 2
	iovec := make([]unix.Iovec, iovSize)
	var itr int
	for paramKey, paramValue := range p {

		// Add Key to IOV
		keyBytes, err := unix.ByteSliceFromString(paramKey)
		if err != nil {
			return nil, err
		}
		pKey := (*byte)(unsafe.Pointer(&keyBytes[0]))

		iovec[itr] = unix.Iovec{
			Base: pKey,
			Len:  uint64(len(keyBytes)),
		}
		itr++

		var valLen uint64
		var valBytes []byte
		var pVal *byte
		// Add Value to IOV
		if paramValue != nil {
			switch v := paramValue.(type) {
			case string:
				valBytes, err = unix.ByteSliceFromString(v)
				if err != nil {
					return nil, err
				}
				pVal = (*byte)(unsafe.Pointer(&valBytes[0]))
				valLen = uint64(len(valBytes))
			case int:
				valBytes = IntToByteArray(int64(v))
				fmt.Printf("%v", valBytes)
				pVal = (*byte)(unsafe.Pointer(&valBytes[0]))
				valLen = uint64(len(valBytes))
			}

		}

		iovec[itr] = unix.Iovec{
			Base: pVal,
			Len:  valLen,
		}
		itr++
	}
	return iovec, nil
}

// getSet performas the given syscall with the params and flags provided
func getSet(call int, iov []unix.Iovec, flags uintptr) error {
	_, _, e1 := unix.Syscall(uintptr(call), uintptr(unsafe.Pointer(&iov[0])), uintptr(len(iov)), flags)
	if e1 != 0 {
		switch call {
		case sysJailGet:
			switch int(e1) {
			case ErrJailGetFaultOutsideOfAllocatedSpace:
				return fmt.Errorf("fault outside of allocated space: %d", e1)
			case enoent:
				return fmt.Errorf("jail referred to either does not exist or is inaccessible: %d", e1)
			case einval:
				return fmt.Errorf("invalid param provided: %d", e1)
			}
		case sysJailSet:
			switch int(e1) {
			case eperm:
				return fmt.Errorf("not allowed or restricted: %d", e1)
			case ErrJailSetFaultOutsideOfAllocatedSpace:
				return fmt.Errorf("fault outside of allocated space: %d", e1)
			case ErrJailSetParamNotExist, ErrJailSetParamWrongSize:
				return fmt.Errorf("invalid param provided: %d", e1)
			case ErrJailSetUpdateFlagNotSet:
				return fmt.Errorf("set update flag not set: %d", e1)
			case ErrJailSetNameTooLong:
				return fmt.Errorf("set name too long: %d", e1)
			case ErrJailSetNoIDsLeft:
				return fmt.Errorf("no JID's left: %d", e1)
			}
		}
	}
	return nil
}

func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}
