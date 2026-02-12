package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/itchyny/gojq"
)

func Print(data json.RawMessage, jqExpr string) error {
	if jqExpr == "" {
		var out bytes.Buffer
		if err := json.Indent(&out, data, "", "  "); err != nil {
			_, err = os.Stdout.Write(data)
			fmt.Fprintln(os.Stdout)
			return err
		}
		out.WriteByte('\n')
		_, err := out.WriteTo(os.Stdout)
		return err
	}

	query, err := gojq.Parse(jqExpr)
	if err != nil {
		return fmt.Errorf("jq parse error: %w", err)
	}

	var input any
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}

	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			return fmt.Errorf("jq error: %w", err)
		}
		switch val := v.(type) {
		case string:
			fmt.Println(val)
		default:
			out, err := json.MarshalIndent(val, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(out))
		}
	}
	return nil
}
