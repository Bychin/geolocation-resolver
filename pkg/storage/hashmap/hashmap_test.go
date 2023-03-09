package hashmap

import (
	"context"
	"net"
	"reflect"
	"testing"

	"geolocation-resolver/pkg/geoloc"
)

const (
	ip1 = "200.106.141.15"
)

func newFilledStorage() *Storage {
	storage := NewStorage()

	storage.storage = map[string]geoloc.Entry{
		ip1: {
			IP: net.ParseIP(ip1),
		},
	}

	return storage
}

func TestWrite(t *testing.T) {
	for i, testcase := range []struct {
		storage  *Storage
		input    *geoloc.Entry
		mustFail bool
	}{
		{
			storage: newFilledStorage(),
			input: &geoloc.Entry{
				IP: net.ParseIP(ip1),
			},
			mustFail: false,
		},
		{
			storage: newFilledStorage(),
			input: &geoloc.Entry{
				IP: net.ParseIP("0.0.0.0"),
			},
			mustFail: false,
		},
		{
			storage:  nil,
			mustFail: true,
		},
	} {
		if err := testcase.storage.Write(context.Background(), testcase.input); err != nil {
			if !testcase.mustFail {
				t.Errorf("#%d: got error on read: %s", i, err)
			}
			continue
		}
		if testcase.mustFail {
			t.Errorf("#%d: expected to fail, but didn't", i)
		}
	}
}

func TestRead(t *testing.T) {
	for i, testcase := range []struct {
		storage        *Storage
		input          net.IP
		expectedResult *geoloc.Entry
		mustFail       bool
	}{
		{
			storage: newFilledStorage(),
			input:   net.ParseIP(ip1),
			expectedResult: &geoloc.Entry{
				IP: net.ParseIP(ip1),
			},
			mustFail: false,
		},
		{
			storage:        newFilledStorage(),
			input:          net.ParseIP("random"),
			expectedResult: nil,
			mustFail:       false,
		},
		{
			storage:        nil,
			input:          net.ParseIP("random"),
			expectedResult: nil,
			mustFail:       true,
		},
	} {
		entry, err := testcase.storage.Read(context.Background(), testcase.input)
		if err != nil {
			if !testcase.mustFail {
				t.Errorf("#%d: got error on read: %s", i, err)
			}
			continue
		}
		if !reflect.DeepEqual(entry, testcase.expectedResult) {
			t.Errorf("#%d: read entry is not equal to expected one, have: %+v, want: %+v",
				i, entry, testcase.expectedResult)
		}
		if testcase.mustFail {
			t.Errorf("#%d: expected to fail, but didn't", i)
		}
	}
}
