# healthcheck


`filter` 和 `template` 对象下的字符串类型配置项支持模板解析，当配置项取值为 `temp: ` 开头时，可以使用下面的模板变量：

```
ID                       string
Node                     string
Address                  string
Datacenter               string
TaggedAddresses          map[string]string
NodeMeta                 map[string]string
ServiceID                string
ServiceName              string
ServiceAddress           string
ServiceTags              []string
ServiceMeta              map[string]string
ServicePort              int
ServiceEnableTagOverride bool
CreateIndex              uint64
ModifyIndex              uint64
Tags                     []string
```

HTTP `success_rule` 支持变量：

```
StatusCode int
Status     string
Body       string
```

PING `success_rule` 支持变量：

```
// PacketsRecv is the number of packets received.
PacketsRecv int
// PacketsSent is the number of packets sent.
PacketsSent int
// PacketsRecvDuplicates is the number of duplicate responses there were to a sent packet.
PacketsRecvDuplicates int
// PacketLoss is the percentage of packets lost.
PacketLoss float64
// IPAddr is the address of the host being pinged.
IPAddr *net.IPAddr
// Addr is the string address of the host being pinged.
Addr string
// Rtts is all of the round-trip times sent via this pinger.
Rtts []time.Duration
// MinRtt is the minimum round-trip time sent via this pinger.
MinRtt time.Duration
// MaxRtt is the maximum round-trip time sent via this pinger.
MaxRtt time.Duration
// AvgRtt is the average round-trip time sent via this pinger.
AvgRtt time.Duration
// StdDevRtt is the standard deviation of the round-trip times sent via
// this pinger.
StdDevRtt time.Duration
```