package candiedyaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"unicode"
)

type Decoder struct {
	parser yaml_parser_t
	event  yaml_event_t
}

type ParserError struct {
	ErrorType   YAML_error_type_t
	Context     string
	ContextMark YAML_mark_t
	Problem     string
	ProblemMark YAML_mark_t
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("yaml: [%s] %s at line %d, column %d", e.Context, e.Problem, e.ProblemMark.line+1, e.ProblemMark.column+1)
}

type UnexpectedEventError struct {
	Value     string
	EventType yaml_event_type_t
	At        YAML_mark_t
}

func (e *UnexpectedEventError) Error() string {
	return fmt.Sprintf("yaml: Unexpect event [%d]: '%s' at line %d, column %d", e.EventType, e.Value, e.At.line+1, e.At.column+1)
}

func Unmarshal(data []byte, v interface{}) error {
	d := NewDecoder(bytes.NewBuffer(data))
	return d.Decode(v)
}

func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{}
	yaml_parser_initialize(&d.parser)
	yaml_parser_set_input_reader(&d.parser, r)
	return d
}

func (d *Decoder) Decode(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		rType := reflect.TypeOf(v)
		msg := "nil"
		if rType != nil {
			msg = rType.String()
		}
		return errors.New("Invalid type: " + msg)
	}

	if d.event.event_type == yaml_NO_EVENT {
		d.nextEvent()

		if d.event.event_type != yaml_STREAM_START_EVENT {
			return errors.New("Invalid stream")
		}

		d.nextEvent()
	}

	d.document(rv)
	return nil
}

func (d *Decoder) error(err error) {
	panic(err)
}

func (d *Decoder) nextEvent() {
	if d.event.event_type == yaml_STREAM_END_EVENT {
		d.error(errors.New("The stream is closed"))
	}

	if !yaml_parser_parse(&d.parser, &d.event) {
		yaml_event_delete(&d.event)

		d.error(&ParserError{
			ErrorType:   d.parser.error,
			Context:     d.parser.context,
			ContextMark: d.parser.context_mark,
			Problem:     d.parser.problem,
			ProblemMark: d.parser.problem_mark,
		})
	}
}

func (d *Decoder) document(rv reflect.Value) {
	if d.event.event_type != yaml_DOCUMENT_START_EVENT {
		d.error(fmt.Errorf("Expected document start - found %d", d.event.event_type))
	}

	d.nextEvent()
	d.parse(rv)

	if d.event.event_type != yaml_DOCUMENT_END_EVENT {
		d.error(fmt.Errorf("Expected document end - found %d", d.event.event_type))
	}

	d.nextEvent()
}

func (d *Decoder) parse(rv reflect.Value) {
	switch d.event.event_type {
	case yaml_SEQUENCE_START_EVENT:
		d.sequence(rv)
	case yaml_MAPPING_START_EVENT:
		d.mapping(rv)
	case yaml_SCALAR_EVENT:
		d.scalar(rv)
	case yaml_ALIAS_EVENT:
		d.alias(rv)
	case yaml_DOCUMENT_END_EVENT:
	default:
		d.error(&UnexpectedEventError{
			Value:     string(d.event.value),
			EventType: d.event.event_type,
			At:        d.event.start_mark,
		})
	}
}

func (d *Decoder) indirect(v reflect.Value) reflect.Value {
	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}

	return v
}

