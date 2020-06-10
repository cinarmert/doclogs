package container

// Status indicates the instantaneous condition of a Session.
type Status int

const (
	Idle Status = iota
	Live
	Terminated
	Errored
)

// String returns string representation of the receiver Status.
func (s Status) String() string {
	switch s {
	case Idle:
		return "Idle"
	case Live:
		return "Live"
	case Terminated:
		return "Terminated"
	case Errored:
		return "Errored"
	default:
		return "Unknown Status"
	}
}
