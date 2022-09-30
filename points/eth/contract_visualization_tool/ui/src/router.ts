import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "home",
    component: () => import("@/views/home.vue"),
    children: [
      {
        path: "/eth/dependence",
        name: "ethDependence",
        component: () => import("@/views/eth/dependence.vue")
      },
      {
        path: "/eth/deploy",
        name: "ethDeploy",
        component: () => import("@/views/eth/deploy.vue")
      },
      {
        path: "/eth/invoke",
        name: "ethManage",
        component: () => import("@/views/eth/manage.vue")
      },
      {
        path: "/cashbox-controller/dependence",
        name: "ccDependence",
        component: () => import("@/views/cashbox_controller/dependence.vue")
      },
      {
        path: "/cashbox-controller/deploy-d",
        name: "ccDeployD",
        component: () => import("@/views/cashbox_controller/deploy_d.vue")
      },
      {
        path: "/cashbox-controller/deploy-c",
        name: "ccDeployC",
        component: () => import("@/views/cashbox_controller/deploy_c.vue")
      },
      {
        path: "/cashbox-controller/invoke",
        name: "ccManage",
        component: () => import("@/views/cashbox_controller/manage.vue")
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'hash',
  base: process.env.BASE_URL,
  routes
})

export default router
