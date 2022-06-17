package tfrecordio

import (
	"context"
	"fmt"
	"hash/fnv"
	"reflect"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/runtime"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/runtime/graphx/schema"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/filesystem"
	"github.com/gonzojive/beam-go-bazel-example/beamgen"
	"github.com/gonzojive/beam-go-bazel-example/tfrecordio/tfrecord"
)

func init() {
	runtime.RegisterType(reflect.TypeOf((*writeFileFn)(nil)).Elem())
	schema.RegisterType(reflect.TypeOf((*writeFileFn)(nil)).Elem())

	runtime.RegisterType(reflect.TypeOf((*assignShardNumberFn)(nil)).Elem())
	schema.RegisterType(reflect.TypeOf((*assignShardNumberFn)(nil)).Elem())
}

func shardNum(data []byte, shardCount int) int {
	h := fnv.New32a()
	h.Write(data)
	return int(h.Sum32()) % shardCount
}

// WriteSharded writes a PCollection<[]byte]> to a file using tfrecord format.
func WriteSharded(s beam.Scope, filenamePrefix string, shardCount int, col beamgen.Collection[[]byte]) {
	type T = []byte
	s = s.Scope("tfrecord.Write")

	if shardCount <= 0 {
		panic(fmt.Errorf("invalid shardCount %d <= 0", shardCount))
	}

	filesystem.ValidateScheme(filenamePrefix)

	// NOTE(BEAM-3579): We may never call Teardown for non-local runners and
	// FinishBundle doesn't have the right granularity. We therefore
	// perform a GBK with a fixed key to get all values in a single invocation.

	// TODO(BEAM-3860) 3/15/2018: use side input instead of GBK.

	pre := beamgen.ParDoKV[T, int, T](s.Scope("AssignShardNumber"), &assignShardNumberFn{shardCount}, col)

	//pre := beamgen.AddFixedKey(s, col)
	post := beamgen.GroupByKey(s, pre)
	beamgen.ParDoGBK0[int, T](s, &writeFileFn{Filename: filenamePrefix, ShardCount: shardCount}, post)
}

type assignShardNumberFn struct {
	ShardCount int
}

func (f *assignShardNumberFn) ProcessElement(ctx context.Context, record []byte, emit func(int, []byte)) error {
	emit(shardNum(record, f.ShardCount), record)
	return nil
}

type writeFileFn struct {
	Filename   string `json:"filename"`
	ShardCount int    `json:"shardCount"`
}

func (w *writeFileFn) ProcessElement(ctx context.Context, shard int, protos func(*[]byte) bool) error {
	fs, err := filesystem.New(ctx, w.Filename)
	if err != nil {
		return err
	}
	defer fs.Close()

	shardName := fmt.Sprintf("%05d-of-%05d", shard+1, w.ShardCount)

	filename := w.Filename + "-" + shardName
	recordWriter, err := tfrecord.NewWriter(filename, &tfrecord.RecordWriterOptions{
		CompressionType: tfrecord.CompressionTypeNone,
	})
	if err != nil {
		return fmt.Errorf("error creating record writer: %w", err)
	}

	var elem []byte
	for protos(&elem) {
		if err := recordWriter.WriteRecord(elem); err != nil {
			return fmt.Errorf("error writing proto to TFRecord file: %w", err)
		}
	}

	if err := recordWriter.Flush(); err != nil {
		return fmt.Errorf("error flushing bytes to TFRecord file: %w", err)
	}

	if err := recordWriter.Close(); err != nil {
		return fmt.Errorf("error closing TFRecord file: %w", err)
	}
	return nil
}
