/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tika

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

// errorServer always response with http.StatusInternalServerError.
var errorServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}))

var errorClient = NewClient(nil, errorServer.URL)

func TestMain(m *testing.M) {
	r := m.Run()
	errorServer.Close()
	os.Exit(r)
}

func TestCallError(t *testing.T) {
	tests := []struct {
		method string
		url    string
	}{
		{"bad method", ""},
		{"GET", "http://unknown_test_url"},
	}
	for _, test := range tests {
		c := NewClient(nil, test.url)
		if _, err := c.call(nil, test.method, "", nil); err == nil {
			t.Errorf("call(%s, %s) got no error, want error", test.method, test.url)
		}

	}
}

func TestParse(t *testing.T) {
	want := "test value"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.Parse(nil)
	if err != nil {
		t.Errorf("Parse returned nil, want %q", want)
	}
	if got != want {
		t.Errorf("Parse got %q, want %q", got, want)
	}
}

func TestParseRecursive(t *testing.T) {
	tests := []struct {
		response string
		want     []string
	}{
		{
			response: `[{"X-TIKA:content":"test 1"}]`,
			want:     []string{"test 1"},
		},
		{
			response: `[{"X-TIKA:content":"test 1"},{"X-TIKA:content":"test 2"}]`,
			want:     []string{"test 1", "test 2"},
		},
		{
			response: `[{"other_key":"other_value"},{"X-TIKA:content":"test"}]`,
			want:     []string{"test"},
		},
		{
			response: `[]`,
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		got, err := c.ParseRecursive(nil)
		if err != nil {
			t.Errorf("ParseRecursive returned an error: %v, want %v", err, test.want)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ParseRecursive(%q) got %v, want %v", test.response, got, test.want)
		}
	}
}

func TestParseRecursiveError(t *testing.T) {
	_, err := errorClient.ParseRecursive(nil)
	if err == nil {
		t.Error("ParseRecursive got no error, want an error")
	}
}

func TestMeta(t *testing.T) {
	want := "test value"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.Meta(nil)
	if err != nil {
		t.Errorf("Meta returned an error: %v, want %q", err, want)
	}
	if got != want {
		t.Errorf("Meta got %q, want %q", got, want)
	}
}

func TestMetaField(t *testing.T) {
	want := "test value"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.MetaField(nil, "")
	if err != nil {
		t.Errorf("MetaField returned an error: %v, want %q", err, want)
	}
	if got != want {
		t.Errorf("MetaField got %q, want %q", got, want)
	}
}

func TestDetect(t *testing.T) {
	want := "test value"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.Detect(nil)
	if err != nil {
		t.Errorf("Detect returned an error: %v, want %q", err, want)
	}
	if got != want {
		t.Errorf("Detect got %q, want %q", got, want)
	}
}

func TestLanguage(t *testing.T) {
	want := "test value"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.Language(nil)
	if err != nil {
		t.Errorf("Language returned an error: %v, want %q", err, want)
	}
	if got != want {
		t.Errorf("Language got %q, want %q", got, want)
	}
}

func TestLanguageString(t *testing.T) {
	want := "test value"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.LanguageString("")
	if err != nil {
		t.Errorf("LanguageString returned an error: %v, want %q", err, want)
	}
	if got != want {
		t.Errorf("LanguageString got %q, want %q", got, want)
	}
}

func TestMetaRecursive(t *testing.T) {
	tests := []struct {
		response string
		want     []map[string][]string
	}{
		{
			response: `[{"X-TIKA:content":"test 1"}]`,
			want: []map[string][]string{
				map[string][]string{
					"X-TIKA:content": []string{"test 1"},
				},
			},
		},
		{
			response: `[{"X-TIKA:content":"test 1"},{"X-TIKA:content":"test 2"}]`,
			want: []map[string][]string{
				map[string][]string{
					"X-TIKA:content": []string{"test 1"},
				},
				map[string][]string{
					"X-TIKA:content": []string{"test 2"},
				},
			},
		},
		{
			response: `[{"other_key":"other_value"},{"X-TIKA:content":"test"}]`,
			want: []map[string][]string{
				map[string][]string{
					"other_key": []string{"other_value"},
				},
				map[string][]string{
					"X-TIKA:content": []string{"test"},
				},
			},
		},
		{
			response: `[{"other_key":["other_value", "other_value2"]}]`,
			want: []map[string][]string{
				map[string][]string{
					"other_key": []string{"other_value", "other_value2"},
				},
			},
		},
		{
			response: `[]`,
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		got, err := c.MetaRecursive(nil)
		if err != nil {
			t.Errorf("MetaRecursive returned an error: %v, want %v", err, test.want)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("MetaRecursive(%q) got %+v, want %+v", test.response, got, test.want)
		}
	}
}
func TestMetaRecursiveError(t *testing.T) {
	tests := []struct {
		name     string
		response string
	}{
		{
			name:     "invalid type",
			response: `[{"other_key":{"test": "fail"}}]`,
		},
		{
			name:     "invalid nested type",
			response: `[{"other_key":["other_value", {"test": "fail"}]}]`,
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		_, err := c.MetaRecursive(nil)
		if err == nil {
			t.Errorf("MetaRecursive(%s) got no error, want an error", test.name)
		}
	}
}

func TestTranslate(t *testing.T) {
	want := "test value"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.Translate(nil, "translator", "src", "dst")
	if err != nil {
		t.Errorf("Translate returned an error: %v, want %q", err, want)
	}
	if got != want {
		t.Errorf("Translate got %q, want %q", got, want)
	}
}

func TestParsers(t *testing.T) {
	tests := []struct {
		response string
		want     Parser
	}{
		{
			response: `{"name":"TestParser"}`,
			want: Parser{
				Name: "TestParser",
			},
		},
		{
			response: `{
				"name":"TestParser",
				"children":[
					{"name":"TestSubParser1"},
					{"name":"TestSubParser2"}
				]
			}`,
			want: Parser{
				Name: "TestParser",
				Children: []Parser{
					Parser{
						Name: "TestSubParser1",
					},
					Parser{
						Name: "TestSubParser2",
					},
				},
			},
		},
		{
			response: `{
				"name":"TestParser",
				"supportedTypes":["test-type"],
				"children":[
					{
						"supportedTypes":["test-type-two"],
						"name":"TestSubParser",
						"decorated":true,
						"composite":false
					}
				],
				"decorated":false,
				"composite":true}`,
			want: Parser{
				Name:           "TestParser",
				Composite:      true,
				SupportedTypes: []string{"test-type"},
				Children: []Parser{
					Parser{
						Name:           "TestSubParser",
						Decorated:      true,
						SupportedTypes: []string{"test-type-two"},
					},
				},
			},
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		got, err := c.Parsers()
		if err != nil {
			t.Errorf("Parsers returned an error: %v, want %+v", err, test.want)
		}
		if !reflect.DeepEqual(*got, test.want) {
			t.Errorf("Parsers got %+v, want %+v", got, test.want)
		}
	}
}

func TestParsersError(t *testing.T) {
	tests := []struct {
		response string
	}{
		{
			response: "invalid",
		},
		{},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		_, err := c.Parsers()
		if err == nil {
			t.Errorf("Parsers(%q) got no error, want an error", test.response)
		}
	}
	if _, err := errorClient.Parsers(); err == nil {
		t.Errorf("Parsers got no error, want an error")
	}
}

func TestVersion(t *testing.T) {
	want := "test value"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.Version()
	if err != nil {
		t.Errorf("Version returned an error: %v, want %q", err, want)
	}
	if got != want {
		t.Errorf("Version got %q, want %q", got, want)
	}
}

func TestMimeTypes(t *testing.T) {
	tests := []struct {
		response string
		want     map[string]MimeType
	}{
		{
			response: `{"empty-mime":{}}`,
			want: map[string]MimeType{
				"empty-mime": MimeType{},
			},
		},
		{
			response: `{"alias-mime":{"alias":["alias1", "alias2"]}}`,
			want: map[string]MimeType{
				"alias-mime": MimeType{
					Alias: []string{"alias1", "alias2"},
				},
			},
		},
		{
			response: `{"empty-mime":{},"super-mime":{"supertype":"super-mime"}}`,
			want: map[string]MimeType{
				"empty-mime": MimeType{},
				"super-mime": MimeType{SuperType: "super-mime"},
			},
		},
		{
			response: `{"super-alias":{"alias":["alias1", "alias2"], "supertype": "super-mime"}}`,
			want: map[string]MimeType{
				"super-alias": MimeType{
					Alias:     []string{"alias1", "alias2"},
					SuperType: "super-mime",
				},
			},
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		got, err := c.MimeTypes()
		if err != nil {
			t.Errorf("MimeTypes returned an error: %v, want %q", err, test.want)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("MimeTypes got %+v, want %+v", got, test.want)
		}
	}
}

func TestMimeTypesError(t *testing.T) {
	tests := []struct {
		response string
	}{
		{
			response: "",
		},
		{
			response: `["test"]`,
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		_, err := c.MimeTypes()
		if err == nil {
			t.Errorf("MimeTypes got no error, want an error")
		}
	}
	if _, err := errorClient.MimeTypes(); err == nil {
		t.Errorf("MimeTypes got no error, want an error")
	}
}

func TestMetaRecursive_BadResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "invalid")
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.MetaRecursive(nil)
	if err == nil {
		t.Errorf("MetaRecursive got %q, want an error", got)
	}
}

func TestMetaRecursive_BadFieldType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"super-alias":{}`)
	}))
	defer ts.Close()
	c := NewClient(nil, ts.URL)
	got, err := c.MetaRecursive(nil)
	if err == nil {
		t.Errorf("MetaRecursive got %q, want an error", got)
	}
}

func TestDetectors(t *testing.T) {
	tests := []struct {
		response string
		want     Detector
	}{
		{
			response: `{"name":"TestDetector"}`,
			want: Detector{
				Name: "TestDetector",
			},
		},
		{
			response: `{
				"name":"TestDetector",
				"children":[
					{"name":"TestSubDetector1"},
					{"name":"TestSubDetector2"}
				]
			}`,
			want: Detector{
				Name: "TestDetector",
				Children: []Detector{
					Detector{
						Name: "TestSubDetector1",
					},
					Detector{
						Name: "TestSubDetector2",
					},
				},
			},
		},
		{
			response: `{
				"name":"TestDetector",
				"children":[
					{
						"name":"TestSubDetector",
						"composite":false
					}
				],
				"composite":true}`,
			want: Detector{
				Name:      "TestDetector",
				Composite: true,
				Children: []Detector{
					Detector{
						Name: "TestSubDetector",
					},
				},
			},
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		got, err := c.Detectors()
		if err != nil {
			t.Errorf("Detectors returned an error: %v, want %+v", err, test.want)
		}
		if !reflect.DeepEqual(*got, test.want) {
			t.Errorf("Detectors got %+v, want %+v", got, test.want)
		}
	}
}

func TestDetectorsError(t *testing.T) {
	tests := []struct {
		response string
	}{
		{
			response: "",
		},
		{
			response: `["test"]`,
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, test.response)
		}))
		defer ts.Close()
		c := NewClient(nil, ts.URL)
		_, err := c.Detectors()
		if err == nil {
			t.Errorf("Detectors got no error, want an error")
		}
	}
	if _, err := errorClient.Detectors(); err == nil {
		t.Errorf("Detectors got no error, want an error")
	}
}