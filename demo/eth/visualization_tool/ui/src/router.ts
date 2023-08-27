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
        path: "/eth/library",
        name: "ethLibrary",
        component: () => import("@/views/eth/library.vue")
      },
      {
        path: "/eth/contract",
        name: "ethDeploy",
        component: () => import("@/views/eth/contract.vue")
      },
      {
        path: "/eth/invoke",
        name: "ethInvoke",
        component: () => import("@/views/eth/invoke.vue")
      },
      {
        path: "/cashbox/library",
        name: "cashboxLibrary",
        component: () => import("@/views/cashbox/library.vue")
      },
      {
        path: "/cashbox/contract-data",
        name: "cashboxDataContract",
        component: () => import("@/views/cashbox/contract_data.vue")
      },
      {
        path: "/cashbox/contract-controller",
        name: "cashboxControllerContract",
        component: () => import("@/views/cashbox/contract_controller.vue")
      },
      {
        path: "/cashbox/invoke",
        name: "cashboxInvoke",
        component: () => import("@/views/cashbox/invoke.vue")
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
