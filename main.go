package main

// #include <fcntl.h>           /* For O_* constants */
// #include <sys/stat.h>        /* For mode constants */
// #include <semaphore.h>
// #include <string.h>
// #include <stdlib.h>
// #include <errno.h>
// #include <stdio.h>
/*
int _errno() {
    return errno;
}
typedef struct {
    sem_t*      val;
    int         err;
} sem_tt;
sem_tt* _sem_open(char* name, int flags) {
    sem_tt* r = (sem_tt*)malloc(sizeof(sem_tt));
    sem_t* sem = sem_open((const char*)name, flags, 0644, 0);
    if (sem == SEM_FAILED) r->err = errno;
    else r->val = sem;
    return r;
}
int _sem_close(void* sem) {
    return sem_close(((sem_tt*)sem)->val) == 0 ? 0 : errno;
}
int _sem_wait(void* sem) {
    return sem_wait(((sem_tt*)sem)->val) == 0 ? 0 : errno;
}
int _sem_trywait(void* sem) {
    return sem_trywait(((sem_tt*)sem)->val) == 0 ? 0 : errno;
}
int _sem_post(void* sem) {
    return sem_post(((sem_tt*)sem)->val) == 0 ? 0 : errno;
}
int _sem_unlink(char* name) {
    return sem_unlink((const char*) name) == 0 ? 0 : errno;
}
int* pInt() {
    int* val = (int*)malloc(sizeof(int));
    return val;
}
const char* cpchar(char* val) {
    return (const char*)val;
}
const struct timespec* new_timespec(time_t sec, long nsec) {
    struct timespec* val = (struct timespec*)malloc(sizeof(struct timespec));
    val->tv_sec = sec;
    val->tv_nsec = nsec;
    return (const struct timespec*)val;
}
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/akutz/goof"
)

type semaphore struct {
	name  string
	cName *C.char
	sema  unsafe.Pointer
}

func openSem(name string, excl bool) (*semaphore, error) {
	name = fmt.Sprintf("%s", name)
	cName := C.CString(name)

	flags := C.O_CREAT
	if excl {
		flags = flags | C.O_EXCL
	}

	sema := C._sem_open(cName, C.int(flags))
	if sema.err != 0 {
		return nil, goof.WithFields(goof.Fields{
			"name":  name,
			"error": sema.err,
		}, "error opening semaphore")
	}

	return &semaphore{
		name:  name,
		cName: cName,
		sema:  unsafe.Pointer(sema),
	}, nil
}

func (s *semaphore) Close() error {
	err := C._sem_close(s.sema)
	if err == 0 {
		return nil
	}
	return goof.WithFields(goof.Fields{
		"name":  s.name,
		"error": int(err),
	}, "error closing semaphore")
}

func (s *semaphore) Wait() error {
	err := C._sem_wait(s.sema)
	if err == 0 {
		return nil
	}
	return goof.WithFields(goof.Fields{
		"name":  s.name,
		"error": int(err),
	}, "error waiting on semaphore")
}

func (s *semaphore) TryWait() error {
	err := C._sem_trywait(s.sema)
	if err == 0 || err == C.EAGAIN {
		return nil
	}
	return goof.WithFields(goof.Fields{
		"name":  s.name,
		"error": int(err),
	}, "error trying wait on semaphore")
}

// func (s *semaphore) TimedWait(t *time.Time) error {
// 	return s.timedWait(t)
// }

func Unlink(name string) error {
	name = fmt.Sprintf("/%s", name)
	cName := C.CString(name)
	err := C._sem_unlink(cName)
	if err == 0 {
		return nil
	}
	return goof.WithFields(goof.Fields{
		"name":  name,
		"error": int(err),
	}, "error unlinking semaphore")
}

func memcpy(dest, src []byte) int {
	n := len(src)
	if len(dest) < len(src) {
		n = len(dest)
	}
	if n == 0 {
		return 0
	}
	C.memcpy(unsafe.Pointer(&dest[0]), unsafe.Pointer(&src[0]), C.size_t(n))
	return n
}

func mapFile(fileName string, base int, size int) (data []byte, err error) {
	mapFile, err := os.OpenFile(fileName, os.O_RDWR, 0755)
	if err != nil {
		return nil, goof.WithFields(goof.Fields{
			"name":  fileName,
			"error": err,
		}, "error opening file")
	}
	defer mapFile.Close()
	data, err = syscall.Mmap(int(mapFile.Fd()), int64(base), size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	return
}
func main() {
	fmt.Println("Hello World")
	mysem, err := openSem("trigger_sem", false)
	if err != nil {
		fmt.Println("Error on openSem: ", err)
		log.Panicln(err)
	}
	fmt.Println("No Error on openSem: ", err)
	defer mysem.Close()

	fmt.Printf("Hello %s/%s\n", runtime.GOOS, runtime.GOARCH)

	regMmap, err := mapFile("/usr/share/client", 0, 2*4096)
	if err != nil {
		log.Panicln(err)
	}
	defer syscall.Munmap(regMmap)

	addr_imaging_units := 0x0F8C
	fmt.Printf("num imaging units %x\n", regMmap[addr_imaging_units])
	var addr_detector_ready = 0x0F60
	*(*int32)(unsafe.Pointer(&regMmap[addr_detector_ready])) = 0
	detector_ready := *(*int)(unsafe.Pointer(&regMmap[addr_detector_ready]))
	fmt.Printf(" addr_detector_ready 0x%08x \r\n", detector_ready)
	*(*int32)(unsafe.Pointer(&regMmap[addr_detector_ready])) = 1
	detector_ready = *(*int)(unsafe.Pointer(&regMmap[addr_detector_ready]))
	fmt.Printf(" addr_detector_ready 0x%08x \r\n", detector_ready)

	mmap2, err := mapFile("/dev/uio3", 0x1000, 0x1000)
	if err != nil {
		fmt.Printf(" Cant map uio3 \r\n")
		log.Panicln(err)
	}
	fmt.Printf(" mapped uio3 at 0x%08x \r\n", unsafe.Pointer(&mmap2[0]))
	defer syscall.Munmap(mmap2)

	for i := 0; i < 4; i++ {
		offset := i * 4
		value := *(*uint32)(unsafe.Pointer(&mmap2[offset]))
		fmt.Printf("Hello %x\n", value)
	}
	__dma_ddr_size_reg := 0x08
	__dma_ddr_base_reg := 0x04
	__dma_ddr_head_reg := 0x2c
	var __DDR_head int64 = int64(*(*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_head_reg])))
	fmt.Printf(" DDR head understood to be at 0x%08x \r\n", __DDR_head)
	var __DDR_base int64 = int64(*(*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_base_reg])))
	fmt.Printf(" DDR base understood to be at 0x%08x \r\n", __DDR_base)
	DDR_size := *(*int)(unsafe.Pointer(&mmap2[__dma_ddr_size_reg]))
	fmt.Printf(" DDR Size understood to be at 0x%08x \r\n", DDR_size)

	//for kernel module map from zero base
	rbMmap, err := mapFile("/dev/mydevice", 0, DDR_size)
	if err != nil {
		fmt.Printf(" cant map module at \r\n")
		log.Panicln(err)
	}
	fmt.Printf(" mapped module at 0x%08x \r\n", unsafe.Pointer(&rbMmap[0]))
	defer syscall.Munmap(rbMmap)

	//TcpServer trigger and addr_detector_ready
	//var addr_detector_ready = 0x0F60
	*(*int32)(unsafe.Pointer(&(regMmap[addr_detector_ready]))) = 0
	detector_ready = *(*int)(unsafe.Pointer(&regMmap[addr_detector_ready]))
	fmt.Printf(" addr_detector_ready 0x%08x \r\n", detector_ready)
	*(*int32)(unsafe.Pointer(&regMmap[addr_detector_ready])) = 1
	detector_ready = *(*int)(unsafe.Pointer(&regMmap[addr_detector_ready]))
	fmt.Printf(" addr_detector_ready 0x%08x \r\n", detector_ready)
	//RingBuffer := (*uint32)(unsafe.Pointer(&rbMmap[0]))
	//DDR_size := *(*int)(unsafe.Pointer(&mmap2[__dma_ddr_size_reg]))
	fmt.Printf(" Awaiting Data 0x%08x \r\n", rbMmap[0])

	//t :=
	mysem.TryWait() //while ...
	//mysem.Wait()

	fmt.Printf(" Sem triggered 0x%08x \r\n")

	rbBase := *(*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_base_reg]))
	pHead := (*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_head_reg]))

	fmt.Printf(" head 0x%08x \r\n", (*pHead)-rbBase)
}

func Readu32(baseAddress int64, offset int64) uint32 {
	var value uint32 = 0xFFFFFFFF
	const bufferSize int = 4096

	file, err := os.Open("/dev/mem")
	if err != nil {
		log.Panicln(err)
	}

	defer file.Close()
	mmap, err := syscall.Mmap(int(file.Fd()), baseAddress, bufferSize, syscall.PROT_READ, syscall.MAP_SHARED)

	if err != nil {
		log.Panicln(err)
	}
	value = *(*uint32)(unsafe.Pointer(&mmap[offset]))
	err = syscall.Munmap(mmap)
	if err != nil {
		log.Panicln(err)
	}
	return value
}
