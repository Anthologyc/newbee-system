# NewBee System - 企业级后台管理系统

NewBee System 是一个基于现代化技术栈构建的全栈后台管理系统，采用 **Go (Gin)** 作为后端，**Vue 3 (TypeScript)** 作为前端，并使用 **Docker** 进行全容器化部署。系统实现了用户认证、题库管理、公告发布及动态仪表盘等核心功能。

## 🚀 技术栈

### 后端 (Backend)
- **语言**: Golang 1.23+
- **框架**: Gin Web Framework
- **ORM**: GORM (支持 MySQL JSON 类型)
- **认证**: JWT (JSON Web Token)

### 前端 (Frontend)
- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **语言**: TypeScript
- **UI 组件库**: Arco Design Vue
- **路由**: Vue Router 4

### 基础设施 (Infrastructure)
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.0
- **网关**: Nginx (反向代理与静态资源托管)
- **部署**: Docker & Docker Compose

## ✨ 核心功能

1.  **用户认证系统**
    - 用户注册与登录
    - JWT Token 签发与验证
    - 密码加密存储 (Bcrypt)

2.  **题库管理模块**
    - 题目的增删改查 (CRUD)
    - 支持单选、多选、判断等多种题型
    - **JSON 数据存储**: 选项与答案使用 MySQL JSON 字段存储，灵活扩展
    - **批量导入**: 支持 JSON 格式题目批量导入功能

3.  **公告管理系统**
    - 发布系统公告
    - 公告列表展示与删除
    - 首页自动展示最新公告

4.  **动态仪表盘 (Dashboard)**
    - 实时统计题目总数与用户数
    - 展示最新系统公告
    - 系统版本信息展示

## 🛠️ 快速开始

### 前置要求
- Docker & Docker Compose

### 安装与运行

1.  **克隆项目**
    ```bash
    git clone https://github.com/Anthologyc/newbee-system.git
    cd newbee-system
    ```

2.  **一键启动 (Docker)**
    使用 Docker Compose 编排所有服务（前端、后端、MySQL、Redis）：
    ```bash
    docker-compose up --build
    ```

3.  **访问系统**
    - 前端页面: [http://localhost](http://localhost)
    - 后端 API: [http://localhost/api](http://localhost/api) (通过 Nginx 转发)
    - 数据库端口: 3307 (宿主机映射)

### 📦 目录结构

```text
.
├── backend/                # Go 后端源码
│   ├── models/             # GORM 数据模型 (User, Question, Announcement)
│   ├── main.go             # 核心业务逻辑与路由
│   └── Dockerfile          # 后端镜像构建文件
├── frontend/               # Vue 前端源码
│   ├── src/
│   │   ├── views/          # 页面组件 (Login, Dashboard, QuestionList...)
│   │   ├── layouts/        # 布局组件 (侧边栏、顶部栏)
│   │   └── router/         # 路由配置
│   ├── nginx.conf          # Nginx 配置文件 (反向代理 /api)
│   └── Dockerfile          # 前端镜像构建文件
├── docker-compose.yml      # 容器编排配置
└── README.md               # 项目说明文档


🔌 API 概览
方法	路径	描述
POST	/register	用户注册
POST	/login	用户登录
GET	/api/dashboard/stats	获取仪表盘统计数据
GET	/api/questions	获取题目列表 (支持筛选)
POST	/api/questions	新增题目
PUT	/api/questions/:id	更新题目
DELETE	/api/questions/:id	删除题目
GET	/api/announcements	获取公告列表
POST	/api/announcements	发布公告