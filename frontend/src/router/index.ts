import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import UserLayout from '../layouts/UserLayout.vue'
import BasicLayout from '../layouts/BasicLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    // ❌ 之前错误的位置：在这里删掉 filter 路由
    {
      path: '/practice',
      component: UserLayout,
      children: [
        {
          path: '',
          name: 'practice-home',
          component: () => import('../views/user/UserHomeView.vue')
        },
        {
          path: 'mode',
          name: 'practice-mode',
          component: () => import('../views/user/PracticeModeView.vue')
        },
        // ✅ 正确的位置：放在 children 里
        {
          path: 'filter',
          name: 'practice-filter',
          component: () => import('../views/user/PracticeFilterView.vue')
        },
        {
          path: 'exam/:mode',
          name: 'practice-exam',
          component: () => import('../views/user/ExamView.vue')
        },
        {
          path: 'exam/config',
          name: 'exam-config',
          component: () => import('../views/user/ExamConfigView.vue')
        },
        {
          path: 'exam/paper',
          name: 'exam-paper',
          component: () => import('../views/user/ExamPaperView.vue')
        },
        {
          path: 'mistakes',
          name: 'mistake-book',
          component: () => import('../views/user/MistakeBookView.vue')
        }
      ]
    },
    // 加在 /practice 下面或者单独一级
    {
      path: '/coop',
      component: UserLayout,
      children: [
        { path: '', name: 'coop-lobby', component: () => import('../views/coop/LobbyView.vue') },
        { path: 'room/:id', name: 'coop-room', component: () => import('../views/coop/RoomView.vue') },
        { path: 'room/:id/exam', name: 'coop-exam', component: () => import('../views/coop/CoopExamView.vue') },
        { path: 'room/:id/result', name: 'coop-result', component: () => import('../views/coop/CoopResultView.vue') }
      ]
    },
    {
      path: '/admin',
      component: BasicLayout,
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('../views/DashboardView.vue')
        },
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
          component: () => import('../views/UserListView.vue')
        }
      ]
    }
  ]
})

// 守卫逻辑保持不变...
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  const role = localStorage.getItem('role');

  if (to.name === 'login') {
    next();
    return;
  }

  if (!token) {
    next({ name: 'login' });
    return;
  }

  if (to.path.startsWith('/admin') && role !== 'admin') {
    next({ name: 'practice-home' });
    return;
  }

  next();
});

export default router;