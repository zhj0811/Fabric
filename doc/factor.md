
北京众享比特科技有限公司
------------

| 产品名称Product name    | 密级Confidentiality level
---------------------- | ---------------------
factor_apiserver       |
产品版本Product version |

# 保理 接口文档

Prepared by拟制   |  杨治彬  | Date日期 | 2017-07-13
-----------------|---------|----------| ------------
Reviewed by评审人 |         | Date日期  |
Approved by 批准  |         | Date日期  |


------

# Revision Record 修订记录
Date 日期 		   |  Revision version 修订版本  | Change description 修改描述 | Author 作者
----------------- |----------------------------|----------------------------| ------------
2017/07/13        | Draft-0.1                  | 概要设计                    |  杨治彬

# 1 基本结构定义

## 1.1 FactorInfo

Field          |  Type     | description
---------------|-----------|-------------
createBy       | String    | 创建者
createTime     | UINT64    | 创建时间
sender         | String    | 发送者
receiver       | []String  | 接收者列表
txData         | String    | 业务数据
lastUpdateTime | UINT64    | 最近一次修改时间
lastUpdateBy   | String    | 最近一次修改者
cryptoFlag     | INT       | 加密标识（0:不加密，1:加密）
cryptoAlgorithm| String    | 加密算法类型
docType        | String    | 业务类型
fabricTxId     | String    | Fabric交易id(uuid)
businessNo     | String    | 业务编号（交易编号）
expand1        | String    | 扩展字段1
expand2        | String    | 扩展字段2
DataVersion    | String    | 数据版本


# 2 接口定义

## 2.1 SaveData接口定义
### 2.1.1 请求 

>* header

Field             |  Type     | description
------------------|-----------|-------------
version           | String    | 0.0.1_snapshot
content-Type      | String    | application/json
trackId           | String    | track-0001919248-1280dj840
language          | String    | zh-CN
www-Authenticate  | String    | do8od-084j-l49f-p49g-fk42-9rk4
signatureAlgorithm| String    | RSA/ECB/PKCS1Padding application/json
authentications   | String    | "ROLE_ADMIN","ROLE_AUDITOR"

>* url     /factor/saveData
>* method  POST
>* body    结构如下结构的json

Field          |  Type     | description
---------------|-----------|-------------
payload        | FactorInfo|   数组

### 2.1.2 响应

>* header

Field          |  Type     | description
---------------|-----------|-------------
version        | String    | 版本号
content-Type   | String    | 类型
trackId        | String    | trackId
language       | String    | 语言
responseStatus | Object    | 响应状态

>* body（为空）


### 2.1.3 通知 （取消websocket通信，采用RabbitMQ消息队列）

>* RabbitMQ消息队列 (queueName=“fatorQueue” 众享方为producer,润和方为consumer)

>* Address（amqp://user:password@IP:5672/)

Field          |  Type     | description
---------------|-----------|-------------
header         |   Header  |   结构如下
contents       |   Contents|   结构如下

>* Header

Field          |  Type     | description
---------------|-----------|-------------
contentDef     | ContentDef| 结构如下
ack            | Ack       | 结构如下
responseStatus | ResponseStatus | 响应状态

>* ContentDef

Field          |  Type     | description
---------------|-----------|-------------
contentType    | String    | 类型
trackId        | String    | trackId
language       | String    | 语言

>* Ack

Field          |  Type     | description
---------------|-----------|-------------
level          | String    | 通知级别
callback       | String    | 通知响应

>* Contents

Field          |  Type     | description
---------------|-----------|-------------
$schema        | String    |   该通知payload对应的schema url
payload        | FactorInfo|   消息内容
command        | CommandObj|   command结构

>* CommandObj

Field          |  Type     | description
---------------|-----------|-------------
uri            |  String   |
action         |  String   |
desc           |  String   |


## 2.2 DslQuery接口定义
### 2.2.1 请求 

>* header

