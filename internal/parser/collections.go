package parser

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"sync"
)

// ============================================================================
// Array
// ============================================================================

// Array represents a PDF array object.
// Arrays are ordered collections of objects: [obj1 obj2 obj3]
// Arrays can contain any PDF objects, including other arrays and dictionaries.
type Array struct {
	elements []PdfObject
	mu       sync.RWMutex // Thread-safe operations
}

// NewArray creates a new empty Array.
func NewArray() *Array {
	return &Array{
		elements: make([]PdfObject, 0),
	}
}

// NewArrayWithCapacity creates a new Array with specified capacity.
func NewArrayWithCapacity(capacity int) *Array {
	return &Array{
		elements: make([]PdfObject, 0, capacity),
	}
}

// NewArrayFromSlice creates an Array from a slice of objects.
func NewArrayFromSlice(objects []PdfObject) *Array {
	arr := &Array{
		elements: make([]PdfObject, len(objects)),
	}
	copy(arr.elements, objects)
	return arr
}

// Len returns the number of elements in the array.
func (a *Array) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.elements)
}

// Get returns the element at index i.
// Returns nil if index is out of bounds.
func (a *Array) Get(i int) PdfObject {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if i < 0 || i >= len(a.elements) {
		return nil
	}
	return a.elements[i]
}

// Set sets the element at index i.
// Returns error if index is out of bounds.
func (a *Array) Set(i int, obj PdfObject) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if i < 0 || i >= len(a.elements) {
		return fmt.Errorf("index %d out of bounds (len=%d)", i, len(a.elements))
	}
	a.elements[i] = obj
	return nil
}

// Append adds an element to the end of the array.
func (a *Array) Append(obj PdfObject) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.elements = append(a.elements, obj)
}

// AppendAll adds multiple elements to the end of the array.
func (a *Array) AppendAll(objects ...PdfObject) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.elements = append(a.elements, objects...)
}

// Insert inserts an element at index i.
func (a *Array) Insert(i int, obj PdfObject) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if i < 0 || i > len(a.elements) {
		return fmt.Errorf("index %d out of bounds (len=%d)", i, len(a.elements))
	}

	// Grow slice
	a.elements = append(a.elements, nil)
	// Shift elements
	copy(a.elements[i+1:], a.elements[i:])
	// Insert new element
	a.elements[i] = obj

	return nil
}

// Remove removes the element at index i.
func (a *Array) Remove(i int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if i < 0 || i >= len(a.elements) {
		return fmt.Errorf("index %d out of bounds (len=%d)", i, len(a.elements))
	}

	a.elements = append(a.elements[:i], a.elements[i+1:]...)
	return nil
}

// Clear removes all elements from the array.
func (a *Array) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.elements = make([]PdfObject, 0)
}

// Elements returns a copy of all elements.
// Returns a new slice to prevent external modification.
func (a *Array) Elements() []PdfObject {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]PdfObject, len(a.elements))
	copy(result, a.elements)
	return result
}

