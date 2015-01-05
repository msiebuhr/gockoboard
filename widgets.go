package gockoboard

import (
    "encoding/json"
)

type valueObject struct {
    Value float64 `json:"value"`
}

// Implements a Geck-O-Meter.
// https://developer.geckoboard.com/#geck-o-meter
type GeckOMeter struct {
    Item float64
    Min  float64
    Max float64
}

// MarshalJSON will marshal the GeckOMeter into JSON.
func (g GeckOMeter) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Item float64 `json:"item"`
        Min valueObject `json:"min"`
        Max valueObject `json:"max"`
    }{
        Item: g.Item,
        Min: valueObject{g.Min},
        Max: valueObject{g.Max},
    })
}
