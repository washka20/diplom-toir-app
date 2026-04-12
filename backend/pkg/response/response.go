package response

// Response представляет стандартный envelope для API ответов.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// Meta содержит метаданные пагинации.
type Meta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

// Success возвращает успешный ответ с данными.
func Success(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// Error возвращает ответ с ошибкой.
func Error(msg string) Response {
	return Response{
		Success: false,
		Error:   msg,
	}
}

// Paginated возвращает успешный ответ с данными и метаданными пагинации.
func Paginated(data interface{}, page, perPage int, total int64) Response {
	return Response{
		Success: true,
		Data:    data,
		Meta: Meta{
			Page:    page,
			PerPage: perPage,
			Total:   total,
		},
	}
}
