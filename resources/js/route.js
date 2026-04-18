import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    name: "home",
    component: () => import("./pages/Home.vue")
  },
  {
    path: "/auth/login",
    name: "login",
    component: () => import("./pages/auth/Login.vue")
  },
  {
    path: "/auth/register",
    name: "register",
    component: () => import("./pages/auth/Register.vue")
  },
  {
    path: "/admin",
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
        component: () => import("./pages/admin/Users.vue")
      },
      {
        path: "roles",
        name: "admin.roles",
        component: () => import("./pages/admin/Roles.vue")
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

export default router;
