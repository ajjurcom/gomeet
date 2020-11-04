## gomeet-backend

FoShan University Conference Room Management System - backend

## API 文档

http://localhost:8203/swagger/index.html

`http://localhost:8203` 为服务网址头，可能会随服务迁移相应改变，注意更换地址

## 技术选型

| 技术          | 说明          | 官网                                     |
| ------------- | ------------- | ---------------------------------------- |
| gin           | Web框架       | https://github.com/gin-gonic/gin         |
| go-ini        | 配置包        | https://github.com/go-ini/ini            |
| mysql         | 关系数据库    | https://github.com/go-sql-driver/mysql   |
| zap           | 日志打印      | https://github.com/uber-go/zap           |
| jwt-go        | 权限验证器    | https://github.com/dgrijalva/jwt-go      |
| swag          | API文档生成器 | https://github.com/swaggo/swag           |
| elasticsearch | 搜索引擎      | https://www.elastic.co/cn/elasticsearch/ |

## 调试

```shell
$ git clone -b back-end git@github.com:mittacy/golog.git
$ cd back-end
# 修改 back-end/config/my.ini 相关配置信息
$ go mod download
$ go run main.go
# 服务将运行在相应的端口，例如：http://localhost:8203/api/v1
```

## Project structrue

```
golog
|-- front-end
|-- back-end
	|-- common	工具类及通用代码
	|-- config	配置包
	|-- router	路由
	|-- controller	API控制器
	|-- repository	数据库操作
	|-- model	模型
	|-- database 数据库连接
	|-- logger	日志
	|-- e	错误码和提示信息
	|-- docs API文档
	main.go	主程序
```

项目文档

```shell
$ godoc -http=:6060
```

打开浏览器：[http://localhost:6060/pkg/com/mittacy/gomeet/](http://localhost:6060/pkg/com/mittacy/meet/) 可查看项目doc文档

## 表结构



