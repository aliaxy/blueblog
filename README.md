# BlueBlog

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

基于 Go + Vue 的社区论坛系统，支持帖子发布、社区分类、投票排序等功能。

## 特性

- JWT 鉴权，令牌桶限流
- 帖子支持按时间 / 分数排序，基于 Redis Sorted Set 实现
- 雪花算法生成分布式唯一 ID
- Zap + Lumberjack 结构化日志 & 自动切分
- Swagger 交互式 API 文档
- Docker Compose 一键部署

## 技术栈

| 层次 | 技术 |
|------|------|
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) |
| 前端框架 | Vue.js |
| 关系数据库 | MySQL 8.0 |
| 缓存 / 排序 | Redis |
| 数据库操作 | [sqlx](https://github.com/jmoiron/sqlx) |
| 配置管理 | [Viper](https://github.com/spf13/viper) |
| 日志 | [Zap](https://github.com/uber-go/zap) + [Lumberjack](https://github.com/natefinch/lumberjack) |
| 认证 | JWT (`dgrijalva/jwt-go`) |
| ID 生成 | Snowflake |
| 限流 | 令牌桶 (`juju/ratelimit`) |
| 参数校验 | `go-playground/validator` |
| API 文档 | Swagger |
| 性能分析 | pprof |
| 容器化 | Docker / Docker Compose |

## 项目结构

```
blueblog/
├── conf/               # 配置文件（config.yaml.example 为模板）
├── controller/         # HTTP 处理层
├── dao/
│   ├── mysql/          # MySQL 数据访问
│   └── redis/          # Redis 数据访问
├── docs/               # Swagger 自动生成文档
├── logger/             # 日志初始化
├── logic/              # 业务逻辑层
├── logs/               # 运行时日志（自动切分，不进 git）
├── middleware/         # 中间件（JWT、限流）
├── models/             # 数据模型 & 请求参数
├── pkg/
│   ├── jwt/            # JWT 工具
│   └── snowflake/      # 雪花 ID 生成
├── router/             # 路由注册
├── scripts/
│   ├── init.sql        # 数据库初始化脚本
│   └── wait.sh         # Docker 服务就绪等待脚本
├── settings/           # 配置加载
├── web/
│   ├── static/         # 前端静态资源
│   └── templates/      # HTML 模板
├── docker-compose.yml.example
├── Dockerfile
└── flake.nix           # Nix 开发环境（可选）
```

## 快速开始

### 前置依赖

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose

### 方式一：Docker Compose（推荐）

```bash
git clone <repo-url>
cd blueblog
cp docker-compose.yml.example docker-compose.yml
cp conf/config.yaml.example conf/config.yaml
docker-compose up -d
```

服务启动后访问 `http://localhost:8888`。

> MySQL 映射端口 `23306`，Redis 映射端口 `26379`，应用端口 `8888`。

### 方式二：本地运行

**前置依赖**：MySQL、Redis

**第一步：进入开发环境（推荐使用 Nix）**

项目提供 [Nix Flake](https://nixos.wiki/wiki/Flakes) 开发环境，自动配置 `go`、`gopls`、`gotools`、`golangci-lint`、`delve`，无需手动安装工具链。

```bash
# 使用 direnv 自动激活（首次需授权）
direnv allow

# 或手动进入
nix develop
```

> 需要安装 [Nix](https://nixos.org/download/) 并开启 `experimental-features = nix-command flakes`。  
> 不使用 Nix 的话，手动安装 Go 1.22+ 即可。

**第二步：初始化数据库**

```bash
mysql -u root -p < scripts/init.sql
```

**第三步：准备配置文件**

```bash
cp conf/config.yaml.example conf/config.yaml
# 按实际环境修改 mysql / redis 连接信息
```

**第四步：编译并运行**

```bash
go build -o blueblog .
./blueblog conf/config.yaml
```

## 配置说明

配置文件位于 `conf/config.yaml`，以 `conf/config.yaml.example` 为模板：

```yaml
name: "blueblog"
mode: "dev"          # dev | release
port: 8888

auth:
  jwt_expire: 8760   # JWT 有效期（小时）

log:
  level: "debug"
  filename: "logs/blueblog.log"
  max_size: 200      # 单文件最大 MB
  max_age: 30        # 保留天数
  max_backups: 7

mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "your_password"
  dbname: "blueblog"

redis:
  host: "127.0.0.1"
  port: 6379
```

## API 接口

启动后交互式文档见 `http://localhost:8888/swagger/index.html`。

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/signup` | 用户注册 |
| POST | `/api/v1/login` | 用户登录，返回 JWT Token |

登录成功后，需鉴权的接口在请求头携带：

```
Authorization: Bearer <token>
```

### 帖子

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|:----:|------|
| GET  | `/api/v1/posts2` | 否 | 帖子列表（支持排序 & 社区筛选） |
| GET  | `/api/v1/posts` | 是 | 帖子列表（分页） |
| GET  | `/api/v1/post/:id` | 是 | 帖子详情 |
| POST | `/api/v1/post` | 是 | 发布帖子 |
| POST | `/api/v1/vote` | 是 | 帖子投票（+1 / -1） |

**`GET /api/v1/posts2` 查询参数**

| 参数 | 类型 | 默认 | 说明 |
|------|------|------|------|
| `page` | int | 1 | 页码 |
| `size` | int | 10 | 每页条数 |
| `order` | string | `time` | 排序：`time` / `score` |
| `community_id` | int | - | 社区 ID，不传则查全部 |

### 社区

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|:----:|------|
| GET  | `/api/v1/community` | 是 | 社区列表 |
| GET  | `/api/v1/community/:id` | 是 | 社区详情 |

### 其他

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/ping` | 健康检查 |
| GET | `/swagger/*any` | API 文档 |
| GET | `/debug/pprof/*` | 性能分析 |

## 数据库设计

| 表名 | 说明 |
|------|------|
| `user` | 用户信息，用户 ID 由 Snowflake 生成 |
| `post` | 帖子，帖子 ID 由 Snowflake 生成 |
| `community` | 社区分类 |

帖子分数排序通过 Redis Sorted Set 实现，投票数据存储于 Redis。

## 开发

### 热重载（Air）

```bash
air
```

### 生成 Swagger 文档

```bash
swag init
```

### 压测

```bash
# 需要安装 wrk
wrk -t8 -c100 -d30s http://localhost:8888/api/v1/posts2
```

### 性能分析

```bash
go tool pprof http://localhost:8888/debug/pprof/profile
```

## License

[MIT](./LICENSE)
