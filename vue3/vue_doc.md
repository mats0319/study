# vue3学习笔记

> 学习时间：2023.9
> 官方文档: https://cn.vuejs.org/guide/introduction.html
> API文档：

vue是一个前端框架，也是一个生态。（可以做服务端渲染，SSR）

每个`.vue`文件包含`template`/`script`/`style`三个部分，分别使用**基于html模板语法的扩展、js/ts、css**实现

我在写前端的时候有一种感觉：html、js、css不是平级的。  
这似乎暗合react的设计哲学：使用jsx，将html、js、css写在一起；  
但实际上，我更喜欢把它们分开写，就像vue：template表示有哪些html元素以及他们的层级关系、script表示功能、style表示样式。  
这对应了两种把**真实事物抽象成代码**的方式。  
举个例子，给你一个网页，你会怎么抽象html元素、功能和样式之间的关系？  
react是以html元素为主，即，一个个html元素**恰好有那样的功能，又恰好是那样的样式**，组合到一起形成了一个网页
(html { js css })；  
vue是以样式布局为主，即，这个网页有一些html元素，它们按照一定的规则排布，然后每个html元素又有其功能(css { html { js } })。  
很难说哪种更好，甚至我都不知道会不会随着年龄的增长，我就自然而然的转向喜欢react的方式了。

## 基础

### v-bind / v-on

`v-bind`用于在html元素中，动态绑定一个属性：`<div v-bind:id="dynamicId"></div>`

简写为`:`：`<div :id="dynamicId"></div>`

`v-on`用于监听dom事件：`<a v-on:click="doSomething"> ... </a>`

简写为`@`：`<a @click="doSomething"> ... </a>`

`v-bind`/`v-on`可以绑定动态属性：`<a :[attributeName]="url"> ... </a>`，例子中的`attribute name`可能是`href`或者其他

v-on通常会绑定一个函数，如果绑定时不给函数后面的括号，则默认带出事件参数event：

```vue 
<script setup lang="ts">
function handleChange(event: Event) {
  console.log((event.target as HTMLInputElement).value)
}
</script>

<template>
  <input type="text" @change="handleChange" /> // 这里函数没有带括号，每次触发时会带上event
</template>
```

如果绑定函数需要参数、无法省略括号，可用`$event`带出事件参数：

```vue
<!-- 使用特殊的 $event 变量 -->
<button @click="warn('Form cannot be submitted yet.', $event)">
  Submit
</button>

<!-- 使用内联箭头函数 -->
<button @click="(event) => warn('Form cannot be submitted yet.', event)">
  Submit
</button>

function warn(message, event) {
  // 这里可以访问原生事件
  if (event) {
    event.preventDefault()
  }
  alert(message)
}
```

### 在模板中绑定动态内容 {{ }}

包括`{{ }}`和`v-xxx`，可以绑定简单的变量，也可以绑定js表达式

每个绑定仅支持单一表达式，也就是一段能够被求值的 JavaScript 代码。一个简单的判断方法是是否可以合法地写在 return 后面。

注意：函数不是实时响应的，它会在组件每次更新时被调用。例如每次新进入页面、html部分重新注入时，会调用函数。

### 修饰符 modifiers (draft)

todo

包括事件修饰符和按键修饰符，似乎是html部分的内容

.stop
.prevent
.self
.capture
.once
.passive

### 变量 ref

主要使用`ref`定义，在`script`部分使用时需要加`.value`，`html`部分不用：

1. html部分也可以使用.value访问，只是大部分情况下vue会自动解包
2. 顶级的ref属性会被解包，或者绑定的内容为变量本身、而不是表达式，也会解包

```vue
import { ref } from 'vue'
import type { Ref } from 'vue'

const year: Ref<string | number> = ref<string | number>('2020')

year.value = 2020 // 成功！
```

```vue
const count = ref<number>(0)
const object = { id: ref(1) }

{{ count + 1 }}
{{ object.id + 1 }} // 结果为：[object Object]1   因为.id不是顶级的ref属性

const { id } = object
{{ id + 1 }} // 可以通过将.id解构为顶级属性的方式，让表达式按照预期工作

{{ object.id }} // 如果绑定的只是简单的变量、不是表达式，也会正常解包
```

### 计算属性 computed

`computed`接收一个`getter`函数，返回一个`ref`，返回值可以和直接使用`ref`定义的变量一样访问，  
例如下方代码中，可以访问`publishedBooksMessage.value`

相较于函数，计算属性有缓存：只要计算属性依赖的内容没有发生变化，多次调用计算属性只会返回缓存结果而不会重新计算：

