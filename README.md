# 用户管理 API

基于 Go + Gin + GORM + MySQL 的用户管理系统，实现了 JWT 认证和完整的用户 CRUD。

## 技术栈

- Go 1.21
- Gin
- GORM
- MySQL
- JWT
- bcrypt

## 功能列表

- POST /register - 用户注册（不需要认证）
- POST /login - 用户登录（不需要认证）
- POST /user - 创建用户（需要认证）
- GET /users - 获取所有用户（需要认证）
- GET /user/:id - 获取单个用户（需要认证）
- PUT /user/:id - 更新用户（需要认证）
- DELETE /user/:id - 删除用户（需要认证）

## 快速开始

1. 克隆项目

git clone https://github.com/Sp1ke-7/gin-demo.git
cd gin-demo

2. 安装依赖

go mod tidy

3. 运行项目

go run main.go

4. 测试注册

curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"username":"test","password":"123"}'

5. 测试登录

curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"test","password":"123"}'

## 项目亮点

- JWT 认证中间件
- bcrypt 密码加密
- 参数校验
- 完整 CRUD

## 联系方式

- GitHub: Sp1ke-7
