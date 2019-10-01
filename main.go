package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/akutz/goof"
	"golang.org/x/exp/mmap"
)

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

type Semaphore interface {

	// Close closes the semaphore.
	Close() error

	// Unlock increments (unlocks) the semaphore pointed to by sem.  If
	// the semaphore's value consequently becomes greater than zero, then
	// another process or thread blocked in a Wait() call will be woken
	// up and proceed to lock the semaphore.
	Unlock() error

	// Wait decrements (locks) the semaphore pointed to by sem.  If
	// the semaphore's value is greater than zero, then the decrement
	// proceeds, and the function returns, immediately.  If the semaphore
	// currently has the value zero, then the call blocks until either it
	// becomes possible to perform the decrement (i.e., the semaphore value
	// ises above zero), or a signal handler interrupts the call.
	Wait() error

	// TryWait is the same as Wait(), except that if the decrement
	// cannot be immediately performed, then call returns an error (errno
	// set to C.EAGAIN) instead of blocking.
	TryWait() error

	// TimedWait is the same as Wait(), except that abs_timeout
	// specifies a limit on the amount of time that the call should block if
	// the decrement cannot be immediately performed.
	TimedWait(timeout *time.Time) error

	// Value returns the current value of the semaphore. If one or more
	// processes or threads are blocked waiting to lock the
	// semaphore with Wait(), POSIX.1 permits two possibilities for the
	// value returned in sval: either 0 is returned; or a negative number
	// whose absolute value is the count of the number of processes and
	// threads currently blocked in Wait().  Linux adopts the former
	// behavior.
	Value() (int, error)
}

// Open creates a new, named semaphore or opens an existing one if one exists
// with the given name.
func Open(name string, excl bool) (Semaphore, error) {
	return open(name, excl)
}

type semaphore struct {
	name  string
	cName *C.char
	sema  unsafe.Pointer
}

func open(name string, excl bool) (Semaphore, error) {
	name = fmt.Sprintf("/%s", name)
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

func main() {
	fmt.Println("Hello World")
	mysem = semaphore.Open("trigger_sem", false)
	//mysem = syscall.sem_open("trigger_sem", O_CREAT, 777, 0)
	var base int64 = 0x35c00000 //1048576 * 768
	//	var c128: complex128 = 0
	var value = Readu32(base, 0)
	fmt.Printf("Hello %x\n", value)
	fmt.Printf("Hello %s/%s\n", runtime.GOOS, runtime.GOARCH)
	//mysem, err := sync.NewSemaphore("trigger_sem", O_CREAT, 777, 0)
	mapped, err := mmap.Open("/usr/share/client")
	if err != nil {
		fmt.Println("Error mmapping: ", err)
	}
	addr_imaging_units := 0x0F8C
	addr_BitsPerPixel := 0x0F98

	for i := 0; i < 4; i++ {
		fmt.Printf("Hello %x\n", mapped.At(addr_BitsPerPixel+i))
	}
	fmt.Printf("Hello %x\n", mapped.At(addr_imaging_units))

	//	__DMA_AND_DATA_SOURCE := "/dev/uio3"
	//  MAIN_MEMORY_ACCESS = /dev/mydevice
	file, err := os.OpenFile("/dev/uio3", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	var baseAddress int64 = 0x1000
	var bufferSize = 4 * 0x100

	mmap2, err := syscall.Mmap(int(file.Fd()), baseAddress, bufferSize, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 32; i++ {
		offset := i * 4
		value := *(*uint32)(unsafe.Pointer(&mmap2[offset]))
		fmt.Printf("Hello %x\n", value) // mapped2.At(0x1000+addr_BitsPerPixel+i))
	}
	__dma_ddr_size_reg := 0x08
	__dma_ddr_base_reg := 0x04
	__DDR_base := *(*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_base_reg]))
	fmt.Printf(" DDR base understood to be at 0x%08x \r\n", __DDR_base)
	DDR_size := *(*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_size_reg]))
	fmt.Printf(" DDR base understood to be at 0x%08x \r\n", DDR_size)

	rbFile, err := os.OpenFile("/dev/mydevice", os.O_RDWR|O_SYNC) //os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	defer rbFile.Close()

	//PROT_READ | PROT_WRITE,  MAP_SHARED , fd_mem, __DDR_offset );
	rbMmap2, err := syscall.Mmap(int(rbfile.Fd()), __DDR_base, DDR_size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}

	RingBuffer := *(*uint32)(unsafe.Pointer(&mmap2[offset]))
}

