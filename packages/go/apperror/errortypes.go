package apperror

type Category int32

const (
	CategoryInternal = iota
	CategoryNotFound
	CategoryInvalidArgument
	CategoryConflict
)

type Error struct {
	Category Category
	Code     string
	Details  any
}

func (e Error) Error() string {
	return e.Code
}

func Internal(code string) Error {
	return Error{Category: CategoryInternal, Code: code}
}

func NotFound(code string) Error {
	return Error{Category: CategoryNotFound, Code: code}
}

func InvalidArgument(code string) Error {
	return Error{Category: CategoryInvalidArgument, Code: code}
}

func InvalidArgumentDetails(code string, details map[string]string) Error {
	return Error{Category: CategoryInvalidArgument, Code: code, Details: details}
}

func Conflict(code string) Error {
	return Error{Category: CategoryConflict, Code: code}
}
