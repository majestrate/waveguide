package oauth

type Annotation struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type Post struct {
	Text        string       `json:"text"`
	Annotations []Annotation `json:"annotations,omitempty"`
}
