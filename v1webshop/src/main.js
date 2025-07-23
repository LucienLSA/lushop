import {createApp} from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import locale from 'element-plus/lib/locale/lang/zh-CN' // lang i18n
import VCharts from 'v-charts'
import 'normalize.css/normalize.css'
import '@/icons' // icon
import '@/permission' // permission control
import router from './router'

createApp(App).use(router).use(ElementPlus, {locale}).use(VCharts).mount('#app')
