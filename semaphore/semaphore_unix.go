// +build darwin

package semaphore

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
	"time"
	"unsafe"

	"github.com/akutz/goof"
)

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

func (s *semaphore) Unlock() error {
	err := C._sem_post(s.sema)
	if err == 0 {
		return nil
	}
	return goof.WithFields(goof.Fields{
		"name":  s.name,
		"error": int(err),
	}, "error unlocking semaphore")
}

func (s *semaphore) Value() (int, error) {
	return s.value()
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

func (s *semaphore) TimedWait(t *time.Time) error {
	return s.timedWait(t)
}

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