Field             |  Type     | description
------------------|-----------|-------------
version           | String    | 0.0.1_snapshot
content-Type      | String    | application/json
trackId           | String    | track-0001919248-1280dj840
language          | String    | zh-CN
www-Authenticate  | String    | do8od-084j-l49f-p49g-fk42-9rk4
signatureAlgorithm| String    | RSA/ECB/PKCS1Padding application/json
authentications   | String    | "ROLE_ADMIN","ROLE_AUDITOR"
page              | Object    | 分页结构如下 

Field          |  Type     | description
---------------|-----------|-------------
currentPage    | INT       | 当前页码 （只有查询用到此字段）
pageSize       | INT       | 页的大小 （只有查询用到此字段）
totalRecords   | INT       | 总记录数

>* url     /factor/dslQuery
>* method  POST
>* body    

Field          |  Type     | description
---------------|-----------|-------------
body           | String    | 查询语句

### 2.2.2 响应

>* header

Field          |  Type     | description
---------------|-----------|-------------
version        | String    | 版本号
content-Type   | String    | 类型
trackId        | String    | trackId
language       | String    | 语言
responseStatus | Object    | 响应状态
page              | Object    | 分页结构如下 

Field          |  Type     | description
---------------|-----------|-------------
currentPage    | INT       | 当前页码 （只有查询用到此字段）
pageSize       | INT       | 页的大小 （只有查询用到此字段）
totalRecords   | INT       | 总记录数

>* body（结构如下）

Field          |  Type     | description
---------------|-----------|-------------
$schema        | String    |   该通知payload对应的schema url
payload        | []FactorInfo|   消息 数组

## 2.3 BlockQuery接口定义
### 2.3.1 请求 

>* header

Field             |  Type     | description
------------------|-----------|-------------
version           | String    | 0.0.1_snapshot
content-Type      | String    | application/json
trackId           | String    | track-0001919248-1280dj840
language          | String    | zh-CN
www-Authenticate  | String    | do8od-084j-l49f-p49g-fk42-9rk4
signatureAlgorithm| String    | RSA/ECB/PKCS1Padding application/json
authentications   | String    | "ROLE_ADMIN","ROLE_AUDITOR"
page              | Object    | 分页结构如下 

Field          |  Type     | description
---------------|-----------|-------------
currentPage    | INT       | 当前页码 （只有查询用到此字段）
pageSize       | INT       | 页的大小 （只有查询用到此字段）
totalRecords   | INT       | 总记录数

>* url     /factor/block/$businessNo (业务编号（交易编号）)
>* method  GET
>* body    为空

### 2.3.2 响应

>* header

Field          |  Type     | description
---------------|-----------|-------------
version        | String    | 版本号
content-Type   | String    | 类型
trackId        | String    | trackId
language       | String    | 语言
responseStatus | Object    | 响应状态
page              | Object    | 分页结构如下 

Field          |  Type     | description
---------------|-----------|-------------
currentPage    | INT       | 当前页码 （只有查询用到此字段）
pageSize       | INT       | 页的大小 （只有查询用到此字段）
totalRecords   | INT       | 总记录数

>* body（结构如下）

Field          |  Type     | description
---------------|-----------|-------------
$schema        | String    |   该通知payload对应的schema url
payload        |   []blockObj |  区块信息结构如下

Field          |  Type      | description
---------------|------------|-------------
txId           | String     | 交易ID
txHash         | String     | 交易hash
blockHash      | String     | 当前区块hash
blockHeight    | UINT       | 当前区块高度
bidbond        | String     | 投标保函唯一编号
bid            | String     | 开立投标保函对应招标需求唯一编号
progress       | String     | 进度描述信息
createBy       | String     | 创建者
createTime     | UINT64     | 创建发布时间
sender         | String     |
receiver       | []String     |
lastUpdateTime | UINT64     | 更新时间
lastUpdateBy   | String     | 更新者
blockData      | String     | 区块信息
remark         | String     |
status         | StateEntity| 结构如下

