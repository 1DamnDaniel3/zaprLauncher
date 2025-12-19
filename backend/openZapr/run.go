package openzapr

import (
	"fmt"
	"syscall"
	"unsafe"
)

func RunZaprAsAdmin(programExe string) {

	fmt.Println("PATH RUN: ", programExe)
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(programExe)

	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExecute := shell32.NewProc("ShellExecuteW")

	shellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		0,
		0,
		1,
	)
}

func IsAdmin() bool {
	h, err := syscall.GetCurrentProcess()
	if err != nil {
		return false
	}

	var token syscall.Token
	err = syscall.OpenProcessToken(h, syscall.TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	var elevation uint32
	var outLen uint32

	err = syscall.GetTokenInformation(
		token,
		syscall.TokenElevation,
		(*byte)(unsafe.Pointer(&elevation)),
		uint32(unsafe.Sizeof(elevation)),
		&outLen,
	)
	return err == nil && elevation != 0
}
