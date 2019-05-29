#WebApi

<!-- TOC -->

- [WebApi](#webapi)
- [对接i1](#对接i1)
    - [SAML](#saml)
- [接口文档](#接口文档)
    - [鉴权](#鉴权)
    - [用户登录](#用户登录)
        - [URL](#url)
        - [请求方式](#请求方式)
        - [请求参数](#请求参数)
        - [请求示例](#请求示例)
        - [返回结果](#返回结果)
    - [用户授权回调地址](#用户授权回调地址)
        - [URL](#url-1)
        - [请求方式](#请求方式-1)
        - [请求参数](#请求参数-1)
        - [请求示例](#请求示例-1)
        - [返回结果](#返回结果-1)
    - [授权信息](#授权信息)
        - [获取+搜索 授权信息](#获取搜索-授权信息)
            - [参数](#参数)
            - [返回值](#返回值)
        - [新增授权信息](#新增授权信息)
            - [参数](#参数-1)
            - [返回值](#返回值-1)
        - [更新授权信息备注](#更新授权信息备注)
            - [参数](#参数-2)
            - [返回值](#返回值-2)
        - [删除授权信息](#删除授权信息)
            - [参数](#参数-3)
            - [返回值](#返回值-3)
        - [启用, 禁用授权信息](#启用-禁用授权信息)
            - [参数](#参数-4)
            - [返回值](#返回值-4)
    - [获取+搜索 域名列表](#获取搜索-域名列表)
            - [参数](#参数-5)
            - [返回值](#返回值-5)
    - [获取+搜索 解析记录列表](#获取搜索-解析记录列表)
            - [参数](#参数-6)
            - [返回值](#返回值-6)

<!-- /TOC -->

# 对接i1

web的登录基于i1的单点登录系统, 所以用户的管理是在i1的那边管理的, 给用户分配权限之后, 用户获得登录该应用的权限.

## SAML

新用户访问时, 单点登录流程如下:

- 判断用户是否登录

  - 根据用户的`cookie`(`user_id` + `token`)判断是否登录 

- 如果登录, 则跳转用户应用界面, 否则跳转`/saml/login`路由, 请求i1的saml认证

- 用户在i1登录和授权

- i1在用户点击授权后, 路由到ns_bridge的`/saml/acs`地址

- 在response中获取用户的信息

# 接口文档 

前端开发关于单点登录的部分不用看, 从`授权信息`看起即可

## 鉴权

只有已登录的用户才有权限访问接口, 权限相关的token和userid保存在cookie中, 过期时间是一个小时

当调用`/web`下接口返回
```json
{
"err_code": 10303,
"err_msg": "user is not login ",
"data": ""
}
```
时, 请使用重定向`c.Redirect(302, "/saml/login?RelayState=" + path)`将网页跳转到登录页面.

`path`: 为用户的点击授权登录后想让他跳转的地址, 一般为用户当前正处于的url, `格式`比如`/index.html`


## 用户登录

用户登录NDS应用, 使用SAML协议, 如不清楚, 可去先熟悉下SAML协议

### URL
$host/saml/login/

### 请求方式
**Get**

### 请求参数

> id_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwOi8vbWYuc2RkZXpuc20uY29tL2xvZ2luIiwiYXZhdGFyIjoiIiwiZGVwYXJ0bWVudCI6IuW-kOWtkOmJtOa1i-ivlee7hOivt-S4jeimgeS9v-eUqF9leG1haWzmtYvor5Xpg6jpl6gxLTEiLCJkb21haW4iOiJvcmcuaTEuZGV2ZWxlbnYuY29tIiwiZW1haWwiOiJtYWZlbmdAaHVhbnFpdS5jb20iLCJleHAiOjE1NTI5OTUwNjIsImZpcnN0bmFtZSI6IuWzsCIsImlzcyI6IkkxIiwibGFzdG5hbWUiOiLpqawiLCJtb2JpbGUiOiIiLCJuYW1lIjoi6ams5bOwIiwibmJmIjoxNTUyOTA4NjYyLCJzdWIiOiJJMSBMb2dpbiIsInRlbGVwaG9uZSI6Ii0tIiwidG9rZW4iOiJ4MHBvamtlNDJjM3YiLCJ1aWQiOiJtYWZlbmcifQ.JrwRjCmtnMEvCcSUnX2FGh8cp8UpJ4Z7Iu7rdgVoG9s

### 返回结果
 302 -> /sso/ack

## 用户授权回调地址
i1单点登录之后的回调地址

### URL
$host/saml/acs


## 授权信息

### 获取+搜索 授权信息

> GET /web/auth

#### 参数
|字段|类型|是否必须|解释|
|:--|:--|:--|:--|
|domain_key|string|no|domainKey模糊值|
|disable|int|no|是否禁用(1禁用,0启用)|
|offset|int|no|偏移量|
|count|int|no|数量|
#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": [
        {
            "id": 7,
            "domain_key": "e8cbc59d-94cc-41d2-8a0d-e5957ead7273", // 公钥
            "domain_secret": "44nA5q1KYKjtWNZ8HppsXYhwOkVPkHTj",  // 私钥
            "remark": "mafeng", // 备注
            "disable": 0,  // 是否禁用
            "create_at": "2019-03-26T15:04:11Z",
            "update_at": "2019-03-26T15:04:11Z"
        },
        ...
    ]
}
```

### 新增授权信息

> POST /web/auth

#### 参数
```json
{"remark":"mafeng"}
```
|字段|类型|是否必须|解释|
|:--|:--|:--|:--|
|remark|string|yes|备注|

#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": null
}
```

### 更新授权信息的备注

> PUT /web/auth

#### 参数
```json
{"remark":"mafeng"}
```
|字段|类型|是否必须|解释|
|:--|:--|:--|:--|
|remark|string|yes|备注|

#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": null
}
```

### 删除授权信息

> DELETE /web/auth/:auth_id

#### 参数
:auth_id

#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": null
}
```

### 启用, 禁用授权信息

> PUT /web/auth/:auth_id

#### 参数
:auth_id

#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": null
}
```

## 获取+搜索 域名列表

> GET /web/domain

#### 参数
|字段|类型|是否必须|解释|
|:--|:--|:--|:--|
|domain|string|no|域名模糊值|
|offset|int|no|偏移量|
|count|int|no|数量|

#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": [
        {
            "id": 9,
            "fone_domain_id": 28,
            "domain_key": "",
            "domain": "dev.newio.cc",
            "name_server": [
                "ns1.newio.cc",
                "ns2.newio.cc"
            ],
            "soa_email": "",
            "remark": "测试，🍰",
            "is_take_over": 1,
            "is_open_key": 1,
            "record_key": "cmVjb3Jk2cab9d2edfd19538",
            "record_secret": "wT3gFaOxhj1nylxJQ1s8",
            "create_at": "2019-03-22T17:11:55Z",
            "update_at": "2019-03-23T14:00:26Z"
        }
    ]
}
```

## 获取+搜索 解析记录列表

> GET /web/record/:domain_id

#### 参数
|字段|类型|是否必须|解释|
|:--|:--|:--|:--|
|domain_id|string|yes|domain_id|
|sub_domain|string|no|主机记录模糊值|
|record_type|string|no|记录类型模糊值|
|value|string|no|记录值模糊值|
|offset|int|no|偏移量|
|count|int|no|数量|
#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": [
        {
            "id": 10,
            "domain_id": 6,  
            "fone_domain_id": 25,
            "fone_record_id": 45,
            "sub_domain": "mafeng",  // 主机记录
            "record_type": "A",  // 记录类型 
            "value": "mf.sddeznsm.com", // 记录值
            "line_id": 2,  // 线路
            "line_name": "Asia", // 线路名
            "ttl": 3,  // ttl
            "unit": "hour", // 时间单位
            "priority": 3,  // 权重
            "disable": 0,  // 是否禁用
            "create_at": "2019-03-23T10:24:07Z",
            "update_at": "2019-03-23T10:24:07Z"
        },
        ...
    ]
}
```

## 获取某个域名的解析记录类型

> GET /web/record/:domain_id/types

#### 参数
:domain_id

#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success ",
    "data": [
        {
            "record_type": "A"
        },
        {
            "record_type": "MX"
        },
        {
            "record_type": "CNAME"
        },
        {
            "record_type": "TXT"
        },
        {
            "record_type": "NS"
        }
    ]
}
```

## 用户退出登录

> GET /web/logout

#### 参数
无
#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success",
    "data": "success"
}
```

## 获取当前登录用户信息

> GET /web/user

#### 参数
无
#### 返回值

```json
{
    "err_code": 0,
    "err_msg": "success ",
    "data": {
        "user_id": "mafeng",
        "user_name": "马XX"
    }
}
```