# 配置结构更新说明

## 更新概述

根据用户需求，将 `editor` 和 `pager` 配置项从 profiles 数组中的每个 profile 移动到 JSON 配置文件的顶层。这样这些配置就是全局的，对所有 profiles 生效，而不是每个 profile 特有的。

## 修改内容

### 1. 配置结构变更

**修改前的结构:**
```json
{
  "current": "default",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your_key",
      "access_key_secret": "your_secret",
      "region_id": "cn-hangzhou",
      "editor": "code",
      "pager": "less -R"
    }
  ]
}
```

**修改后的结构:**
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

### 2. 代码修改

#### `internal/config/config.go`
- 从 `ConfigProfile` 结构体中移除了 `Editor` 和 `Pager` 字段
- 在 `AliyunConfig` 结构体中添加了顶层的 `Editor` 和 `Pager` 字段
- 更新了 `LoadAliyunConfig()` 函数，现在从配置文件顶层读取这些字段

#### 文档更新
- 更新了 `EDITOR_PAGER_FEATURE.md` 中的配置示例
- 更新了 `MODIFICATION_SUMMARY.md` 中的说明

### 3. 测试验证

添加了配置结构测试到 `test_features.go`，验证：
- `GetEditor()` 函数正常工作
- `GetPager()` 函数正常工作
- 配置优先级正确（配置文件 → 环境变量 → 默认值）

## 优势

1. **全局配置**: `editor` 和 `pager` 现在是全局设置，不需要为每个 profile 重复配置
2. **简化配置**: 减少了配置文件的冗余
3. **更符合逻辑**: 编辑器和分页器通常是用户的个人偏好，不应该与特定的云账号 profile 绑定
4. **向后兼容**: 如果配置文件中没有这些字段，系统会回退到环境变量或默认值

## 配置优先级

### 编辑器选择优先级
1. 配置文件顶层的 `editor` 字段
2. `VISUAL` 环境变量
3. `EDITOR` 环境变量
4. 默认使用 `vim`

### 分页器选择优先级
1. 配置文件顶层的 `pager` 字段
2. `PAGER` 环境变量
3. 默认使用 `less`

## 迁移指南

对于现有用户：

1. **无需立即修改**: 现有配置文件仍然可以正常工作，系统会忽略 profile 级别的 `editor` 和 `pager` 字段
2. **推荐迁移**: 建议将 `editor` 和 `pager` 字段移动到配置文件顶层
3. **环境变量**: 如果不想修改配置文件，也可以使用环境变量 `VISUAL`、`EDITOR` 和 `PAGER`

## 测试状态

- ✅ 编译成功
- ✅ 配置加载测试通过
- ✅ GetEditor() 和 GetPager() 函数测试通过
- ✅ 向后兼容性保持
- ✅ 所有现有功能正常工作

## 示例配置

```json
{
  "current": "default",
  "editor": "code --wait",
  "pager": "less -R",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your_access_key",
      "access_key_secret": "your_secret_key",
      "region_id": "cn-hangzhou"
    },
    {
      "name": "production",
      "mode": "AK",
      "access_key_id": "prod_access_key",
      "access_key_secret": "prod_secret_key",
      "region_id": "cn-shanghai"
    }
  ]
}
```

这样配置后，无论切换到哪个 profile，都会使用相同的编辑器和分页器设置。 