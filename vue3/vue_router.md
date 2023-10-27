# vue router4

> 学习时间：2023.9
> 文档：https://router.vuejs.org/zh/guide/

## 动态路由匹配

指路由中带参数的情况，例如：  
路由规则：`/users/:username/posts/:postId`  
实际：`/users/eduardo/posts/123`  
`$route.params`：`{ username: 'eduardo', postId: '123' }`

使用动态路由导航时（例如`/users/:id`，从`/users/mario`导航到`/users/mats9693`），组件实例将被重复使用，而不是销毁-重新创建；  
这导致组件的声明周期钩子函数不会重新触发，如果你想要响应动态路由的导航，可以`watch $route.params`，或者使用**导航守卫**

## 捕获所有路由

用于捕获在此之前没有匹配成功的所有路由

```vue 
// 将匹配所有内容并将其放在 `$route.params.pathMatch` 下
{ path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
// 将匹配以 `/user-` 开头的所有内容，并将其放在 `$route.params.afterUser` 下
{ path: '/user-:afterUser(.*)', component: UserGeneric },
```

## 路由匹配语法

当我们写`/users/:id`时，`:id`的含义等于`([^/]+)`，即一直到下一个`/`之前的内容，  
这个正则是可以自定义的，例如仅匹配数字、重复匹配多次(`/:id+`可以匹配`/0/1/2`)、可选参数(`/:id?`)

路由默认**不区分大小写**，且**可以匹配末尾带有`/`或不带有`/`的情况**，可以使用`sensitive & strict`明确设置

[路由调试](https://paths.esm.dev/?p=AAMeJSyAwR4UbFDAFxAcAGAIJXMAAA..#)

## 编程式导航

导航到不同位置，可以拼动态参数，[参考](https://router.vuejs.org/zh/guide/essentials/navigation.html)

## 嵌套路由

`children`属性

子路由列表中，空的`path`将被默认渲染到`router view`的位置

## 命名路由和视图 view

一个路由需要有path、component，可以有name，name会带来好处：

1. 没有硬编码的 URL
2. params 的自动编码/解码。
3. 防止你在 url 中出现打字错误。
4. 绕过路径排序（如显示一个）

如果你想要在同一个层级展示多个视图(view)，可以先给它们起个名字：

```vue
<router-view class="view left-sidebar" name="LeftSidebar"></router-view>
<router-view class="view main-content"></router-view>
<router-view class="view right-sidebar" name="RightSidebar"></router-view>

{
  path: '/',
  components: {
    default: Home,
    // LeftSidebar: LeftSidebar 的缩写
    LeftSidebar,
    // 它们与 `<router-view>` 上的 `name` 属性匹配
    RightSidebar,
  },
},
```

命名的视图可以嵌套

## 重定向和别名

`redirect`/`alias`

重定向会将一个url导航到另一个url或路由，可以接收函数作为参数：

```vue
{ path: '/home', redirect: '/' }
{ path: '/home', redirect: { name: 'homepage' } }
{
  // /search/screens -> /search?q=screens
  path: '/search/:searchText',
  redirect: to => {
    // 方法接收目标路由作为参数
    // return 重定向的字符串路径/路径对象
    return { path: '/search', query: { q: to.params.searchText } }
  },
},
{
  // 将总是把/users/123/posts重定向到/users/123/profile。
  path: '/users/:id/posts',
  redirect: to => {
    // 该函数接收目标路由作为参数
    // 相对位置不以`/`开头
    // 或 { path: 'profile'}
    return 'profile'
  },
},
// 因为导航守卫仅应用在目标路由上，所以在以上例子中编写beforeEnter导航守卫不会生效
```

别名可以为一个路由匹配多个url

## 通过路由组件传参

`$route.params.id`，这样可以在组件中获取路由里的动态参数，然而使用$route会让组件和路由紧密耦合，可以使用prop：

```vue
const User = {
  template: '<div>User {{ $route.params.id }}</div>'
}
const routes = [{ path: '/user/:id', component: User }]

 ↓

const User = {
  // 请确保添加一个与路由参数完全相同的 prop 名
  props: ['id'],
  template: '<div>User {{ id }}</div>'
}
const routes = [{ path: '/user/:id', component: User, props: true }]
```

## history mode

hash模式、history模式

两种模式的主要区别是：history模式无法应对形如：`https://example.com/user/id` 的**直接**请求，需要服务端配合：

```nginx
location / {
    // root ...
    try_files $uri $uri/ /index.html;
}
```

此时服务端不会返回404，而是对所有路径都返回index.html，这需要程序里添加一个404页面：`{ path: '*', component: NotFoundComponent }`

1. 为什么history模式无法应对默认首页以外的直接请求？  
   `https://example.com/user/id` 发到服务端，服务端会去文件夹下面找`user/id`文件，找不到就返回404
2. try_files是什么意思？
   `try_files file ... uri;`，至少需要2个参数，会依次匹配root中对应的文件或文件夹；最后一个参数在前面全都没有匹配上的时候生效
   举个例子：`http://.../home (try_files $uri $uri/ /index.html;)`，`$uri`表示`/home`，  
   这个请求发到服务端，会先找home文件，然后去找home文件夹（第二个参数），最后都没找到，返回index.html

## 导航守卫

导航守卫，主要通过**跳转**和**取消**的方式守卫导航，其注册方式可以是全局的、单个路由的、单个组件的

[语法参考](https://router.vuejs.org/zh/guide/advanced/navigation-guards.html)

导航解析流程：

1. 导航被触发。
2. 在失活的组件里调用 beforeRouteLeave 守卫。
3. 调用全局的 beforeEach 守卫。
4. 在重用的组件里调用 beforeRouteUpdate 守卫(2.2+)。
5. 在路由配置里调用 beforeEnter。
6. 解析异步路由组件。
7. 在被激活的组件里调用 beforeRouteEnter。
8. 调用全局的 beforeResolve 守卫(2.5+)。
9. 导航被确认。
10. 调用全局的 afterEach 钩子。
11. 触发 DOM 更新。
12. 调用 beforeRouteEnter 守卫中传给 next 的回调函数，创建好的组件实例会作为回调函数的参数传入。

## 路由元参数 meta

```vue
{
  path: 'new',
  component: PostsNew,
  // 只有经过身份验证的用户才能创建帖子
  meta: { requiresAuth: true }
},
```

在ts中使用meta参数，需要声明：

```vue
// typings.d.ts or router.ts
import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    // 是可选的
    isAdmin?: boolean
    // 每个路由都必须声明
    requiresAuth: boolean
  }
}
```

## 数据获取

数据获取的时间节点：

1. 导航完成前：在导航守卫中获取数据，然后执行导航
2. 导航完成后：在组件的声明周期钩子函数中获取数据

两种获取方式主要的区别在于用户体验

## 组合式api

通常我们使用this.$router/this.$route访问路由，但是在vue setup中无法访问this：（template部分可以访问$router/$route）

```vue
import { useRouter, useRoute } from 'vue-router'

export default {
  setup() {
    const router = useRouter()
    const route = useRoute()

    function pushWithQuery(query) {
      router.push({
        name: 'search',
        query: {
          ...route.query,
          ...query,
        },
      })
    }
  },
}
```

## 滚动行为

vue-router可以让你在切换到新路由时，滚动到页面顶部、保持当前滚动位置、滚动到锚点等

## 路由懒加载

`import UserDetails from './views/UserDetails.vue'`  
替换成：  
`() => import('./views/UserDetails.vue')`
