# mqtgo [![Build Status](https://travis-ci.org/arthurkiller/mqtgo.svg?branch=master)](https://travis-ci.org/arthurkiller/mqtgo) [![Go Report Card](https://goreportcard.com/badge/github.com/arthurkiller/mqtgo)](https://goreportcard.com/report/github.com/arthurkiller/mqtgo) [![codecov](https://codecov.io/gh/arthurkiller/mqtgo/branch/master/graph/badge.svg)](https://codecov.io/gh/arthurkiller/mqtgo) [![GoDoc](https://godoc.org/github.com/arthurkiller/mqtgo/packets?status.svg)](https://godoc.org/github.com/arthurkiller/mqtgo/packets)
mqtgo is a pure MQTT protocol packets in go focus on performence

this project was based on eclipse.paho mqtt packet

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

## TODO
* try to redesign the protocol

``` golang
package mqtt

type Unmarshaler interface {
    UnmarshalMqtt() error
}
type Marshaler interface {
    MarshalMqtt()([]byte, error)
}

func Marshal(v ControlPacket) ([]byte, error)
func Unmarshal(data []byte, v ControlPacket) error

type Decoder struct
    func NewDecoder(r io.Reader) *Decoder
    func (d *Decoder) Decode() ControlPacket, error
    func (d *Decoder) DecodeFixedHeader() *FixedHeader, error
    func (d *Decoder) DecodeControlPacket() ControlPacket, error

type Encoder struct
    func NewEncoder(w io.Writer) *Encoder
    func (e *Encoder) Encode(v ControlPacket) error
    func (e *Encoder) EncodeFixedHeader(fh *FixedHeader) error
    func (e *Eecoder) EncodeControlPacket(cp ControlPacket) error
```

``` golang
type XXXPacket struct
    func NewXXXPacket() *XXXPacket
    func NewXXXPacketWithHeader(fh *FixedHeader) *XXXPacket
    func (m XXXPacket) MarshalMqtt() ([]byte, error)
    func (m *XXXPacket) UnmarshalMqtt(data []byte) error

type XXXPacketDecoder struct
    func NewXXXPacketDecoder(r io.Reader) *XXXPacketDecoder
    func (d XXXPacketDecoder)Decode() *NewXXXPacket,error
```

