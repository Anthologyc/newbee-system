import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ArcoVue from '@arco-design/web-vue';
import '@arco-design/web-vue/dist/arco.css'; // 引入样式

const app = createApp(App)

app.use(router)
app.use(ArcoVue) // 注册 Arco
app.mount('#app')