// Package testutils provides utilities used in tests
package testutils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func checkEqual(a interface{}, b interface{}, message string) (bool, string) {
	if a == b {
		return true, ""
	}

	var m string
	if len(message) == 0 {
		m = fmt.Sprintf("%v != %v", a, b)
	} else {
		m = message
	}
	errorMessage := fmt.Sprintf("%s. Actual: %+v. Expected: %+v.", m, a, b)

	return false, errorMessage
}

// AssertEqual errors a test if the actual does not match the expected
func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if !ok {
		t.Error(m)
	}
}

// AssertEqualf fails a test if the actual does not match the expected
func AssertEqualf(t *testing.T, a interface{}, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if !ok {
		t.Fatal(m)
	}
}

// AssertNotEqual fails a test if the actual matches the expected
func AssertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v == %v", a, b)
	}
	t.Errorf("%s. Actual: %+v. Expected: %+v.", message, a, b)
}

// AssertDeepEqual fails a test if the actual does not deeply equal the expected
func AssertDeepEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if reflect.DeepEqual(a, b) {
		return
	}

	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Errorf("%s.\nActual:   %+v.\nExpected: %+v.", message, a, b)
}

// ReadJSON reads JSON fixture to the struct at the destination address
func ReadJSON(path string, destination interface{}) {
	var dat []byte
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, "Failed to load fixture payload"))
	}
	if err := json.Unmarshal(dat, destination); err != nil {
		panic(errors.Wrap(err, "Failed to get event"))
	}
}

// IsEqualJSON deeply compares two JSON byte slices
func IsEqualJSON(s1, s2 []byte) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	if err := json.Unmarshal(s1, &o1); err != nil {
		return false, errors.Wrap(err, "unmarshalling first  JSON")
	}
	if err := json.Unmarshal(s2, &o2); err != nil {
		return false, errors.Wrap(err, "unmarshalling second JSON")
	}

	return reflect.DeepEqual(o1, o2), nil
}

// MustExec executes the given SQL query and fails a test if an error occurs
func MustExec(t *testing.T, message string, db *sql.DB, query string, args ...interface{}) sql.Result {
	result, err := db.Exec(query, args...)
	if err != nil {
		t.Fatal(errors.Wrap(errors.Wrap(err, "executing sql"), message))
	}

	return result
}

// MustMarshalJSON marshalls the given interface into JSON.
// If there is any error, it fails the test.
func MustMarshalJSON(t *testing.T, v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("%s: marshalling data", t.Name())
	}

	return b
}

// MustUnmarshalJSON marshalls the given interface into JSON.
// If there is any error, it fails the test.
func MustUnmarshalJSON(t *testing.T, data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		t.Fatalf("%s: unmarshalling data", t.Name())
	}
}
