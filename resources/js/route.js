import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    component: () => import("./pages/PublicLayout.vue"),
    children: [
      {
        path: "",
        name: "home",
        component: () => import("./pages/auth/Home.vue")
      },
      {
        path: "auth/login",
        name: "login",
        component: () => import("./pages/auth/Login.vue")
      },
      {
        path: "auth/register",
        name: "register",
        component: () => import("./pages/auth/Register.vue")
      },
      {
        path: "auth/forgot-password",
        name: "forgot_password",
        component: () => import("./pages/auth/ForgotPassword.vue")
      },
      {
        path: "auth/unauthenticated",
        name: "unauthenticated",
        component: () => import("./pages/auth/Unauthenticated.vue")
      },
      {
        path: "auth/unauthorized",
        name: "unauthorized",
        component: () => import("./pages/auth/Unauthorized.vue")
      },
      {
        path: "auth/not-found",
        name: "not_found",
        component: () => import("./pages/auth/NotFound.vue")
      },
      {
        path: ":pathMatch(.*)*",
        name: "notfound",
        component: () => import("./pages/auth/NotFound.vue")
      }
    ]
  },
  {
    path: "/admin",
    meta: { requiresAuth: true },
    component: () => import("./pages/Layout.vue"),
    children: [
      {
        path: "",
        name: "admin.dashboard",
        component: () => import("./pages/admin/Dashboard.vue")
      },
      {
        path: "users",
        name: "admin.users",
        component: () => import("./pages/users/Index.vue")
      },
      {
        path: "users/new",
        name: "admin.users.create",
        component: () => import("./pages/users/AddCreate.vue")
      },
      {
        path: "users/:id/edit",
        name: "admin.users.edit",
        component: () => import("./pages/users/AddCreate.vue")
      },
      {
        path: "roles",
        name: "admin.roles",
        component: () => import("./pages/roles/Index.vue")
      },
      {
        path: "roles/new",
        name: "admin.roles.create",
        component: () => import("./pages/roles/AddEdit.vue")
      },
      {
        path: "roles/:id/edit",
        name: "admin.roles.edit",
        component: () => import("./pages/roles/AddEdit.vue")
      },
      {
        path: "posts",
        name: "admin.posts",
        component: () => import("./pages/posts/Index.vue")
      },
      {
        path: "posts/new",
        name: "admin.posts.create",
        component: () => import("./pages/posts/AddCreate.vue")
      },
      {
        path: "posts/:id/edit",
        name: "admin.posts.edit",
        component: () => import("./pages/posts/AddCreate.vue")
      },
      {
        path: "settings",
        name: "admin.settings",
        component: () => import("./pages/admin/Settings.vue")
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
});

function getToken() {
  try {
    return localStorage.getItem("token") || "";
  } catch {
    return "";
  }
}

router.beforeEach((to) => {
  const requiresAuth = to.matched.some((record) => record.meta && record.meta.requiresAuth);
  if (!requiresAuth) return true;

  const hasToken = getToken().trim().length > 0;
  if (hasToken) return true;

  return {
    name: "login",
    query: {
      redirect: to.fullPath
    }
  };
});

export default router;
