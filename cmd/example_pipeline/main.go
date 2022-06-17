package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
	"github.com/golang/glog"
	"github.com/gonzojive/beam-go-bazel-example/beamgen"
	"github.com/gonzojive/beam-go-bazel-example/tfrecordio"
)

var (
	recordsOutput = flag.String("output", "", "output tfrecords prefix")
	shardCount    = flag.Int("shard-count", 5, "number of output shards")
	recordSize    = flag.Int("record-bytes", 1024*1024*3, "output tfrecords prefix")
	recordCount   = flag.Int("record-count", 1000, "output tfrecords prefix")
)

func init() {
	beam.RegisterDoFn(reflect.TypeOf(randomBytesGeneratorFn{}))
}

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("runtime error: %v", err)
	}
}

func run(ctx context.Context) error {
	// Needed to transmit the worker binary to the server without relying on `go
	// build`, which typically doesn't work with bazel builds.
	if err := setWorkerBinaryFlag(); err != nil {
		return fmt.Errorf("failed to set --worker_binary: %w", err)
	}
	beam.Init()

	p := beam.NewPipeline()
	s := p.Root()
	var counts beamgen.Collection[int] = beamgen.Create[int](s, *recordCount)
	var records beamgen.Collection[[]byte] = beamgen.ParDo1[int, []byte](s, randomBytesGeneratorFn{
		Seed:           42,
		BytesPerRecord: *recordSize,
	}, counts)

	tfrecordio.WriteSharded(s, *recordsOutput, *shardCount, records)

	if err := beamx.Run(ctx, p); err != nil {
		return fmt.Errorf("failed to execute job: %w", err)
	}
	fmt.Println("Pipeline completed successfully.")
	return nil
}

func setWorkerBinaryFlag() error {
	binPath, err := os.Executable()
	if err != nil {
		return err
	}
	flag.Set("worker_binary", binPath)
	return nil
}

type randomBytesGeneratorFn struct {
	Seed           int64 `json:"seed"`
	BytesPerRecord int   `json:"size"`
}

func (fn randomBytesGeneratorFn) ProcessElement(_ context.Context, count int, emit func([]byte)) error {
	r := rand.New(rand.NewSource(fn.Seed))
	for i := 0; i < count; i++ {
		buf := make([]byte, fn.BytesPerRecord)
		if _, err := r.Read(buf); err != nil {
			return fmt.Errorf("error generating record %d: %w", i, err)
		}
		emit(buf)
	}
	return nil
}
