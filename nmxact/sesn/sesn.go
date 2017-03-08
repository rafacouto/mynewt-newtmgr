package sesn

import (
	"fmt"
	"time"

	"mynewt.apache.org/newt/nmxact/nmp"
	"mynewt.apache.org/newt/nmxact/xport"
)

type TxOptions struct {
	Timeout time.Duration
	Tries   int
}

func (opt *TxOptions) AfterTimeout() <-chan time.Time {
	if opt.Timeout == 0 {
		return nil
	} else {
		return time.After(opt.Timeout)
	}
}

// Represents a communication session with a specific peer.  The particulars
// vary according to protocol and transport. Several Sesn instances can use the
// same Xport.
type Sesn interface {
	// Initiates communication with the peer.  For connection-oriented
	// transports, this creates a connection.
	Open() error

	// Ends communication with the peer.  For connection-oriented transports,
	// this closes the connection.
	Close() error

	// Retrieves the maximum data payload for outgoing NMP requests.
	MtuOut() int

	// Retrieves the maximum data payload for incoming NMP responses.
	MtuIn() int

	// Transmits a single NMP message and listens for the response.  Blocking.
	TxNmpOnce(msg *nmp.NmpMsg, opt TxOptions) (nmp.NmpRsp, error)

	// Stops a receive operation in progress.  This must be called from a
	// separate thread, as sesn receive operations are blocking.
	AbortRx(nmpSeq uint8) error
}

// Represents an NMP timeout; request sent, but no response received.
type TimeoutError struct {
	Text string
}

func NewTimeoutError(text string) *TimeoutError {
	return &TimeoutError{
		Text: text,
	}
}

func FmtTimeoutError(format string, args ...interface{}) *TimeoutError {
	return NewTimeoutError(fmt.Sprintf(format, args...))
}

func (e *TimeoutError) Error() string {
	return e.Text
}

func IsTimeout(err error) bool {
	_, ok := err.(*TimeoutError)
	return ok
}

// Represents an NMP timeout; request sent, but no response received.
type DisconnectError struct {
	Text string
}

func NewDisconnectError(text string) *DisconnectError {
	return &DisconnectError{
		Text: text,
	}
}

func FmtDisconnectError(format string, args ...interface{}) *DisconnectError {
	return NewDisconnectError(fmt.Sprintf(format, args...))
}

func (e *DisconnectError) Error() string {
	return e.Text
}

func IsDisconnect(err error) bool {
	_, ok := err.(*DisconnectError)
	return ok
}

func TxNmp(s Sesn, m *nmp.NmpMsg, o TxOptions) (nmp.NmpRsp, error) {
	retries := o.Tries - 1
	for i := 0; ; i++ {
		r, err := s.TxNmpOnce(m, o)
		if err == nil {
			return r, nil
		}

		if (!IsTimeout(err) && !xport.IsTimeout(err)) || i >= retries {
			return nil, err
		}
	}
}
