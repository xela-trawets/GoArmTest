package semaphore

import (
	"time"
)

// Semaphore enables processes and threads to synchronize their actions.
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
