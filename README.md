# Todo List

一个基于 Gin + Gorm 的待办事项（Todo List）RESTful API 项目，用于学习与实践 Web 后端分层架构、JWT 鉴权、配置化与容器化部署。

> 本项目在 B 站 UP 主「小生凡一」的 TodoList 项目基础上做了改进与扩展，原项目地址：https://github.com/CocaineCong/TodoList/tree/main

---

## 特性

- 基于 Gin 的 RESTful API 路由与中间件
- 基于 Gorm 的 MySQL ORM（包含自动迁移与软删除）
- JWT 登录鉴权与中间件保护
- CORS 支持
- 分层架构：Router / API / Service / DAO / Model / Types / Middleware / Utils
- Viper 配置管理（支持多环境：本地与 prod）
- Docker 与 docker-compose 部署示例

## 技术栈

- 语言：Go（go.mod 指定 `go 1.25`）
- Web：Gin v1.11
- ORM：Gorm v1.31 + gorm.io/driver/mysql
- 配置：Viper v1.21
- 鉴权：golang-jwt/jwt v5
- 其他：CORS、x/crypto

## 目录结构

```
cmd/server/main.go           # 程序入口
config/                      # 配置文件目录（由 Viper 读取）
internal/
	api/v1/                    # HTTP Handler（Controller）
	middleware/                # 中间件（JWT、CORS）
	repository/db/             # 数据访问
		init.go                  # 数据库初始化（MySQL、Gorm 设置与自动迁移）
		dao/                     # DAO 接口与实现（task、user）
		model/                   # Gorm Model 定义（task、user）
	routes/router.go           # 路由注册
	service/                   # 业务服务层（task、user）
pkg/
	conf/                      # 配置加载（Viper）
	ctl/                       # 统一响应包装
	e/                         # 错误码与消息
	utils/                     # 工具（JWT、随机名）
types/                       # 请求/响应 DTO
Dockerfile                   # Docker 多阶段构建
deploy/docker-compose.yml    # MySQL、Redis、应用容器编排
```

## 配置说明

项目使用 Viper 从工作目录下的 `config/` 读取配置：

- 环境变量 `APP_ENV=prod` 时读取 `config.prod.yaml`，否则读取 `config.yaml`
- 配置结构位于 `pkg/conf/config.go`：
	- `service.app_mode`：`debug` 或 `release`
	- `service.http_port`：服务端口（当前主程序默认监听 `:8080`）
	- `database`：MySQL 连接信息（用户、密码、主机、库名、字符集、时区等）
	- `redis`：Redis 连接信息

JWT 密钥通过环境变量 `JWT_SECRET` 提供（必填，生产环境必须设置为高强度随机串）。

## 快速开始

### 方式一：本地运行

前置条件：Go 1.25、MySQL、Redis。

1) 准备 MySQL 和 Redis（可使用 docker-compose 仅启动依赖）：

```bash
# 仅启动 mysql 与 redis（在项目根目录执行）
docker compose -f deploy/docker-compose.yml up -d mysql redis
```

2) 修改 `config/config.yaml` 中的数据库连接信息，确保能连上你的 MySQL。

3) 设置 JWT 密钥并运行：

```bash
export JWT_SECRET="your-strong-secret"
GO111MODULE=on go run ./cmd/server/main.go
# 或者使用 Makefile：
# make run
```

服务默认监听 `http://localhost:8080`。

### 方式二：容器化运行（包含应用）

使用预置的镜像与编排文件一键启动：

```bash
docker compose -f deploy/docker-compose.yml up -d
```

- 将启动 `mysql`、`redis` 与 `todolist-app` 三个容器
- 其中应用容器镜像为 `hzhenj/todolist-server:latest`
- 如需使用本地构建镜像，可用下述命令替换 compose 文件中的镜像名

### 可选：本地构建镜像

```bash
# 在项目根目录构建镜像
docker build -t todolist-server:v1.0 .

# 运行（与 compose 创建的网络相连仅作示例，按需调整）
docker run -d --name todolist-app -p 8080:8080 \
	-e APP_ENV=prod -e JWT_SECRET=your-strong-secret \
	todolist-server:v1.0
```

## API 速览

