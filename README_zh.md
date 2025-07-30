# 🤖 iFlow CLI GitHub Action

一个 GitHub Action，使您能够在 GitHub 工作流中运行 [iFlow CLI](https://github.com/iflow-ai/iflow-cli) 命令。这个基于 Docker 的操作预装了 Node.js 22 和 npm 以实现最佳性能，并使用 iFlow CLI 执行您指定的命令。

> 文档站点（使用 iFlow CLI GitHub Action 生成）：[https://vibe-ideas.github.io/iflow-cli-action/](https://vibe-ideas.github.io/iflow-cli-action/)

## 功能特性

- ✅ 基于 Docker 的操作，预装 Node.js 22 和 npm
- ✅ 可配置的 iFlow API 认证
- ✅ 支持自定义模型和 API 端点
- ✅ 灵活的命令执行和超时控制
- ✅ 可在任何工作目录中运行
- ✅ 使用 Go 构建，快速可靠
- ✅ **GitHub Actions 摘要集成**：在 PR 摘要中提供丰富的执行报告

## 使用方法

### 基础示例

```yaml
name: iFlow CLI 示例
on: [push]

jobs:
  analyze-code:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: 运行 iFlow CLI
        uses: vibe-ideas/iflow-cli-action@v1.1.0
        with:
          prompt: "分析此代码库并提出改进建议"
          api_key: ${{ secrets.IFLOW_API_KEY }}
```

### 高级示例

```yaml
name: 高级 iFlow CLI 用法
on: 
  pull_request:
    types: [opened, synchronize]

jobs:
  code-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: 初始化项目分析
        uses: vibe-ideas/iflow-cli-action@v1.1.0
        with:
          prompt: "/init"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "Qwen3-Coder"
          timeout: "600"
          working_directory: "."
      
      - name: 生成技术文档
        uses: vibe-ideas/iflow-cli-action@v1.1.0
        with:
          prompt: "根据代码库分析生成技术文档"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          base_url: "https://apis.iflow.cn/v1"
          model: "DeepSeek-V3"
        id: docs
      
      - name: 显示结果
        run: |
          echo "生成的文档："
          echo "${{ steps.docs.outputs.result }}"
```

### 多命令示例

```yaml
name: 多步骤 iFlow 分析
on: [workflow_dispatch]

jobs:
  comprehensive-analysis:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: 项目概览
        uses: vibe-ideas/iflow-cli-action@v1.1.0
        with:
          prompt: |
            分析项目结构并提供：
            1. 主要架构组件
            2. 关键依赖及其用途
            3. 潜在的安全考虑
          api_key: ${{ secrets.IFLOW_API_KEY }}
          timeout: "900"
      
      - name: 代码质量评估
        uses: vibe-ideas/iflow-cli-action@v1.1.0
        with:
          prompt: "审查代码以了解最佳实践、潜在错误和性能改进"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "Kimi-K2"
```

## 输入参数

| 输入 | 描述 | 必需 | 默认值 |
|-------|-------------|----------|---------|
| `prompt` | 要使用 iFlow CLI 执行的提示 | ✅ 是 | - |
| `api_key` | 用于认证的 iFlow API 密钥 | ✅ 是 | - |
| `settings_json` | 完整的 iFlow settings.json 内容（JSON 字符串）。如果提供，将覆盖其他配置选项。 | ❌ 否 | - |
| `base_url` | iFlow API 的自定义基础 URL | ❌ 否 | `https://apis.iflow.cn/v1` |
| `model` | 要使用的模型名称 | ❌ 否 | `Qwen3-Coder` |
| `working_directory` | 运行 iFlow CLI 的工作目录 | ❌ 否 | `.` |
| `timeout` | iFlow CLI 执行超时时间（秒） | ❌ 否 | `300` |

## 输出参数

| 输出 | 描述 |
|--------|-------------|
| `result` | iFlow CLI 执行的输出 |
| `exit_code` | iFlow CLI 执行的退出代码 |

## 认证

### 获取 iFlow API 密钥

1. 在 [iflow.cn](https://iflow.cn) 注册 iFlow 账户
2. 转到您的个人资料设置或[点击这里](https://iflow.cn/?open=setting)
3. 在弹出对话框中点击"重置"以生成新的 API 密钥
4. 将 API 密钥添加到您的 GitHub 仓库 secrets 中，命名为 `IFLOW_API_KEY`

### 可用模型

- `Qwen3-Coder`（默认）- 适用于代码分析和生成
- `Kimi-K2` - 适用于通用 AI 任务和更长的上下文
- `DeepSeek-V3` - 高级推理和问题解决
- 支持通过 OpenAI 兼容 API 的自定义模型

## 自定义配置

### 自定义 iFlow 配置

```yaml
- name: 自定义 iFlow 配置
  uses: vibe-ideas/iflow-cli-action@v1.1.0
  with:
    prompt: "使用自定义配置分析此代码库"
    api_key: ${{ secrets.IFLOW_API_KEY }}  # 仍需要用于基本验证
    settings_json: |
      {
        "theme": "Dark",
        "selectedAuthType": "iflow",
        "apiKey": "${{ secrets.IFLOW_API_KEY }}",
        "baseUrl": "https://custom-api.example.com/v1",
        "modelName": "custom-model",
        "searchApiKey": "${{ secrets.SEARCH_API_KEY }}",
        "customField": "customValue"
      }
```

## 常见用例

### 代码分析和审查

```yaml
- name: 代码审查
  uses: vibe-ideas/iflow-cli-action@v1.1.0
  with:
    prompt: "审查此拉取请求的代码质量、安全问题和最佳实践"
    api_key: ${{ secrets.IFLOW_API_KEY }}
```

### 文档生成

```yaml
- name: 生成文档
  uses: vibe-ideas/iflow-cli-action@v1.1.0
  with:
    prompt: "/init && 生成全面的 API 文档"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "600"
```

### 自动化测试建议

```yaml
- name: 测试策略
  uses: vibe-ideas/iflow-cli-action@v1.1.0
  with:
    prompt: "分析代码库并建议全面的测试策略和具体测试用例"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    model: "DeepSeek-V3"
```

### 架构分析

```yaml
- name: 架构审查
  uses: vibe-ideas/iflow-cli-action@v1.1.0
  with:
    prompt: "分析系统架构并提出可扩展性和可维护性的改进建议"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "900"
```

## 要求

- **运行器**：基于 Linux 的 GitHub Actions 运行器（推荐 ubuntu-latest）
- **权限**：操作需要互联网访问权限以下载依赖项
- **资源**：足够的命令执行超时时间（根据复杂性调整）

## 故障排除

### 常见问题

**命令超时：** 为复杂操作增加 `timeout` 值
```yaml
timeout: "900"  # 15 分钟
```

**API 认证失败：** 验证您的 API 密钥是否正确设置在仓库 secrets 中

**工作目录未找到：** 确保路径存在且使用了 checkout 操作
```yaml
- uses: actions/checkout@v4  # 使用 iFlow 操作前必需
```

### 调试模式

通过设置环境变量启用详细日志记录：
```yaml
env:
  ACTIONS_STEP_DEBUG: true
```

## 贡献

欢迎贡献！请随时提交问题和拉取请求。

## 许可证

该项目根据 MIT 许可证授权 - 有关详细信息，请参见 [LICENSE](LICENSE) 文件。

## 相关链接

- [iFlow CLI](https://github.com/iflow-ai/iflow-cli) - 底层 CLI 工具
- [iFlow 平台](https://docs.iflow.cn/en/docs) - 官方文档
- [GitHub Actions 文档](https://docs.github.com/en/actions)