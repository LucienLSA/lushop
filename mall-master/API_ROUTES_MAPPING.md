# Lushop 前后端API路由映射

## 基础配置

### 前端配置
- 开发环境API地址: `http://172.25.6.248:8101/v2/`
- 认证方式: JWT Bearer Token
- 请求头: `Authorization: Bearer <token>`

### 后端配置
- 服务地址: `http://172.25.6.248:8101`
- API版本: `v2`
- 认证中间件: JWT Auth

## API路由映射

### 1. 商品管理 (Goods)

| 前端路由 | 后端API | 方法 | 描述 |
|---------|---------|------|------|
| `/product` | `/v2/goods` | GET | 商品列表 |
| `/addProduct` | `/v2/goods` | POST | 创建商品 |
| `/updateProduct/:id` | `/v2/goods/:id` | PUT | 更新商品 |
| `/product/:id` | `/v2/goods/:id` | DELETE | 删除商品 |
| `/product/:id` | `/v2/goods/:id` | GET | 商品详情 |
| - | `/v2/goods/:id/stocks` | GET | 获取库存 |

### 2. 品牌管理 (Brands)

| 前端路由 | 后端API | 方法 | 描述 |
|---------|---------|------|------|
| `/brand` | `/v2/brands` | GET | 品牌列表 |
| `/addBrand` | `/v2/brands` | POST | 创建品牌 |
| `/updateBrand/:id` | `/v2/brands/:id` | PUT | 更新品牌 |
| `/brand/:id` | `/v2/brands/:id` | DELETE | 删除品牌 |

### 3. 分类管理 (Categories)

| 前端路由 | 后端API | 方法 | 描述 |
|---------|---------|------|------|
| `/productCate` | `/v2/categorys` | GET | 分类列表 |
| `/addProductCate` | `/v2/categorys` | POST | 创建分类 |
| `/updateProductCate/:id` | `/v2/categorys/:id` | PUT | 更新分类 |
| `/productCate/:id` | `/v2/categorys/:id` | DELETE | 删除分类 |
| - | `/v2/categorybrands/:id` | GET | 获取分类品牌 |

### 4. 订单管理 (Orders)

| 前端路由 | 后端API | 方法 | 描述 |
|---------|---------|------|------|
| `/oms/order` | `/v2/order` | GET | 订单列表 |
| `/oms/orderDetail/:id` | `/v2/order/:id` | GET | 订单详情 |
| - | `/v2/order` | POST | 创建订单 |
| - | `/v2/shopcart` | GET | 购物车列表 |
| - | `/v2/shopcart` | POST | 添加到购物车 |
| - | `/v2/shopcart/:id` | PUT | 更新购物车 |
| - | `/v2/shopcart/:id` | DELETE | 删除购物车商品 |

### 5. 用户管理 (Users)

| 前端路由 | 后端API | 方法 | 描述 |
|---------|---------|------|------|
| `/user/user` | `/v2/user/list` | GET | 用户列表 |
| `/login` | `/v2/user/pwd_login` | POST | 用户登录 |
| - | `/v2/user/register` | POST | 用户注册 |
| - | `/v2/user/detail` | GET | 用户详情 |
| - | `/v2/user/update` | PATCH | 更新用户信息 |
| `/user/address` | `/v2/address` | GET | 地址列表 |
| - | `/v2/address` | POST | 创建地址 |
| - | `/v2/address/:id` | PUT | 更新地址 |
| - | `/v2/address/:id` | DELETE | 删除地址 |
| `/user/rotation` | `/v2/banners` | GET | 轮播图列表 |
| - | `/v2/banners` | POST | 创建轮播图 |
| - | `/v2/banners/:id` | PUT | 更新轮播图 |
| - | `/v2/banners/:id` | DELETE | 删除轮播图 |

### 6. 文件上传

| 前端路由 | 后端API | 方法 | 描述 |
|---------|---------|------|------|
| - | `/v2/oss/upload` | POST | 文件上传 |

## 认证流程

1. **登录获取Token**
   ```
   POST /v2/user/pwd_login
   {
     "mobile": "手机号",
     "password": "密码"
   }
   ```

2. **请求头携带Token**
   ```
   Authorization: Bearer <access_token>
   ```

3. **Token刷新**
   ```
   GET /v2/user/refresh
   ```

## 响应格式

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    // 具体数据
  }
}
```

## 错误处理

- `401`: 未授权，需要重新登录
- `403`: 权限不足
- `404`: 资源不存在
- `500`: 服务器内部错误

## 启动命令

### 后端启动
```bash
cd app/lushop_api
go run main.go
```

### 前端启动
```bash
cd mall-master
npm run dev
```

### 一键启动
```bash
start_dev.bat
``` 