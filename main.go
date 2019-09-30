package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/exp/mmap"
)

//#include <string.h>
import "C"

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
	//mysem = semaphore.Open("trigger_sem", O_CREAT, 777, 0)
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

	//addr_total_onboard_ram_lower := 0x1080
	//__dma_ddr_size_reg := 0x08
	//__dma_ddr_head_reg := 0x2c
	//__dma_ddr_base_reg := 0x04
	for i := 0; i < 300; i++ {
		//fmt.Printf("Hello %x\n", mapped.At(i*4))
	}

	for i := 0; i < 4; i++ {
		fmt.Printf("Hello %x\n", mapped.At(addr_BitsPerPixel+i))
	}
	fmt.Printf("Hello %x\n", mapped.At(addr_imaging_units))

	//	__DMA_AND_DATA_SOURCE := "/dev/uio3"
	//  MAIN_MEMORY_ACCESS = /dev/mydevice
	file, err := os.Open("/dev/uio3")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	var baseAddress int64 = 0x1000
	var bufferSize = 4 * 0x100

	mmap2, err := syscall.Mmap(int(file.Fd()), baseAddress, bufferSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}
	// mapped2, err := mmap.Open("/dev/uio3")
	// if err != nil {
	// 	fmt.Println("Error mmapping 2: ", err)
	// }
	//__dma_ddr_base_reg := 0x04
	for i := 0; i < 32; i++ {
		offset := i * 4
		value := *(*uint32)(unsafe.Pointer(&mmap2[offset]))
		fmt.Printf("Hello %x\n", value) // mapped2.At(0x1000+addr_BitsPerPixel+i))
	}
	//__DDR_base = __ptrDMA[__dma_ddr_base_reg/4]

	// fd_dma = open(__DMA_AND_DATA_SOURCE, O_RDWR);
	// if (fd_dma < 1) {
	// ptr = 0;
	// printf("Invalid UIO device file\n\r");
	// }
	// else {
	// /* mmap the UIO device */
	// ptr = (pvui)mmap(NULL, UIO_MAP_SIZE, PROT_READ|PROT_WRITE, MAP_SHARED, fd_dma, offset);
	// }

	//*((uint32_t*)(ptrRegBank + addr_detector_ready )) = 1;
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
