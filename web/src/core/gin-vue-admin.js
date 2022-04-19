/*
 * kubespace web框架组
 *
 * */
// 加载网站配置文件夹
import { register } from './global'

export default {
  install: (app) => {
    register(app)
    console.log(`
       欢迎使用 Kubespace
       当前版本: v1.0.0
       默认swagger文档地址: http://127.0.0.1:${import.meta.env.VITE_SERVER_PORT}/swagger/index.html
       默认前端运行地址: http://127.0.0.1:${import.meta.env.VITE_CLI_PORT}
       默认后端运行地址: http://127.0.0.1:8888
    `)
  }
}
