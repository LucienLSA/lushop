<template>
  <div class="test-api">
    <el-card>
      <div slot="header">
        <span>API连接测试</span>
      </div>
      
      <el-form :model="loginForm" label-width="80px">
        <el-form-item label="手机号">
          <el-input v-model="loginForm.mobile" placeholder="请输入手机号"></el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="testLogin">测试登录</el-button>
          <el-button @click="testBrands">测试品牌列表</el-button>
          <el-button @click="testGoods">测试商品列表</el-button>
        </el-form-item>
      </el-form>

      <el-divider></el-divider>
      
      <div v-if="testResult">
        <h4>测试结果:</h4>
        <pre>{{ JSON.stringify(testResult, null, 2) }}</pre>
      </div>
    </el-card>
  </div>
</template>

<script>
import { login, fetchList as fetchBrands } from '@/api/lushop/user'
import { fetchList as fetchGoods } from '@/api/lushop/product'
import { fetchList as fetchBrandList } from '@/api/lushop/brand'

export default {
  name: 'TestApi',
  data() {
    return {
      loginForm: {
        mobile: '13800138000',
        password: '123456'
      },
      testResult: null
    }
  },
  methods: {
    async testLogin() {
      try {
        this.testResult = null
        const result = await login(this.loginForm)
        this.testResult = {
          type: '登录测试',
          success: true,
          data: result
        }
        this.$message.success('登录测试成功')
      } catch (error) {
        this.testResult = {
          type: '登录测试',
          success: false,
          error: error.message || '登录失败'
        }
        this.$message.error('登录测试失败')
      }
    },
    async testBrands() {
      try {
        this.testResult = null
        const result = await fetchBrandList({ pn: 1, pnum: 10 })
        this.testResult = {
          type: '品牌列表测试',
          success: true,
          data: result
        }
        this.$message.success('品牌列表测试成功')
      } catch (error) {
        this.testResult = {
          type: '品牌列表测试',
          success: false,
          error: error.message || '获取品牌列表失败'
        }
        this.$message.error('品牌列表测试失败')
      }
    },
    async testGoods() {
      try {
        this.testResult = null
        const result = await fetchGoods({ pn: 1, pnum: 10 })
        this.testResult = {
          type: '商品列表测试',
          success: true,
          data: result
        }
        this.$message.success('商品列表测试成功')
      } catch (error) {
        this.testResult = {
          type: '商品列表测试',
          success: false,
          error: error.message || '获取商品列表失败'
        }
        this.$message.error('商品列表测试失败')
      }
    }
  }
}
</script>

<style scoped>
.test-api {
  padding: 20px;
}
pre {
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
}
</style> 