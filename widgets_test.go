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

		// Leaderboard
		{
			NewLeaderboard(
				LeaderboardItem{Label: "First"},
				LeaderboardItem{Label: "Second"},
			),
			`{"item":[{"label":"First"},{"label":"Second"}]}`,
		},
		{
			NewLeaderboard(
				LeaderboardItem{"First", 123, 7},
			),
			`{"item":[{"label":"First","value":123,"previous_rank":7}]}`,
		},
		{
			Leaderboard{
				Items:  []LeaderboardItem{LeaderboardItem{Label: "x"}},
				Format: LEADERBOARD_FORMAT_DECIMAL,
			},
			`{"item":[{"label":"x"}],"format":"decimal"}`,
		},
		{
			Leaderboard{
				Items:  []LeaderboardItem{LeaderboardItem{Label: "x"}},
				Format: LEADERBOARD_FORMAT_PERCENT,
			},
			`{"item":[{"label":"x"}],"format":"percent"}`,
		},
		{
			Leaderboard{
				Items:  []LeaderboardItem{LeaderboardItem{Label: "x"}},
				Format: LEADERBOARD_FORMAT_CURRENCY,
				Unit:   "DKK",
			},
			`{"item":[{"label":"x"}],"format":"currency","unit":"DKK"}`,
		},

		// Monitoring
		{
			Monitoring{Status: MONITORING_UP},
			`{"status":"Up"}`,
		},
		{
			Monitoring{Status: MONITORING_DOWN},
			`{"status":"Down"}`,
		},
		{
			Monitoring{
				Status:       "Whatever",
				Downtime:     "down",
				Responsetime: "0",
			},
			`{"status":"Whatever","downTime":"down","responseTime":"0"}`,
		},
		// RAG
		{
			RAG{
				Red:   &RAGItem{1, "11"},
				Amber: &RAGItem{2, "22"},
				Green: &RAGItem{3, "33"},
			},
			`{"item":[{"value":1,"text":"11"},{"value":2,"text":"22"},{"value":3,"text":"33"}]}`,
		},
		{
			RAG{
				Red:     &RAGItem{1, "11"},
				Amber:   &RAGItem{2, "22"},
				Prefix:  "@",
				Reverse: true,
			},
			`{"item":[{"value":1,"text":"11"},{"value":2,"text":"22"}],"prefix":"@","reverse":true}`,
		},

		// Text widget
		{
			Text{TextPage{Text: "1 2 3"}},
			`{"item":[{"text":"1 2 3"}]}`,
		},
		{
			Text{
				TextPage{"none", TEXT_TYPE_NONE},
				TextPage{"alert", TEXT_TYPE_ALERT},
				TextPage{"info", TEXT_TYPE_INFO},
			},
			`{"item":[{"text":"none"},{"text":"alert","type":1},{"text":"info","type":2}]}`,
		},
		{
			NewSimpleText("1", "2", "3"),
			`{"item":[{"text":"1"},{"text":"2"},{"text":"3"}]}`,
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

		// RAG
		{
			RAG{Red: &RAGItem{}},
			"json: error calling MarshalJSON for type gockoboard.RAG: Amber is required.",
		},
		{
			RAG{Amber: &RAGItem{}},
			"json: error calling MarshalJSON for type gockoboard.RAG: Red is required.",
		},

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
		j, err := json.Marshal(tt.in)

		if err == nil {
			t.Fatalf("Expected error when Marshal()'ing:\n\t`%v`\nGot JSON:\n]t`%v`", tt.in, j)
		}

		if err.Error() != tt.err {
			t.Fatalf("Expected error\n\t`%v`,\nGot\n\t`%v`.", tt.err, err.Error())
		}
	}
}
