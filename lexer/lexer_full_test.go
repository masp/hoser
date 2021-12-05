package lexer

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestExamples(t *testing.T) {
	err := os.MkdirAll("examples/want", os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}

	items, _ := ioutil.ReadDir("examples/")
	for _, item := range items {
		if !item.IsDir() && strings.HasSuffix(item.Name(), "hos") {
			example := item.Name()
			fmt.Println(example)
			t.Run(fmt.Sprintf("example: %v", example), func(t *testing.T) {
				exampleSrc, err := readExampleSrc(example)
				if err != nil {
					t.Errorf("readExampleSrc() = %v", err)
					return
				}

				want, err := readWantSrc(example, exampleSrc)
				if err != nil {
					t.Errorf("readWantSrc() = %v", err)
					return
				}

				got, err := ScanAll(exampleSrc)
				if err != nil {
					t.Errorf("ScanAll(examples/%v) = %v", example, err)
					return
				}

				if !reflect.DeepEqual(got, want) {
					t.Errorf("ScanAll(examples/%v) = ScanAll(examples/want/%v), got %v, want %v", example, example, got, want)
				}
			})
		}
	}
}

func readExampleSrc(example string) (string, error) {
	exampleFile, err := os.Open(fmt.Sprintf("examples/%s", example))
	if err != nil {
		return "", err
	}

	exampleSrc, err := ioutil.ReadAll(exampleFile)
	if err != nil {
		return "", err
	}
	return string(exampleSrc), nil
}

func readWantSrc(example string, exampleSrc string) ([]Token, error) {
	var want []Token
	wantFile, err := os.Open(fmt.Sprintf("examples/want/%s", example))
	if err != nil {
		if os.IsNotExist(err) {
			want, err = generateWantFile(example, exampleSrc)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		want, err = UnmarshalFrom(wantFile)
		if err != nil {
			return nil, err
		}
	}
	return want, nil
}

// generateWantFile will create the example test file if it does not exist. This is useful
// when behavior changes and you want to recreate the examples (you can just do a `rm *.hos`).
func generateWantFile(example string, exampleSrc string) ([]Token, error) {
	wantFile, err := os.Create(fmt.Sprintf("examples/want/%s", example))
	if err != nil {
		return nil, err
	}

	tokens, err := ScanAll(exampleSrc)
	if err != nil {
		return nil, err
	}

	// Write to want file as well so that we avoid generating it new each time
	err = MarshalTo(tokens, wantFile)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
