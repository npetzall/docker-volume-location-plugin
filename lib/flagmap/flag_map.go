package flagmap

import (
	"fmt"
	"strings"
)

type FlagMap map[string]string

func (f *FlagMap) String() string {
	return fmt.Sprint(*f)
}

func (f *FlagMap) Set(value string) error {
	keyAndValue := strings.Split(value, "=")
	if len(keyAndValue) == 1 {
		(*f)["default"] = value
	} else if len(keyAndValue) != 2 {
		fmt.Printf("Invalid location: %s\n", value)
	} else {
		(*f)[keyAndValue[0]] = keyAndValue[1]
	}
	return nil
}
