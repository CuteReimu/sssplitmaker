import {createApp} from 'vue'
import App from './App.vue'
import router from './router.js'
import './style.css';
import 'element-plus/dist/index.css';
import 'element-plus/theme-chalk/dark/css-vars.css';

createApp(App).use(router).mount('#app')
