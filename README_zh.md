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
        uses: vibe-ideas/iflow-cli-action@v1.3.0
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
        uses: vibe-ideas/iflow-cli-action@v1.3.0
        with:
          prompt: "/init"
          api_key: ${{ secrets.IFLOW_API_KEY }}
          model: "Qwen3-Coder"
          timeout: "600"
          working_directory: "."
      
      - name: 生成技术文档
        uses: vibe-ideas/iflow-cli-action@v1.3.0
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
        uses: vibe-ideas/iflow-cli-action@v1.3.0
        with:
          prompt: |
            分析项目结构并提供：
            1. 主要架构组件
            2. 关键依赖及其用途
            3. 潜在的安全考虑
          api_key: ${{ secrets.IFLOW_API_KEY }}
          timeout: "900"
      
      - name: 代码质量评估
        uses: vibe-ideas/iflow-cli-action@v1.3.0
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
| `timeout` | iFlow CLI 执行超时时间（秒） | ❌ 否 | `86400` |
| `extra_args` | 传递给 iFlow CLI 的附加命令行参数（空格分隔的字符串） | ❌ 否 | `` |
| `precmd` | 在运行 iFlow CLI 之前执行的 Shell 命令（例如 "npm install", "git fetch"） | ❌ 否 | `` |

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

### 使用附加参数

`extra_args` 输入允许您直接向 iFlow CLI 传递附加的命令行参数。这提供了灵活性，可以使用未作为专用操作输入公开的高级 iFlow CLI 功能。

```yaml
- name: 带自定义参数的 iFlow
  uses: vibe-ideas/iflow-cli-action@v1.3.0
  with:
    prompt: "使用调试输出分析此代码库"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    extra_args: "--debug --max-tokens 3000"
```

#### 附加参数示例

- `--debug` - 启用 iFLOW CLI 调试模式

### 使用预执行命令

`precmd` 输入允许您在执行 iFlow CLI 之前运行 Shell 命令。这对于设置环境或安装 iFlow 命令所需的依赖项非常有用。

```yaml
- name: 带预执行命令的 iFlow
  uses: vibe-ideas/iflow-cli-action@v1.3.0
  with:
    prompt: "在安装依赖项后分析此代码库"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    precmd: |
      npm install
      git fetch origin main
```

#### 多行命令

您可以通过用换行符分隔来指定多个命令：

```yaml
precmd: |
  npm ci
  npm run build
```

#### 带引号的参数

对于包含空格的参数，请使用引号：

```yaml
extra_args: '--debug'
```

### 使用自定义设置

对于需要完全控制 iFlow 配置的高级用户，您可以直接提供自定义的 `settings.json`：

```yaml
- name: 自定义 iFlow 配置
  uses: vibe-ideas/iflow-cli-action@v1.3.0
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

当提供 `settings_json` 时，它优先于单个配置输入（`base_url`、`model` 等）。这允许您：

- 使用自定义认证类型
- 配置输入中不可用的附加字段
- 在多个工作流运行中维护复杂配置
- 支持自定义 API 端点和模型

**注意：** 仍需要 `api_key` 输入进行验证，但实际使用的 API 密钥将是您在 `settings_json` 中指定的密钥。

## 使用 MCP 服务器

[MCP (Model Context Protocol)](https://modelcontextprotocol.io) 允许 iFlow CLI 连接到外部工具和服务，扩展其超越 AI 模型交互的能力。您可以在工作流中配置 MCP 服务器，以启用代码搜索、数据库查询或自定义工具集成等功能。

### 示例：使用 DeepWiki MCP 服务器

以下示例演示了如何配置和使用 DeepWiki MCP 服务器以增强代码搜索功能：

```yaml
- name: 带 MCP 服务器的 iFlow CLI
  uses: vibe-ideas/iflow-cli-action@v1.3.0
  with:
    prompt: "使用 @deepwiki 搜索如何使用 Skynet 构建游戏"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    settings_json: |
      {
        "selectedAuthType": "iflow",
        "apiKey": "${{ secrets.IFLOW_API_KEY }}",
        "baseUrl": "https://apis.iflow.cn/v1",
        "modelName": "Qwen3-Coder",
        "searchApiKey": "${{ secrets.IFLOW_API_KEY }}",
        "mcpServers": {
          "deepwiki": {
            "command": "npx",
            "args": ["-y", "mcp-deepwiki@latest"]
          }
        }
      }
    model: "Qwen3-Coder"
    timeout: "1800"
    extra_args: "--debug"