// String returns a string representation of the array.
func (a *Array) String() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, elem := range a.elements {
		if i > 0 {
			buf.WriteByte(' ')
		}
		if elem != nil {
			buf.WriteString(elem.String())
		} else {
			buf.WriteString("null")
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

// WriteTo writes the PDF representation of the array to w.
func (a *Array) WriteTo(w io.Writer) (int64, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var total int64

	// Write opening bracket
	n, err := w.Write([]byte("["))
	total += int64(n)
	if err != nil {
		return total, err
	}

	// Write elements
	for i, elem := range a.elements {
		// Space before element (except first)
		if i > 0 {
			n, err := w.Write([]byte(" "))
			total += int64(n)
			if err != nil {
				return total, err
			}
		}

		// Write element
		if elem != nil {
			written, err := elem.WriteTo(w)
			total += written
			if err != nil {
				return total, err
			}
		} else {
			n, err := w.Write([]byte("null"))
			total += int64(n)
			if err != nil {
				return total, err
			}
		}
	}

	// Write closing bracket
	n, err = w.Write([]byte("]"))
	total += int64(n)
	return total, err
}

// Clone creates a deep copy of the array.
func (a *Array) Clone() *Array {
	a.mu.RLock()
	defer a.mu.RUnlock()

	cloned := NewArrayWithCapacity(len(a.elements))
	for _, elem := range a.elements {
		cloned.elements = append(cloned.elements, Clone(elem))
	}
	return cloned
}

// ============================================================================
// Dictionary
// ============================================================================

// Dictionary represents a PDF dictionary object.
// Dictionaries are associative tables: << /Key1 value1 /Key2 value2 >>
// Keys are always Name objects, values can be any PDF object.
type Dictionary struct {
	entries map[string]PdfObject
	keys    []string // Maintains insertion order
	mu      sync.RWMutex
}

// NewDictionary creates a new empty Dictionary.
func NewDictionary() *Dictionary {
	return &Dictionary{
		entries: make(map[string]PdfObject),
		keys:    make([]string, 0),
	}
}

// NewDictionaryWithCapacity creates a new Dictionary with specified capacity.
func NewDictionaryWithCapacity(capacity int) *Dictionary {
	return &Dictionary{
		entries: make(map[string]PdfObject, capacity),
		keys:    make([]string, 0, capacity),
	}
}

// Len returns the number of entries in the dictionary.
func (d *Dictionary) Len() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.entries)
}

// Has checks if a key exists in the dictionary.
func (d *Dictionary) Has(key string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.entries[key]
	return exists
}

// Get returns the value for a key.
// Returns nil if key doesn't exist.
func (d *Dictionary) Get(key string) PdfObject {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.entries[key]
}

// GetName is a convenience method to get a Name value.
// Returns nil if key doesn't exist or value is not a Name.
func (d *Dictionary) GetName(key string) *Name {
	obj := d.Get(key)
	if name, ok := obj.(*Name); ok {
		return name
	}
	return nil
}

// GetInteger is a convenience method to get an Integer value.
// Returns 0 if key doesn't exist or value is not an Integer.
func (d *Dictionary) GetInteger(key string) int64 {
	obj := d.Get(key)
	if i, ok := obj.(*Integer); ok {
		return i.Value()
	}
	return 0
}

// GetReal is a convenience method to get a Real value.
// Returns 0.0 if key doesn't exist or value is not a Real.
func (d *Dictionary) GetReal(key string) float64 {
	obj := d.Get(key)
	if r, ok := obj.(*Real); ok {
		return r.Value()
	}
	return 0.0
}

// GetBoolean is a convenience method to get a Boolean value.
// Returns false if key doesn't exist or value is not a Boolean.
func (d *Dictionary) GetBoolean(key string) bool {
	obj := d.Get(key)
	if b, ok := obj.(*Boolean); ok {
		return b.Value()
	}
	return false
}

// GetString is a convenience method to get a String value.
// Returns empty string if key doesn't exist or value is not a String.
func (d *Dictionary) GetString(key string) string {
	obj := d.Get(key)
	if s, ok := obj.(*String); ok {
		return s.Value()
	}
	return ""
}

// GetArray is a convenience method to get an Array value.
// Returns nil if key doesn't exist or value is not an Array.
func (d *Dictionary) GetArray(key string) *Array {
	obj := d.Get(key)
	if arr, ok := obj.(*Array); ok {
		return arr
	}
	return nil
}

// GetDictionary is a convenience method to get a Dictionary value.
// Returns nil if key doesn't exist or value is not a Dictionary.
func (d *Dictionary) GetDictionary(key string) *Dictionary {
	obj := d.Get(key)
	if dict, ok := obj.(*Dictionary); ok {
		return dict
	}
	return nil
}

// Set sets a key-value pair in the dictionary.
// If key already exists, its value is replaced.
func (d *Dictionary) Set(key string, value PdfObject) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// If key doesn't exist, add to keys slice
	if _, exists := d.entries[key]; !exists {
		d.keys = append(d.keys, key)
	}

	d.entries[key] = value
}

