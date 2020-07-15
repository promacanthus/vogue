---
title: "15 uintptr与unsafe.Pointer"
date: 2020-07-15T09:46:08+08:00
draft: true
---

- [0.1. `uintptr`](#01-uintptr)
- [0.2. `unsafe.Pointer`](#02-unsafepointer)
- [0.3. 区别](#03-区别)

Golang有两个无类型的指针：

- `uintptr`
- `unsafe.Pointer`

这两者可以互相转换。

## 0.1. `uintptr`

`uintptr` 是**整数**，不是引用。

将 `Pointer` 转换为 `uintptr` 会创建一个没有指针语义的整数值。

> 即使 `uintptr` 持有某个对象的地址，如果对象移动，垃圾收集器并不会更新 `uintptr` 的值，`uintptr` 也无法阻止该对象被回收。

## 0.2. `unsafe.Pointer`

`unsafe`包中`type Pointer`的说明如下：

`Pointer`表示指向任意类型的指针。`Pointer`类型有四种特殊操作，而其他类型则没有：

- 任何类型的指针值都可以转换为`Pointer`
- `Pointer`可以转换为任何类型的指针值
- 可以将`uintptr`转换为`Pointer`
- 可以将`Pointer`转换为`uintptr`

`Pointer`允许程序绕过类型系统并读取任意的内存地址。**使用时要格外小心**。

尽管 `unsafe.Pointer` 是通用指针，但 Go 垃圾收集器知道它们指向 Go 对象；换句话说，它们是真正的 **Go指针**。

垃圾收集器使用`Pointer`来防止活动对象被回收并发现更多活动对象（如果`unsafe.Pointer`指向的对象自身持有指针）。

## 0.3. 区别

因此，对 `unsafe.Pointer` 的合法操作上的许多限制归结为“在任何时候，它们都必须指向真正的 Go 对象”。如果创建的 `unsafe.Pointer` 并不符合，即使很短的时间，Go 垃圾收集器也可能会在该时刻扫描，然后由于发现了无效的 Go 指针而崩溃。

相比之下，`uintptr` 只是一个数字。这种特殊的垃圾收集魔法机制并不适用于 `uintptr` 所“引用”的对象，因为它仅仅是一个数字，一个 `uintptr` 不会引用任何东西。

反过来，这导致在将 `unsafe.Pointer` 转换为 `uintptr`，对其进行操作然后再将其转回的各种方式上存在许多微妙的限制。

基本要求是以这种方式进行操作，使编译器和运行时可以屏蔽不安全的指针的临时非指针性，使其免受垃圾收集器的干扰，因此这种临时转换对于垃圾收集将是原子的。

从 Go 1.8 开始，即使当时没有运行垃圾回收，所有 Go 指针必须始终有效（包括 `unsafe.Pointer`）。如果在变量或字段中存储了无效的指针，则仅通过将字段更新为包括 nil 在内的完全有效的值即可使代码崩溃。
