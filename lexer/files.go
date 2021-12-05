package lexer

import (
	"bufio"
	"encoding/json"
	"io"
)

// MarshalTo writes a JSON lines representation of tokens to out
//
// A token has the form {"Kind": }
func MarshalTo(tokens []Token, out io.Writer) error {
	for _, token := range tokens {
		data, err := json.Marshal(token)
		if err != nil {
			return err
		}

		_, err = out.Write(data)
		if err != nil {
			return err
		}

		_, err = out.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalFrom will read in a JSON lines token stream representation and unmarshal it to
// a slice of tokens. It is the inverse of MarshalTo.
func UnmarshalFrom(in io.Reader) ([]Token, error) {
	ret := make([]Token, 0)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		var token Token
		err := json.Unmarshal(scanner.Bytes(), &token)
		if err != nil {
			return nil, err
		}
		ret = append(ret, token)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}
