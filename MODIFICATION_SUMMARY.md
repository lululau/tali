# 编辑器和分页器配置功能 - 修改总结

## 修改概述

根据用户需求，实现了可配置的编辑器和分页器功能，替换了原来硬编码的nvim编辑器。配置项位于配置文件的顶层，对所有 profiles 生效。

## 修改的文件

### 1. `internal/config/config.go`

**结构体修改:**
- 从 `ConfigProfile` 结构体中移除了 `Editor` 和 `Pager` 字段
- 在 `AliyunConfig` 结构体顶层添加了 `Editor` 和 `Pager` 字段
- `Config` 结构体保持 `Editor` 和 `Pager` 字段不变

**新增函数:**
- `GetEditor()` - 获取编辑器命令，按优先级：配置文件顶层 → VISUAL环境变量 → EDITOR环境变量 → vim
- `GetPager()` - 获取分页器命令，按优先级：配置文件顶层 → PAGER环境变量 → less

**加载逻辑修改:**
- `LoadAliyunConfig()` 函数现在从配置文件顶层读取 `editor` 和 `pager` 字段

### 2. `internal/ui/utils.go`

**新增函数:**
- `OpenInEditor()` - 使用配置的编辑器打开JSON数据
- `OpenInPager()` - 使用配置的分页器查看JSON数据

**功能特点:**
- 支持带参数的命令（如 `"less -R"`）
- 自动创建和清理临时文件
- 正确处理tview应用的挂起和恢复

### 3. `internal/ui/components.go`

**修改函数:**
- `CreateInteractiveJSONDetailViewWithSearch()` - 添加了 `onPager` 参数
- 添加了 `v` 键的输入处理
- 更新了帮助文本，将 "edit in nvim" 改为 "edit"，并添加了 "v: view in pager"

**更新快捷键说明:**
- 所有详情页面的快捷键说明都更新为包含新的 `v` 键功能

### 4. `internal/app/navigation.go`

**全面更新:**
- 所有 `CreateInteractiveJSONDetailViewWithSearch` 调用都添加了 `onPager` 回调函数
- 所有 `OpenInNvimWithSuspend` 调用都替换为 `OpenInEditor`
- 为每个详情页面添加了分页器功能

## 新增按键绑定

- **`e` 键**: 使用配置的编辑器打开JSON数据
- **`v` 键**: 使用配置的分页器查看JSON数据

## 配置结构

### 新的配置文件结构
```json
{
  "current": "default",
  "editor": "code",
  "pager": "less -R",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your_key",
      "access_key_secret": "your_secret",
      "region_id": "cn-hangzhou"
    }
  ]
}
```

### 配置优先级

#### 编辑器选择
1. 配置文件顶层的 `editor` 字段
2. `VISUAL` 环境变量
3. `EDITOR` 环境变量
4. 默认 `vim`

#### 分页器选择
1. 配置文件顶层的 `pager` 字段
2. `PAGER` 环境变量
3. 默认 `less`

## 向后兼容性

- 如果配置文件中没有顶层的 `editor` 和 `pager` 字段，功能仍然正常工作
- 现有的配置文件无需修改即可继续使用
- 保持了所有原有的功能和按键绑定
- 旧的 profile 级别的 `editor` 和 `pager` 字段会被忽略

## 测试状态

- ✅ 编译成功
- ✅ 所有linter错误已修复
- ✅ 保持了原有功能的完整性
- ✅ 新功能按预期工作

## 使用示例

```json
{
  "current": "default",
  "editor": "code",
  "pager": "less -R",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your_key",
      "access_key_secret": "your_secret",
      "region_id": "cn-hangzhou"
    }
  ]
}
```

用户现在可以：
1. 在详情页面按 `e` 使用自定义编辑器编辑JSON
2. 在详情页面按 `v` 使用自定义分页器查看JSON
3. 通过配置文件顶层或环境变量自定义编辑器和分页器
4. 编辑器和分页器配置对所有 profiles 生效 