package liquid

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/osteele/liquid/expressions"
	"github.com/osteele/liquid/generics"
	"github.com/russross/blackfriday"
)

func (e *Wrapper) addJekyllFilters() {
	// arrays
	e.engine.DefineFilter("array_to_sentence_string", arrayToSentenceStringFilter)
	// TODO neither Liquid nor Jekyll docs this, but it appears to be present
	e.engine.DefineFilter("filter", func(values []map[string]interface{}, key string) []interface{} {
		out := []interface{}{}
		for _, value := range values {
			if _, ok := value[key]; ok {
				out = append(out, value)
			}
		}
		return out
	})
	// sort overrides the Liquid filter with one that takes parameters
	e.engine.DefineFilter("sort", sortFilter)
	e.engine.DefineFilter("where", whereFilter) // TODO test case
	e.engine.DefineFilter("where_exp", whereExpFilter)
	e.engine.DefineFilter("xml_escape", xml.Marshal)

	e.engine.DefineFilter("push", func(array []interface{}, item interface{}) interface{} {
		return append(array, generics.MustConvertItem(item, array))
	})
	e.engine.DefineFilter("unshift", func(array []interface{}, item interface{}) interface{} {
		return append([]interface{}{generics.MustConvertItem(item, array)}, array...)
	})

	// dates
	e.engine.DefineFilter("date_to_rfc822", func(date time.Time) string {
		return date.Format(time.RFC822)
		// Out: Mon, 07 Nov 2008 13:07:54 -0800
	})
	e.engine.DefineFilter("date_to_string", func(date time.Time) string {
		return date.Format("02 Jan 2006")
		// Out: 07 Nov 2008
	})
	e.engine.DefineFilter("date_to_long_string", func(date time.Time) string {
		return date.Format("02 January 2006")
		// Out: 07 November 2008
	})
	e.engine.DefineFilter("date_to_xmlschema", func(date time.Time) string {
		return date.Format("2006-01-02T15:04:05-07:00")
		// Out: 2008-11-07T13:07:54-08:00
	})

	// strings
	e.engine.DefineFilter("absolute_url", func(s string) string {
		return e.AbsoluteURL + e.BaseURL + s
	})
	e.engine.DefineFilter("relative_url", func(s string) string {
		return e.BaseURL + s
	})
	e.engine.DefineFilter("jsonify", json.Marshal)
	e.engine.DefineFilter("markdownify", blackfriday.MarkdownCommon)
	// e.engine.DefineFilter("normalize_whitespace", func(s string) string {
	// 	wsPattern := regexp.MustCompile(`(?s:[\s\n]+)`)
	// 	return wsPattern.ReplaceAllString(s, " ")
	// })
	e.engine.DefineFilter("to_integer", func(n int) int { return n })
	e.engine.DefineFilter("number_of_words", func(s string) int {
		wordPattern := regexp.MustCompile(`\w+`) // TODO what's the Jekyll spec for a word?
		m := wordPattern.FindAllStringIndex(s, -1)
		if m == nil {
			return 0
		}
		return len(m)
	})

	// string escapes
	// e.engine.DefineFilter("uri_escape", func(s string) string {
	// 	parts := strings.SplitN(s, "?", 2)
	// 	if len(parts) > 0 {
	// TODO PathEscape is the wrong function
	// 		parts[len(parts)-1] = url.PathEscape(parts[len(parts)-1])
	// 	}
	// 	return strings.Join(parts, "?")
	// })
	e.engine.DefineFilter("xml_escape", func(s string) string {
		// TODO can't handle maps
		// eval https://github.com/clbanning/mxj
		// adapt https://stackoverflow.com/questions/30928770/marshall-map-to-xml-in-go
		buf := new(bytes.Buffer)
		if err := xml.EscapeText(buf, []byte(s)); err != nil {
			panic(err)
		}
		return buf.String()
	})
}

func arrayToSentenceStringFilter(value []string, conjunction interface{}) string {
	conj, ok := conjunction.(string)
	if !ok {
		conj = "and "
	}
	rt := reflect.ValueOf(value)
	ar := make([]string, rt.Len())
	for i, v := range value {
		ar[i] = v
		if i == rt.Len()-1 {
			ar[i] = conj + v
		}
	}
	return strings.Join(ar, ", ")
}

func sortFilter(in []interface{}, key interface{}, nilFirst interface{}) []interface{} {
	nf, ok := nilFirst.(bool)
	if !ok {
		nf = true
	}
	out := make([]interface{}, len(in))
	copy(out, in)
	if key == nil {
		generics.Sort(out)
	} else {
		generics.SortByProperty(out, key.(string), nf)
	}
	return out
}

func whereExpFilter(in []interface{}, name string, expr expressions.Closure) ([]interface{}, error) {
	rt := reflect.ValueOf(in)
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return in, nil
	}
	out := []interface{}{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i).Interface()
		value, err := expr.Bind(name, item).Evaluate()
		if err != nil {
			return nil, err
		}
		if value != nil && value != false {
			out = append(out, item)
		}
	}
	return out, nil
}

func whereFilter(in []map[string]interface{}, key string, value interface{}) []interface{} {
	rt := reflect.ValueOf(in)
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil
	}
	out := []interface{}{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i)
		if item.Kind() == reflect.Map && item.Type().Key().Kind() == reflect.String {
			attr := item.MapIndex(reflect.ValueOf(key))
			if attr.IsValid() && (value == nil || attr.Interface() == value) {
				out = append(out, item.Interface())
			}
		}
	}
	return out
}
