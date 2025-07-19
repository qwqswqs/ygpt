# 接口文档
管理端：admin / 123456
租户端：yu233 /123456

## 发布管理相关接口

### 算力产品管理接口
#### 添加算力产品
方法类型：POST

路径：/product/computing/addComputingInfo

请求：
```json5
{
    "name":"A100GPU",
    "type":3, // “1.元素”、“2.容器”、“3.云主机”、“4.裸金属”或“5.存储”
    "productDesc":"{\"name\":\"服务1\",\"age\":18,\"sex\":\"男\"}",
    "status": 1, // 状态(1.正常，2.下架)
    "payWay":1 // 支付方式 1微信 2支付宝
}
```
响应：
```json5
{
  "code": 0,
  "data": {
    "ID": 1,
    "CreatedAt": "0001-01-01T00:00:00Z",
    "UpdatedAt": "0001-01-01T00:00:00Z"
  },
  "msg": "添加成功"
}
```
#### 修改算力产品
方法类型：PUT

路径：/product/computing/updateComputingInfo

请求：
```json5
{
  "id":1,
  "name":"A100GPU",
  "type":3,
  "productDesc":"{\"name\":\"服务1\",\"age\":18,\"sex\":\"男\"}",
  "status": 1,
  "payWay":1
}
```
响应：
```json5
{
  "code": 0,
  "data": {
    "updateAt": 1734056664
  },
  "msg": "修改成功"
}
```
#### 获取算力产品列表
方法类型：POST

路径：/product/computing/queryComputingInfo

请求：
```json5
{
  "page":1,
  "pageSize":10
}
```
响应：
```json5
{
  "code": 0,
  "data": {
    "list": [
      {
        "ID": 1,
        "CreatedAt": "2024-12-13T10:13:47.394+08:00",
        "UpdatedAt": "2024-12-13T10:24:24.122+08:00",
        "name": "A100GPU",
        "type": 3,
        "productDesc": "{\"age\": 18, \"sex\": \"男\", \"name\": \"服务1\"}",
        "status": 1,
        "payWay": 1
      },
      {
        "ID": 2,
        "CreatedAt": "2024-12-13T10:25:04.763+08:00",
        "UpdatedAt": "2024-12-13T10:25:04.763+08:00",
        "name": "A100GPU",
        "type": 3,
        "productDesc": "{\"age\": 18, \"sex\": \"男\", \"name\": \"服务1\"}",
        "status": 1,
        "payWay": 1
      }
    ],
    "total": 2,
    "page": 1,
    "pageSize": 10
  },
  "msg": "获取成功"
}
```
#### 删除算力产品
方法类型：Delete

路径：/product/computing/deleteComputingInfo

请求：
```json5
{
  "id":1,
}
```
响应：
```json5
{
  "code": 0,
  "data": {
    "deleteAt": 1734057195
  },
  "msg": "删除成功"
}
```
### 模型管理和数据集管理接口
### 供给和需求

#### 添加供需信息
方法类型：POST

路径：/product/computing/addComputingInfo

请求：
```json5
{
  "title":"100张A100待出售",
  "type":1, // 1.供给 2.需求
  "supplyType":1, // 1.算力”、“2.机柜”、“3.带宽”
  "description":"100张", 
  "specDesc": "{\"name\":\"服务1\",\"age\":18,\"sex\":\"男\"}",
  "status":1
}
```
响应：
```json5
{
  "code": 0,
  "data": {
    "ID": 1,
    "CreatedAt": "0001-01-01T00:00:00Z",
    "UpdatedAt": "0001-01-01T00:00:00Z"
  },
  "msg": "添加成功"
}
```
#### 修改供需信息
方法类型：PUT

路径：/product/supply/updateSupplyInfo

请求：
```json5
{
  "id":1,
  "title": "100张A100待出售",
  "type": 1,
  "supplyType": 1,
  "description": "100张",
  "specDesc": "{\"name\":\"服务1\",\"age\":18,\"sex\":\"男\"}",
  "status": 1
}
```
响应：
```json5
{
  "code": 0,
  "data": {
    "updateAt": 1734056664
  },
  "msg": "修改成功"
}
```
#### 获取供需列表
方法类型：POST

路径：/product/supply/querySupplyInfo

请求：
```json5
{
  "pageIndex":1,
  "pageSize":10
}
```
响应：
```json5
{
  "code": 0,
  "data": {
    "list": [
      {
        "ID": 1,
        "CreatedAt": "2024-12-13T10:56:19.302+08:00",
        "UpdatedAt": "2024-12-13T11:00:49.927+08:00",
        "title": "",
        "type": 1,
        "supplyType": 1,
        "description": "100张",
        "specDesc": "{\"age\": 18, \"sex\": \"男\", \"name\": \"服务1\"}",
        "status": 1
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  },
  "msg": "获取成功"
}
```
#### 删除供需信息
方法类型：Delete

路径：/product/supply/deleteSupplyInfo

请求：
```json5
{
  "id":1,
}
```
响应：
```json5
{
  "code": 0,
  "data": {
    "deleteAt": 1734057195
  },
  "msg": "删除成功"
}
```