Field          |  Type     | description
---------------|-----------|-------------
changeEvent    | String    | 变更事件
preState       | String    | 上一状态
currState      | String    | 当前状态

## 2.4 BlockQuery接口定义
### 2.4.1 请求 

>* header

Field             |  Type     | description
------------------|-----------|-------------
version           | String    | 0.0.1_snapshot
content-Type      | String    | application/json
trackId           | String    | track-0001919248-1280dj840
language          | String    | zh-CN
www-Authenticate  | String    | do8od-084j-l49f-p49g-fk42-9rk4
signatureAlgorithm| String    | RSA/ECB/PKCS1Padding application/json
authentications   | String    | "ROLE_ADMIN","ROLE_AUDITOR"
page              | Object    | 分页结构如下 

Field          |  Type     | description
---------------|-----------|-------------
currentPage    | INT       | 当前页码 （只有查询用到此字段）
pageSize       | INT       | 页的大小 （只有查询用到此字段）
totalRecords   | INT       | 总记录数

>* url     /factor/blockQuery/$businessNo (业务编号（交易编号）)
>* method  GET
>* body    为空

### 2.4.2 响应

>* header

Field          |  Type     | description
---------------|-----------|-------------
version        | String    | 版本号
content-Type   | String    | 类型
trackId        | String    | trackId
language       | String    | 语言
responseStatus | Object    | 响应状态
page              | Object    | 分页结构如下 

Field          |  Type     | description
---------------|-----------|-------------
currentPage    | INT       | 当前页码 （只有查询用到此字段）
pageSize       | INT       | 页的大小 （只有查询用到此字段）
totalRecords   | INT       | 总记录数


区块信息结构查询返回body
>* body（结构如下）

Field          |  Type      | description
---------------|------------|-------------
$schema        | String     |   该通知payload对应的schema url:    /schema/blockList.json
payload        | []blockDate |  区块信息结构如下

## blockDate
Field             |  Type     | description
------------------|-----------|-------------  
blockHash            String         当前区块hash
blockHeight          INT            当前区块高度
previousHash         String         前一个区块哈希
events               []Events       生成的事件结构如下

##  Events
Field             |   Type     | description
------------------|------------|-------------
chaincodeId          String        链码ID
txId                 String        交易ID
eventName            String        事件名
payload              FactorInfo    业务数据


## 2.5 keepaliveQuery接口定义
### 2.4.1 请求 

>* header

Field             |  Type     | description
------------------|-----------|-------------
version           | String    | 0.0.1_snapshot
content-Type      | String    | application/json
trackId           | String    | track-0001919248-1280dj840
language          | String    | zh-CN
www-Authenticate  | String    | do8od-084j-l49f-p49g-fk42-9rk4
signatureAlgorithm| String    | RSA/ECB/PKCS1Padding application/json
authentications   | String    | "ROLE_ADMIN","ROLE_AUDITOR"
page              | Object    | 分页结构如下 

Field          |  Type     | description
---------------|-----------|-------------
currentPage    | INT       | 当前页码 （只有查询用到此字段）
pageSize       | INT       | 页的大小 （只有查询用到此字段）
totalRecords   | INT       | 总记录数

>* url     /factor//keepaliveQuery
>* method  HEAD
>* body    为空

### 2.4.2 响应
>* header

Field          |  Type     | description
---------------|-----------|-------------
version        | String    | 版本号
content-Type   | String    | 类型
trackId        | String    | trackId
language       | String    | 语言
responseStatus | Object    | 响应状态
page              | Object    | 分页结构如下 

Field          |  Type     | description
---------------|-----------|-------------
currentPage    | INT       | 当前页码 （只有查询用到此字段）
pageSize       | INT       | 页的大小 （只有查询用到此字段）
totalRecords   | INT       | 总记录数

>* body    为空