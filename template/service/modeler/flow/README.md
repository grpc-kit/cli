## 工作流

规范：

```text
# 用于定义 argo workflow 流程编排
template/

# 用于实现工作流脚本编码
script/
```

每个工作流根据类型分组在各自文件夹下，如：

```text
template/example/
script/example/
```

当创建新目录时，需要在 `embed.go` 中添加对应目录的文件引用，如：

```text
//go:embed example/*.sh
//go:embed example/*.py
```

所有脚本大小尽量避免超过 100kb 否则有可能无法运行成功。
