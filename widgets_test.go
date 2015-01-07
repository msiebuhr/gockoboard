package gockoboard

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestMarshalGeckOMeter(t *testing.T) {
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

func TestMarshalNumberAndTrendline(t *testing.T) {
	tl := NumberAndTrendline{
		Text:  "t",
		Value: 42,
		Trend: []float64{1, 2, 3},
	}

	geckoJson, err := json.Marshal(tl)

	if err != nil {
		t.Fatalf("Unexpected error when Marshal()'ing: %v", err)
	}

	expectedJson := `{"item":[{"text":"t","value":42},[1,2,3]]}`

	if string(geckoJson) != expectedJson {
		t.Fatalf("Expected '%v', got '%v'.", expectedJson, string(geckoJson))
	}
}

func ExampleNumberAndTrendline() {
	tl := NumberAndTrendline{
		Text:  "Monthly new users",
		Value: 32,
		Trend: []float64{2, 4, 8, 16},
	}

	b, err := json.Marshal(tl)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// {"item":[{"text":"Monthly new users","value":32},[2,4,8,16]]}
}