1. 计算属性的缓存特性适用于变量本身较大的情况，例如下方代码中`books`是一个很大的数组
2. 计算属性和函数的刷新节点不同：  
   例如在另一个函数中改变了`books`，计算属性刷新、函数不刷新；  
   而一些非响应式依赖则会在调用函数时刷新、计算属性不刷新，如`const now = computed(() => Date.now())`

```vue
<script setup>
import { reactive, computed } from 'vue'

const author = reactive({
  name: 'John Doe',
  books: ['Vue 2 - Advanced Guide']
})

// 一个计算属性 ref
const publishedBooksMessage = computed<string>(() => {
  return author.books.length > 0 ? 'Yes' : 'No'
})
</script>

<template>
  <p>Has published books:</p>
  <span>{{ publishedBooksMessage }}</span>
</template>
```

计算属性默认是只读的，如果想要修改计算属性，推荐修改它依赖的变量，而不是直接set计算属性

计算属性的getter应该只是做计算，而不修改什么（参考solidity中应用view描述的函数），包括但不限于异步请求、更改dom，这些可以在watch中实现。

### 动态绑定样式 :class

基本语法：`<div :class="{ 'active': isActive }"></div>`

1. active表示样式名，例如`text-alive`，一般写在style部分
2. isActive表示是否绑定该样式，类型是bool
3. 动态绑定可以和静态绑定同时存在

无论是动态还是静态绑定样式，如果**子组件有多个顶级html元素**，就需要`$attrs`属性来指定：  
当父组件调用子组件时，如果也给了样式，应该绑定到子组件的哪个元素上，示例代码如下：

```vue 
<MyComponent class="baz" /> // 父组件

// 以下是子组件顶级html元素
<p :class="$attrs.class">Hi!</p> // 渲染结果：<p class="baz">Hi!</p>
<span>This is a child component</span>
```

内联样式也可以这样绑定(`<div :style="{ color: activeColor, fontSize: fontSize + 'px' }"></div>`)，  
此时`:`前面依然表示样式名，而后面则表示样式的值，例如`50vw`

我个人更喜欢把内联样式也写到style部分，所以当需要绑定内联样式时，我可能会在style部分写一个样式，然后使用class绑定

### 条件渲染 v-if/v-show

`v-if`/`v-else`/`v-else-if`，后两个要求前面有一个v-if且只能绑定在同级html元素上

可以使用`template`组合多个html元素，统一控制其是否渲染

相较于`v-show`：

1. v-if会销毁、重新生成dom元素，v-show只是控制其display属性
2. v-show不支持用在`template`上，也不支持v-else
3. v-if切换开销更大，但v-if只有在值第一次为true时才初始化dom

v-if与v-for同级时，v-if优先级更高，这意味着v-if无法访问v-for中的迭代对象。不建议这么做，可以在v-if外层套一个`template v-for`：

```vue
<ul>
  <template v-for="user in users" :key="user.id">
    <li v-if="user.isActive">
      {{ user.name }}
    </li>
  </template>
</ul>
```

### 列表渲染 v-for

基本语法：`v-for="(item, index) in items"`/`v-for="(value, key， index) in myObject"`/`v-for="n in 10"`

1. item表示迭代变量名
2. index表示当前迭代序号
3. in可以换成of
4. 可以用于遍历一个对象的所有属性，顺序根据`Object.keys()`
5. 可以用于渲染固定次数：`n次，n∈[1,10]`

可以像v-if一样，使用`template`组合多个html元素

v-for通常和`:key`一起使用，用于标识列表中的不同dom并管理它们的状态，同一父元素下的key不能重复。  
实际上，key属性作为vue内置属性，作用是标识vnode，进而在更新dom树时，如果key改变了，对应dom会销毁-重新创建。

v-for可以用在组件上，但不会传值进去，因为作用域。如果需要传值，可以为组件指定属性：

```vue
<MyComponent
  v-for="(item, index) in items"
  :item="item"
  :index="index"
  :key="item.id"
/>
```

vue能响应数组变更方法，例如`push`/`shift`/`sort`等，无法响应非变更方法，例如`filter()`
，它不改变原数组，而是返回一个新的数组。  
如果你想要vue响应filter方法，应将方法返回的数组重新赋值给原始数组，参考`go slice append`的通常做法

### 表单输入绑定 v-model

因为js代码导致的变量的值发生改变，vue会帮我们反映到html元素上。  
但如果我们想将html元素中的内容绑定到变量，需要这样写：

```vue
<input
  :value="text"
  @input="event => text = event.target.value">
```

v-model指令可以简化这种写法：`<input v-model="text">`

v-model 还可以用于各种不同类型的输入，<textarea>、<select> 元素。它会根据所使用的元素自动使用对应的 DOM 属性和事件组合：

