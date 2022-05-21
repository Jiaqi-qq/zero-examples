# fx

​		`fx`是一个完整的流处理组件。 它与 `MapReduce` 类似，`fx` 也有一个并发处理函数：`Parallel(fn, options)`。但同时，它又不仅仅是并发处理。`From(chan)`, `Map(fn)`, `Filter(fn)`, `Reduce(fn)`等，从数据源读入流，处理流数据，最后聚合流数据。

| API                  | 函数                                                         |
| -------------------- | ------------------------------------------------------------ |
| `Distinct(fn)`       | 在fn中选择一个特定的项目类型，并将其去掉。                   |
| `Filter(fn, option)` | fn指定特定的规则，符合规则的 `element` 被传递到下一个 `stream`。 |
| `Group(fn)`          | 根据fn，`stream` 中的元素被分为不同的组。                    |
| `Head(num)`          | 取出 `stream` 中的前n个元素，生成一个新的 `stream`。         |
| `Map(fn, option)`    | 将每个元素转换为另一个对应的元素，并将其传递给下一个 `stream`。 |
| `Merge()`            | 将所有的ele合并成一个 `slice`，并生成一个新的 `stream`。     |
| `Reverse()`          | 反转 `stream` 中的元素。[使用双指针]                         |
| `Sort(fn)`           | 根据 fn 对 `stream` 中的元素进行排序。                       |
| `Tail(num)`          | 取出 `stream` 的最后 n 个元素，生成一个新的 `stream`。[使用一个双链表] |
| `Walk(fn, option)`   | 将fn应用于 `source` 的每个元素。生成一个新的 `stream`        |

不再生成一个新的`stream`，做最后的评估操作。

| API                    | 函数                                                         |
| ---------------------- | ------------------------------------------------------------ |
| `ForAll(fn)`           | 根据fn处理`stream`，不再生成 `stream` [评估操作]             |
| `ForEach(fn)`          | 对 `stream` 中的所有元素进行fn[求值操作] !                   |
| `Parallel(fn, option)` | 同时对每个 `element` 应用给定的fn和给定数量的worker[求值操作] |
| `Reduce(fn)`           | 直接处理 `stream` [评估操作]                                 |
| `Done()`               | 不做任何事情，等待所有操作完成                               |



## 用法

### `Distinct`

```go
func TestDistinct(t *testing.T) {
    fx.Just(1, 3, 5, 7, 1, 3, 5, 5, 5, 5).Distinct(func(item interface{}) interface{} {
        // 通过key进行去重，相同key只保留一个
        return item
    }).ForEach(func(item interface{}) {
        fmt.Println(item)
    })
}
/*
1
3
5
7
*/
```

### `Group`

```go
func TestGroup(t *testing.T) {
    fx.Just(1, 3, 5, 7, 1, 3, 5, 5, 5, 5).Group(func(item interface{}) interface{} {
        return item
    }).ForEach(func(item interface{}) {
        fmt.Println(item)
    })
}
/*
[1 1]
[3 3]
[5 5 5 5 5]
[7]
*/
```

### `Map`

```go
func TestMap(t *testing.T) {
    fx.Just(1, 3, 5).Map(func(item interface{}) interface{} {
        return strconv.Itoa(item.(int)) + " - abc"
    }).ForEach(func(item interface{}) {
        fmt.Println(item)
    })
}
/*
5 - abc
1 - abc
3 - abc
*/
```

### `Merge`

```go
func TestMerge(t *testing.T) {
	fx.Just("abc", 3, "dtf", 5).Merge().ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}
/*
[abc 3 dtf 5]
*/
```

### Parallel

```go
func TestParallel(t *testing.T) {
	fx.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Group(func(item interface{}) interface{} {
		return item.(int) % 3
	}).Parallel(func(item interface{}) {
		fmt.Println(item)
		time.Sleep(time.Second * 2)
	})
}
/* 三行同时输出 并发数与元素个数相等
[2 5 8]
[1 4 7 10]
[3 6 9]
*/
```

### Reduce

```go
func TestReduce(t *testing.T) {
	reduce, err := fx.Just(1, 2, 3, 4).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for x := range pipe {
			time.Sleep(time.Second)
			fmt.Println(x)
		}
		return "你以为这是什么", nil
	})
	fmt.Println(reduce, err)
}
/*
1
2
3
4
你以为这是什么 <nil>
--- PASS: TestReduce (4.04s)
*/
```

