// Package beamgen contains utilites for working with beam pipelines using
// generics.
package beamgen

import (
	"context"
	"reflect"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/runtime"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/runtime/graphx/schema"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/samber/lo"
)

// Collection is a typed PCollection.
type Collection[T any] struct {
	underlying beam.PCollection
}

// PCollection returns a generic PCollection from a typed PCollection.
func (c Collection[T]) PCollection() beam.PCollection { return c.underlying }

// AssertType returns a Collection[T] from a PCollection that asserts it has a
// given type.
func AssertType[T any](coll beam.PCollection) Collection[T] {
	return Collection[T]{coll}
}

// Create returns a PCollection from a set of in-memory objects.
//
// Create inserts a fixed non-empty set of values into the pipeline. The values must
// be of the same type 'A' and the returned PCollection is of type A.
//
// The returned PCollections can be used as any other PCollections. The values
// are JSON-coded. Each runner may place limits on the sizes of the values and
// Create should generally only be used for small collections.
func Create[T any](scope beam.Scope, obj ...T) Collection[T] {
	asInterfaces := lo.Map(obj, func(elem T, _ int) any { return elem })
	untyped := beam.Create(scope, asInterfaces...)
	return Collection[T]{untyped}
}

//func ParDoFunc[InT, OutT any](scope beam.Scope, in Collection[InT], fn func(value InT))

// ParDoUnsafe is like beam.ParDo in that it accepts an `any` dofn, but the
// input and output collections are typed.
func ParDoUnsafe[InT, OutT any](scope beam.Scope, dofn any, inCol Collection[InT], opts ...beam.Option) Collection[OutT] {
	return Collection[OutT]{
		beam.ParDo(scope, dofn, inCol.PCollection(), opts...),
	}
}

type DoFnInterfaceStruct1[InT, OutT1 any] interface {
	ProcessElement(ctx context.Context, value InT, emit func(out OutT1)) error
}

// ParDo1 is like beam.ParDo in that it accepts an `any` dofn, but the
// input and output collections are typed.
func ParDo1[InT, OutT any](scope beam.Scope, dofn DoFnInterfaceStruct1[InT, OutT], inCol Collection[InT], opts ...beam.Option) Collection[OutT] {
	return Collection[OutT]{
		beam.ParDo(scope, dofn, inCol.PCollection(), opts...),
	}
}

type DoFnInterfaceStruct2[InT, OutT1, OutT2 any] interface {
	ProcessElement(ctx context.Context, value InT, emit1 func(OutT1), emit2 func(OutT2)) error
}

// ParDo2 is like beam.ParDo in that it accepts an `any` dofn, but the
// input and output collections are typed.
func ParDo2[InT, OutT1, OutT2 any](scope beam.Scope, dofn DoFnInterfaceStruct2[InT, OutT1, OutT2], inCol Collection[InT], opts ...beam.Option) (Collection[OutT1], Collection[OutT2]) {
	c1, c2 := beam.ParDo2(scope, dofn, inCol.PCollection(), opts...)
	return Collection[OutT1]{c1}, Collection[OutT2]{c2}
}

// ParDo1Func is like beam.ParDo in that it accepts an `any` dofn, but the
// input and output collections are typed.
func ParDo1Func[InT, OutT any](scope beam.Scope, dofn func(ctx context.Context, in InT, emit OutT) error, inCol Collection[InT], opts ...beam.Option) Collection[OutT] {
	return Collection[OutT]{
		beam.ParDo(scope, dofn, inCol.PCollection(), opts...),
	}
}

// ParDoGBKFunc is used to handle the result of a call to GroupByKey.
func ParDoGBKFunc[InK, InV, OutT any](
	scope beam.Scope,
	dofn func(ctx context.Context, in InK, next func(*InV) bool, emit func(OutT)) error,
	inCol Collection[GroupedByKey[InK, InV]],
	opts ...beam.Option) Collection[OutT] {
	return Collection[OutT]{
		beam.ParDo(scope, dofn, inCol.PCollection(), opts...),
	}
}

// ParDoGBK is like beam.ParDo in that it accepts an `any` dofn, but the
// input and output collections are typed.
func ParDoGBK[InK, InV, OutT any](
	scope beam.Scope,
	dofn interface {
		ProcessElement(ctx context.Context, in InK, next func(*InV) bool, emit func(OutT)) error
	},
	inCol Collection[GroupedByKey[InK, InV]],
	opts ...beam.Option) Collection[OutT] {
	return Collection[OutT]{
		beam.ParDo(scope, dofn, inCol.PCollection(), opts...),
	}
}

// ParDoGBK0 is like beam.ParDo in that it accepts an `any` dofn, but the
// input and output collections are typed.
func ParDoGBK0[InK, InV any](
	scope beam.Scope,
	dofn DoFnInterfaceGBK0[InK, InV],
	inCol Collection[GroupedByKey[InK, InV]],
	opts ...beam.Option) {
	beam.ParDo0(scope, dofn, inCol.PCollection(), opts...)
}

