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

// The primary kind of Number-widget (there's NumberWithText as well)
type Number struct {
	Value         float64
	Text          string
	Prefix        string
	Type          string
	SecondaryStat interface{}
}

func (n Number) MarshalJSON() ([]byte, error) {
	encodedObj, err := json.Marshal(struct {
		Value  float64 `json:"value"`
		Text   string  `json:"text,omitempty"`
		Prefix string  `json:"prefix,omitempty"`
		Type   string  `json:"type,omitempty"`
	}{
		Value:  n.Value,
		Text:   n.Text,
		Prefix: n.Prefix,
		Type:   n.Type,
	})

	if err != nil {
		return encodedObj, err
	}

	// No secondary stat?
	if n.SecondaryStat == nil {
		return []byte(fmt.Sprintf(`{"item":[%s]}`, encodedObj)), nil
	}

	secondaryStat, err := json.Marshal(n.SecondaryStat)

	if err != nil {
		return secondaryStat, err
	}

	return []byte(fmt.Sprintf(`{"item":[%s,%s]}`, encodedObj, secondaryStat)), nil
}

// Trendline implements the geckoboard trendline widget
// https://developer.geckoboard.com/#trendline-example
type TrendlineSecondary []float64
