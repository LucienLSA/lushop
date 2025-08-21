import request from '@/utils/request'

// 商品列表
export function fetchList(params) {
  return request({
    url: 'goods',
    method: 'get',
    params: params
  })
}

// 创建商品
export function createProduct(data) {
  return request({
    url: 'goods',
    method: 'post',
    data: data
  })
}

// 更新商品
export function updateProduct(id, data) {
  return request({
    url: `goods/${id}`,
    method: 'put',
    data: data
  })
}

// 删除商品
export function deleteProduct(id) {
  return request({
    url: `goods/${id}`,
    method: 'delete'
  })
}

// 获取商品详情
export function getProduct(id) {
  return request({
    url: `goods/${id}`,
    method: 'get'
  })
}

// 更新商品状态
export function updateProductStatus(id, data) {
  return request({
    url: `goods/${id}`,
    method: 'patch',
    data: data
  })
}

// 获取商品库存
export function getProductStocks(id) {
  return request({
    url: `goods/${id}/stocks`,
    method: 'get'
  })
} 