1. 文本类型的 <input> 和 <textarea> 元素会绑定 value property 并侦听 input 事件；
2. `<input type="checkbox">` 和 `<input type="radio">` 会绑定 checked property 并侦听 change 事件；
3. `<select>` 会绑定 value property 并侦听 change 事件。

v-model指令可以添加一些修饰符：

1. .lazy：`<input v-model.lazy="msg" />`，默认是监听input事件，即输入一个字符就会触发一次同步，现在会在change事件之后统一同步
2. .number：`<input v-model.number="age" />`，将用户输入处理成数字，  
   如果该值无法被 parseFloat() 处理，那么将返回原始值。  
   number 修饰符会在输入框有 type="number" 时自动启用。
3. .trim：`<input v-model.trim="msg" />`，默认去除用户输入内容收尾的空格

### 生命周期图示

https://cn.vuejs.org/guide/essentials/lifecycle.html#lifecycle-diagram

### watch

正如前面提到的，计算属性更像是**一个变量基于某种规则产生的新变量**，不应该有复杂操作，像是异步请求、修改dom之类的。

如果需要执行一些复杂操作，可以在watch中实现。

```vue
<script setup>
import { ref, watch } from 'vue'

const question = ref('')
const answer = ref('Questions usually contain a question mark. ;-)')

// 可以直接侦听一个 ref
watch(question, async (newQuestion, oldQuestion) => {
  if (newQuestion.indexOf('?') > -1) {
    answer.value = 'Thinking...'
    try {
      const res = await fetch('https://yesno.wtf/api')
      answer.value = (await res.json()).answer
    } catch (error) {
      answer.value = 'Error! Could not reach the API. ' + error
    }
  }
})
</script>

<template>
  <p>
    Ask a yes/no question:
    <input v-model="question" />
  </p>
  <p>{{ answer }}</p>
</template>
```

通常我们直接watch一个ref，如果你想watch一个对象的某个属性，或者其他什么数据源，[参考](https://cn.vuejs.org/guide/essentials/watchers.html#basic-example)

像是C语言的do...while循环，watch本身的含义是数据源变化了，执行回调，如果我们想创建watch时立刻执行一次回调，可以使用`immediate`：

```vue
watch(source, (newValue, oldValue) => {
  // 立即执行，且当 `source` 改变时再次执行
}, { immediate: true })
```

watch只追踪数据源的变化，如果有多个数据源会导致变化，可以使用watchEffect

现在你改了一个变量，它可能同时触发html元素的修改和watch的回调，默认会**优先调用watch的回调**，  
如果想要在html元素修改之后调用，可以使用`flush: 'post'`：`watch(source, callback, { flush: 'post' })`

### 模板引用 ref in html

部分场景下，可能需要直接操作dom，vue提供了一种获取dom的方法：

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'

// 要求变量与html元素中的ref同名
const el = ref<HTMLInputElement | null>(null)

onMounted(() => {
  el.value?.focus()
})
</script>

<template>
  <input ref="el" />
</template>
```

在v-for中使用模板引用，会获得一个dom数组

在组件上使用模板引用，会获得一个组件实例

### 组件基础

父组件向子组件传值（props）：

```vue
<!-- BlogPost.vue -->
<script setup lang="ts">
const props = defineProps<{
  title: string
  bar?: number
}>()
</script>

<template>
  <h4>{{ title }}</h4>
</template>

// use
<BlogPost title="My journey with Vue" />
```

子组件向父组件传递消息（emit）：其中子组件可以直接在html部分调用`$emit`触发事件，但我更喜欢把他们分开

```vue
// 子组件，假设叫Child.vue
<script setup>
const emit = defineEmits<{
  self-define-event: [id: number]
  update: [value: string]
}>()

emit("self-define-event", 1)
</script>

// 父组件
<Child @self-define-event="onEvent" />
```

slot：

```vue
<AlertBox>
  Something bad happened.
</AlertBox>

<template>
  <div class="alert-box">
    <strong>This is an Error for Demo Purposes</strong>
    <slot />
  </div>
</template>
```

动态组件：(`:is`)

```vue
<script setup>
import Home from './Home.vue'
import Posts from './Posts.vue'
import Archive from './Archive.vue'
import { ref } from 'vue'
 
const currentTab = ref('Home')

const tabs = {
  Home,
  Posts,
  Archive
}
</script>

<template>
  <div class="demo">
    <button
       v-for="(_, tab) in tabs"
       :key="tab"
       :class="['tab-button', { active: currentTab === tab }]"
       @click="currentTab = tab"
     >
      {{ tab }}
    </button>
	  <component :is="tabs[currentTab]"></component>
  </div>
</template>
```
