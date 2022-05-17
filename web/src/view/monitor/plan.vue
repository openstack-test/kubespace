<template>
  <div class="plan">
    <div class="search-bar">
      <div class="title">
        <el-button class="add-cluster-button-import" type="primary" @click="showCreateRule()" size="small">新建</el-button>
        <el-input class="title-input" v-model="queryData.rule_labels" placeholder="输入组名称" style="width: 200px;"></el-input>
        <el-button class="title-search" type="primary" icon="el-icon-search" @click="pageCurrentChange(1)" size="small">搜索</el-button>
      </div>
    </div>
    <el-card class="table-card">
      <div v-loading="loadingTable">
        <el-table
          :height="this.pageHeight-200"
          highlight-current-row
          :sort-orders="['ascending', 'descending']"
          @expand-change="expandChange"
          :data="tableData.data">
          <el-table-column type="expand">
            <template slot-scope="props">
	            <el-table
	                :data="props.row.childrenData"
	                v-loading="props.row.loading">
                <el-table-column
                  align="center"
                  show-overflow-tooltip
                  label="告警时间段">
                  <template slot-scope="scope">
                    <div>
                      {{scope.row.start_time}}~{{scope.row.end_time}}
                    </div>
                  </template>
                </el-table-column>
                <el-table-column
                  prop="start"
                  align="center"
                  show-overflow-tooltip
                  label="告警延迟">
                </el-table-column>
                <el-table-column
                  prop="period"
                  align="center"
                  show-overflow-tooltip
                  label="持续周期">
                </el-table-column>
                <el-table-column
                  prop="group"
                  align="center"
                  show-overflow-tooltip
                  label="告警组">
                </el-table-column>
                <el-table-column
                  prop="expression"
                  align="center"
                  show-overflow-tooltip
                  label="告警过滤表达式">
                </el-table-column>
                <el-table-column
                  prop="method"
                  align="center"
                  show-overflow-tooltip
                  label="通知方式">
                </el-table-column>
                <el-table-column prop="handle">
                  <template slot-scope="scope">
                    <span class="buton-css" @click="showUpdReceiver(props.row.id,scope.row)">编辑 |</span>
                    <span class="buton-css" @click="delReceiver(props.row.id,scope.row)"> 删除</span>
                  </template>
                  <template slot="header">
                    <span>操作 | </span>
                    <!-- props.row.id 将告警组的ID传到子项 -->
                    <span class="buton-css"  @click="showPlanReceiver(props.row.id)">添加</span>
                  </template>
                </el-table-column>
              </el-table>
            </template>
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="编号"
            width="80"
            align="center"
            prop="id">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="名称"
            align="center"
            width="auto"
            prop="rule_labels"
            >
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="描述"
            align="center"
            width="auto"
            prop="description">
          </el-table-column>
          <el-table-column
            label="操作"
            align="center"
            width="140">
            <template slot-scope="scope">
              <div>
                <el-button type="primary" @click="showUpdatePlan(scope.row)" size="mini">
                  编辑<i class=" el-icon--right"></i>
                </el-button>
                <el-button type="primary" @click="delPlan(scope.row)" size="mini">
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

    <!-- 添加告警计划组 -->
    <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="450px" :close="handCreateDialogClose">
      <el-form :model="formFields" ref="formFields" :rules="formRules" label-width="100px" class="demo-form-inline">
        <el-form-item label="告警组名称" prop="rule_labels" >
          <el-input v-model.trim="formFields.rule_labels" placeholder="输入告警组名称"  />
        </el-form-item>
        <el-form-item label="描述" prop="description" >
          <el-input v-model.trim="formFields.description" placeholder="请输入描述信息." />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="createPlan('formFields')" :disabled="isSubmiting" :loading="isSubmiting"><span v-show="!isSubmiting">确 定</span><span v-show="isSubmiting">创建中</span></el-button>
      </div>
    </el-dialog>

    <!-- 编辑告警计划组 -->
    <el-dialog :title="editdialogFormTitle" :visible.sync="editdialogFormVisible" width="450px" :close="handEditDialogClose">
      <el-form :model="formEdit" ref="formEdit" :rules="formRules" label-width="100px" class="demo-form-inline">
        <el-form-item label="告警组名称" prop="rule_labels" >
          <el-input v-model.trim="formEdit.rule_labels" placeholder="输入告警组名称"  />
        </el-form-item>
        <el-form-item label="描述" prop="description" >
          <el-input v-model.trim="formEdit.description" placeholder="请输入描述信息." />
        </el-form-item>
      </el-form>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="editdialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="editPlan('formEdit')" :disabled="isSubmiting" :loading="isSubmiting"><span v-show="!isSubmiting">更 新</span><span v-show="isSubmiting">创建中</span></el-button>
      </div>
    </el-dialog>

    <!-- 添加告警策略 -->
    <el-dialog :title="planReceiverFormVisibleTitle" :visible.sync="planReceiverFormVisible" width="550px" :close="handCreateDialogClose">
      <el-form :model="formPlanReceiver" ref="formPlanReceiver" :rules="formPlanReceiverRules" label-width="100px" class="demo-form-inline">
        <el-form-item label="起始时间" prop="startTime">
          <el-time-select
            placeholder="起始时间"
            v-model="formPlanReceiver.start_time"
            :picker-options="{
              start: '00:00',
              step: '00:03',
              end: '23:59'
            }">
          </el-time-select>
        </el-form-item>
        <el-form-item label="结束时间" prop="endTime">
          <el-time-select
            placeholder="结束时间"
            v-model="formPlanReceiver.end_time"
            filterable
            :picker-options="{
              start: '00:00',
              step: '00:03',
              end: '23:59',
              minTime: formPlanReceiver.start_time
            }">
          </el-time-select> 
        </el-form-item>
        <el-form-item label="报警延迟" prop="start" >
          <el-input v-model.number="formPlanReceiver.start" placeholder="请输入告警延迟时间,单位分钟." />
        </el-form-item>
        <el-form-item label="触发周期" prop="period" >
          <el-input v-model.number="formPlanReceiver.period" placeholder="请输入触发周期次数" />
        </el-form-item>
        <el-form-item label="告警组" prop="group">
          <el-select style="width:50px" v-model="formPlanReceiver.group" filterable multiple placeholder="请选择通知告警组">
            <el-option
              v-for="(item,index) in userList"
              :key="index"
              :label="item.name"
              :value="item.name">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="告警过滤" prop="expression" >
          <el-input v-model.trim="formPlanReceiver.expression" placeholder="请输入告警过滤规则,例如:job=web." />
        </el-form-item>
        <el-form-item label="通知媒介" prop="method" >
          <el-radio v-model="formPlanReceiver.method" label="WorkChat">企微</el-radio>
          <el-radio v-model="formPlanReceiver.method" label="Email">邮件</el-radio>
        </el-form-item>
        <el-form-item label="通知地址" prop="call_url" >
          <el-select style="width:60px" v-model="formPlanReceiver.call_url" filterable multiple placeholder="通知机器人,仅限第三方接口">
            <el-option
              v-for="(item,index) in methodList"
              :key="index"
              :label="item.name"
              :value="item.name">
            </el-option>
          </el-select>
        </el-form-item>          
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="planReceiverFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="createReceiver('formPlanReceiver')" :disabled="isSubmiting" :loading="isSubmiting"><span v-show="!isSubmiting">确 定</span><span v-show="isSubmiting">创建中</span></el-button>
      </div>
    </el-dialog>

    <!-- 更新告警策略 -->
    <el-dialog :title="editplanReceiverFormTitle" :visible.sync="editplanReceiverFormVisible" width="550px" :close="handCreateDialogClose">
      <el-form :model="formPlanReceiverEdit" ref="formPlanReceiver" :rules="formPlanReceiverRules" label-width="100px" class="demo-form-inline">
        <el-form-item label="起始时间" prop="startTime">
          <el-time-select
            placeholder="起始时间"
            v-model="formPlanReceiverEdit.start_time"
            :picker-options="{
              start: '00:00',
              step: '00:03',
              end: '23:59'
            }">
          </el-time-select>
        </el-form-item>
        <el-form-item label="结束时间" prop="endTime">
          <el-time-select
            placeholder="结束时间"
            v-model="formPlanReceiverEdit.end_time"
            :picker-options="{
              start: '00:00',
              step: '00:03',
              end: '23:59',
              minTime: formPlanReceiverEdit.start_time
            }">
          </el-time-select> 
        </el-form-item>
        <el-form-item label="报警延迟" prop="start" >
          <el-input v-model.number="formPlanReceiverEdit.start" placeholder="请输入告警延迟时间,单位分钟." />
        </el-form-item>
        <el-form-item label="触发周期" prop="period" >
          <el-input v-model.number="formPlanReceiverEdit.period" placeholder="请输入触发周期次数" />
        </el-form-item>
        <el-form-item label="告警组" prop="group">
          <el-select style="width: 60px" v-model="formPlanReceiverEdit.group" filterable multiple placeholder="请选择通知告警组">
            <el-option
              v-for="(item,index) in userList"
              :key="index"
              :label="item.name"
              :value="item.name">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="告警过滤" prop="expression" >
          <el-input v-model.trim="formPlanReceiverEdit.expression" placeholder="请输入告警过滤规则,例如:job=web." />
        </el-form-item>
        <el-form-item label="通知媒介" prop="method" >
          <el-radio v-model="formPlanReceiverEdit.method" label="WorkChat">企微</el-radio>
          <el-radio v-model="formPlanReceiverEdit.method" label="Email">邮件</el-radio>
        </el-form-item>
        <el-form-item label="通知地址" prop="call_url" >
          <el-select style="width:60px" v-model="formPlanReceiverEdit.call_url" filterable multiple placeholder="通知机器人,仅限第三方接口">
            <el-option
              v-for="(item,index) in methodList"
              :key="index"
              :label="item.name"
              :value="item.name">
            </el-option>
          </el-select>
        </el-form-item>        
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="editplanReceiverFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="editReceiver('formPlanReceiver')" :disabled="isSubmiting" :loading="isSubmiting"><span v-show="!isSubmiting">确 定</span><span v-show="isSubmiting">创建中</span></el-button>
      </div>
    </el-dialog>    
  </div>
