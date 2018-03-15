# mqtgo
mqtgo is a pure MQTT protocol packets in go focus on performence

this project was based on eclipse.paho mqtt packet

test coverage: 83.8% of statements

## benchmark
* goos: darwin
* goarch: amd64

```
BenchmarkNewControlPacket-4             20000000              85.5 ns/op               0 B/op          0 allocs/op
BenchmarkParallelNewControlPacket-4     10000000               169 ns/op               0 B/op          0 allocs/op
BenchmarkReadPacket-4                     500000              2340 ns/op             930 B/op          2 allocs/op
```

compare with paho mqtt packets
```
benchmark                               old ns/op     new ns/op     delta
BenchmarkNewControlPacket-4             90.9          85.5          -5.94%
BenchmarkParallelNewControlPacket-4     178           169           -5.06%
BenchmarkReadPacket-4                   2707          2340          -13.56%

benchmark                               old allocs     new allocs     delta
BenchmarkNewControlPacket-4             1              0              -100.00%
BenchmarkParallelNewControlPacket-4     1              0              -100.00%
BenchmarkReadPacket-4                   8              2              -75.00%

benchmark                               old bytes     new bytes     delta
BenchmarkNewControlPacket-4             44            0             -100.00%
BenchmarkParallelNewControlPacket-4     44            0             -100.00%
BenchmarkReadPacket-4                   2116          930           -56.05%
```


