package ginutil

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	// Meta    *Meta      `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Internal bool   `json:"internal"`
	Key      string `json:"key"`
}

// type Meta struct {
// 	Page       int `json:"page,omitempty"`
// 	PerPage    int `json:"per_page,omitempty"`
// 	Total      int `json:"total,omitempty"`
// 	TotalPages int `json:"total_pages,omitempty"`
// }