基础前缀：`/api/v1`

### 用户

- 注册：`POST /api/v1/user/register`

请求示例：
```json
{
	"email": "test@example.com",
	"password": "123456"
}
```

- 登录：`POST /api/v1/user/login`

响应示例：
```json
{
	"code": 200,
	"msg": "ok",
	"data": {
		"token": "<jwt>",
		"user": {
			"id": 1,
			"username": "u_xxxxxx",
			"email": "test@example.com",
			"nickname": "u_xxxxxx"
		}
	}
}
```

- 登出：`POST /api/v1/user/logout`（需要 `Authorization: Bearer <token>`）

### 任务

- 创建：`POST /api/v1/task/create`（需登录）

请求示例：
```json
{
	"title": "Read GORM docs",
	"content": "focus on associations",
	"status": 0,
	"category": "study"
}
```

- 列表：`GET /api/v1/task/list?page=1&pageSize=10[&status=0]`（需登录）

响应示例：
```json
{
	"code": 200,
	"msg": "ok",
	"data": {
		"items": [
			{"id": 1, "title": "Read GORM docs", "status": 0, "user_id": 1},
			{"id": 2, "title": "Write tests", "status": 1, "user_id": 1}
		],
		"total": 2
	}
}
```

- 更新：`PUT /api/v1/task/:id`（需登录）

请求示例（字段可选，未提供的不更新）：
```json
{
	"title": "Read GORM docs v2",
	"content": "add preload examples",
	"status": 1,
	"category": "study"
}
```

- 删除：`DELETE /api/v1/task/:id`（需登录）

说明：使用 Gorm 软删除（保留审计信息）。

### 统一响应格式

成功：
```json
{"code":200, "msg":"ok", "data":{...}}
```
失败：
```json
{"code":<ErrCode>, "msg":"<ErrMsg>", "data":null}
```

常见错误码（节选）：

| Code  | Msg                | 含义                 |
|-------|--------------------|----------------------|
| 200   | ok                 | 成功                 |
| 400   | 请求参数错误        | 参数校验失败         |
| 10001 | Token鉴权失败       | 未携带或格式错误     |
| 10002 | Token已超时         | 令牌已过期           |
| 20001 | 用户不存在          | 登陆/鉴权失败        |
| 20002 | 用户名已存在        | 注册时邮箱已占用     |
| 20003 | 密码错误            | 登陆失败             |
| 30001 | 任务不存在          | 资源未找到           |

## 开发说明

- 分层职责：
	- `routes` 仅做路由分发
	- `api/v1` 负责参数绑定与响应封装
	- `service` 承载业务逻辑
	- `repository/db/dao` 负责数据库读写
	- `repository/db/model` 定义数据结构
	- `types` 定义请求与响应结构体
- CORS：见 `internal/middleware/cors.go`
- JWT：见 `internal/middleware/jwt.go` 与 `pkg/utils/jwt.go`
- 自动迁移：`internal/repository/db/init.go` 在启动时对 `User` 与 `Task` 执行 `AutoMigrate`

## 本地测试示例

登录获取 Token 后，带上请求头 `Authorization: Bearer <token>` 即可请求受保护接口：

```bash
# 创建任务
curl -X POST "http://localhost:8080/api/v1/task/create" \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer <your_jwt_token>" \
	-d '{"title":"Read GORM docs","content":"...","status":0,"category":"study"}'

# 查询任务
curl "http://localhost:8080/api/v1/task/list?page=1&pageSize=10" \
	-H "Authorization: Bearer <your_jwt_token>"

# 更新任务
curl -X PUT "http://localhost:8080/api/v1/task/1" \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer <your_jwt_token>" \
	-d '{"status":1}'

# 删除任务
curl -X DELETE "http://localhost:8080/api/v1/task/1" \
	-H "Authorization: Bearer <your_jwt_token>"
```

## 备注

- 删除任务为软删除，如需硬删除可在 DAO 层改为 `Unscoped().Delete(...)`
- 生产环境请务必：
	- 设置强随机的 `JWT_SECRET`
	- 使用 `APP_ENV=prod` 并提供 `config/config.prod.yaml`
	- 对外仅暴露必要端口，做好网络与数据安全

