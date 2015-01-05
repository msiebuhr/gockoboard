package gockoboard

import (
	"encoding/json"
	"fmt"
	"os"
)

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
