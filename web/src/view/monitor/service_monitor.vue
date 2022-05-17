<template>
  <div class="serviceMonitor">
    <div class="yamlInput">
      <codemirror v-model="code" :options="cmOption" @input="onCmCodeChange" />
    </div>
    <div class="item-footer">
      <div class="item-cluster">
        请选择集群: 
        <el-select v-model="cluster" filterable placeholder="全部" style="width: 120px;">
          <el-option
            v-for="key in Object.keys(cls_list)"
            :key="key"
            :label="cls_list[key]"
            :value="cls_list[key]">
          </el-option>
        </el-select>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="reStore">重 置</el-button>
        <el-button type="primary" @click="applyMonit"><span>创 建</span></el-button>
      </div>
    </div>
  </div>
</template>

<script>
  import dedent from 'dedent';
  import {request} from "../../utils/rquestes";
  import { codemirror } from 'vue-codemirror';
  // language
  import 'codemirror/mode/yaml/yaml.js';
  // theme css
  import 'codemirror/theme/eclipse.css';
  // require active-line.js
  import'codemirror/addon/selection/active-line.js';
  // autoCloseTags
  import'codemirror/addon/edit/closetag.js';
  
  export default {
    name: 'codemirror-example-html',
    title: 'Mode: text/x-yaml & Theme: eclipse',
    components: {
      codemirror
    },
    data() {
      return {
        code: dedent`
kind: ServiceMonitor
metadata:
    annotations:
    labels:
      k8s-app: mk-qw-account-svc
    name: mk-qw-account-svc-monitor
    namespace: monitoring
  spec:
    endpoints:
    - interval: 30s
    port: http
    namespaceSelector:
      matchNames:
      - mk-release
    selector:
      matchLabels:
        svc: mk-qw-account-svc
        `,
        cluster: null,
        cls_list: [],
        cmOption: {
          tabSize: 4,
          styleActiveLine: true,
          lineNumbers: true,
          autoCloseTags: true,
          matchBrackets: true,
          lineWrapping: true,
          line: true,
          indentUnit: 1,
          mode: 'text/x-yaml',
          theme: 'eclipse'
        }
      }
    },
    beforeMount() {
      this.getClusterList();      
    },
    methods: {
      onCmCodeChange(newCode) {
        console.log('this is new code', 11)
        this.code = newCode
      },
      reStore() {
        this.code = '';
      },
      async applyMonit() {
        if (this.checkYaml(this.code)) {
          if (this.cluster == null) {
            this.$message.error("请选择集群！");
            return
          }
            // 正确的yaml, 提交到后端
          let Base64 = require('js-base64').Base64;
          let res = Base64.encode(this.code);
          try{
            let param = {"config":res, "name": this.cluster};
            await request.post("/monitor/addClusterResource", param);
            this.$message.success("提交成功！");
          }catch (e) {
            // this.$message.error("提交失败！");
            console.log(e);
          }finally {
          }
        } else {
          console.log("state is false");
        }
      },
      // yaml 检查函数
      checkYaml(cm) {
        let state = false;
        try {
          let jsyaml = require("js-yaml");
          state = !!jsyaml.load(cm);
        } catch(e) {
          console.log(e);
          this.$message.error("yaml格式不正确,请检查格式!");
          return
        }
        return state
      },

      // 从后端获取可用集群列表
      async getClusterList() {
        let res = await request.get("/monitor/getClusterList");
        this.cls_list = res.data;
      },
    },
    computed: {
      codemirror() {
        return this.$refs.cmEditor.codemirror
      }
    },
  }
</script>

<style scope>
  .CodeMirror * {
    font-family: monospace;
    font-size: 16px;
  }
  .yamlInput {
    width: 100%;
    height: 80%;
  }
  .CodeMirror {
    height: 60% !important;
  }
  .item-footer {
    margin-top: 50px;

  }
  .dialog-footer  {
    float: left;
    margin-left: 30px;
  }
  .item-cluster {
    float: left;    
  }
</style>
