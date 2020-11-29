module github.com/qydysky/other/proxylist

go 1.14

require (
	github.com/golang/protobuf v1.4.3
	github.com/klauspost/compress v1.10.11 // indirect
	github.com/qydysky/part v0.0.0-20200908005332-f8509bab9fa5
	github.com/shirou/gopsutil v2.20.8+incompatible // indirect
	v2ray.com/core v4.15.0+incompatible
)

replace v2ray.com/core => ../../qqqaadd/v2ray-core