func (d *Decoder) sequence(v reflect.Value) {
	if d.event.event_type != yaml_SEQUENCE_START_EVENT {
		d.error(fmt.Errorf("Expected sequence start - found %d", d.event.event_type))
	}

	pv := d.indirect(v)
	v = pv

	// Check type of target.
	switch v.Kind() {
	case reflect.Interface:
		if v.NumMethod() == 0 {
			// Decoding into nil interface?  Switch to non-reflect code.
			v.Set(reflect.ValueOf(d.sequenceInterface()))
			return
		}
		// Otherwise it's invalid.
		fallthrough
	default:
		d.error(errors.New("array: invalid type: " + v.Type().String()))
	case reflect.Array:
	case reflect.Slice:
		break
	}

	d.nextEvent()

	i := 0
	for {
		if d.event.event_type == yaml_SEQUENCE_END_EVENT {
			break
		}

		// Get element of array, growing if necessary.
		if v.Kind() == reflect.Slice {
			// Grow slice if necessary
			if i >= v.Cap() {
				newcap := v.Cap() + v.Cap()/2
				if newcap < 4 {
					newcap = 4
				}
				newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
				reflect.Copy(newv, v)
				v.Set(newv)
			}
			if i >= v.Len() {
				v.SetLen(i + 1)
			}
		}

		if i < v.Len() {
			// Decode into element.
			d.parse(v.Index(i))
		} else {
			// Ran out of fixed array: skip.
			d.parse(reflect.Value{})
		}
		i++
	}

	if i < v.Len() {
		if v.Kind() == reflect.Array {
			// Array.  Zero the rest.
			z := reflect.Zero(v.Type().Elem())
			for ; i < v.Len(); i++ {
				v.Index(i).Set(z)
			}
		} else {
			v.SetLen(i)
		}
	}
	if i == 0 && v.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}

	d.nextEvent()
}

func (d *Decoder) mapping(v reflect.Value) {
	pv := d.indirect(v)
	v = pv

	// Decoding into nil interface?  Switch to non-reflect code.
	if v.Kind() == reflect.Interface && v.NumMethod() == 0 {
		v.Set(reflect.ValueOf(d.mappingInterface()))
		return
	}

	// Check type of target: struct or map[X]Y
	switch v.Kind() {
	case reflect.Struct:
		d.mappingStruct(v)
		return
	case reflect.Map:
	default:
		d.error(errors.New("mapping: invalid type: " + v.Type().String()))
	}

	mapt := v.Type()
	if v.IsNil() {
		v.Set(reflect.MakeMap(mapt))
	}

	d.nextEvent()

	keyt := mapt.Key()
	valuet := mapt.Elem()

	for {
		if d.event.event_type == yaml_MAPPING_END_EVENT {
			break
		}

		key := reflect.New(keyt)
		d.parse(key.Elem())

		value := reflect.New(valuet)
		d.parse(value.Elem())

		v.SetMapIndex(key.Elem(), value.Elem())
	}

	d.nextEvent()
}

func (d *Decoder) mappingStruct(v reflect.Value) {

	structt := v.Type()
	var f *field
	fields := cachedTypeFields(structt)

	d.nextEvent()

	for {
		if d.event.event_type == yaml_MAPPING_END_EVENT {
			break
		}
		key := ""
		d.parse(reflect.ValueOf(&key))

		// Figure out field corresponding to key.
		var subv reflect.Value

		for i := range fields {
			ff := &fields[i]
			if ff.name == key {
				f = ff
				break
			}
			if f == nil && strings.EqualFold(ff.name, key) {
				f = ff
			}
		}

		if f != nil {
			subv = v
			for _, i := range f.index {
				if subv.Kind() == reflect.Ptr {
					if subv.IsNil() {
						subv.Set(reflect.New(subv.Type().Elem()))
					}
					subv = subv.Elem()
				}
				subv = subv.Field(i)
			}
		}
		d.parse(subv)
	}

	d.nextEvent()
}

func (d *Decoder) scalar(v reflect.Value) {
	pv := d.indirect(v)

	v = pv

	err := resolve(d.event, v)
	if err != nil {
		d.error(err)
	}

	d.nextEvent()
}

func (d *Decoder) alias(v interface{}) {
	d.nextEvent()
}

// arrayInterface is like array but returns []interface{}.
func (d *Decoder) sequenceInterface() []interface{} {
	var v = make([]interface{}, 0)

	d.nextEvent()
	for {
		if d.event.event_type == yaml_SEQUENCE_END_EVENT {
			break
		}

		v = append(v, d.valueInterface())
	}

	d.nextEvent()
	return v
}

