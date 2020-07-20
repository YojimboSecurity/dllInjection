package src

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/contester/runlib/win32"
)

const PROCESS_ALL_ACCESS = 0x1F0FFF

// const PAGE_READWRITE = 0x04
const Virtual_Mem = 0x1000 | 0x2000

func DLLInjection(pid int16, dll string) {
	fmt.Println("Start")

	dllLen := uint32(len(dll))
	unsPointer := unsafe.Pointer(&dll)

	handle, err := syscall.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pid))
	defer syscall.CloseHandle(handle)
	if err != nil {
		log.Fatalf("OpenProcess %v", err)
	}
	argAddress, err := win32.VirtualAllocEx(handle, 0, dllLen, Virtual_Mem, syscall.PAGE_READWRITE)
	if err != nil {
		log.Fatalf("VirtualAllocEx %v", err)
	}
	_, err = win32.WriteProcessMemory(handle, argAddress, unsPointer, dllLen)
	if err != nil {
		log.Fatalf("WriteProcessMemory %v", err)
	}

	k32uint := syscall.StringToUTF16Ptr("kernel32.dll")

	hKernel32, err := win32.GetModuleHandle(k32uint)
	if err != nil {
		log.Fatalf("GetModuleHandle %v", err)
	}
	hLoadLib, err := syscall.GetProcAddress(hKernel32, "LoadLibraryA")
	if err != nil {
		log.Fatalf("GetProcAddress %v", err)
	}
	_, threadID, err := win32.CreateRemoteThread(handle, nil, 0, hLoadLib, argAddress, 0)
	if err != nil {
		log.Fatalf("CreateRemoteThread %v", err)
	}
	fmt.Println("Remote thread with id %v created", threadID)
}
