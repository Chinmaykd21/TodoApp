package serverErrors

type ErrorCode int

const (
	BodyParse ErrorCode = iota + 1
	ParseInt
	InvalidID
	InsertError
	RetreivalError
	RecordsNotFound
)

func (c ErrorCode) String() string {
	switch c {
	case BodyParse:
		return "Invalid Todo"
	case ParseInt:
		return "Cannot parse the ID"
	case InvalidID:
		return "Invalid ID"
	case InsertError:
		return "Cannot insert document to collection"
	case RetreivalError:
		return "Error while getting the documents from the collection"
	case RecordsNotFound:
		return "Could not find the records to delete"
	default:
		return "An Invalid ErrorCode"
	}
}
