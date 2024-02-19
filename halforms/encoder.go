package halforms

import "encoding/json"

type Encoder interface {
	Encode(document Document) ([]byte, error)
}

type jsonEncoder struct{}

// NewJSONEncoder creates a JSON encoder
func NewJSONEncoder() Encoder {
	return new(jsonEncoder)
}

// Encode generates a HAL-FORMS document from provided Document.
func (enc *jsonEncoder) Encode(document Document) ([]byte, error) {
	namedMap := document.ToMap()
	return json.Marshal(namedMap.Content)
}
