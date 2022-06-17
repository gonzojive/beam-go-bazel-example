## Reproduction of OOM issue

This is a reproduction of a beam pipeline crashing due to OOM.
https://github.com/apache/beam/issues/21817

Start up the Flink server

```shell
bazel run :flink-job-server-driver -- \
  --flink-conf-dir $PWD/flink-conf/ \
  --job-port 8099
```

```shell
mkdir tmp

bazel run //cmd/example_pipeline -- \
  --alsologtostderr \
  --runner flink \
  --environment_type LOOPBACK \
  --endpoint localhost:8099 \
  --output "$PWD/tmp/records" \
  --record-bytes 25000000 \
  --record-count 10000
```