// pvui open_dma( off_t offset)
// {

// pvui ptr;

// 	/* Open the UIO device file */
// 	fd_dma = open(__DMA_AND_DATA_SOURCE, O_RDWR);
// 	if (fd_dma < 1) {
// 	ptr = 0;
// 	printf("Invalid UIO device file\n\r");
// 	}
// 	else {
// 	/* mmap the UIO device */
// 	//ptr = (pvui)mmap(NULL, UIO_MAP_SIZE, PROT_READ|PROT_WRITE, MAP_SHARED, fd_dma, offset);
// 	ptr = (pvui)mmap(NULL, UIO_MAP_SIZE, PROT_READ|PROT_WRITE, MAP_SHARED, fd_dma, offset);
// 	//*((unsigned *)(ptr + 4)) = 1344;
// 	}
// 	return ptr;
// }
// void close_dma(pvui ptr)
// {

//     munmap((void*)ptr, UIO_MAP_SIZE);
//     close(fd_dma);
// }

func Readu32(baseAddress int64, offset int64) uint32 {
	var value uint32 = 0xFFFFFFFF
	const bufferSize int = 4096
	// 	int __DMA_base = 0x40001000;
	// int __DDR_base = 0x20000000;
	// int __dma_ddr_size_reg = 0x08;
	// int __dma_ddr_head_reg = 0x2c;
	// int __dma_ddr_top_reg = 0x28;
	// int __dma_ddr_base_reg = 0x04;

	// __dma_ddr_base_reg = 0x40001000
	// DDR_base = 0x35c00000
	//__dma_ddr_base_reg = iniparser_getint(ini, "offsets:dma_ddr_base_reg", 0x04)
	//__dma_ddr_size_reg = iniparser_getint(ini, "offsets:dma_ddr_size_reg", 0x08)
	//__dma_ddr_head_reg = iniparser_getint(ini, "offsets:dma_ddr_head_reg", 0x2c)
	//__dma_ddr_top_reg = iniparser_getint(ini, "offsets:dma_ddr_top_reg", 0x28)

	file, err := os.Open("/dev/mem")
	if err != nil {
		log.Fatal(err)
	}

	// fdSharedFile = open("/usr/share/client", O_RDWR);
	// if (fdSharedFile < 1) {
	//   perror("no shared file");

	// } else {
	//   /* mmap the device into memory */

	//   ptrSharedFile = mmap(NULL, page_size * 2, PROT_READ | PROT_WRITE, MAP_SHARED , fdSharedFile, 0x00);

	//   fprintf(stdout,"ASSEMBLER: Mapped %d bytes (2 pages) from /usr/share/client at offset %d \n", page_size * 2, 0x00);
	// }

	// ptrRegBank = ptrSharedFile;

	//   initDevices();
	// //ptrRegBank += 4;

	//   *((uint32_t*)(ptrRegBank + addr_assembler_version )) =  version2();
	//   *((uint32_t*)(ptrRegBank + addr_usleep_send )) = 10;
	//   *((uint32_t*)(ptrRegBank + addr_detector_ready )) = 0;

	// 	__DDR_base = __ptrDMA[__dma_ddr_base_reg / 4];
	// 	fprintf(stdout, " DDR base understood to be at 0x%08x \r\n", __DDR_base );
	// DDR_size = __ptrDMA[__dma_ddr_size_reg / 4];
	// 	fprintf(stdout, " DDR base understood to be at 0x%08x \r\n", DDR_size );

	// 	if(cverbose) fprintf(stdout, " Mapping ring buffer @ 0x%08x 0x%08x \r\n", __DDR_base , DDR_size );

	// 	__DDR_offset = 0;
	// 	fprintf(stdout, "Hamdan: %d\n", __DDR_offset);
	// 	ptrBaton =  (pvui)mmap(NULL, DDR_size  , PROT_READ | PROT_WRITE,  MAP_SHARED , fd_mem, __DDR_offset ); // see detector config or 0x00000

	//	0x38100000 0x07d00000
	defer file.Close()
	mmap, err := syscall.Mmap(int(file.Fd()), baseAddress, bufferSize, syscall.PROT_READ, syscall.MAP_SHARED)

	if err != nil {
		log.Fatal(err)
	}
	value = *(*uint32)(unsafe.Pointer(&mmap[offset]))
	err = syscall.Munmap(mmap)
	if err != nil {
		log.Fatal(err)
	}
	return value
}
