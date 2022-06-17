package tfrecordio

import (
	"bufio"
	"context"
	"fmt"
	"hash/fnv"
	"reflect"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/runtime"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/runtime/graphx/schema"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/filesystem"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/log"
	"github.com/gonzojive/beam-go-bazel-example/beamgen"
	"github.com/gonzojive/beam-go-bazel-example/tfrecordio/tfrecord"
	"google.golang.org/protobuf/proto"
)

func init() [
	runtime.RegisterType(reflect.TypeOf((*writeFileFn)(nil)).Elem())
	schema.RegisterType(reflect.TypeOf((*writeFileFn)(nil)).Elem())

	runtime.RegisterType(reflect.TypeOf((*assignShardNumberFn[T])(nil)).Elem())
	schema.RegisterType(reflect.TypeOf((*assignShardNumberFn[T])(nil)).Elem())
]

func shardNum(data []byte, shardCount int) int {
	h := fnv.New32a()
	h.Write(data)
	return int(h.Sum32()) % shardCount
}

// WriteShartded writes a PCollection<[]byte]> to a file using tfrecord format.
func WriteShartded(s beam.Scope, filenamePrefix string, shardCount int, col beamgen.Collection[[]byte]) {
	type T = []byte
	s = s.Scope("textio.Write")

	if shardCount <= 0 {
		panic(fmt.Errorf("invalid shardCount %d <= 0", shardCount))
	}

	filesystem.ValidateScheme(filenamePrefix)

	// NOTE(BEAM-3579): We may never call Teardown for non-local runners and
	// FinishBundle doesn't have the right granularity. We therefore
	// perform a GBK with a fixed key to get all values in a single invocation.

	// TODO(BEAM-3860) 3/15/2018: use side input instead of GBK.

	// var shardNumbers []int
	// for i := 0; i < shardCount; i++ {
	// 	shardNumbers = append(shardNumbers, i)
	// }
	//shardNumbersCol := beamgen.Create(s.Scope("ShardNumbers"), shardNumbers...)

	//beamgen.ParDoKV[InT any, OutK any, OutV any](scope beam.Scope, dofn DoFnInterfaceKVStruct[InT, OutK, OutV], inCol Collection[InT], opts ...beam.Option)
	pre := beamgen.ParDoKV[T, int, T](s.Scope("AssignShardNumber"), &assignShardNumberFn[T]{shardCount}, col)

	//pre := beamgen.AddFixedKey(s, col)
	post := beamgen.GroupByKey(s, pre)
	beamgen.ParDoGBK0[int, T](s, &writeFileFn[T]{Filename: filenamePrefix, ShardCount: shardCount}, post)
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

func (w *writeFileFn) ProcessElement(ctx context.Context, shard int, protos func(*T) bool) error {
	fs, err := filesystem.New(ctx, w.Filename)
	if err != nil {
		return err
	}
	defer fs.Close()

	shardName := fmt.Sprintf("%05d-of-%05d", shard+1, w.ShardCount)

	fd, err := fs.OpenWrite(ctx, w.Filename+"-"+shardName)
	if err != nil {
		return err
	}
	buf := bufio.NewWriterSize(fd, 1000*1000*5) // use 5MB buffer
	recordWriter := tfrecord.NewWriter(buf, nil)

	messageName := zeroT.ProtoReflect().Descriptor().FullName()

	log.Infof(ctx, "Writing riegeli record of %s protos to %v", messageName, w.Filename)

	var elem T
	for protos(&elem) {
		if err := recordWriter.PutProto(elem); err != nil {
			return fmt.Errorf("error writing proto to riegeli file: %w", err)
		}
	}

	if err := recordWriter.Flush(); err != nil {
		return fmt.Errorf("error flushing bytes to riegeli file: %w", err)
	}

	if err := buf.Flush(); err != nil {
		return fmt.Errorf("error flushing bytes to riegeli file: %w", err)
	}
	if err := fd.Close(); err != nil {
		return fmt.Errorf("error closing riegeli file: %w", err)
	}
	return nil
}
