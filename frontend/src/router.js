import { createRouter, createWebHashHistory } from 'vue-router';
import Home from './Home.vue';
import Translate from './Translate.vue';

export default createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', component: Home },
    { path: '/translate', component: Translate },
  ],
});
