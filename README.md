[![wercker status](https://app.wercker.com/status/0811a08698b2f39f79350758d325cc85/s/ "wercker status")](https://app.wercker.com/project/bykey/0811a08698b2f39f79350758d325cc85)

### Persistent Queue on File System

This project is created to test the performance of various distributed file systems via FUSE interface.


### Running Tests

```bash
cd /path/to/fsqueue
go test

# To run benchmarks
go test -bench=.
```


### Benchmark Results

```
BenchmarkPush   100000       30519 ns/op
testing: BenchmarkPush left GOMAXPROCS set to 4
BenchmarkPop    500000       10410 ns/op
testing: BenchmarkPop left GOMAXPROCS set to 4
BenchmarkPushPop    500000       78938 ns/op
testing: BenchmarkPushPop left GOMAXPROCS set to 1
```