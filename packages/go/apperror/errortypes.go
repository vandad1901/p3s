package apperror

type Category int32

const (
	CategoryNotFound = iota + 1
	CategoryInvalidArgument
	CategoryConflict
	CategoryUnauthenticated
)

type Error struct {
	Category Category
	Code     string
	Details  any
}

func (e Error) Error() string {
	return e.Code
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

func Unauthenticated(code string) Error {
	return Error{Category: CategoryUnauthenticated, Code: code}
}
