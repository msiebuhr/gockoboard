package gockoboard

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestMarshalGeckOMeter(t *testing.T) {
	t.Parallel()
	g := GeckOMeter{
		Item: 123.4,
		Min:  20,
		Max:  400,
	}

	geckoJson, err := json.Marshal(g)

	if err != nil {
		t.Fatalf("Unexpected error when Marshal()'ing: %v", err)
	}

	expectedJson := `{"item":123.4,"min":{"value":20},"max":{"value":400}}`

	if string(geckoJson) != expectedJson {
		t.Fatalf("Expected '%v', got '%v'.", expectedJson, string(geckoJson))
	}
}

func TestMarshalNumberAndSecondaryTrendline(t *testing.T) {
	t.Parallel()
	tl := Number{
		Value:         1,
		SecondaryStat: TrendlineSecondary{1, 2, 3},
	}

	geckoJson, err := json.Marshal(tl)

	if err != nil {
		t.Fatalf("Unexpected error when Marshal()'ing: %v", err)
	}

	expectedJson := `{"item":[{"value":1},[1,2,3]]}`

	if string(geckoJson) != expectedJson {
		t.Fatalf("Expected '%v', got '%v'.", expectedJson, string(geckoJson))
	}
}

func ExampleNumberAndTrendlineSecondary() {
	tl := Number{
		Value:         32,
		Text:          "Monthly new users",
		SecondaryStat: TrendlineSecondary{2, 4, 8, 16},
	}

	b, err := json.Marshal(tl)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// {"item":[{"value":32,"text":"Monthly new users"},[2,4,8,16]]}
}

func TestMarshalWidgets(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		in  interface{}
		out string
	}{
		// Number widgets
		{Number{Value: 123}, `{"item":[{"value":123}]}`},
		{Number{Value: 42, Prefix: "%"}, `{"item":[{"value":42,"prefix":"%"}]}`},
		{Number{Value: 42, Text: "HG2G"}, `{"item":[{"value":42,"text":"HG2G"}]}`},
		{Number{Value: 42, Type: "reverse"}, `{"item":[{"value":42,"type":"reverse"}]}`},

		// Text widget
		{
			Text{TextPage{Text: "1 2 3"}},
			`{"item":[{"text":"1 2 3"}]}`,
		},
		{
			Text{
				TextPage{"1", TEXT_TYPE_NONE},
				TextPage{"2", TEXT_TYPE_INFO},
				TextPage{"2", TEXT_TYPE_ALERT},
			},
			`{"item":[{"text":"1"},{"text":"2","type":1},{"text":"2","type":2}]}`,
		},
	}

	for _, tt := range tests {
		geckoJson, err := json.Marshal(tt.in)

		if err != nil {
			t.Fatalf("Unexpected error when Marshal()'ing: %v", err)
		}

		if string(geckoJson) != tt.out {
			t.Fatalf("Expected\n\t`%v`,\nGot\n\t`%v`.", tt.out, string(geckoJson))
		}
	}
}

func TestMarshalWidgetErrorss(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		in  interface{}
		err string
	}{
		// Number widgets
		//{Number{}, ``},

		// 11 pages of text (only 10 allowed)
		{
			Text{
				TextPage{}, TextPage{}, TextPage{}, TextPage{}, TextPage{},
				TextPage{}, TextPage{}, TextPage{}, TextPage{}, TextPage{},
				TextPage{},
			},
			`json: error calling MarshalJSON for type gockoboard.Text: Text widget support at most 10 entries.`,
		},
	}

	for _, tt := range tests {
		_, err := json.Marshal(tt.in)

		if err == nil {
			t.Fatalf("Expected error when Marshal()'ing:\n\t`%v`\nGot JSON:\n]t`%v`", tt.in)
		}

		if err.Error() != tt.err {
			t.Fatalf("Expected error\n\t`%v`,\nGot\n\t`%v`.", tt.err, err.Error())
		}
	}
}