// SetName is a convenience method to set a Name value.
func (d *Dictionary) SetName(key, value string) {
	d.Set(key, NewName(value))
}

// SetInteger is a convenience method to set an Integer value.
func (d *Dictionary) SetInteger(key string, value int64) {
	d.Set(key, NewInteger(value))
}

// SetReal is a convenience method to set a Real value.
func (d *Dictionary) SetReal(key string, value float64) {
	d.Set(key, NewReal(value))
}

// SetBoolean is a convenience method to set a Boolean value.
func (d *Dictionary) SetBoolean(key string, value bool) {
	d.Set(key, NewBoolean(value))
}

// SetString is a convenience method to set a String value.
func (d *Dictionary) SetString(key, value string) {
	d.Set(key, NewString(value))
}

// Remove removes a key from the dictionary.
func (d *Dictionary) Remove(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.entries[key]; exists {
		delete(d.entries, key)

		// Remove from keys slice
		for i, k := range d.keys {
			if k == key {
				d.keys = append(d.keys[:i], d.keys[i+1:]...)
				break
			}
		}
	}
}

// Keys returns all keys in insertion order.
func (d *Dictionary) Keys() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]string, len(d.keys))
	copy(result, d.keys)
	return result
}

// KeysSorted returns all keys in alphabetical order.
func (d *Dictionary) KeysSorted() []string {
	keys := d.Keys()
	sort.Strings(keys)
	return keys
}

// Clear removes all entries from the dictionary.
func (d *Dictionary) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.entries = make(map[string]PdfObject)
	d.keys = make([]string, 0)
}

// String returns a string representation of the dictionary.
func (d *Dictionary) String() string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var buf bytes.Buffer
	buf.WriteString("<<")

	for i, key := range d.keys {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteByte('/')
		buf.WriteString(key)
		buf.WriteByte(' ')

		value := d.entries[key]
		if value != nil {
			buf.WriteString(value.String())
		} else {
			buf.WriteString("null")
		}
	}

	buf.WriteString(">>")
	return buf.String()
}

// WriteTo writes the PDF representation of the dictionary to w.
func (d *Dictionary) WriteTo(w io.Writer) (int64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var total int64

	// Write opening
	n, err := w.Write([]byte("<<"))
	total += int64(n)
	if err != nil {
		return total, err
	}

	// Write key-value pairs
	for i, key := range d.keys {
		// Space before entry (except first)
		if i > 0 {
			n, err := w.Write([]byte(" "))
			total += int64(n)
			if err != nil {
				return total, err
			}
		}

		// Write key (as Name)
		name := NewName(key)
		written, err := name.WriteTo(w)
		total += written
		if err != nil {
			return total, err
		}

		// Space between key and value
		n, err := w.Write([]byte(" "))
		total += int64(n)
		if err != nil {
			return total, err
		}

		// Write value
		value := d.entries[key]
		if value != nil {
			written, err := value.WriteTo(w)
			total += written
			if err != nil {
				return total, err
			}
		} else {
			n, err := w.Write([]byte("null"))
			total += int64(n)
			if err != nil {
				return total, err
			}
		}
	}

	// Write closing
	n, err = w.Write([]byte(">>"))
	total += int64(n)
	return total, err
}

// Clone creates a deep copy of the dictionary.
func (d *Dictionary) Clone() *Dictionary {
	d.mu.RLock()
	defer d.mu.RUnlock()

	cloned := NewDictionaryWithCapacity(len(d.entries))
	for _, key := range d.keys {
		cloned.Set(key, Clone(d.entries[key]))
	}
	return cloned
}

// Merge merges another dictionary into this one.
// Existing keys are overwritten.
func (d *Dictionary) Merge(other *Dictionary) {
	if other == nil {
		return
	}

	for _, key := range other.Keys() {
		d.Set(key, other.Get(key))
	}
}
