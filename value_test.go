package sabre_test

import "testing"

import "github.com/spy16/sabre"

import "reflect"

func TestValues_First(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		vals := sabre.Values{}

		want := sabre.Nil{}
		got := vals.First()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("First() want=%#v, got=%#v", want, got)
		}
	})

	t.Run("Nil", func(t *testing.T) {
		vals := sabre.Values(nil)

		want := sabre.Nil{}
		got := vals.First()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("First() want=%#v, got=%#v", want, got)
		}
	})

	t.Run("NonEmpty", func(t *testing.T) {
		vals := sabre.Values{sabre.Int64(10)}

		want := sabre.Int64(10)
		got := vals.First()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("First() want=%#v, got=%#v", want, got)
		}
	})
}

func TestValues_Next(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		vals []sabre.Value
		want sabre.Seq
	}{
		{
			name: "Nil",
			vals: []sabre.Value(nil),
			want: nil,
		},
		{
			name: "Empty",
			vals: []sabre.Value{},
			want: nil,
		},
		{
			name: "SingleItem",
			vals: []sabre.Value{sabre.Int64(10)},
			want: nil,
		},
		{
			name: "MultiItem",
			vals: []sabre.Value{sabre.Int64(10), sabre.String("hello"), sabre.Bool(true)},
			want: sabre.Values{sabre.String("hello"), sabre.Bool(true)},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := sabre.Values(tt.vals).Next()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Next() want=%#v, got=%#v", tt.want, got)
			}
		})
	}
}

func TestValues_Cons(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		vals []sabre.Value
		item sabre.Value
		want sabre.Seq
	}{
		{
			name: "Nil",
			vals: []sabre.Value(nil),
			item: sabre.Int64(10),
			want: sabre.Values{sabre.Int64(10)},
		},
		{
			name: "Empty",
			vals: []sabre.Value{},
			item: sabre.Int64(10),
			want: sabre.Values{sabre.Int64(10)},
		},
		{
			name: "SingleItem",
			vals: []sabre.Value{sabre.Int64(10)},
			item: sabre.String("hello"),
			want: sabre.Values{sabre.String("hello"), sabre.Int64(10)},
		},
		{
			name: "MultiItem",
			vals: []sabre.Value{sabre.Int64(10), sabre.String("hello")},
			item: sabre.Bool(true),
			want: sabre.Values{sabre.Bool(true), sabre.Int64(10), sabre.String("hello")},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := sabre.Values(tt.vals).Cons(tt.item)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Next() want=%#v, got=%#v", tt.want, got)
			}
		})
	}
}
