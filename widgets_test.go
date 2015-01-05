package gockoboard

import (
    "testing"
    "encoding/json"
    "fmt"
)

func TestMarshalGeckOMeter(t *testing.T) {
    g := GeckOMeter{
        Item: 123.4,
        Min: 20,
        Max: 400,
    }

    geckoJson, err := json.Marshal(g);

    if err != nil{
        t.Fatalf("Unexpected error when Marshal()'ing:", err)
    }

    expectedJson := `{"item":123.4,"min":{"value":20},"max":{"value":400}}`

    if string(geckoJson) != expectedJson {
        t.Fatalf("Expected '%v', got '%v'.", expectedJson, string(geckoJson))
    }
}
