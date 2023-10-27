# vue生态

## pinia

> 学习时间：2023.9
> 文档：https://pinia.vuejs.org/zh/core-concepts/

pinia基于vuex一个升级方向的探索，它们都是用来管理全局状态的。

基础语法：（使用setup store格式）

```ts
import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', () => {
  const count = ref(0)
  function increment() {
    count.value++
  }

  return { count, increment }
})
```

解构会破坏响应性，推荐不定义新的变量、直接使用：

```ts
<script setup>
const store = useCounterStore()
// ❌ 这将不起作用，因为它破坏了响应性
// 这就和直接解构 `props` 一样
const { name, doubleCount } = store 
name // 将始终是 "Eduardo" 
doubleCount // 将始终是 0 
setTimeout(() => {
  store.increment()
}, 1000)
// ✅ 这样写是响应式的
// 💡 当然你也可以直接使用 `store.doubleCount`
const doubleValue = computed(() => store.doubleCount)
</script>
```

## vite配置文档

> 学习时间：2023.9
> 文档：https://cn.vitejs.dev/config/shared-options.html

路径别名：

path包需要安装依赖：@types/node

```ts
resolve: {
    alias: {
        '@': path.resolve(__dirname, './src'),
    }
},
```

不清屏：`clearScreen: false`

server.host：如果将此设置为 `0.0.0.0` 或者 `true` 将监听所有地址，包括局域网和公网地址，可能响应的是其他服务器而不是 Vite。

## ts编译选项 draft

> 学习时间：
> ts编译选项文档(tsconfig.json)：https://www.typescriptlang.org/docs/handbook/compiler-options.html
