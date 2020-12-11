package jailcgo

/*
#cgo CFLAGS: -I /usr/lib
#cgo LDFLAGS: -L. -ljail
#include <stdlib.h>
#include <jail.h>

*/
import "C"
import (
	"fmt"
	"unsafe"
	"strconv"
)

type JailParams struct {
	Params    []C.struct_jailparam
	MapParams map[string]string
}

// TODO: Maybe allow updating also per jail name?
func (jailparams JailParams) Update(jailID int) (int, error) {
	if jailID < 1 {
		return 0, fmt.Errorf("Invalid JailID specified for update: %d", jailID)
	}
	jailid, err := jailparams.SetOrUpdate(2, jailID)
	if err != nil {
		return -1, fmt.Errorf("Error Updating the Jail: %v", err)
	}
	return jailid, nil
}

func (jailparams JailParams) Set() (int, error) {
	jailid, err := jailparams.SetOrUpdate(1, -1)
	if err != nil {
		return -1, fmt.Errorf("Error Updating the Jail: %v", err)
	}
	return jailid, nil
}

func (jailparams JailParams) SetOrUpdate(operation int, jailid int) (int, error) {
	
	var paramsLen int
	pointers := unsafePointers{}

	if operation == 2 {
		paramsLen = len(jailparams.MapParams) + 1 // Need extra one for the jailID
	} else {
		paramsLen = len(jailparams.MapParams) // Need extra one for the jailID
	}
	jailparams.Params = make([]C.struct_jailparam, paramsLen)
	
	// The safest thing is: let's free everything if anything goes wrong :)
	defer C.jailparam_free(&jailparams.Params[0], C.uint(len(jailparams.Params)))
	defer pointers.cleanUnsafePointers()

	var idx int = 0
	
	// If Update, JailID must be the first parameter
	if operation == 2 {
		err := jailparams.addJailParam("jid", strconv.Itoa(jailid), pointers, idx)
		if err != nil {
			return -1, fmt.Errorf("Failed to update the Jail Parameters for Jail %d: %v", jailid, err)
		}
		idx++
	}
	for paramName, paramValue := range jailparams.MapParams {
		err := jailparams.addJailParam(paramName, paramValue, pointers, idx)
		if err != nil {
			return -1, fmt.Errorf("Failed to create the Jail Parameters: %v", err)
		}
		idx++
	}

	paramsSize := C.uint(len(jailparams.Params))
	// TODO: Replace C.int() here for a const :)
	jailID := C.jailparam_set(&jailparams.Params[0], paramsSize, C.int(operation))

	if int(jailID) < 0 {
		return -1, fmt.Errorf("Unexpected error")
	}

	return int(jailID), nil
}


func (jailparams JailParams) addJailParam(key, value string, pointers unsafePointers, index int) error {
	keyPointer := C.CString(key)
	valuePointer := C.CString(value)
	pointers = append(pointers, unsafe.Pointer(keyPointer), unsafe.Pointer(valuePointer))

	retinit := C.jailparam_init(&jailparams.Params[index], keyPointer)
	if int(retinit) < 0 {
		return fmt.Errorf("Failed to add parameter: invalid parameter %s", key)
	}

	retval := C.jailparam_import(&jailparams.Params[index], valuePointer)
	if int(retval) < 0 {
		return fmt.Errorf("Failed to add parameter %s: invalid value %s", key, value)
	}

	return nil
}

func (pointers unsafePointers) cleanUnsafePointers() {
	for _, pointer := range pointers {
		C.free(pointer)
	}
}
