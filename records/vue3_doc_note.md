# vue3官方文档阅读笔记

1. 不要在选项property或回调上使用箭头函数，例如`created: () => console.log(this.a)`或`vm.$watch('a', newValue => this.myMethod())`，  
   因为箭头函数没有this，this会被视为变量，一直向上级词法作用域查找，直到找到，所以容易出现`Uncaught TypeError: Cannot read property of undefined`  
   或`Uncaught TypeError: this.myMethod is not a function`之类的错误。
    1. 其他地方也有该问题，总之就是.vue文件中少用箭头函数
2. 指令：`v-`开头的特殊属性
    1. v-bind/v-on可以在后面用`:`
       接收一个参数，这个参数可以是[动态的](https://v3.cn.vuejs.org/guide/template-syntax.html#%E5%8A%A8%E6%80%81%E5%8F%82%E6%95%B0)
3. vue官方推荐实现方式：不要在模板（html部分）放入太多逻辑。  
   仅在html部分绑定数据(`{{ data }}`)，或编写简单的表达式(`{{ data+1 }}`)；而对于任何包含响应式数据的复杂逻辑，应编写成计算属性（`{{ computedData }}`，以下称为computed)
   ```txt
   computed: {
     // 计算属性的 getter
     computedData() {
       // `this` 指向 vm 实例
       return this.author.books.length > 0 ? 'Yes' : 'No'
     }
   }   
   ```
    2. computed将基于它们的响应依赖关系缓存，computed只会在相关响应式依赖发生改变时重新求值（例如上例中的`this.author.books`）  
       而像`Date.now()`这种，则不会重新计算，因为`Date.now()`不是响应式依赖
    3. 编写方法(methods)也能实现**将复杂逻辑转移到模板之外**，但结合上一条，如果我们需要缓存这个逻辑的执行结果，可以使用computed；如果不需要，可以使用方法
    4. watch常用于异步操作，例如发送http请求等，如果是简单的行为，应使用computed（例如a改变了，b根据a变化）
4. 动态绑定class
    1. 语法1：`:class="{ [class name]: [is active] }"`，解释：如果`is active`的结果为true，则html渲染结果为`class=[class name]`
    2. 语法2：`:style="{ color: [color data] }"`，解释：假设`color data`为`red`，则html渲染结果为`style={ color: 'red' }`，个人建议，仅了解
    3. 可以传入更多k-v，都在大括号中，以`,`分隔
    4. 结合computed使用：
       ```txt 
       <div :class="classObject"></div>
       data() {
         return {
           isActive: true,
           error: null
         }
       },
       computed: {
         classObject() {
           return {
             active: this.isActive && !this.error,
             'text-danger': this.error && this.error.type === 'fatal'
           }
         }
       }
       ```
    5. 使用对象（大括号）要求为每一个`class name`绑定一个`bool`变量，如果不想绑定`bool`变量，可以使用数组，数组里还可以使用对象，个人建议，仅了解
5. v-if还有v-else/v-else-if，后两种还需要注意结构关系，感觉没什么必要，可以直接上v-if
    1. v-if是条件渲染，它可以保证在切换过程中，事件监听器和子组件适当地被销毁和重建；v-show始终渲染组件，只是切换元素的`display`css属性
    2. 如果需要频繁切换，适合使用v-show；如果运行时条件很少改变，适合使用v-if
    3. v-if与v-for一起使用，v-if优先级更高；当然，vue不推荐一起使用（优先级更高，意味着v-if无法访问v-for中的变量）
    4. 遇到逻辑上需要v-if、v-for同级的场景，可以将v-for提升一级，形如：
       ```txt 
       <template v-for="item in items">
         <li v-if="item.ok" />
       </template>
       ```
6. v-for遍历
    1. 遍历数组，语法：`v-for="(item,index) in items"`，item表示数组元素，index表示其索引
    2. 遍历对象，语法：`v-for="(value,key,index) in object"`，value表示对象内一条属性的值，key表示value对应的键，index表示索引
    3. 重复遍历，语法：`v-for="n in 10"`，表示将该模板重复10次
    4. v-for通常建议提供key，这涉及到元素刷新问题，除非你想要使用默认的元素刷新规则做事情
7. 数组更新检测，用于检测数组的更新，即使在模板（html部分）使用数组元素作为绑定值，也可以像使用字符串、数字等基础类型的绑定值一样，触发视图更新
    1. 变更方法，包括`push()`/`pop()`/`shift()`/`unshift()`/`splice()`/`sort()`/`reverse()`
    2. 替换数组，相较于上一条，还有一些方法不是原地修改，而是整体替换，例如`filter()`/`concat()`/`slice()`
8. 事件处理，v-on，通常简写为`@`，可以绑多个方法，例如`@click="one,two"`
    1. 事件修饰符，用于处理dom事件细节，例如`.stop`/`.prevent`/`.self`/`.once`等
    2. $event参数，可以传递原始dom事件
9. 选择框，vue推荐将第一个值设置为空且禁用，因为v-model的初始值如果没能匹配一个选项，select将被渲染为*未选中*状态，在ios中，用户无法选择第一个选项，因为不会触发change事件

## 问题

1. 如何注册全局方法？——在app.vue里编写方法，然后将这个vue实例公开
2. vue原生script格式：

```ts 
import {defineComponent} from 'vue';
import HelloWorld from '@/components/HelloWorld.vue'; // @ is an alias to /src

export default defineComponent({
  name: 'HomeView',
  components: {
    HelloWorld,
  },

  data() {
    return {
      welcome: "Welcome to vue3 !"
    }
  },

  mounted() {
    // placeholder
  }
});

// 全部可用属性
declare interface LegacyOptions<Props, D, C extends ComputedOptions, M extends MethodOptions, Mixin extends ComponentOptionsMixin, Extends extends ComponentOptionsMixin> {
  compatConfig?: CompatConfig;

  [key: string]: any;

  data?: (this: CreateComponentPublicInstance<Props, {}, {}, {}, MethodOptions, Mixin, Extends>, vm: CreateComponentPublicInstance<Props, {}, {}, {}, MethodOptions, Mixin, Extends>) => D;
  computed?: C;
  methods?: M;
  watch?: ComponentWatchOptions;
  provide?: ComponentProvideOptions;
  inject?: ComponentInjectOptions;
  filters?: Record<string, Function>;
  mixins?: Mixin[];
  extends?: Extends;

  beforeCreate?(): void;

  created?(): void;

  beforeMount?(): void;

  mounted?(): void;

  beforeUpdate?(): void;

  updated?(): void;

  activated?(): void;

  deactivated?(): void;

  /** @deprecated use `beforeUnmount` instead */
  beforeDestroy?(): void;

  beforeUnmount?(): void;

  /** @deprecated use `unmounted` instead */
  destroyed?(): void;

  unmounted?(): void;

  renderTracked?: DebuggerHook;
  renderTriggered?: DebuggerHook;
  errorCaptured?: ErrorCapturedHook;
  /**
   * runtime compile only
   * @deprecated use `compilerOptions.delimiters` instead.
   */
  delimiters?: [string, string];
  /**
   * #3468
   *
   * type-only, used to assist Mixin's type inference,
   * typescript will try to simplify the inferred `Mixin` type,
   * with the `__differentiator`, typescript won't be able to combine different mixins,
   * because the `__differentiator` will be different
   */
  __differentiator?: keyof D | keyof C | keyof M;
}
```