```

在此示例中：

- `mcpServers` 配置定义了一个名为 `deepwiki` 的服务器
- 服务器通过 `npx -y mcp-deepwiki@latest` 执行
- 提示中使用 `@deepwiki` 引用服务器以利用其功能
- `searchApiKey` 用于 DeepWiki 服务的认证

### 何时使用 MCP 服务器

当您需要以下功能时，MCP 服务器特别有用：

- 增强的代码搜索和文档查找功能
- 与外部工具和服务的集成
- 访问专业知识库或数据库
- 扩展 iFlow CLI 功能的自定义工具

## 示例工作流程

该仓库包含多个示例工作流程，展示了 iFlow CLI GitHub Action 的各种用例。您可以在 [.github/workflows](.github/workflows) 目录中找到这些示例。

### 定时生成主页

[deploy-homepage.yml](.github/workflows/deploy-homepage.yml) 工作流程演示了如何使用 iFlow CLI 自动从您的 README 生成 Reveal.js 演示文稿，并按计划部署到 GitHub Pages。

主要特性：
- 每天 UTC 时间 02:00 运行，或在主分支的 README 更改时运行
- 使用 iFlow CLI 将 README.md 转换为 Reveal.js 演示文稿
- 将生成的演示文稿部署到 GitHub Pages

### 使用 MCP 服务器

[iflow-with-mcp.yml](.github/workflows/iflow-with-mcp.yml) 工作流程展示了如何使用 iFlow CLI 与 MCP (Model Context Protocol) 服务器集成以扩展功能。

主要特性：
- 演示与 DeepWiki MCP 服务器的集成
- 展示如何在提示中使用 `@deepwiki` 来访问增强的搜索功能
- 使用自定义 settings_json 配置 MCP 服务器

### 自动化问题实现

[issue-killer.yml](.github/workflows/issue-killer.yml) 工作流程通过使用 iFlow CLI 和 GitHub MCP 集成来实现自动化问题处理。

主要特性：
- 通过在问题上评论 `@iflow-cli /issue-killer` 触发
- 自动生成实现计划和进度更新
- 使用 GitHub MCP 工具创建拉取请求和添加评论
- 直接根据问题描述实现功能

### 拉取请求审查自动化

该仓库包含多个用于 PR 审查自动化的流程：
- [pr-review.yml](.github/workflows/pr-review.yml) - 带有详细反馈的自动化 PR 审查
- [pr-review-killer.yml](.github/workflows/pr-review-killer.yml) - 基于审查评论的 PR 修改

### Common Use Cases

### Code Analysis and Review

```yaml
- name: Code Review
  uses: vibe-ideas/iflow-cli-action@v1.3.0
  with:
    prompt: "Review this pull request for code quality, security issues, and best practices"
    api_key: ${{ secrets.IFLOW_API_KEY }}
```

### Documentation Generation

```yaml
- name: Generate Documentation
  uses: vibe-ideas/iflow-cli-action@v1.3.0
  with:
    prompt: "/init && Generate comprehensive API documentation"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "600"
```

### Automated Testing Suggestions

```yaml
- name: Test Strategy
  uses: vibe-ideas/iflow-cli-action@v1.3.0
  with:
    prompt: "Analyze the codebase and suggest a comprehensive testing strategy with specific test cases"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    model: "DeepSeek-V3"
```

### Architecture Analysis

```yaml
- name: Architecture Review
  uses: vibe-ideas/iflow-cli-action@v1.3.0
  with:
    prompt: "Analyze the system architecture and suggest improvements for scalability and maintainability"
    api_key: ${{ secrets.IFLOW_API_KEY }}
    timeout: "900"
```

## Troubleshooting

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
- [Gemini CLI GitHub Action](https://github.com/google-github-actions/run-gemini-cli)