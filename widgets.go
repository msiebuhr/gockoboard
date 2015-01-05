package gockoboard

import (
	"encoding/json"
	"fmt"
)

type valueObject struct {
	Value float64 `json:"value"`
}

// Implements a Geck-O-Meter.
// https://developer.geckoboard.com/#geck-o-meter
type GeckOMeter struct {
	Item float64
	Min  float64
	Max  float64
}

// MarshalJSON will marshal the GeckOMeter into JSON.
func (g GeckOMeter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Item float64     `json:"item"`
		Min  valueObject `json:"min"`
		Max  valueObject `json:"max"`
	}{
		Item: g.Item,
		Min:  valueObject{g.Min},
		Max:  valueObject{g.Max},
	})
}

// Trendline implements the geckoboard trendline widget
// https://developer.geckoboard.com/#trendline-example
type Trendline struct {
	Text  string
	Value float64
	Trend []float64
}

// MarshalJSON will marshal the Trendline into JSON.
//
// {"item": [{text and value}, [trendline numbers]]}
func (t Trendline) MarshalJSON() ([]byte, error) {
	encodedObj, err := json.Marshal(struct {
		Text  string  `json:"text"`
		Value float64 `json:"value"`
	}{Text: t.Text, Value: t.Value})

	if err != nil {
		return encodedObj, err
	}

	encodedArray, err := json.Marshal(t.Trend)

	if err != nil {
		return encodedArray, err
	}

	return []byte(fmt.Sprintf(`{"item":[%s,%s]}`, encodedObj, encodedArray)), nil
}
