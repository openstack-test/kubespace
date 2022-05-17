<template>
  <div class="cluster">
    <div class="search-bar">
      <div class="title">
        <el-button class="add-cluster-button-import" type="primary" @click="showCreateMethod()" size="small">新建</el-button>
        <el-input class="title-input" v-model="queryData.name" placeholder="输入方法名称" style="width: 200px;"></el-input>
        <el-button class="title-search" type="primary" icon="el-icon-search" @click="pageCurrentChange(1)" size="small">搜索</el-button>
      </div>
    </div>
    <el-card class="table-card">
      <div v-loading="loadingTable">
        <el-table
          :height="this.pageHeight-200"
          highlight-current-row
          :sort-orders="['ascending', 'descending']"
          :data="tableData.data">
          <el-table-column
            show-overflow-tooltip
            label="编号"
            width="80"
            align="center"
            prop="id">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="方法名称"
            width="auto"
            align="center"
            prop="name">
          </el-table-column>          
          <el-table-column
            show-overflow-tooltip
            label="通知地址"
            align="center"
            prop="call_url"
            width="auto">
          </el-table-column>
          <el-table-column
            label="操作"
            align="center"
            width="140">
            <template slot-scope="scope">
              <div>
                <el-button type="primary" @click="showUpdateMethod(scope.row)" size="mini">
                  编辑<i class=" el-icon--right"></i>
                </el-button>
                <el-button type="primary" @click="delMethod(scope.row)" size="mini">
                  删除<i class=" el-icon--right"></i>
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          class="center"
          ref="pagination"
          :pager-count="5"
          :small="true"
          @size-change="pageSizeChange"
          :current-page="queryData.current_page"
          :page-sizes="[10, 20, 30, 40, 50]"
          :page-size="queryData.page_size"
          layout="total ,sizes, prev, pager, next, jumper"
          :total="tableData.total">
        </el-pagination>
      </div>
    </el-card>

    <!-- 创建告警组 -->
    <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="450px" :close="handCreateDialogClose">
      <el-form :model="formFields" ref="formFields" :rules="formRules" label-width="100px" class="demo-form-inline">
        <el-form-item label="方法名称" prop="name" >
          <el-input v-model.trim="formFields.name" placeholder="输入方法名称"  />
        </el-form-item>      
        <el-form-item label="通知地址" prop="call_url">
          <el-input v-model="formFields.call_url"  placeholder="请输入通知地址" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="createMethod('formFields')" :disabled="isSubmiting" :loading="isSubmiting"><span v-show="!isSubmiting">确 定</span><span v-show="isSubmiting">创建中</span></el-button>
      </div>
    </el-dialog>

    <!--更新告警组 -->
    <el-dialog :title="editdialogFormTitle" :visible.sync="editdialogFormVisible" width="450px" :close="handEditDialogClose">
      <el-form :model="formEdit" ref="formEdit" :rules="formRules" label-width="100px" class="demo-form-inline">
        <el-form-item label="方法名称" prop="name" >
          <el-input v-model.trim="formEdit.name" placeholder="输入方法名称"  />
        </el-form-item>      
        <el-form-item label="通知地址" prop="call_url">
          <el-input  v-model="formEdit.call_url" placeholder="请输入通知地址" />
        </el-form-item>
      </el-form>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="editdialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="editMethod('formEdit')" :disabled="isSubmiting" :loading="isSubmiting"><span v-show="!isSubmiting">更 新</span><span v-show="isSubmiting">创建中</span></el-button>
      </div>
    </el-dialog>    
  </div>
</template>

