package collection

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spy16/sabre/sabre/core"
)

var (
	_ core.Value = (*HashMap)(nil)
	_ core.Map   = (*HashMap)(nil)
)

// HashMap represents a container for key-value pairs implemented using Go
// native hashmap.
type HashMap struct {
	core.Position
	data map[core.Value]core.Value
}

// Eval evaluates all keys and values and returns a new HashMap containing
// the evaluated values.
func (hm *HashMap) Eval(env core.Env) (core.Value, error) {
	res := &HashMap{data: map[core.Value]core.Value{}}
	for k, v := range hm.data {
		kvRes, err := core.EvalAll(env, []core.Value{k, v})
		if err != nil {
			return nil, err
		}

		key, val := kvRes[0], kvRes[1]
		if !isHashable(key) {
			return nil, fmt.Errorf("value of type '%s' is not hashable", reflect.TypeOf(key))
		}

		res.data[key] = val
	}

	return res, nil
}

func (hm *HashMap) String() string {
	lst := core.Seq(&core.List{})
	for k, v := range hm.data {
		lst = lst.Conj(k, v)
	}
	return core.SeqString(lst, "{", "}", " ")
}

// Count returns the number of entries in the map.
func (hm *HashMap) Count() int { return len(hm.data) }

// Get returns the value associated with the given key if found. Returns error
// otherwise.
func (hm *HashMap) Get(key core.Value) (core.Value, error) {
	if !isHashable(key) {
		return nil, core.ErrNotFound
	}

	v, found := hm.data[key]
	if !found {
		return nil, nil
	}
	return v, nil
}

// Assoc sets/updates the value associated with the given key.
func (hm *HashMap) Assoc(key, val core.Value) (core.Map, error) {
	if !isHashable(key) {
		return nil, fmt.Errorf("value of type '%s' is not hashable", reflect.TypeOf(key))
	}

	if hm.data == nil {
		hm.data = map[core.Value]core.Value{}
	}

	hm.data[key] = val
	return hm, nil
}

// Dissoc removes the entry with given key from the hash map and returns
// the new map.
func (hm *HashMap) Dissoc(key core.Value) (core.Map, error) {
	if !isHashable(key) {
		return nil, fmt.Errorf("value of type '%s' is not hashable", reflect.TypeOf(key))
	}

	if hm.data != nil {
		delete(hm.data, key)
	}
	return hm, nil
}

// HasKey returns true if the given key exists in the hash map.
func (hm *HashMap) HasKey(key core.Value) bool {
	if !isHashable(key) {
		return false
	}
	_, found := hm.data[key]
	return found
}

// Keys returns all the keys in the hashmap.
func (hm *HashMap) Keys() core.Seq {
	lst := core.Seq(&core.List{})
	for k := range hm.data {
		lst = lst.Conj(k)
	}
	return lst
}

// Vals returns all the values in the hashmap.
func (hm *HashMap) Vals() core.Seq {
	lst := core.Seq(&core.List{})
	for _, v := range hm.data {
		lst = lst.Conj(v)
	}
	return lst
}

// HashMapReader implements reader macro for reading a hash map from source.
func HashMapReader(rd *core.Reader, _ rune) (core.Value, error) {
	const mapEnd = '}'

	pi := rd.Position()
	forms, err := rd.Container(mapEnd, "hash-map")
	if err != nil {
		return nil, err
	}

	if len(forms)%2 != 0 {
		return nil, errors.New("expecting even number of forms within {}")
	}

	m := core.Map(&HashMap{Position: pi})

	for i := 0; i < len(forms); i += 2 {
		if m.HasKey(forms[i]) {
			return nil, fmt.Errorf("duplicate key: %v", forms[i])
		}

		m, err = m.Assoc(forms[i], forms[i+1])
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func isHashable(v core.Value) bool {
	switch v.(type) {
	case core.String, core.Int64, core.Float64, core.Nil,
		core.Char, core.Keyword, core.Symbol:
		return true

	default:
		return false
	}
}
