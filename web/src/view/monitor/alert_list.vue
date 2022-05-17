<template>
  <div class="alert_list">
    <div class="search-bar">
      <el-form :model="queryData" ref="searchForm" :inline="true">
        <el-form-item label="ID">
          <el-input v-model="queryData.id" placeholder="ID" style="width: 80px;"></el-input>
        </el-form-item>      
        <el-form-item label="标题">
          <el-input v-model="queryData.summary" placeholder="标题" style="width: 100px;"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryData.status" filterable placeholder="全部" style="width: 100px;" @change="selQueryStatusChange">
            <el-option
              v-for="key in Object.keys(alert_status)"
              :key="key"
              :label="alert_status[key]"
              :value="key">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="触发时间">
          <el-date-picker
            style="width:210px;"
            :editable="false"
            v-model="queryData.start_time"
            :picker-options="startDatePicker"
            type="datetime"
            value-format="yyyy-MM-dd HH:mm:ss"
            format="yyyy-MM-dd HH:mm:ss"
            placeholder="选择开始日期" />
        </el-form-item>
        <el-form-item>
          <div>~</div>
        </el-form-item>
        <el-form-item >
          <el-date-picker
            style="width:210px;"
            :editable="false"
            v-model="queryData.end_time"
            :picker-options="endDatePicker"
            type="datetime"
            value-format="yyyy-MM-dd HH:mm:ss"
            format="yyyy-MM-dd HH:mm:ss"
            default-time="['23:59:59']"
            placeholder="选择结束日期"
            />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="pageCurrentChange(1)" size="small">查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-card class="table-card">
      <div v-loading="loadingTable">
        <el-table
          tooltip-effect="dark"
          :height="this.pageHeight-200"
          ref="singleTable"
          highlight-current-row
          border stripe
          :sort-orders="['ascending', 'descending']"
          :data="tableData.data">
          <el-table-column
            show-overflow-tooltip
            label="ID"
            prop="id"
            width="80">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="Rule ID"
            prop="rule_id"
            width="60">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="当前状态"
            width="80">
            <template slot-scope="scope">
                <p>{{ alert_status[scope.row.status] }}</p>
            </template>
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="异常分钟数"
            width="100"
            prop="count">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            prop="summary"
            label="summary"
            min-width="150"
            width="auto">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            prop="title"
            label="标题"
            min-width="150"
            width="auto">
          </el-table-column>          
          <el-table-column
            show-overflow-tooltip
            label="label"
            width="auto">
            <template slot-scope="scope">
                <div v-if="scope.row.labels && Object.prototype.toString.call(scope.row.labels) == '[object Object]'">
                  <div v-for="key in Object.keys(scope.row.labels)">
                    <el-tag type="success" effect="dark" style={{ marginTop: '5px' }}>
                      {{key}}:{{scope.row.labels[key]}}
                    </el-tag>
                  </div>
                </div>
                <div v-else>
                  <span>-</span>
                </div>
            </template>            
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="描述"
            prop="description"
            width="auto">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            prop="confirmed_by"
            label="确认人"
            width="auto">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            prop="fired_at"
            label="触发时间"
            width="160">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            prop="confirmed_at"
            label="确认时间"
            width="150">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="确认截止时间"
            prop="confirmed_before"
            width="auto">
          </el-table-column>
          <el-table-column
            show-overflow-tooltip
            label="恢复时间"
            prop="resolved_at"
            width="auto">
          </el-table-column>
        </el-table>
        <!-- background -->
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
  </div>
</template>

<script>
import  {AlERT_STATUS} from "../../utils/vm-cfg";
import {deepCopy, isEmpty} from "../../utils/common";
import {request} from "../../utils/rquestes";

export default {
  name: "alert_list",
  data: function () {
    return {
      //时间选择变量
      startDatePicker: this.beginDate(),
      endDatePicker: this.processDate(),
      loadingTable: false,
      pageHeight: document.body.scrollHeight,
      tableData: [],
      queryData: {
        id: null,
        summary: null,
        status: null,
        start_time: null,
        end_time: null,
        current_page: 1,
        page_size: 20,
      },
      alert_status: AlERT_STATUS,
      queryEnvData: [],
      queryZoneData: [],  
    }
  },
  mounted: function () {
    this.getPageList();
  },
  methods: {
    beginDate() {
      const self = this;
      return {
        disabledDate(time) {
          if (self.end_time) {
            //如果结束时间不为空，则小于结束时间
            return new Date(self.end_time).getTime() < time.getTime();
          } else {
            // return time.getTime() > Date.now()//开始时间不选时，结束时间最大值小于等于当天
          }
        }
      };
    },
    processDate() {
      const self = this;
      return {
        disabledDate(time) {
          if (self.start_time) {
            //如果开始时间不为空，则结束时间大于开始时间
            return new Date(self.start_time).getTime() > time.getTime();
          } else {
            // return time.getTime() > Date.now()//开始时间不选时，结束时间最大值小于等于当天
          }
        }
      };
    },
    // 查询方法
    pageCurrentChange(val) {
      this.queryData.current_page = val;
      this.getPageList();
    },
    async getPageList() {
        this.loadingTable=true;
        let filters = deepCopy(this.queryData);
        if(filters["status"] === "全部"){
          filters["status"] = '';
        }
        this.tableData = await request.get('/monitor/getAlerts', filters);
        this.loadingTable=false;
    },
    selQueryStatusChange(selVal) {
      if (selVal === "全部") {
        selVal = null;
      }
      this.queryData.status = selVal;
      this.pageCurrentChange(1);
    },
    //分页
    pageSizeChange(val) {
      this.queryData.page_size = val;
      this.getPageList();
    },    
  }
}
</script>

<style scope>
  .alert_list {
    width: 100%;
    height: 100%;
    padding: 10px;
    box-sizing: border-box;
    display: flex;
    flex-wrap: wrap;
  }
  .el-pagination {
      text-align: center; 
  }  
</style>