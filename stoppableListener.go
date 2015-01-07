package martini

import (
	"errors"
	"net"
	"time"
)

// wraps a Listener and adds a channel to allow user to indicate the desire to shut down
type stoppableListener struct {
	*net.TCPListener
	// This channel only exists to indicate the need to shut down (by closing it).
	// No messages are actually passed over it.
	stop chan int
}

var stoppedError = errors.New("Listener stopped.")

// wrap Listener in a StoppableListener
func newStoppableListener(l net.Listener) (*stoppableListener, error) {
	tcpL, ok := l.(*net.TCPListener)
	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}
	retval := &stoppableListener{}
	retval.TCPListener = tcpL
	retval.stop = make(chan int)

	return retval, nil
}

// Hide the listener's Accept() method to incorporate stop checks
func (sl *stoppableListener) Accept() (net.Conn, error) {
	for {
		// set a time out for our acccept operation
		sl.SetDeadline(time.Now().Add(time.Second))

		// accept
		newConn, err := sl.TCPListener.Accept()

		// check to see if we've been told to stop
		select {
		case <-sl.stop:
			return nil, stoppedError
		default:
			// channel still open, continue as normal
		}
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}
		return newConn, err
	}
}

func (sl *stoppableListener) Stop() {
	close(sl.stop)
}
