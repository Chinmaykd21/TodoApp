package serverErrors

type ErrorCode int

const (
	BodyParse ErrorCode = iota + 1
	ParseInt
	InvalidID
)

func (c ErrorCode) String() string {
	switch c {
	case BodyParse:
		return "Invalid Todo"
	case ParseInt:
		return "Cannot parse the ID"
	case InvalidID:
		return "Invalid ID"
	default:
		return "An Invalid ErrorCode"
	}
}
