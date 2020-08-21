## 认证

采用[Basic Auth](https://swagger.io/docs/specification/authentication/basic-authentication/).

在使用本项目任何API时都需要携带 **Authorization Header** （**Basic sid:password**）。


## API
### 查询余额

| Method | URL                   | Header        |
| ------ | --------------------- | ------------- |
| GET    | /api/card/v1/balance/ | Authorization |

**RESPONSE Data:**
```json
{
    "code": 0,
    "message": "OK",
    "data": {
        "balance": 104.67,
        "status": "在用"
    }
}
```

### 查询流水

| Method | URL                  | Header        |
| ------ | -------------------- | ------------- |
| GET    | /api/card/v2/account | Authorization |

**URL Params:**
```
    limit: string
    page: string // 页码，默认为1
    start : string  // 开始日期，格式：2018-01-01，默认为当天
    end : string // 结束日期
```

**RESPONSE Data:**
```json
{
    "code": 0,
    "message": "OK",
    "data": {
        "count": 1,
        "list": [
            {
                "dealName": "消费",
                "orgName": "华中师范大学/后勤集团/商贸中心/超市/学子超市",
                "transMoney": 11.7,                   // 交易金额
                "dealDate": "2020-01-16 13:03:35",
                "outMoney": 112.17                    // 剩余余额
            }
        ]
    }
}
```
