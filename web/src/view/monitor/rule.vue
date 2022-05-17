<template>
  <div class="cluster">
    <div class="search-bar">
      <el-form :model="queryData" ref="searchForm" :inline="true">
        <el-form-item> <!--class="add-cluster-button-import" -->
          <el-button type="primary" @click="showCreateRule()" size="small">新建</el-button>
        </el-form-item>
        <el-form-item label="环境">
          <el-select v-model="queryData.prom_id" filterable placeholder="全部" style="width: 100px;" @change="selQueryStatusChange">
            <el-option
              v-for="(key,index) in prom_id_list"
              :key="index"
              :label="key.name"
              :value="key.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input class="title-input" v-model="queryData.summary" placeholder="输入规则标题" style="width: 200px;"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button class="title-search" type="primary" icon="el-icon-search" @click="pageCurrentChange(1)" size="small">搜索</el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-card class="table-card">
      <div v-loading="loadingTable">
        <el-table
          ref="singleTable"
          :height="this.pageHeight-200"
          highlight-current-row
          border stripe
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
            label="标题"
            width="80"
            align="center"
            prop="title">
          </el-table-column>          
          <el-table-column
            show-overflow-tooltip
            label="表达式"
            align="center"
            width="330">
            <template slot-scope="scope">
              <div>{{scope.row.expr}} {{scope.row.op}} {{scope.row.value}}</div>
            </template>
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="持续时间"
            align="center"
            width="180"
            prop="for">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="告警类别"
            align="center"
            width="80"
            prop="monitor_type">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="告警级别"
            align="center"
            width="80"
            prop="alarm_level">
          </el-table-column> 
          <el-table-column
            show-overflow-tooltip
            label="Summary"
            align="center"
            width="200"
            prop="summary">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="描述"
            align="center"
            width="auto"
            prop="description">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="Prometheus数据源"
            align="center"
            width="200">
            <template slot-scope="scope">
              <div>{{prom_id_list[scope.row.prom_id].name}}</div>
            </template>
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="告警计划"
            align="center"
            width="180">
            <template slot-scope="scope">
              <div>
                {{plan_id_list[scope.row.plan_id].rule_labels}}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            align="center"
            width="140">
            <template slot-scope="scope">
              <div>
                <el-button type="primary" @click="showUpdateRule(scope.row)" size="mini">
                  编辑<i class=" el-icon--right"></i>
                </el-button>
                <el-button type="primary" @click="delRules(scope.row)" size="mini">
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
          @current-change="pageCurrentChange"
          :current-page="queryData.current_page"
          :page-sizes="[10, 20, 30, 40, 50]"
          :page-size="queryData.page_size"
          layout="total ,sizes, prev, pager, next, jumper"
          :total="tableData.total">
        </el-pagination>
      </div>
    </el-card>

    <!-- 创建告警规则 -->
    <el-dialog :title="dialogFormTitle" :visible.sync="dialogFormVisible" width="450px" :close="handCreateDialogClose">
      <el-form :model="formFields" ref="formFields" :rules="formRules" label-width="100px" class="demo-form-inline">
        <el-form-item label="Title" prop="title" >
          <el-input v-model.trim="formFields.title" placeholder="输入监控标题"  />
        </el-form-item>      
        <el-form-item label="监控指标" prop="expr" >
          <el-input v-model.trim="formFields.expr" placeholder="输入监控metric;例如:node_load1"  />
        </el-form-item>
        <el-form-item label="监控阈值" class="descClass">
          <el-select style="width:40px" v-model="formFields.op" filterable placeholder=">">
            <el-option
              v-for="(key,index) in operator"
              :key="index"
              :label="key"
              :value="key">
            </el-option>
          </el-select>
          <el-input style="width: 0px" v-model.trim="formFields.value" placeholder="输入阈值"  />
        </el-form-item>
        <el-form-item label="Label">
          <div v-module="formFields.labels">
            <div v-for="(module,index) in modules" :key="index">
              <span class="descClass1">
                <el-input style="width:80px" v-model.trim="module.name"></el-input>
                <el-input style="width:0px" v-model.trim="module.value"></el-input>
              </span> 
            </div>
            <el-button type="primary" @click="add">添加标签</el-button>
          </div>
        </el-form-item>
        <el-form-item label="告警类别" prop="monitor_type" >
          <el-input v-model.trim="formFields.monitor_type" placeholder="告警分类"   />
        </el-form-item>
        <el-form-item label="告警级别" prop="alarm_level" >
          <el-input v-model.trim="formFields.alarm_level" placeholder="告警级别"   />
        </el-form-item>       
        <el-form-item label="持续时间" prop="for" >
          <el-input v-model.trim="formFields.for" placeholder="持续时间,默认单位:s"   />
        </el-form-item>
        <el-form-item label="Summary" prop="summary" >
          <el-input v-model.trim="formFields.summary" placeholder="规则标题,该标题会在告警时发出." />
        </el-form-item>
        <el-form-item label="描述" prop="description" >
          <el-input v-model.trim="formFields.description" placeholder="请输入规则描述信息." />
        </el-form-item>
        <el-form-item label="数据源" prop="prom_id">
          <el-select v-model="formFields.prom_id" filterable placeholder="请选择">
            <el-option
              v-for="(key,index) in prom_id_list"
              :key="index"
              :label="key.name"
              :value="key.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="告警计划" prop="plan_id">
          <el-select v-model="formFields.plan_id" filterable placeholder="请选择">
            <el-option
              v-for="(key,index) in plan_id_list"
              :key="index"
              :label="key.rule_labels"
              :value="key.id">
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="createRule('formFields')" :disabled="isSubmiting" :loading="isSubmiting"><span v-show="!isSubmiting">确 定</span><span v-show="isSubmiting">创建中</span></el-button>
      </div>
    </el-dialog>

    <!--更新告警规则 -->
    <el-dialog :title="editdialogFormTitle" :visible.sync="editdialogFormVisible" width="450px" :close="handEditDialogClose">
      <el-form :model="formEdit" ref="formEdit" :rules="formRules" label-width="100px" class="demo-form-inline">
        <el-form-item label="监控标题" prop="title" >
          <el-input v-model.trim="formEdit.title" placeholder="输入监控标题" />
        </el-form-item>      
        <el-form-item label="监控指标" prop="expr" >
          <el-input v-model.trim="formEdit.expr" placeholder="输入监控metric;例如:node_load1"  />
        </el-form-item>
        <el-form-item label="监控阈值" class="descClass">
          <el-select style="width:40px" v-model="formEdit.op" filterable placeholder=">">
            <el-option
              v-for="(key,index) in operator"
              :key="index"
              :label="key"
              :value="key">
            </el-option>
          </el-select>
          <el-input style="width: 0px" v-model.trim="formEdit.value" placeholder="输入阈值"  />
        </el-form-item>
        <el-form-item label="Label">
          <div v-model="formEdit.labels">
            <div v-for="(module,index) in modules" :key="index">
              <span class="descClass1">
                <el-input style="width:80px" v-model.trim="module.name"></el-input>
                <el-input style="width:0px" v-model.trim="module.value"></el-input>
              </span> 
            </div>
            <el-button type="primary" @click="add">添加标签</el-button>
          </div>
        </el-form-item>
        <el-form-item label="告警类别" prop="monitor_type" >
          <el-input v-model.trim="formEdit.monitor_type" placeholder="告警分类"   />
        </el-form-item>
        <el-form-item label="告警级别" prop="alarm_level" >
          <el-input v-model.trim="formEdit.alarm_level" placeholder="告警级别"   />
        </el-form-item>            
        <el-form-item label="持续时间" prop="for" >
          <el-input v-model.trim="formEdit.for" placeholder="持续时间,默认单位:s"   />
        </el-form-item>
        <el-form-item label="监控标题" prop="summary" >
          <el-input v-model.trim="formEdit.summary" placeholder="规则标题,该标题会在告警时发出." />
        </el-form-item>
        <el-form-item label="描述" prop="description" >
          <el-input v-model.trim="formEdit.description" placeholder="请输入规则描述信息." />
        </el-form-item>
        <el-form-item label="数据源" prop="prom_id">
          <el-select v-model="formEdit.prom_id" filterable placeholder="请选择">
            <el-option
              v-for="(key,index) in prom_id_list"
              :key="index"
              :label="key.name"
              :value="key.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="告警计划" prop="plan_id">
          <el-select v-model="formEdit.plan_id" filterable placeholder="请选择">
            <el-option
              v-for="(key,index) in plan_id_list"
              :key="index"
              :label="key.rule_labels"
              :value="key.id">
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="editdialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="editRules('formEdit')" :disabled="isSubmiting" :loading="isSubmiting"><span v-show="!isSubmiting">更 新</span><span v-show="isSubmiting">创建中</span></el-button>
      </div>
    </el-dialog>    
  </div>
