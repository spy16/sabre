package sabre

// Runtime implementation provides ways to allocate collection types and other
// facilities such as Analyze. Runtime acts as a means to customize and extend
// Sabre.
type Runtime interface {
	EmptySeq() Seq
	EmptyMap() Map
	EmptyVec() Vector
	EmptySet() Set
}

// Analyzer can be implemented by Runtime implementations to provide support for
// analyzing custom forms during evaluation.
type Analyzer interface {
	Runtime

	// Analyze should analyze the form given in context of Sabre instance, return
	// an expr that can be evaluated.
	Analyze(s *Sabre, form Value) (Expr, error)
}

// Value represents data/forms in sabre. This includes those emitted by Reader,
// values obtained as result of an evaluation etc..
type Value interface {
	// String should return the LISP representation of the value.
	String() string
}

// Invokable represents a value that can be invoked when it appears as the first
// entry in a list.
type Invokable interface {
	Value

	// Invoke is called when this value appears as first item in a list. Remaining
	// items of the list will be passed un-evaluated as arguments.
	Invoke(env *Sabre, args ...Value) (Value, error)
}

// Seq implements a sequence of values (e.g., List) that may be realized lazily.
type Seq interface {
	Value

	// First returns the first value of the sequence if not empty. Returns 'nil'
	// if empty.
	First() Value

	// Next returns the remaining sequence when the first value of the sequence
	// is excluded. 'nil' if the sequence is empty or has single item.
	Next() Seq

	// Conj returns a new sequence which includes values from this sequence and
	// the arguments. Position of conjoined values is not part of the contract.
	Conj(vals ...Value) Seq

	// Count returns the number of items in the map.
	Count() int
}

// Seqable is any value that can be converted to a sequence.
type Seqable interface {
	Value

	// Seq returns the implementing value as a sequence.
	Seq() Seq
}

// Vector represents a container for values that provides fast index lookups and
// iterations.
type Vector interface {
	Seqable

	// Count returns the number of elements in the vector.
	Count() int

	// EntryAt returns the item at given index. Returns error if the index
	// is out of range.
	EntryAt(index int) (Value, error)

	// Conj returns a new vector with items appended.
	Conj(items ...Value) Vector

	// Assoc returns a new vector with the value at given index updated.
	// Returns error if the index is out of range.
	Assoc(index int, val Value) (Vector, error)
}

// Map represents any value that can store key-value pairs and provide fast
// lookups.
type Map interface {
	Seqable

	// Keys returns all the keys in the map as a sequence.
	Keys() Seq

	// Vals returns all the values in the map as a sequence.
	Vals() Seq

	// HasKey returns true if the map contains the given key.
	HasKey(key Value) bool

	// EntryAt returns the value associated with the given key. Returns nil
	// if the key is not found or is not hashable.
	EntryAt(key Value) Value

	// Assoc should return a new map which contains all the current values
	// with the given key-val pair added.
	Assoc(key, val Value) (Map, error)

	// Dissoc should return a new map which contains all the current entries
	// except the one with given key.
	Dissoc(key Value) (Map, error)
}

// Set represents a container for storing unique values.
type Set interface {
	Seqable

	// HasKey returns true if the key is present in the set.
	HasKey(key Value) bool

	// Conj returns a new set with the vals conjoined.
	Conj(vals ...Value) Set

	// Disj returns a new set with the vals dis-joined.
	Disj(vals ...Value) Set
}
