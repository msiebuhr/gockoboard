package gockoboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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

// Leaderboard implements a ranked list of items with optional values
// https://developer.geckoboard.com/#leaderboard
type Leaderboard struct {
	Items  []LeaderboardItem `json:"item"`
	Format LeaderboardFormat `json:"format,omitempty"`
	Unit   string            `json:"unit,omitempty"`
}

func NewLeaderboard(boards ...LeaderboardItem) Leaderboard {
	return Leaderboard{Items: boards}
}

// Find the previous rank of each element and set it in the leaderboard
func (l *Leaderboard) CalculatePreviousRanks(old Leaderboard) {
	for i, item := range l.Items {
		item.PreviousRank = 0
		// Find it's previous rank
		for j, jtem := range old.Items {
			if item.Label == jtem.Label {
				l.Items[i].PreviousRank = float64(j + 1)
				break
			}
		}
	}
}

// Implement sort.Interface, to make Leaderboards sortable
func (l Leaderboard) Len() int { return len(l.Items) }

// Implement sort.Interface, to make Leaderboards sortable
func (l Leaderboard) Swap(i, j int) {
	l.Items[i], l.Items[j] = l.Items[j], l.Items[i]
}

// Implement sort.Interface, to make Leaderboards sortable
func (l Leaderboard) Less(i, j int) bool {
	if l.Items[i].Value == l.Items[j].Value {
		return sort.StringsAreSorted([]string{
			l.Items[i].Label,
			l.Items[j].Label,
		})
	}
	return l.Items[i].Value < l.Items[j].Value
}

type LeaderboardFormat string

const (
	LEADERBOARD_FORMAT_DEFAULT  LeaderboardFormat = ""
	LEADERBOARD_FORMAT_DECIMAL  LeaderboardFormat = "decimal"
	LEADERBOARD_FORMAT_PERCENT  LeaderboardFormat = "percent"
	LEADERBOARD_FORMAT_CURRENCY LeaderboardFormat = "currency"
)

type LeaderboardItem struct {
	Label        string  `json:"label"`
	Value        float64 `json:"value,omitempty"`
	PreviousRank float64 `json:"previous_rank,omitempty"`
}

// Monitoring will show if a service is up or down
type Monitoring struct {
	Status       string `json:"status"`
	Downtime     string `json:"downTime,omitempty"`
	Responsetime string `json:"responseTime,omitempty"`
}

// Default monitoring values from Geckoboard; arbitrary values are
// permitted.
const (
	MONITORING_UP   string = "Up"
	MONITORING_DOWN string = "Down"
)

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

// RAG implements the Red/Amber/Green widget.
// Note the use of pointers (to make fields optional)
// https://developer.geckoboard.com/#rag
type RAG struct {
	Red     *RAGItem
	Amber   *RAGItem
	Green   *RAGItem
	Prefix  string
	Reverse bool
}

func (r RAG) MarshalJSON() ([]byte, error) {
	// Require Red and Amber to be present
	if r.Red == nil {
		return []byte{}, errors.New("Red is required.")
	}
	if r.Amber == nil {
		return []byte{}, errors.New("Amber is required.")
	}

	items := []RAGItem{*r.Red, *r.Amber}

	if r.Green != nil {
		items = append(items, *r.Green)
	}

	// Pass custom object to regular JSON marshaller
	return json.Marshal(struct {
		Item    []RAGItem `json:"item"`
		Prefix  string    `json:"prefix,omitempty"`
		Reverse bool      `json:"reverse,omitempty"`
	}{
		Item:    items,
		Prefix:  r.Prefix,
		Reverse: r.Reverse,
	})
}

type RAGItem struct {
	Value float64 `json:"value"`
	Text  string  `json:"text"`
}

// Trendline implements the geckoboard trendline widget
// https://developer.geckoboard.com/#trendline-example
type TrendlineSecondary []float64

// Text implements the Text widget
// https://developer.geckoboard.com/#text
type Text []TextPage

// Create a new Text-widget from the given strings
func NewSimpleText(args ...string) Text {
	ret := make(Text, len(args))
	for i, arg := range args {
		ret[i] = TextPage{Text: arg}
	}
	return ret
}

func (t Text) MarshalJSON() ([]byte, error) {
	// Bail if we have too many texts
	if len(t) > 10 {
		return []byte{}, errors.New("Text widget support at most 10 entries.")
	}

	// Pass custom object to regular JSON marshaller
	return json.Marshal(struct {
		Item []TextPage `json:"item"`
	}{
		Item: t,
	})
}

// TextPage is a page of text that a Text-widget will cycle between
type TextPage struct {
	Text string   `json:"text"`
	Type TextType `json:"type,omitempty"`
}

// TextType tells if the text widget will have any special ornamentation
type TextType byte

const (
	TEXT_TYPE_NONE  TextType = iota // No type
	TEXT_TYPE_ALERT                 // Type 1: Alert
	TEXT_TYPE_INFO                  // Type 2: Info
)
