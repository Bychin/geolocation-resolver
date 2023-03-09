package geoloc

import (
	"net"
	"reflect"
	"testing"
)

func TestNewEntryFromStringSlice(t *testing.T) {
	for i, testcase := range []struct {
		input          []string
		expectedResult *Entry
		mustFail       bool
	}{
		{
			input: []string{
				"200.106.141.15",
				"SI",
				"Nepal",
				"DuBuquemouth",
				"-84.87503094689836",
				"7.206435933364332",
			},
			expectedResult: &Entry{
				IP:          net.ParseIP("200.106.141.15"),
				CountryCode: "SI",
				Country:     "Nepal",
				City:        "DuBuquemouth",
				Latitude:    -84.87503094689836,
				Longitude:   7.206435933364332,
			},
			mustFail: false,
		},
		{
			input: []string{
				"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
				"SI",
				"Nepal",
				"DuBuquemouth",
				"-84.87503094689836",
				"7.206435933364332",
			},
			expectedResult: &Entry{
				IP:          net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
				CountryCode: "SI",
				Country:     "Nepal",
				City:        "DuBuquemouth",
				Latitude:    -84.87503094689836,
				Longitude:   7.206435933364332,
			},
			mustFail: false,
		},
		{
			input: []string{
				"invalid_ip",
				"SI",
				"Nepal",
				"DuBuquemouth",
				"-84.87503094689836",
				"7.206435933364332",
			},
			mustFail: true,
		},
		{
			input: []string{
				"256.256.256.256", // invalid IPv4
				"SI",
				"Nepal",
				"DuBuquemouth",
				"-84.87503094689836",
				"7.206435933364332",
			},
			mustFail: true,
		},
		{
			input: []string{
				"200.106.141.15",
			},
			mustFail: true,
		},
		{
			input: []string{
				"200.106.141.15",
				"SI",
				"Nepal",
				"DuBuquemouth",
				"invalid_lat",
				"7.206435933364332",
			},
			mustFail: true,
		},
		{
			input: []string{
				"200.106.141.15",
				"SI",
				"Nepal",
				"DuBuquemouth",
				"7.206435933364332",
				"invalid_lon",
			},
			mustFail: true,
		},
	} {
		parsedEntry, err := NewEntryFromStringSlice(testcase.input)
		if err != nil {
			if !testcase.mustFail {
				t.Errorf("#%d: got error when parsing entry: %s", i, err)
			}
			continue
		}
		if !reflect.DeepEqual(parsedEntry, testcase.expectedResult) {
			t.Errorf("#%d: parsed entry is not equal to expected one, have: %+v, want: %+v",
				i, parsedEntry, testcase.expectedResult)
		}
		if testcase.mustFail {
			t.Errorf("#%d: expected to fail, but didn't", i)
		}
	}
}

func TestValidate(t *testing.T) {
	for i, testcase := range []struct {
		input    *Entry
		mustFail bool
	}{
		{
			input: &Entry{
				IP:          net.ParseIP("200.106.141.15"),
				CountryCode: "SI",
				Country:     "Nepal",
				City:        "DuBuquemouth",
				Latitude:    -84.87503094689836,
				Longitude:   7.206435933364332,
			},
			mustFail: false,
		},
		{
			input: &Entry{
				IP:          nil, // empty IP
				CountryCode: "SI",
				Country:     "Nepal",
				City:        "DuBuquemouth",
				Latitude:    -84.87503094689836,
				Longitude:   7.206435933364332,
			},
			mustFail: true,
		},
		{
			input: &Entry{
				IP:          net.ParseIP("200.106.141.15"),
				CountryCode: "CODE", // invalid code
				Country:     "Nepal",
				City:        "DuBuquemouth",
				Latitude:    -84.87503094689836,
				Longitude:   7.206435933364332,
			},
			mustFail: true,
		},
		{
			input: &Entry{
				IP:          net.ParseIP("200.106.141.15"),
				CountryCode: "SI",
				Country:     "", // no value
				City:        "DuBuquemouth",
				Latitude:    -84.87503094689836,
				Longitude:   7.206435933364332,
			},
			mustFail: true,
		},
		{
			input: &Entry{
				IP:          net.ParseIP("200.106.141.15"),
				CountryCode: "SI",
				Country:     "Nepal",
				City:        "", // no value
				Latitude:    -84.87503094689836,
				Longitude:   7.206435933364332,
			},
			mustFail: true,
		},
		{
			input: &Entry{
				IP:          net.ParseIP("200.106.141.15"),
				CountryCode: "SI",
				Country:     "Nepal",
				City:        "DuBuquemouth",
				Latitude:    0, // missing both coordinates
				Longitude:   0,
			},
			mustFail: true,
		},
		{
			input:    &Entry{},
			mustFail: true,
		},
	} {
		if err := testcase.input.Validate(); err != nil {
			if !testcase.mustFail {
				t.Errorf("#%d: got validation error: %s", i, err)
			}
			continue
		}
		if testcase.mustFail {
			t.Errorf("#%d: expected to fail, but didn't", i)
		}
	}
}