</template>

<script>
  import {request} from "../../utils/rquestes";
  import {deepCopy, isEmpty} from "../../utils/common";
  export default {
    name: "plan_list",
    data: function () {
      return {
        loadingTable: true,
        pageHeight: document.body.scrollHeight,
        dialogFormVisible: false,
        dialogFormTitle: null,

        editdialogFormVisible: false,
        editdialogFormTitle: null,
        editdialogFormUrl: null,
        //告警策略
        planReceiverFormVisible: false,
        planReceiverFormVisibleTitle: null,
        planReceiverFormUrl: null,

        editplanReceiverFormTitle: null,
        editplanReceiverFormUrl: null,
        editplanReceiverFormVisible: false,

        dialogFormUrl: null,
        userList: [],
        methodList: [],
        isSubmiting: false,
        formFields:{},
        formPlanReceiver: {},
        formEdit: {},
        formPlanReceiverEdit: {},
        prom_id_list: [],
        plan_id_list: [],
        operator: ["==",">","<",">=","<=","!="],
        queryData: {
          current_page: 1,
          page_size: 20,
        },
        tableData: {},
        formPlanReceiverRules: {
          start_time: [
            {required: true, message: '请选择开始时间', trigger: 'blur'},
          ],
          end_time: [
            {required: true, message: '请选择结束时间', trigger: 'blur'},
          ],
          start:[
            {required: true, message: '请输入延迟时间', trigger: 'blur'},
          ],
          period:[
            {required: true, message: '请周期次数', trigger: 'blur'},
          ],
          // user:[
          //   {required: true, message: '请输入告警通知用户', trigger: 'blur'},
          // ],
          // group:[
          //   {required: true, message: '请输入告警通知用户组', trigger: 'blur'},
          // ],
          expression:[
            {required: true, message: '请输入告警过滤规则', trigger: 'blur'},
          ],
          // method:[
          //   {required: true, message: '请输入告警通知媒介', trigger: 'blur'},
          // ],
        },
        formRules: {
          rule_labels: [
            {required: true, message: '请输入告警组名称', trigger: 'blur'},
          ],
        },
      }
    },
    beforeMount() {
      this.getPromList();
      this.getPlanList();
      this.getUserGroupList();
      this.getMethodList();
    },
    mounted: function () {
      this.getPageList();
    },
    methods: {
      // 创建规则
      showCreateRule: function() {
        this.dialogFormTitle = "添加告警组"
        this.dialogFormVisible = true
        this.dialogFormUrl = "/monitor/addPlan"
      },
      createPlan: function(formName){
        this.$refs[formName].validate(async (valid) => {
          if (valid) {
            this.isSubmiting = true;
            try{
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
      showUpdatePlan: function(evl) {
        this.editdialogFormTitle = "编辑告警组";
        this.editdialogFormVisible = true;
        this.editdialogFormUrl = "/monitor/updatePlan";
        this.formEdit = deepCopy(evl);
      },
      //更新告警策略
      showUpdReceiver: function(id,evl) {
        this.editplanReceiverFormTitle = "编辑告警规则";
        this.editplanReceiverFormUrl = "/monitor/updateReceiver/" + id;
        this.editplanReceiverFormVisible = true;
        this.formPlanReceiverEdit = deepCopy(evl);
        this.changeSelect(evl);
        console.log("");
      },
      // 创建PlanReceiver
      showPlanReceiver: function(evl) {
        this.planReceiverFormVisible = true;
        this.planReceiverFormVisibleTitle = "添加告警策略";
        this.planReceiverFormUrl = "/monitor/addReceiver/"+evl; //添加告警策略,需要明确是加到那个告警组下;此处根据group的ID
      },
      //  创建告警通知
      createReceiver: function(formName) {
        this.$refs[formName].validate(async (valid) => {
          if (valid) {
            this.isSubmiting = true;
            try {
              this.formPlanReceiver.group = this.formPlanReceiver.group.toString(); //转字符串
              this.formPlanReceiver.call_url = this.formPlanReceiver.call_url.toString(); //转字符串              
              await request.post(this.planReceiverFormUrl, this.formPlanReceiver)
              this.pageCurrentChange(1);
              this.planReceiverFormVisible = false;
              this.$message.success("提交成功!!");
              this.formPlanReceiver = {}; 
            } catch(e) {
            }finally {
              this.isSubmiting = false;
            }
          } else {
            console.log('error submit!!');
            return false;
          }
        });
      },
      // 编辑告警策略
      editReceiver: function(formName) {
        this.$refs[formName].validate(async (valid) => {
          if (valid) {
            this.isSubmiting = true;
            try {
              this.formPlanReceiverEdit.group = this.formPlanReceiverEdit.group.toString();
              this.formPlanReceiverEdit.call_url = this.formPlanReceiverEdit.call_url.toString();
              await request.post(this.editplanReceiverFormUrl, this.formPlanReceiverEdit);
              this.pageCurrentChange(1);
              this.editplanReceiverFormVisible = false;
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
      // 编辑告警组
      editPlan: function(formName) {
        this.$refs[formName].validate(async (valid) => {
          if (valid) {
            this.isSubmiting = true;
            try {
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
      //  删除告警组
      delPlan(evl) {
        this.$confirm('是否要删除告警组:\n'+evl.rule_labels+' ? 删除告警组会将所有子计划一起删除,且不可找回!!!!', '提示', {
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
          await request.post('/monitor/deletePlan', evl);
          this.pageCurrentChange(1);
          this.$message.success("操作成功！");
        }).catch(() => {
        });
      },
      // 处理数据回显
      changeSelect(data) {
        if ( data.call_url !== '') {
          let Methods = data.call_url.toString();
          let methData = Methods.split(',');          
          this.formPlanReceiverEdit.call_url = methData;
        } else {
          this.formPlanReceiverEdit.call_url = [];
        };
        if (data.group !== '') {
          let UserIds = data.group.toString();
          let peodata = UserIds.split(',');          
          this.formPlanReceiverEdit.group = peodata;
        } else {
          this.formPlanReceiverEdit.group = [];
        }        
      },

      // 删除告警策略
      delReceiver(id,evl) {
        this.$confirm('是否要删除策略? 删除策略不可找回!!!!', '提示', {
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
          await request.post('/monitor/deleteReceiver/' +id, evl);
          this.pageCurrentChange(1);
          this.$message.success("操作成功！");
        }).catch(() => {
        });
      },
      // 获取用户组
      async getUserGroupList() {
        let res = await request.get('/monitor/getGroups');
        this.userList = res.data;
      },
      // 获取告警媒介
      async getMethodList() {
        let res = await request.get('/monitor/getMethods');
        this.methodList = res.data;
      },
      // 展开表格事件
      expandChange(row, expandRows) {
        if (!row.childrenData) {
          //通过$set设置loading实现动态加载表格
          this.$set(row, 'loading',true)
        }
      },

      //  分页获取
      async getPageList() {
        this.loadingTable=true;
        let filters = deepCopy(this.queryData);
        this.tableData = await request.get('/monitor/getPlans', filters);
        this.loadingTable=false;
      },
      // 获取prom数据源
      async getPromList() {
        let res = await request.get('/monitor/getProm');
        this.prom_id_list = res.data;
      },
      // 获取告警计划源
      async getPlanList() {
        let res = await request.get('/monitor/getPlan');
        this.plan_id_list = res.data;
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
  .plan .el-dialog__body .descClass .el-input__inner {
    width: 220px !important;
  }
  .buton-css {
    cursor: pointer;
    color:#1e90ff;
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
  .plan {
    width: 100%;
    height: 100%;
    padding: 10px;
    box-sizing: border-box;
    display: flex;
    flex-wrap: wrap;
  }
  .plan .el-dialog__body .el-form-item{
    margin-right: 0;
    width: 49%;
  }
  .plan .el-dialog__body .el-form-item .el-input__inner {
    width: 260px;
  }

  .plan .demo-table-expand {
    font-size: 0;
  }

  .plan .demo-table-expand label {
    width: 100px;
    color: #99a9bf;
  }

  .plan .demo-table-expand .el-form-item {
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