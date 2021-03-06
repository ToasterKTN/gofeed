package shared

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	xpp "github.com/mmcdole/goxpp"
)

var (
	emailNameRgx = regexp.MustCompile(`^([^@]+@[^\s]+)\s+\(([^@]+)\)$`)
	nameEmailRgx = regexp.MustCompile(`^([^@]+)\s+\(([^@]+@[^)]+)\)$`)
	nameOnlyRgx  = regexp.MustCompile(`^([^@()]+)$`)
	emailOnlyRgx = regexp.MustCompile(`^([^@()]+@[^@()]+)$`)

	TruncatedEntity         = errors.New("truncated entity")
	InvalidNumericReference = errors.New("invalid numeric reference")
)

const CDATA_START = "<![CDATA["
const CDATA_END = "]]>"

// ParseText is a helper function for parsing the text
// from the current element of the XMLPullParser.
// This function can handle parsing naked XML text from
// an element.
func ParseText(p *xpp.XMLPullParser) (string, error) {
	var text struct {
		Type     string `xml:"type,attr"`
		InnerXML string `xml:",innerxml"`
	}

	err := p.DecodeElement(&text)
	if err != nil {
		return "", err
	}

	result := text.InnerXML
	result = strings.TrimSpace(result)

	if strings.Contains(result, CDATA_START) {
		prestring, _ := DecodeEntities(result[:strings.Index(result, CDATA_START)])
		cdatastring := StripCDATA(result)
		poststring, _ := DecodeEntities(result[strings.Index(result, CDATA_END)+3:])
		return prestring + cdatastring + poststring, nil
	}

	return DecodeEntities(result)
}

// StripCDATA removes CDATA tags from the string
// content outside of CDATA tags is passed via DecodeEntities
func StripCDATA(str string) string {
	buf := bytes.NewBuffer([]byte{})

	curr := 0

	for curr < len(str) {

		start := indexAt(str, CDATA_START, curr)

		if start == -1 {
			dec, _ := DecodeEntities(str[curr:])
			buf.Write([]byte(dec))
			return buf.String()
		}

		end := indexAt(str, CDATA_END, start)

		if end == -1 {
			dec, _ := DecodeEntities(str[curr:])
			buf.Write([]byte(dec))
			return buf.String()
		}

		buf.Write([]byte(str[start+len(CDATA_START) : end]))

		curr = curr + end + len(CDATA_END)
	}

	return buf.String()
}

// DecodeEntities decodes escaped XML entities
// in a string and returns the unescaped string
func DecodeEntities(str string) (string, error) {
	data := []byte(str)
	buf := bytes.NewBuffer([]byte{})

	for len(data) > 0 {
		// Find the next entity
		idx := bytes.IndexByte(data, '&')
		if idx == -1 {
			buf.Write(data)
			break
		}

		buf.Write(data[:idx])
		data = data[idx:]

		// If there is only the '&' left here
		if len(data) == 1 {
			buf.Write(data)
			return buf.String(), nil
		}

		// Find the end of the entity
		end := bytes.IndexByte(data, ';')
		if end == -1 {
			// it's not an entitiy. just a plain old '&' possibly with extra bytes
			buf.Write(data)
			return buf.String(), nil
		}

		// Check if there is a space somewhere within the 'entitiy'.
		// If there is then skip the whole thing since it's not a real entity.
		if strings.Contains(string(data[1:end]), " ") {
			buf.Write(data)
			return buf.String(), nil
		} else {
			if data[1] == '#' {
				// Numerical character reference
				var str string
				base := 10

				if len(data) > 2 && data[2] == 'x' {
					str = string(data[3:end])
					base = 16
				} else {
					str = string(data[2:end])
				}

				i, err := strconv.ParseUint(str, base, 32)
				if err != nil {
					return "", InvalidNumericReference
				}

				buf.WriteRune(rune(i))
			} else {
				// Predefined entity
				name := string(data[1:end])

				var c byte
				switch name {
				case "lt":
					c = '<'
				case "gt":
					c = '>'
				case "quot":
					c = '"'
				case "apos":
					c = '\''
				case "amp":
					c = '&'
				default:
					return "", fmt.Errorf("unknown predefined "+
						"entity &%s;", name)
				}

				buf.WriteByte(c)
			}
		}

		// Skip the entity
		data = data[end+1:]
	}

	return buf.String(), nil
}

// ParseNameAddress parses name/email strings commonly
// found in RSS feeds of the format "Example Name (example@site.com)"
// and other variations of this format.
func ParseNameAddress(nameAddressText string) (name string, address string) {
	if nameAddressText == "" {
		return
	}

	if emailNameRgx.MatchString(nameAddressText) {
		result := emailNameRgx.FindStringSubmatch(nameAddressText)
		address = result[1]
		name = result[2]
	} else if nameEmailRgx.MatchString(nameAddressText) {
		result := nameEmailRgx.FindStringSubmatch(nameAddressText)
		name = result[1]
		address = result[2]
	} else if nameOnlyRgx.MatchString(nameAddressText) {
		result := nameOnlyRgx.FindStringSubmatch(nameAddressText)
		name = result[1]
	} else if emailOnlyRgx.MatchString(nameAddressText) {
		result := emailOnlyRgx.FindStringSubmatch(nameAddressText)
		address = result[1]
	}
	return
}

func indexAt(str, substr string, start int) int {
	idx := strings.Index(str[start:], substr)
	if idx > -1 {
		idx += start
	}
	return idx
}
