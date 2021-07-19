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

`success_rule` 支持变量：

```
StatusCode int
Status     string
Body       string
```