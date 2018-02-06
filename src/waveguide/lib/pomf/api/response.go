package api

type Response struct {
	Success     bool       `json:"success"`
	ErrorCode   int        `json:"errorcode"`
	Description string     `json:"description"`
	Files       []Resource `json:"files"`
}

/** SetError sets this response error code and message */
func (r *Response) SetError(code int, err error) {
	if err != nil {
		r.Success = false
		r.ErrorCode = code
		r.Description = err.Error()
	}
}

func (r *Response) AddFile(file Resource) {
	r.Files = append(r.Files, file)
}
