import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import BasicLayout from '../layouts/BasicLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    {
      path: '/admin',
      component: BasicLayout, // 这里使用布局组件
      // 所有的后台页面都放在 children 里
      children: [
        {
          path: 'dashboard', // 访问路径 /admin/dashboard
          name: 'dashboard',
          component: () => import('../views/DashboardView.vue')
        },
        // 先占位，防止菜单点击报错，后面我们再创建这些文件
        {
          path: 'questions',
          name: 'question-list',
          component: () => import('../views/QuestionListView.vue')
        },
        {
          path: 'questions/create',
          name: 'question-create',
          component: () => import('../views/QuestionCreateView.vue')
        },
        {
          path: 'questions/edit/:id',
          name: 'question-edit',
          component: () => import('../views/QuestionCreateView.vue')
        },
{
          path: 'announcements',
          name: 'announcements',
          component: () => import('../views/AnnouncementView.vue')
        },
         {
          path: 'users',
          name: 'users',
          component: () => import('../views/DashboardView.vue') // 暂时先用 Dashboard 顶替
        }
      ]
    }
  ]
})

export default router