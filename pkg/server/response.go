package server

type (
	// JSONResponse structure
	JSONResponse struct {
		Data  interface{}            `json:"data"`
		Links map[string]string      `json:"links,omitempty"`
		Meta  map[string]interface{} `json:"meta,omitempty"`
	}
)

// Response factory
func Response(data interface{}) *JSONResponse {
	return &JSONResponse{
		Data: data,
	}
}

// SetMeta helper
func (r *JSONResponse) SetMeta(meta map[string]interface{}) *JSONResponse {
	r.Meta = meta
	return r
}

// SetLinks helper
func (r *JSONResponse) SetLinks(links map[string]string) *JSONResponse {
	r.Links = links
	return r
}
