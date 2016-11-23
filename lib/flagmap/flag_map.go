package flagmap

import (
	"fmt"
	"strings"
)

//FlagMap is used with flag.Var to handle map like structure.
type FlagMap map[string]string

//String returns a string representation of flagMap.
func (f *FlagMap) String() string {
	return fmt.Sprint(*f)
}

//Set method from Value.Set.
//Will plit string on '='
//Key in map will be default is singel value (no '=')
//If multiple '=' exists it will be ignored
//For usage with flag.Var()
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
