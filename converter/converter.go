package converter

import (
	"github.com/liuzl/gocc"
	"log"
)

func ConvertString(in string, rule string) (string, error) {
	t2s, err := gocc.New(rule)
	if err != nil {
		return "", err
	}
	out, err := t2s.Convert(in)
	if err != nil {
		log.Fatal(err)
	}
	return out, nil
}
