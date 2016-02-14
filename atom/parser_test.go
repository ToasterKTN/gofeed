package atom_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/atom"
	"github.com/stretchr/testify/assert"
)

func TestAtomParser_ParseFeed_WellFormed(t *testing.T) {
	files, _ := filepath.Glob("../testdata/wellformed/atom/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/wellformed/atom/%s.xml", name)
		f, _ := ioutil.ReadFile(ff)

		// Parse actual feed
		fp := &atom.Parser{}
		actual, _ := fp.ParseFeed(string(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/wellformed/atom/%s.json", name)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &atom.Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, actual, expected, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

func TestAtomParser_ParseFeed_IllFormed(t *testing.T) {
	files, _ := filepath.Glob("../testdata/illformed/atom/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/illformed/atom/%s.xml", name)
		f, _ := ioutil.ReadFile(ff)

		// Parse actual feed
		fp := &atom.Parser{}
		actual, _ := fp.ParseFeed(string(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/illformed/atom/%s.json", name)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &atom.Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, actual, expected, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}