</template>

<script>
  import {request} from "../../utils/rquestes";
  import {deepCopy, isEmpty} from "../../utils/common";
  export default {
    name: "rule_list",
    data: function () {
      return {
        loadingTable: true,
        pageHeight: document.body.scrollHeight,
        dialogFormVisible: false,
        dialogFormTitle: null,
        modules:[{name:'severity',value: 'critical'}],
        editdialogFormVisible: false,
        editdialogFormTitle: null,
        editdialogFormUrl: null,

        dialogFormUrl: null,
        isSubmiting: false,
        formFieldsAdd:{
          title: "",
          expr: "",
          op: "==",
          value: "",
          for: "0s",
          labels: {},
          summary: "",
          monitor_type: "",
          alarm_level: "",
          description: "",
          prom_id: "",
          plan_id: "",
        },
        formFields: {},
        formEdit: {},
        prom_id_list: [],
        plan_id_list: [],
        operator: ["==",">","<",">=","<=","!="],
        queryData: {
          current_page: 1,
          page_size: 20,
        },
        tableData: {},
        formRules: {
          title: [
            {required: true, message: '请输入告警标题', trigger: 'blur'},
          ],
          expr: [
            {required: true, message: '请输入rule规则', trigger: 'blur'},
          ],
          op: [
            {required: true, message: '请输入value比较方式', trigger: 'blur'},
          ],
          value: [
            {required: true, message: '请输入告警阈值', trigger: 'blur'},
          ],
          for: [
            {required: true, message: '请输入告警持续时间,单位s', trigger: 'blur'},
          ],
          monitor_type: [
            {required: true, message: '请输入告警类别', trigger: 'blur'},
          ],
          alarm_level: [
            {required: true, message: '请输入告警级别', trigger: 'blur'},
          ],
          summary: [
            {required: true, message: '请输入告警标题', trigger: 'blur'},
          ],
          prom_id: [
            {required: true, message: '请选择prometheus数据源.', trigger: 'blur'},
          ],
          plan_id: [
            {required: true, message: '请选择告警计划.', trigger: 'blur'},
          ],
        },
      }
    },
    beforeMount() {
      this.getPromList();
      this.getPlanList();
    },
    mounted: function () {
      this.getPageList();
      // for(var item of this.list) {
      //   this.modules.push({text:item});
      // };
    },
    methods: {
      add() {
        this.modules.push({name:'',value: ''});  // 每点一下，push一次
      },
      // 创建规则
      showCreateRule: function() {
        this.dialogFormTitle = "添加监控规则"
        this.dialogFormVisible = true
        this.dialogFormUrl = "/monitor/addRules"
        this.formFields = deepCopy(this.formFieldsAdd);
      },
      createRule: function(formName){
        this.$refs[formName].validate(async (valid) => {
          if (valid) {
            this.isSubmiting = true;
            try{
              let Param = {};
              for (var item of this.modules) {
                if (item.name !== "" && item.value !== "") {
                  Param[item.name] = item.value;
                };
              };
              this.formFields.labels = Param;
              await request.post(this.dialogFormUrl, this.formFields);
              this.pageCurrentChange(1);
              this.dialogFormVisible = false;
              this.$message.success("提交成功！");
            }catch (e) {
              console.log(e);
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
      showUpdateRule: function(evl) {
        this.editdialogFormTitle = "编辑监控规则";
        this.editdialogFormVisible = true;
        this.editdialogFormUrl = "/monitor/updateRules";

        // 处理label
        let Lab = JSON.parse(evl.labels);
        let LabList = [];
        for (var key in Lab) {
          let map = {name:'',value: ''};
          map.name = key;
          map.value = Lab[key];
          LabList.push(map);
        };
        this.modules = LabList;
        this.formEdit = deepCopy(evl);
      },
      editRules: function(formName) {
        this.$refs[formName].validate(async (valid) => {
          if (valid) {
            this.isSubmiting = true;
            try {
              let Param = {};
              for (var item of this.modules) {
                if (item.name !== "" && item.value !== "") {
                  Param[item.name] = item.value;
                };
              };
              this.formEdit.labels = Param;    
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
      delRules(evl) {
        this.$confirm('是否要删除规则:\n'+evl.expr+' ? 删除规则不可找回!!!!', '提示', {
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
          // 处理label
          let Lab = JSON.parse(evl.labels);
          evl.labels = Lab;
          await request.post('/monitor/deleteRules', evl);
          this.pageCurrentChange(1);
          this.$message.success("操作成功！");
        }).catch(() => {
        });
      },
      //  分页获取
      async getPageList() {
        this.loadingTable=true;
        let filters = deepCopy(this.queryData);
        this.tableData = await request.get('/monitor/getRules', filters);
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
      // 根据数据源筛选告警.
      selQueryStatusChange(selVal) {
        if (selVal === "全部") {
          selVal = null;
        }
        this.queryData.prom_id = selVal;
        this.pageCurrentChange(1);
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

  .cluster .el-dialog__body .descClass1 .el-input__inner {
    width: 100px !important;
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