type DoFnInterfaceGBK0[InK, InV any] interface {
	ProcessElement(ctx context.Context, in InK, next func(*InV) bool) error
}

// DoFnInterfaceKVStruct is a DoFn that outputs a PCollection<KV<OutK, OutV>> from a PCollection<InT>
type DoFnInterfaceKVStruct[InT, OutK, OutV any] interface {
	ProcessElement(ctx context.Context, value InT, emit func(key OutK, value OutV)) error
}

// ParDoKV is a key/value version of ParDo
func ParDoKV[InT, OutK, OutV any](scope beam.Scope, dofn DoFnInterfaceKVStruct[InT, OutK, OutV], inCol Collection[InT], opts ...beam.Option) Collection[KV[OutK, OutV]] {
	return Collection[KV[OutK, OutV]]{
		beam.ParDo(scope, dofn, inCol.PCollection(), opts...),
	}
}

// A type for holding a key and value, mostly used only as a type.
type KV[K, V any] struct {
	// Key   K `json:"key"`
	// Value V `json:"value"`
}

type GroupedByKey[K, V any] struct{}

// GroupByKey is a PTransform that takes a PCollection of type KV<A,B>,
// groups the values by key and windows, and returns a PCollection of type
// GBK<A,B> representing a map from each distinct key and window of the
// input PCollection to an iterable over all the values associated with
// that key in the input per window. Each key in the output PCollection is
// unique within each window.
func GroupByKey[K, V any](scope beam.Scope, col Collection[KV[K, V]]) Collection[GroupedByKey[K, V]] {
	return Collection[GroupedByKey[K, V]]{
		beam.GroupByKey(scope, col.underlying),
	}
}

// type KeyExtractor[V, K any] struct {
// 	fn func
// }

// func KeyBySimple[K, V](scope beam.Scope, col Collection[V], keyFn func(elem V) K) Collection

// RemoveDuplicates takes a PCollection with duplicate entries and returns a
// PCollection with the duplicates removed.
func RemoveDuplicates[T any](scope beam.Scope, col Collection[T]) Collection[T] {
	scope = scope.Scope("RemoveDuplcates")
	var table Collection[KV[T, T]] = ParDoKV[T, T, T](scope, &xToKVXXFn[T]{}, col)
	var grouped Collection[GroupedByKey[T, T]] = GroupByKey(scope, table)
	return ParDoGBK[T, T, T](scope, &keysOfGBKFn[T, T]{}, grouped)
}

func Reshuffle[T any](scope beam.Scope, col Collection[T]) Collection[T] {
	return Collection[T](AssertType[T](beam.Reshuffle(scope, col.PCollection())))
}

func RemoveDuplicatesInit[T any]() {
	runtime.RegisterType(reflect.TypeOf((*keysOfGBKFn[string, string])(nil)).Elem())
	schema.RegisterType(reflect.TypeOf((*keysOfGBKFn[string, string])(nil)).Elem())
	runtime.RegisterType(reflect.TypeOf((*xToKVXXFn[string])(nil)).Elem())
	schema.RegisterType(reflect.TypeOf((*xToKVXXFn[string])(nil)).Elem())
}

type xToKVXXFn[T any] struct{}

func (f *xToKVXXFn[T]) ProcessElement(_ context.Context, in T, emit func(T, T)) error {
	emit(in, in)
	return nil
}

type keysOfGBKFn[K, V any] struct{}

func (f *keysOfGBKFn[K, V]) ProcessElement(_ context.Context, key K, nextValue func(*V) bool, emit func(K)) error {
	emit(key)
	return nil
}

func keysOfGBKFnOld[K, V any](_ context.Context, key K, nextValue func(*V) bool, emit func(K)) error {
	emit(key)
	return nil
}

func TextIOWrite(scope beam.Scope, filename string, col Collection[string]) {
	textio.Write(scope, filename, col.PCollection())
}

// IterToSlice returns a slice from a beam iterator.
func IterToSlice[T any](next func(*T) bool) []T {
	var out []T
	for {
		var value T
		if !next(&value) {
			break
		}
		out = append(out, value)
	}
	return out
}

// IterForEachErr returns a slice from a beam iterator.
func IterForEachErr[T any](next func(*T) bool, fn func(t T) error) error {
	for {
		var value T
		if !next(&value) {
			return nil
		}
		if err := fn(value); err != nil {
			return err
		}
	}
}

// AddFixedKey adds a fixed key (0) to every element.
func AddFixedKey[T any](scope beam.Scope, col Collection[T]) Collection[KV[int, T]] {
	return Collection[KV[int, T]]{
		beam.AddFixedKey(scope, col.PCollection()),
	}
}

// func GroupByKey(scope beam.Scope, a PCollection) PCollection {

// }

// TODO
//type DoFn[InT, OutT any] struct {
//	asInterface any
//}
