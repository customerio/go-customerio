package customerio

import (
	"fmt"
	"log"
)

// CustomerIOError is returned by any method that fails at the API level
type CustomerIOError struct {
	status int
	url    string
	body   []byte
}

func (e *CustomerIOError) Error() string {
	return fmt.Sprintf("%v: %v %v", e.status, e.url, string(e.body))
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal("action failed: ", err)
	}
}
