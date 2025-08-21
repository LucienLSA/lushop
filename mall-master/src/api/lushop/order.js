import request from '@/utils/request'

// 订单列表
export function fetchList(params) {
  return request({
    url: 'order',
    method: 'get',
    params: params
  })
}

// 获取订单详情
export function getOrderDetail(id) {
  return request({
    url: `order/${id}`,
    method: 'get'
  })
}

// 创建订单
export function createOrder(data) {
  return request({
    url: 'order',
    method: 'post',
    data: data
  })
}

// 购物车列表
export function getShopCartList() {
  return request({
    url: 'shopcart',
    method: 'get'
  })
}

// 添加到购物车
export function addToShopCart(data) {
  return request({
    url: 'shopcart',
    method: 'post',
    data: data
  })
}

// 更新购物车
export function updateShopCart(id, data) {
  return request({
    url: `shopcart/${id}`,
    method: 'put',
    data: data
  })
}

// 删除购物车商品
export function deleteShopCart(id) {
  return request({
    url: `shopcart/${id}`,
    method: 'delete'
  })
} 