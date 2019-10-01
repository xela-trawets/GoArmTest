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

func main() {
	fmt.Println("Hello World")
	mysem, err := openSem("trigger_sem", false)
	if err != nil {
		fmt.Println("Error on openSem: ", err)
		log.Fatal(err)
	}
	fmt.Println("No Error on openSem: ", err)
	defer mysem.Close()
	//mysem = syscall.sem_open("trigger_sem", O_CREAT, 777, 0)
	var base int64 = 0x35c00000 //1048576 * 768
	//	var c128: complex128 = 0
	var value = Readu32(base, 0)
	fmt.Printf("Hello %x\n", value)
	fmt.Printf("Hello %s/%s\n", runtime.GOOS, runtime.GOARCH)
	//mysem, err := sync.NewSemaphore("trigger_sem", O_CREAT, 777, 0)
	regMapFile, err := os.OpenFile("/usr/share/client", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer regMapFile.Close()
	regMmap, err := syscall.Mmap(int(regMapFile.Fd()), 0, 2*4096, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Munmap(regMmap)
	addr_imaging_units := 0x0F8C
	// addr_BitsPerPixel := 0x0F98

	// for i := 0; i < 4; i++ {
	// 	fmt.Printf("Hello %x\n", regMmap[addr_BitsPerPixel+i])
	// }
	fmt.Printf("Hello %x\n", regMmap[addr_imaging_units])
	var addr_detector_ready = 0x0F60
	//*(*int32)(unsafe.Pointer(&(regMmap[addr_detector_ready]))) = 0
	detector_ready := *(*int)(unsafe.Pointer(&regMmap[addr_detector_ready]))
	fmt.Printf(" addr_detector_ready 0x%08x \r\n", detector_ready)
	*(*int32)(unsafe.Pointer(&regMmap[addr_detector_ready])) = 1
	detector_ready = *(*int)(unsafe.Pointer(&regMmap[addr_detector_ready]))
	fmt.Printf(" addr_detector_ready 0x%08x \r\n", detector_ready)

	//	__DMA_AND_DATA_SOURCE := "/dev/uio3"
	//  MAIN_MEMORY_ACCESS = /dev/mydevice
	file, err := os.OpenFile("/dev/uio3", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" opened uio3 at 0x%08x \r\n", file.Fd)

	defer file.Close()
	var baseAddress int64 = 0x1000
	var bufferSize = 1 * 0x1000

	mmap2, err := syscall.Mmap(int(file.Fd()), baseAddress, bufferSize, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Printf(" Cant map uio3 \r\n")
		log.Fatal(err)
	}
	fmt.Printf(" mapped uio3 at 0x%08x \r\n", unsafe.Pointer(&mmap2[0]))
	defer syscall.Munmap(mmap2)

	for i := 0; i < 4; i++ {
		offset := i * 4
		value := *(*uint32)(unsafe.Pointer(&mmap2[offset]))
		fmt.Printf("Hello %x\n", value) // mapped2.At(0x1000+addr_BitsPerPixel+i))
		//log.Fatal(err)
	}
	__dma_ddr_size_reg := 0x08
	__dma_ddr_base_reg := 0x04
	__dma_ddr_head_reg := 0x2c
	var __DDR_base int64 = int64(*(*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_base_reg])))
	fmt.Printf(" DDR base understood to be at 0x%08x \r\n", __DDR_base)
	DDR_size := *(*int)(unsafe.Pointer(&mmap2[__dma_ddr_size_reg]))
	fmt.Printf(" DDR Size understood to be at 0x%08x \r\n", DDR_size)

	rbFile, err := os.OpenFile("/dev/mydevice", os.O_RDWR|os.O_SYNC, 0755)
	if err != nil {
		fmt.Printf(" My Device Not Opened 0x%08x \r\n", rbFile.Fd)
		log.Fatal(err)
	}

	fmt.Printf(" My Device Opened 0x%08x \r\n", rbFile.Fd)
	defer rbFile.Close()

	//PROT_READ | PROT_WRITE,  MAP_SHARED , fd_mem, __DDR_offset );
	//rbMmap2, err := syscall.Mmap(int(rbFile.Fd()), __DDR_base, DDR_size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	//using kernel module
	rbMmap, err := syscall.Mmap(int(rbFile.Fd()), 0, DDR_size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Printf(" cant map module at \r\n")
		log.Fatal(err)
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

	//mysem.TryWait()//while ...
	mysem.Wait()

	fmt.Printf(" Sem triggered 0x%08x \r\n", rbMmap[0])

	rbBase := *(*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_base_reg]))
	pHead := (*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_head_reg]))

	fmt.Printf(" head 0x%08x \r\n", *pHead-rbBase)
	//*(*uint32)(unsafe.Pointer(&mmap2[__dma_ddr_base_reg])))
	//mysem.Close()
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