// objectInterface is like object but returns map[string]interface{}.
func (d *Decoder) mappingInterface() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})

	d.nextEvent()

	for {
		if d.event.event_type == yaml_MAPPING_END_EVENT {
			break
		}

		key := d.valueInterface()

		// Read value.
		m[key] = d.valueInterface()
	}

	d.nextEvent()
	return m
}

func (d *Decoder) valueInterface() interface{} {
	switch d.event.event_type {
	case yaml_SEQUENCE_START_EVENT:
		return d.sequenceInterface()
	case yaml_MAPPING_START_EVENT:
		return d.mappingInterface()
	case yaml_SCALAR_EVENT:
		return d.scalarInterface()
	case yaml_ALIAS_EVENT:
		d.error(errors.New("alias interface??"))
	case yaml_DOCUMENT_END_EVENT:
	}

	d.error(&UnexpectedEventError{
		Value:     string(d.event.value),
		EventType: d.event.event_type,
		At:        d.event.start_mark,
	})

	panic("unreachable")
}

func (d *Decoder) scalarInterface() interface{} {
	v := resolveInterface(d.event)

	d.nextEvent()
	return v
}

// A field represents a single field found in a struct.
type field struct {
	name      string
	tag       bool
	index     []int
	typ       reflect.Type
	omitEmpty bool
	quoted    bool
}

// byName sorts field by name, breaking ties with depth,
// then breaking ties with "name came from json tag", then
// breaking ties with index sequence.
type byName []field

func (x byName) Len() int { return len(x) }

func (x byName) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byName) Less(i, j int) bool {
	if x[i].name != x[j].name {
		return x[i].name < x[j].name
	}
	if len(x[i].index) != len(x[j].index) {
		return len(x[i].index) < len(x[j].index)
	}
	if x[i].tag != x[j].tag {
		return x[i].tag
	}
	return byIndex(x).Less(i, j)
}

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

// typeFields returns a list of fields that JSON should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func typeFields(t reflect.Type) []field {
	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.PkgPath != "" { // unexported
					continue
				}
				tag := sf.Tag.Get("yaml")
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					fields = append(fields, field{name, tagged, index, ft,
						opts.Contains("omitempty"), opts.Contains("string")})
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	sort.Sort(byName(fields))

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 { // Only one field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))

	return fields
}

// dominantField looks through the fields, all of which are known to
// have the same name, to find the single field that dominates the
// others using Go's embedding rules, modified by the presence of
// JSON tags. If there are multiple top-level fields, the boolean
// will be false: This condition is an error in Go and we skip all
// the fields.
func dominantField(fields []field) (field, bool) {
	// The fields are sorted in increasing index-length order. The winner
	// must therefore be one with the shortest index length. Drop all
	// longer entries, which is easy: just truncate the slice.
	length := len(fields[0].index)
	tagged := -1 // Index of first tagged field.
	for i, f := range fields {
		if len(f.index) > length {
			fields = fields[:i]
			break
		}
		if f.tag {
			if tagged >= 0 {
				// Multiple tagged fields at the same level: conflict.
				// Return no field.
				return field{}, false
			}
			tagged = i
		}
	}
	if tagged >= 0 {
		return fields[tagged], true
	}
	// All remaining fields have the same length. If there's more than one,
	// we have a conflict (two fields named "X" at the same level) and we
	// return no field.
	if len(fields) > 1 {
		return field{}, false
	}
	return fields[0], true
}

var fieldCache struct {
	sync.RWMutex
	m map[reflect.Type][]field
}

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field {
	fieldCache.RLock()
	f := fieldCache.m[t]
	fieldCache.RUnlock()
	if f != nil {
		return f
	}

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = typeFields(t)
	if f == nil {
		f = []field{}
	}

	fieldCache.Lock()
	if fieldCache.m == nil {
		fieldCache.m = map[reflect.Type][]field{}
	}
	fieldCache.m[t] = f
	fieldCache.Unlock()
	return f
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
