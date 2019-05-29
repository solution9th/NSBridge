# API

<!-- TOC -->

- [API](#api)
    - [HTTP API 鉴权](#http-api-鉴权)
    - [grpc鉴权](#grpc鉴权)
        - [前提](#前提)
        - [域名接口鉴权](#域名接口鉴权)
        - [解析记录接口鉴权](#解析记录接口鉴权)
    - [域名操作](#域名操作)
        - [域名列表](#域名列表)
            - [HTTP 调用](#http-调用)
            - [GRPC 调用](#grpc-调用)
        - [添加域名 [域名key]](#添加域名-域名key)
            - [HTTP 调用](#http-调用-1)
            - [GRPC 调用](#grpc-调用-1)
        - [删除域名 [域名key]](#删除域名-域名key)
            - [HTTP 调用](#http-调用-2)
            - [GRPC 调用](#grpc-调用-2)
        - [检查域名托管状态 [记录key]](#检查域名托管状态-记录key)
            - [HTTP 调用](#http-调用-3)
            - [GRPC 调用](#grpc-调用-3)
    - [解析记录 [记录key]](#解析记录-记录key)
        - [获取解析记录列表](#获取解析记录列表)
            - [HTTP 调用](#http-调用-4)
            - [GRPC 调用](#grpc-调用-4)
        - [获取解析记录详情](#获取解析记录详情)
            - [HTTP 调用](#http-调用-5)
            - [GRPC 调用](#grpc-调用-5)
        - [新增解析记录](#新增解析记录)
            - [HTTP 调用](#http-调用-6)
            - [GRPC 调用](#grpc-调用-6)
        - [更新解析记录](#更新解析记录)
            - [HTTP 请求](#http-请求)
            - [GRPC 调用](#grpc-调用-7)
        - [删除解析记录](#删除解析记录)
            - [HTTP 调用](#http-调用-7)
            - [GRPC 调用](#grpc-调用-8)
        - [暂停、启动解析记录](#暂停启动解析记录)
            - [HTTP 调用](#http-调用-8)
            - [GRPC 调用](#grpc-调用-9)
    - [获取线路信息](#获取线路信息)
        - [获取记录类型列表](#获取记录类型列表)
            - [HTTP 调用](#http-调用-9)
            - [GRPC 调用](#grpc-调用-10)
        - [获取所有线路对应关系](#获取所有线路对应关系)
            - [HTTP调用](#http调用)
            - [GRPC 调用](#grpc-调用-11)
        - [获取洲线路](#获取洲线路)
        - [获取ISP线路](#获取isp线路)
        - [获取各个国家线路](#获取各个国家线路)
        - [获取国内各个省线路](#获取国内各个省线路)
        - [获取国外首都线路](#获取国外首都线路)
        - [获取省与isp组合线路](#获取省与isp组合线路)

<!-- /TOC -->

## HTTP API 鉴权

分为两种 `key-secret` 组合。

- 一种为操作域名的 `ks`，从网站后台获取，可以`添加域名`、`删除域名`，简称为 `域名key`。

- 一种为操作域名解析记录的 `ks`，添加域名后会出现在返回值中，可以用来`添加记录`、`删除记录`等与**该域名相关**的所有解析记录操作，简称为 `记录 key`。


在调用 API 时，不同接口需要不同种的 key 值，但是调用方法是一致的，通过 HTTP 头传递所需的参数，**5分钟**之内请求有效

|HTTP 头字段|说明|
|:--|:--|
|X-API-KEY|对应 Key|
|X-API-TIMESTAMP|当前时间戳|
|X-API-HMAC|通过MD5获得的哈希|

- API-HMAC 算法：`md5(请求方法+X-API-KEY+请求URL+请求参数体+X-API-TIMESTAMP+SECRET-KEY)`
- 接口请求URL 是接口请求的 URL 字符串，例如获取指定 host_id 的解析记录列表的 URL 是 `/api/ping?host_id=XXX&offset=0&row_num=30`
- 请求参数体是指 post 和 put 请求时的 JSON 的请求参数体，例如添加域名的 JSON 请求 参数体是`{"domain":"XXXX.net"}`，建议是 json 保持一行，不要有缩进。
- 请求方法大写 `GET`

以下是 `curl` 示例

```shell
$ curl -X GET \
  'http://127.0.0.1:8080/api/v1/domains?start=0&count=19' \
  -H 'X-API-TIMESTAMP: 1553569617' \
  -H 'X-API-HMAC: D0EC8715CA84A40D12612D4576B418EC' \
  -H 'X-API-KEY: dodomain'
```

key 不仅用来鉴权，还用来指定具体操作的哪个域名，所以操作解析记录时都需要该域名对应的 key

## grpc鉴权

### 前提

使用`TLS`方式连接grpc服务端.
证书 在 keys/grpc目录下

下面的域名鉴权和解析记录的鉴权参数和其他业务参数是在同一级的, 参照proto文件即可.

### 域名接口鉴权

|鉴权字段|类型|说明|
|:--|:--|:--|
|api_key|string|用户中心获得的 API-KEY|

### 解析记录接口鉴权

|鉴权字段|类型|说明|
|:--|:--|:--|
|record_key|string|通过添加域名获取的操作该域名的key|

## 域名操作

### 域名列表

#### HTTP 调用

> GET  /api/v1/domains

分页 `?start=0&count=10`  count大于100时会只返回50个

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": {
        "total": 8,
        "list": [
          {
            "id": 8 , //域名id
            "fone_domain_id": 123,
            "domain_key": "mougeren", // 创建这个域名的key
            "domain": "dev.com",
            "name_server":["ns1.newio.cc","ns2.newio.cc"],
            "soa_email": "uuu@newio.cc", // soa_email
            "remark": "这是一个测试", // 备注
            "is_take_over": 1, // 接管状态
            "is_open_key": 1, // ks 是否启用
            "record_key": "02a6da7989c3e84edd6916cf86d10afe",
            "record_secret": "Dep6QmiLsHwFmofSbmkH",
            "create_at": "2019-03-22T11:07:48Z",
            "update_at": "2019-03-22T11:07:48Z",
          }
        ]
    }
}
```

#### GRPC 调用

`rpc DomainsList (RequestDomainsList) returns (ResponseDomainList)`

### 获取授权Key下的域名列表

#### HTTP 调用

> GET  /api/v1/domains/own

分页 `?start=0&count=10`  count大于100时会只返回50个

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": {
        "total": 8,
        "list": [
          {
            "id": 8 , //域名id
            "fone_domain_id": 123,
            "domain_key": "mougeren", // 创建这个域名的key
            "domain": "dev.com",
            "name_server":["ns1.newio.cc","ns2.newio.cc"],
            "soa_email": "uuu@newio.cc", // soa_email
            "remark": "这是一个测试", // 备注
            "is_take_over": 1, // 接管状态
            "is_open_key": 1, // ks 是否启用
            "record_key": "02a6da7989c3e84edd6916cf86d10afe",
            "record_secret": "Dep6QmiLsHwFmofSbmkH",
            "create_at": "2019-03-22T11:07:48Z",
            "update_at": "2019-03-22T11:07:48Z",
          }
        ]
    }
}
```

#### GRPC 调用

`rpc OwnDomainsList (RequestOwnDomainsList) returns (ResponseDomainList)`

### 添加域名 [域名key]

#### HTTP 调用

> POST /api/v1/domains

[HTTP API 鉴权](#http-api-鉴权)

```json
{
  "domain":"dev.com",
  "remark":"" // 备注
}
```

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": {
        "id": 9,
        "fone_domain_id": 28,
        "domain_key": "123123", // 添加这个域名的 key
        "domain": "hh.com",
        "name_server": [
            "ns1.newio.cc",
            "ns2.newio.cc"
        ],
        "soa_email": "",
        "remark": "测试",  // 备注
        "is_take_over": 0, // 是否接管
        "is_open_key": 1, // key，secret 是否可以启用
        "record_key": "cmVjb3Jk2cab9d2edfd19538",
        "record_secret": "wT3gFaOxhj1nylxJQ1s8",
        "create_at": "2019-03-22T17:11:55Z",
        "update_at": "2019-03-22T17:11:56Z"
    }
}
```

#### GRPC 调用

`rpc DomainCreate (RequestDomainCreate) returns (ResponseDomainCreate)`

### 删除域名 [域名key]

只能用创建域名时的 key 删除此域名

#### HTTP 调用

> DELETE /api/v1/domain/:id

[HTTP API 鉴权](#http-api-鉴权)

```json
{
    "err_code": 0,
    "err_msg": "success"
}
```

#### GRPC 调用

`rpc DomainDelete (RequestDomainDelete) returns (ResponseDomainDelete)`

### 检查域名托管状态 [记录key]

#### HTTP 调用

由于是通过根节点递归查询，所以即使显示已托管也有可以因为 `114.114.114.114` 的 TTL，导致不能域名不能正确访问

> GET /api/v1/status

[HTTP API 鉴权](#http-api-鉴权)

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": {
        "is_take_over": 1 // 0 为未托管， 1 为托管
    }
}
```

#### GRPC 调用

`rpc DomainStatus (RequestDomainStatus) returns (ResponseDomainStatus)`


## 解析记录 [记录key]

[HTTP API 鉴权](#http-api-鉴权)

### 查看record_key的域名 [记录key]

#### HTTP 调用

> GET /api/v1/record/domain

[HTTP API 鉴权](#http-api-鉴权)

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": {
        "id": 9,
        "fone_domain_id": 28,
        "domain_key": "123123", // 添加这个域名的 key
        "domain": "hh.com",
        "name_server": [
            "ns1.newio.cc",
            "ns2.newio.cc"
        ],
        "soa_email": "",
        "remark": "测试",  // 备注
        "is_take_over": 0, // 是否接管
        "is_open_key": 1, // key，secret 是否可以启用
        "record_key": "cmVjb3Jk2cab9d2edfd19538",
        "record_secret": "wT3gFaOxhj1nylxJQ1s8",
        "create_at": "2019-03-22T17:11:55Z",
        "update_at": "2019-03-22T17:11:56Z"
    }
}
```

#### GRPC 调用

`rpc RecordDomainOfRK (RequestDomainOfRK) returns (ResponseDomainOfRK)`

### 获取解析记录列表

#### HTTP 调用

> GET /api/v1/records

分页 ?start=0&count=10

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data":{
      "total": 8,
      "list": [
        {
            "id": 13,  // 记录 id
            "domain_id": 6, // 域名 id
            "fone_domain_id": 25, // fone中的id
            "fone_record_id": 48, // fone中记录的id
            "sub_domain": "mafeng", // 主机值
            "record_type": "TXT", // 记录类型
            "value": "eznsm.com", // 记录值
            "line_id": 2, // 线路id
            "ttl": 3, // ttl
            "unit": "hour", // ttl 单位
            "priority": 3, // 优先级
            "disable": 0, // 是否关闭
            "create_at": "2019-03-23T15:21:22Z",
            "update_at": "2019-03-23T15:21:22Z"
        },
      ]
    }
}
```

#### GRPC 调用

`rpc RecordList (RequestRecordList) returns (ResponseRecordList)`

### 获取解析记录详情

#### HTTP 调用

> GET /api/v1/record/:recordid

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": {
        "id": 13,
        "domain_id": 6,
        "fone_domain_id": 25,
        "fone_record_id": 48,
        "sub_domain": "mn1234g",
        "record_type": "MX", // 记录类型
        "value": "mf.m.com",  // 记录值
        "line_id": 2, // 线路 id
        "ttl": 6,
        "unit": "min",
        "priority": 3, // 优先级
        "disable": 0, // 0 为正常，1为暂停
        "create_at": "2019-03-23T15:21:22Z",
        "update_at": "2019-03-23T15:50:29Z"
    }
}
```

#### GRPC 调用

`rpc RecordInfo (RequestRecordInfo) returns (ResponseRecordInfo)`

### 新增解析记录

#### HTTP 调用

> POST /api/v1/records

|字段|类型|说明|
|:--|:--|:--|
|sub_domain|string|主机值，长度小于63|
|record_type|string|记录类型[A,CNAME,MX,TXT,NS]|
|value|string|记录，根据type进行区分|
|line_id|int|线路id(如果没有特殊需求, 请使用默认线路id)|
|priority|int|优先级，当type为MX时选填，默认值为5|
|ttl|int|生存时间，配合unit参数使用|
|unit|string|ttl单位[sec,min,hour,day]|

```json
{
  "sub_domain":"aaa",
  "record_type":"A",
  "value":"1.1.1.1",
  "line_id":0, // 线路 id
  "priority": 1, // 优先级,范围 1-100.当记录类型是 MX/AX/CNAMEX 时有效并且必选
  "ttl": 5,
  "unit": "min" // ttl 单位
}
```

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": {
        "id": 46,
        "domain_id": 6,
        "fone_domain_id": 25,
        "fone_record_id": 119,
        "sub_domain": "mao",
        "record_type": "TXT",
        "value": "_9999_",
        "line_id": 0,
        "ttl": 5,
        "unit": "min",
        "priority": 0,
        "disable": 0,
        "create_at": "2019-04-09T10:40:14Z",
        "update_at": "2019-04-09T10:40:15Z"
    }
}
```

#### GRPC 调用

`rpc RecordCreate (RequestRecordCreate) returns (ResponseRecordCreate)`

### 更新解析记录

#### HTTP 请求

> PUT /api/v1/record/:recordid

```json
{
  "sub_domain":"nihao",
  "record_type": "TXT",
  "value":"_9999_",
  "line_id":0,
  "ttl": 5,
  "unit":"min",
  "priority": 1 // 优先级,范围 1-100.当记录类型是 MX 时有效，默认是5
}
```

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": "success"
}
```

#### GRPC 调用

`rpc RecordUpdate (RequestRecordUpdate) returns (ResponseRecordUpdate)`

### 删除解析记录

#### HTTP 调用

> DELETE /api/v1/record/:recordid

```json
{
    "err_code": 0,
    "err_msg": "success"
}
```

#### GRPC 调用

`rpc RecordDelete (RequestRecordDelete) returns (ResponseRecordDelete)`

### 暂停、启动解析记录

#### HTTP 调用

> PATCH /api/v1/record/:recordid

```json
{
  "disable": true
}
```

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": "success"
}
```

#### GRPC 调用

`rpc RecordDisable (RequestRecordDisable) returns (ResponseRecordDisable)`

## 获取线路信息


### 获取记录类型列表

#### HTTP 调用

获取可以添加的记录类型

> GET /api/v1/types

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": [
        "A",
        "AAAA",
        "MX",
        "CNAME",
        "TXT",
        "NS"
    ]
}
```

#### GRPC 调用

`rpc Types (google.protobuf.Empty) returns (google.api.HttpBody)`

### 获取所有线路对应关系

#### HTTP调用

> GET /api/v1/lines/all

```json
{
  "err_code": 0,
  "err_msg": 0,
  "data": [
      {
        "id": 0,
        "name": "default"
      },
      {
        "id": 1,
        "name": "Africa"
      }
    ]
}
```

#### GRPC 调用

`rpc LineAll (google.protobuf.Empty) returns (google.api.HttpBody)`

### 获取洲线路

> GET /api/v1/lines/continental

> rpc LineContinental (google.protobuf.Empty) returns (google.api.HttpBody)

同上

### 获取ISP线路

> GET /api/v1/lines/isp

> rpc LineISP (google.protobuf.Empty) returns (google.api.HttpBody)

同上


### 获取各个国家线路

> GET /api/v1/lines/country

> rpc LineCountry (google.protobuf.Empty) returns (google.api.HttpBody)

同上

### 获取国内各个省线路

> GET /api/v1/lines/province

> rpc LineProvince (google.protobuf.Empty) returns (google.api.HttpBody)

同上

### 获取国外首都线路

> GET /api/v1/lines/outcity

> rpc LineOutCity (google.protobuf.Empty) returns (google.api.HttpBody)

同上

### 获取省与isp组合线路

> GET /api/v1/lines/cityisp

> rpc LineCityIsp (google.protobuf.Empty) returns (google.api.HttpBody)

同上