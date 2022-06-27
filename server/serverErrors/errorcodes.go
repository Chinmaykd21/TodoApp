package serverErrors

type ErrorCode int

const (
	BodyParse ErrorCode = iota + 1
	ParseInt
	InvalidID
	InsertError
	RetreivalError
	RecordsNotFound
	UpdateError
)

func (c ErrorCode) String() string {
	switch c {
	case BodyParse:
		return "Invalid Todo"
	case ParseInt:
		return "Cannot parse ID"
	case InvalidID:
		return "Invalid ID"
	case InsertError:
		return "Cannot insert document to mongoDB collection"
	case RetreivalError:
		return "Error while getting the documents from mongoDB collection"
	case RecordsNotFound:
		return "Could not find the docuemtns to delete in mongoDB collection"
	case UpdateError:
		return "Error while updating todo document in mongoDB collection"
	default:
		return "An Invalid ErrorCode"
	}
}
