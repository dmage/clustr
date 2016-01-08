package daemon

import "unsafe"

// #include <stdio.h>
// #include <libproc.h>
import "C"

func getProcExe(pid int) (string, error) {
	var buf [C.PROC_PIDPATHINFO_MAXSIZE]C.char
	n, err := C.proc_pidpath(C.int(pid), unsafe.Pointer(&buf[0]), C.PROC_PIDPATHINFO_MAXSIZE)
	return C.GoStringN(&buf[0], n), err
}
