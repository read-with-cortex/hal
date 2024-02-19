package hal

import (
	"encoding/json"
)

// Encoder to encode a Resource into a valid HAL document.
type Encoder interface {
	Encode(resource Resource) ([]byte, error)
}

type jsonEncoder struct{}

// NewJSONEncoder creates a JSON encoder
func NewJSONEncoder() Encoder {
	return new(jsonEncoder)
}

// Encode generates a HAL document from provided Resource.
func (enc *jsonEncoder) Encode(resource Resource) ([]byte, error) {
	namedMap := resource.ToMap()
	return json.Marshal(namedMap.Content)
}
