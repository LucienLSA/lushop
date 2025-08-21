<template>
  <div class="app-container">
    <el-card class="operate-container" shadow="never">
      <i class="el-icon-tickets"></i>
      <span>数据列表</span>
      <el-button
        class="btn-add"
        @click="addBrand()"
        size="mini">
        添加
      </el-button>
    </el-card>
    <div class="table-container">
      <el-table ref="brandTable"
                :data="list"
                style="width: 100%"
                @selection-change="handleSelectionChange"
                v-loading="listLoading"
                border>
                  <el-table-column label="品牌id" align="center">
          <template slot-scope="scope">{{scope.row.id}}</template>
        </el-table-column>
        <el-table-column label="品牌名称" align="center">
          <template slot-scope="scope">{{scope.row.name}}</template>
        </el-table-column>
       
        <el-table-column label="图片" align="center">
          <template slot-scope="scope">
            <img :src="scope.row.logo" class="imgs" alt="">
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center">
          <template slot-scope="scope">
            <el-button
              size="mini"
              @click="handleUpdate(scope.$index, scope.row)">编辑
            </el-button>
            <el-button
              size="mini"
              type="danger"
              @click="handleDelete(scope.$index, scope.row)">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
     <div class="pagination-container">
      <el-pagination
        background
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        layout="total, sizes,prev, pager, next,jumper"
        :page-size="listQuery.pnum"
        :current-page.sync="pageNum"
        :total="total">
        <!-- :page-sizes="[5,10,15]" -->
      </el-pagination>
    </div>
  </div>
</template>
<script>
  import {fetchList, deleteBrand} from '@/api/lushop/brand'

  export default {
    name: 'brandList',
    data() {
      return {
        operates: [
          {
            label: "显示品牌",
            value: "showBrand"
          },
          {
            label: "隐藏品牌",
            value: "hideBrand"
          }
        ],
        operateType: null,
        pageNum:0,
        listQuery: {
          pn: 1,
          pnum: 10
        },
        list: [],
        total: null,
        listLoading: true,
        multipleSelection: []
      }
    },
    created() {
      this.getList();
    },
    methods: {
      getList() {
        this.listLoading = true;
        fetchList(this.listQuery).then(response => {
          this.listLoading = false;
          // 适配lushop_api的响应格式
          if (response.code === 200) {
            this.list = response.data.list || [];
            this.total = response.data.total || 0;
          } else {
            this.$message.error(response.msg || '获取品牌列表失败');
          }
        }).catch(error => {
          this.listLoading = false;
          this.$message.error('获取品牌列表失败');
        });
      },
      handleSelectionChange(val) {
        this.multipleSelection = val;
      },
      handleSizeChange(val) {
        this.listQuery.pnum = val;
        this.getList();
      },
      handleCurrentChange(val) {
        this.listQuery.pn = val;
        this.getList();
      },
      addBrand() {
        this.$router.push({path: '/pms/addBrand'})
      },
      handleUpdate(index, row) {
        this.$router.push({path: '/pms/updateBrand', query: {id: row.id}})
      },
      handleDelete(index, row) {
        this.$confirm('是否要删除该品牌?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          deleteBrand(row.id).then(response => {
            this.$message({
              type: 'success',
              message: '删除成功!'
            });
            this.getList();
          }).catch(error => {
            this.$message.error('删除失败');
          });
        });
      }
    }
  }
</script>
<style scoped>
  .operate-container {
    margin-top: 0;
  }
  .operate-container .btn-add {
    float: right;
  }
  .table-container {
    margin-top: 20px;
  }
  .pagination-container {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
  .imgs {
    height: 40px;
    width: 40px;
  }
</style>
</style>


