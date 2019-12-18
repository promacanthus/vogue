# 06-Hooks

Helm 提供了一个 hook 机制，允许在 release 的生命周期中的某些点进行干预。例如，可以使用 hooks 来|

1. 在加载任何其他 chart 之前加载 ConfigMap 或 Secret。
2. 在安装新 chart 之前执行job1以备份数据库，在升级完成后执行job2以恢复数据。
3. 在删除 release 之前执行job，以便在删除 release 之前优雅地停止服务。

Hooks 像常规模板一样工作，但它们具有**特殊的注释**，可以使 Helm 以不同的方式使用它们。

## 可用的Hooks

|功能|参数|描述|
---|---|---
安装前| pre-install|在模板渲染后，资源加载到Kubernetes前，执行
安装后| post-install|在所有资源加载到Kubernetes后，执行
删除前| pre-delete|删除任何release资源前，执行
删除后| post-delete|删除所有release资源后，执行
升级前| pre-upgrade|在模板渲染后，资源加载到Kubernetes前，执行
升级后| post-upgrade|在所有资源加载到Kubernetes后，执行
回滚前| pre-rollback|在渲染模板后，资源加载到Kubernetes前，执行
回滚后| post-rollback|在所有资源加载到Kubernetes后，执行

## Hooks与release生命周期

Hooks 让 chart 开发人员有机会在 release 的生命周期中的**关键点**执行操作。

默认情况下，`release-a`的生命周期如下：

1. 用户运行`helm install chart-a`
2. chart-a被加载到Tiller中
3. 经过验证后，Tiller渲染chart-a中的模板
4. Tiller将渲染后产生的YAML文件加载到Kubernetes中
5. Tiller将release-a名称和其他数据返回给客户端
6. 客户端接收数据并退出

假设在上述release-a的生命周期中定义两个hook：`pre-install`和`post-install`，新的生命周期如下:

1. 用户运行`helm install chart-a`
2. chart-a被加载到Tiller中
3. 经过验证后，Tiller渲染chart-a中的模板
4. **Tiller准备执行`pre-install` hook （将hook资源加载到kubernetes中）**
5. **Tiller根据权重（默认分配权重为0）对hook进行排序，相同权重hook按升序排列**
6. **Tiller加载最低权重的hook（从负到正）**
7. **Tiller等待直到hook操作完成（如果设置`--wait`标志，Tiller等待直到所有资源都处于就绪状态，并且在准备就绪前不会运行`post-install` hook）**
8. Tiller将渲染后产生的YAML文件加载到Kubernetes中
9. **Tiller执行`post-install` hook（将hook资源加载到kubernetes中）**
10. **Tiller根据权重（默认分配权重为0）对hook进行排序，相同权重hook按升序排**
11. **Tiller加载最低权重的hook（从负到正）**
12. **Tiller等待直到hook操作完成**
13. Tiller将release-a名称和其他数据返回给客户端
14. 客户端接收数据并退出

> 加粗的步骤为执行hook操作的步骤。添加Hook权重是比较好的做法，如果权重不重要则设置为0，默认也为0。

等到 hook 准备就绪取决于在 hook 中声明的资源：

1. 如果资源是 **Job**，Tiller 将等到作业成功完成。如果作业失败，则发布失败。这是一个阻塞操作，所以 Helm 客户端会在 Job 运行时暂停。
2. 对于**其他类型**，只要 Kubernetes 将资源标记为加载（添加或更新），资源就被视为 “就绪”。当一个 hook 声明了很多资源时，这些资源将被串行执行。如果有 hook 权重，按照加权顺序执行。否则，顺序不被保证（在 Helm 2.3.0 及之后的版本中，按字母顺序排列）。

## 注意

Hook创建的资源不作为release的一部分进行跟踪或管理。一旦Tiller验证Hook已经达到其就绪状态，它将Hook资源放在一边。

这意味着在 Hook 中创建资源，则不能依赖于 `helm delete` 删除资源。要销毁这些资源，需要编写代码在 `pre-delete` 或 `post-delete` Hook中执行此操作，或者将 `"helm.sh/hook-delete-policy"` 注释添加到 Hook 模板文件。

## 例子

Hook也是Kubernetes的manifest文件，只是在metadata部分有**特殊注释** 。Hook也是模板文件，可以使用模板的所有功能，包括读取 `.Values`，`.Release` 和 `.Template`。

创建一个Hook文件存放在 `templates/post-install-job.yaml`文件中，将其声明为在`post-install`阶段运行：

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{.Release.Name}}"
  labels:
    app.kubernetes.io/managed-by: {{.Release.Service | quote}}
    app.kubernetes.io/instance: {{.Release.Name | quote}}
    helm.sh/chart: "{{.Chart.Name}}-{{.Chart.Version}}"
  annotations:
    # 在这里添加注释，将资源定义为Hook，没有这些注释，这个资源将被任务是release的一部分

    "helm.sh/hook": post-install，post-upgrade      # 部署为多个Hook

    "helm.sh/hook-weight": "-5"                     # 设置Hook的权重，确定执行顺序，
                                                    # 权重可以是正数或负数，但必须表示为字符串

    "helm.sh/hook-delete-policy": hook-succeeded    # 定义删除Hook资源的时间和策略
                                                    # hook-succeeded：在执行成功后删除hook
                                                    # hook-failed：在执行失败后删除hook
                                                    # before-hook-creation：创建新hook之前删除旧hook
    "helm.sh/resource-policy": keep                 # 指示Tiller在helm delete操作过程中跳过此资源（将变成孤儿）
spec:
  template:
    metadata:
      name: "{{.Release.Name}}"
      labels:
      app.kubernetes.io/managed-by: {{.Release.Service | quote}}
      app.kubernetes.io/instance: {{.Release.Name | quote}}
      helm.sh/chart: "{{.Chart.Name}}-{{.Chart.Version}}"
    spec:
      restartPolicy: Never
      containers:
      - name: post-install-job
        image: "alpine:3.3"
        command: ["/bin/sleep","{{default"10".Values.sleepyTime}}"]
```

- 实现一个给定的Hook的不同种类资源数量没有限制。例如，可以将secret和config map声明为`per-install` Hook。
- 子chart声明Hook时，Tiller也会渲染这些Hook。顶级chart无法禁用子chart所声明的Hook。