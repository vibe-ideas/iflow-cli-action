# Bot Name 全局变量重构

## 修改摘要

将 GitHub Actions 工作流中的 `gemini-cli` 硬编码字符串完全抽取为可配置的 repository variable `BOT_NAME`。

## 修改内容

### 1. 添加全局环境变量（使用 repository variable）
在工作流顶部添加了全局环境变量定义：
```yaml
env:
  BOT_NAME: ${{ vars.BOT_NAME || 'gemini-cli' }}
```

### 2. 替换条件表达式中的硬编码值
使用 `format()` 函数和 `vars.BOT_NAME` 替换所有条件表达式中的硬编码值：
- `contains(github.event.issue.body, '@gemini-cli')` → `contains(github.event.issue.body, format('@{0}', vars.BOT_NAME || 'gemini-cli'))`
- `contains(github.event.comment.body, '@gemini-cli')` → `contains(github.event.comment.body, format('@{0}', vars.BOT_NAME || 'gemini-cli'))`
- `contains(github.event.review.body, '@gemini-cli')` → `contains(github.event.review.body, format('@{0}', vars.BOT_NAME || 'gemini-cli'))`

### 3. 替换步骤中的硬编码值
- 在 `Get context from event` 步骤中，将 `sed 's/.*@gemini-cli//'` 替换为 `sed "s/.*@${BOT_NAME}//"`
- 在 `Set up git user for commits` 步骤中，将：
  - `'gemini-cli[bot]'` 替换为 `"${BOT_NAME}[bot]"`
  - `'gemini-cli[bot]@users.noreply.github.com'` 替换为 `"${BOT_NAME}[bot]@users.noreply.github.com"`

## 配置方法

### 设置 Repository Variable
1. 进入 GitHub 仓库设置
2. 导航到 **Settings** > **Secrets and variables** > **Actions**
3. 在 **Variables** 标签页中，点击 **New repository variable**
4. 添加变量：
   - Name: `BOT_NAME`
   - Value: `your-bot-name` (例如 `my-custom-cli`)

### 默认行为
如果没有设置 `BOT_NAME` repository variable，工作流将使用默认值 `gemini-cli`。

## 使用方法

要更改 bot 名称：
1. 在 GitHub 仓库中设置 `BOT_NAME` repository variable
2. 无需修改工作流文件！所有引用将自动更新

## 优势

- **完全可配置**：无需修改代码即可更改 bot 名称
- **集中化管理**：通过 repository variables 统一配置
- **向后兼容**：未设置变量时使用默认值
- **更加灵活**：可以为不同环境设置不同的 bot 名称
- **安全性**：敏感配置信息不暴露在代码中

## 技术实现

使用了 GitHub Actions 的以下特性：
- **Repository Variables** (`vars.BOT_NAME`) 用于条件表达式
- **Environment Variables** (`env.BOT_NAME`) 用于步骤执行
- **Format Function** (`format('@{0}', vars.BOT_NAME)`) 用于字符串拼接
- **默认值语法** (`vars.BOT_NAME || 'gemini-cli'`) 提供回退选项