<script>
  import {request} from "../../utils/rquestes";
  import {deepCopy, isEmpty} from "../../utils/common";
  export default {
    name: "method",
    data: function () {
      return {
        loadingTable: true,
        pageHeight: document.body.scrollHeight,
        dialogFormVisible: false,
        dialogFormTitle: null,
        editdialogFormVisible: false,
        editdialogFormTitle: null,
        editdialogFormUrl: null,
        dialogFormUrl: null,
        isSubmiting: false,
        formFields:{
          name: "",
          call_url: "",
        },
        formEdit: {},
        userList: [],
        queryData: {
          current_page: 1,
          page_size: 20,
        },
        tableData: {},
        formRules: {
          name: [
            {required: true, message: '请输入告警组', trigger: 'blur'},
          ],
          call_url: [
            {required: true, message: '请输入通知地址', trigger: 'blur'},
          ],
        },
      }
    },
    // beforeMount() {
      // this.getPromList();
      // this.getUserList();
    // },
    mounted: function () {
      this.getPageList();
    },
    methods: {
      // 创建规则
      showCreateMethod: function() {
        this.dialogFormTitle = "添加告警媒介"
        this.dialogFormVisible = true
        this.dialogFormUrl = "/monitor/addMethod"
      },
      createMethod: function(formName){
        this.$refs[formName].validate(async (valid) => {
          if (valid) {
            this.isSubmiting = true;
            try{
              // this.formFields.user = this.formFields.user.toString(); //转字符串
              await request.post(this.dialogFormUrl, this.formFields);
              this.pageCurrentChange(1);
              this.dialogFormVisible = false;
              this.$message.success("提交成功！");
            }catch (e) {
              // console.log(e);
            }finally {
              this.isSubmiting = false;
            }
          } else {
            console.log('error submit!!');
            return false;
          }
        });
      },
      // 更新规则
      showUpdateMethod: function(evl) {
        this.editdialogFormTitle = "编辑告警媒介";
        this.editdialogFormVisible = true;
        this.editdialogFormUrl = "/monitor/putMethod";
        this.formEdit = deepCopy(evl);
        // this.changeSelect(evl);
      },
      editMethod: function(formName) {
        this.$refs[formName].validate(async (valid) => {
          if (valid) {
            this.isSubmiting = true;
            try {
              // this.formEdit.user = this.formEdit.user.toString(); //转字符串
              await request.post(this.editdialogFormUrl, this.formEdit);
              this.pageCurrentChange(1);
              this.editdialogFormVisible = false;
              this.$message.success("修改成功!");
            } catch (e) {
            } finally {
              this.isSubmiting = false;
            }
          } else {
            console.log("error submit!!");
            return false;
          }
        });
      },
      delMethod(evl) {
        this.$confirm('是否要媒介:\n'+evl.name+' ? 删除不可找回!!!!', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          beforeClose: (action, instance, done) => {
            if (action === 'confirm') {
              this.loadingTable = true;
              done();
            } else {
              done();
            }
          }
        }).then(async () => {
          await request.post('/monitor/deleteMethod', evl);
          this.pageCurrentChange(1);
          this.$message.success("操作成功！");
        }).catch(() => {
        });
      },
      //  分页获取
      async getPageList() {
        this.loadingTable=true;
        let filters = deepCopy(this.queryData);
        this.tableData = await request.get('/monitor/getMethod', filters);
        this.loadingTable=false;
      },
      pageCurrentChange(val) {
        this.queryData.current_page = val;
        this.getPageList();
      },
      //分页
      pageSizeChange(val) {
        this.queryData.page_size = val;
        this.getPageList();
      },
      handCreateDialogClose(){
        this.isSubmiting = false;
      },
      handEditDialogClose(){
        this.handEditDialogClose = false;
      },      
    },
  };
</script>

<style scope>
  .cluster .el-dialog__body .descClass .el-input__inner {
    width: 221px !important;
  }

  .el-icon-arrow-up:before {
      content: '';
  }
  .el-pagination {
      text-align: center; 
  }
  .title {
    text-align: right;
  }
  .title-input {
    vertical-align: middle;
  }
  .title-search {
    margin-left: 12px;
    vertical-align: middle;
  }
  .add-cluster-button {
    /* margin-right: 70%;*/
    vertical-align: middle;
  }
  .add-cluster-button-import {
    margin-right: 78%;
  }
  .cluster {
    width: 100%;
    height: 100%;
    padding: 10px;
    box-sizing: border-box;
    display: flex;
    flex-wrap: wrap;
  }
  .cluster .el-dialog__body .el-form-item{
    margin-right: 0;
    width: 49%;
  }
  .cluster .el-dialog__body .el-form-item .el-input__inner {
    width: 260px;
  }

  .cluster .demo-table-expand {
    font-size: 0;
  }

  .cluster .demo-table-expand label {
    width: 100px;
    color: #99a9bf;
  }

  .cluster .demo-table-expand .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
    width: 33%;
  }
  .row div {
    margin-button: 20px;
    margin-left: 40px;
    width: 20%;
    display: inline-block;
  }
  .cut-img {
    width: 15px;
    height: 15px;
    /* animation:changDeg 2s linear 0.2s infinite; */
  }
	@keyframes changDeg{
			0%{
				transform: rotate(0deg);
			}
			100%{
				transform: rotate(360deg);
			}
		}
</style>