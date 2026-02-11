<div align="center">
  <img src="media/logo-autocache.png" alt="AutoCache Logo" width="400"/>

# Autocache

  **带有 ROI 分析的智能 Anthropic API 缓存代理**

  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
  [![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
  [![Tests](https://img.shields.io/badge/tests-passing-brightgreen)]()

</div>

Autocache 是一个智能代理服务器，可自动向 Anthropic Claude API 请求中注入 `cache-control` 字段，从而将成本降低高达 **90%**，延迟降低高达 **85%**，同时通过响应头提供详细的 ROI（投资回报率）分析。

## 动机

现代 AI 代理平台（如 **n8n**、**Flowise**、**Make.com**）甚至流行的框架（如 **LangChain** 和 **LlamaIndex**）目前大多不支持 Anthropic 的提示词缓存（Prompt Caching）——尽管用户正在构建日益复杂的代理，这些代理通常包含：

- 📝 **大型系统提示词** (1,000-5,000+ tokens)
- 🛠️ **10 个以上的工具定义** (5,000-15,000+ tokens)
- 🔄 **重复的代理交互**（相同的上下文，不同的查询）

### 问题

当你在 **n8n** 中构建一个带有详细系统提示词和多个工具的复杂代理时，每次 API 调用都会重新发送完整的上下文——这导致成本比必要的高出 10 倍。例如：

- **不使用缓存**：15,000 token 的代理 → 每次请求 $0.045
- **使用缓存**：相同的代理 → 首次调用后每次请求 $0.0045（节省 90%）

### 真实用户痛点

AI 社区一直在请求此功能：

- 🔗 [n8n GitHub Issue #13231](https://github.com/n8n-io/n8n/issues/13231) - "Anthropic 模型未缓存系统提示词"
- 🔗 [Flowise Issue #4289](https://github.com/FlowiseAI/Flowise/issues/4289) - "支持 Anthropic 提示词缓存"
- 🔗 [n8n 社区请求](https://community.n8n.io/t/request-prompt-caching-support-for-claude/101941) - 多个关于缓存支持的请求
- 🔗 [LangChain Issue #26701](https://github.com/langchain-ai/langchain/issues/26701) - 实现困难

### 解决方案

**Autocache** 作为一个透明代理工作，它会自动分析您的请求，并在最佳断点处注入缓存控制头——**无需更改代码**。只需将现有的 n8n/Flowise/Make.com 工作流指向 Autocache，而不是直接指向 Anthropic 的 API。

**结果**：相同的代理，成本降低 90%，延迟降低 85% —— 自动实现。

## 替代方案与对比

有几种工具提供提示词缓存支持，但 Autocache 在结合 **零配置透明代理** 与 **智能 ROI 分析** 方面是独一无二的：

### 现有解决方案

| 解决方案                                                                             | 类型    | 自动注入          | 智能程度                    | ROI 分析       | 适用于 n8n/Flowise |
| ------------------------------------------------------------------------------------ | ------- | ----------------------- | ------------------------------- | ------------------- | ----------------------- |
| **Autocache**                                                                  | 代理   | ✅ 完全自动      | ✅ Token 分析 + ROI 评分 | ✅ 响应头 | ✅ 是                  |
| [LiteLLM](https://docs.litellm.ai/docs/tutorials/prompt_caching)                        | 代理   | ⚠️ 需要配置    | ❌ 基于规则                   | ❌ 否               | ✅ 是                  |
| [langchain-smart-cache](https://github.com/imranarshad/langchain-anthropic-smart-cache) | 库 | ✅ 完全自动      | ✅ 基于优先级               | ✅ 统计信息       | ❌ 仅限 LangChain       |
| [anthropic-cost-tracker](https://github.com/Supgrade/anthropic-API-cost-tracker)        | 库 | ❓ 不明确              | ❓ 未知                      | ✅ 仪表板        | ❌ 仅限 Python          |
| OpenRouter                                                                           | 服务 | ⚠️ 取决于提供商 | ❌ 否                           | ❌ 否               | ✅ 是                  |
| AWS Bedrock                                                                          | 云   | ✅ 基于机器学习             | ✅ 是                          | ✅ 仅限 AWS         | ❌ 仅限 AWS             |

### 为什么选择 Autocache

**Autocache 是唯一结合了以下特性的解决方案：**

1. 🔄 **透明代理** - 适用于任何工具（n8n, Flowise, Make.com），无需更改代码。
2. 🧠 **智能分析** - 自动 Token 计数、ROI 评分和最佳缓存断点放置。
3. 📊 **实时 ROI** - 每次响应头中都包含成本节省和盈亏平衡分析。
4. 🏠 **私有部署** - 无外部依赖，无云厂商锁定。
5. ⚙️ **零配置** - 开箱即用，支持多种策略（保守/中等/激进）。

**其他解决方案**通常需要繁琐的配置 (LiteLLM)、框架锁定 (langchain-smart-cache)，或者不为代理构建者提供透明代理功能。

## 功能特性

✨ **无缝替换**：只需更改 API URL 即可获得自动缓存。
📊 **ROI 分析**：通过响应头提供详细的成本节省和盈亏平衡分析。
🎯 **智能缓存**：使用多种策略智能放置缓存断点。
⚡ **高性能**：同时支持流式（streaming）和非流式请求。
🔧 **可配置**：多种缓存策略和可自定义的阈值。
🐳 **Docker 就绪**：使用 Docker 和 docker-compose 轻松部署。
📋 **全面日志**：具有结构化输出的详细请求/响应日志。

## 使用 Docker 快速开始

使用 Autocache 最快的方法是使用 GitHub Container Registry 发布的 Docker 镜像：

### 快速启动 (30 秒)

**1. 运行容器：**

```bash
# 选项 A：在环境变量中设置 API 密钥
docker run -d -p 8080:8080 \
  -e ANTHROPIC_API_KEY=sk-ant-... \
  --name autocache \
  ghcr.io/montevive/autocache:latest

# 选项 B：不设置 API 密钥（在每个请求的 Header 中传递）
docker run -d -p 8080:8080 \
  --name autocache \
  ghcr.io/montevive/autocache:latest
```

**2. 验证运行状态：**

```bash
curl http://localhost:8080/health
# {"status":"healthy","version":"1.0.1","strategy":"moderate"}
```

**3. 发送请求测试：**

```bash
# 如果使用选项 A（环境变量）：
curl http://localhost:8080/v1/messages \
  -H "Content-Type: application/json" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-3-5-haiku-20241022",
    "max_tokens": 50,
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# 如果使用选项 B（无环境变量），在 Header 中传递 API 密钥：
curl http://localhost:8080/v1/messages \
  -H "Content-Type: application/json" \
  -H "x-api-key: sk-ant-..." \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-3-5-haiku-20241022",
    "max_tokens": 50,
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

**4. 检查响应头中的缓存元数据：**

```bash
X-Autocache-Injected: true
X-Autocache-Cache-Ratio: 0.750
X-Autocache-ROI-Percent: 85.2
X-Autocache-Savings-100req: $1.75
```

### 可用的 Docker 标签

- `latest` - 最新的稳定版本（推荐）
- `v1.0.1` - 特定版本标签
- `1.0.1`, `1.0`, `1` - 语义化版本别名

### Docker 镜像详情

- **注册表**: `ghcr.io/montevive/autocache`
- **架构**: `linux/amd64`, `linux/arm64`
- **大小**: ~29 MB (基于优化的 Alpine 镜像)
- **源代码**: https://github.com/montevive/autocache

### 下一步

- 关于使用 docker-compose 的**生产部署**，请参阅下文的 [快速入门](#快速入门) 章节。
- 关于**配置选项**，请参阅 [配置](#配置)。
- 关于 **n8n 集成**，请参阅我们的 [n8n 设置指南](N8N_TESTING.md)。

## 快速入门

### 使用 Docker Compose (推荐)

1. **克隆并配置**：

```bash
git clone <repository-url>
cd autocache
cp .env.example .env
# 编辑 .env 文件，设置 ANTHROPIC_API_KEY（可选 - 也可以在 Header 中传递）
```

2. **启动代理**：

```bash
docker-compose up -d
```

3. **在应用中使用**：

```bash
# 将您的 API 基础 URL 从：
# https://api.anthropic.com
# 修改为：
# http://localhost:8080
```

### 直接运行

1. **构建并运行**：

```bash
go mod download
go build -o autocache ./cmd/autocache
# 选项 1：通过环境变量设置 API 密钥（可选）
ANTHROPIC_API_KEY=sk-ant-... ./autocache
# 选项 2：不带 API 密钥运行（在请求头中传递）
./autocache
```

2. **配置您的客户端**：

```python
# Python 示例 - 在 Header 中传递 API 密钥
import anthropic

client = anthropic.Anthropic(
    api_key="sk-ant-...",  # 这将被转发给 Anthropic
    base_url="http://localhost:8080"  # 指向 autocache
)
```

## 配置

### 环境变量

| 变量 | 默认值 | 描述 |
| ------------------------- | ------------ | -------------------------------------------------------------- |
| `PORT` | `8080` | 服务器端口 |
| `ANTHROPIC_API_KEY` | - | 您的 Anthropic API 密钥（如果在请求头中传递则可选） |
| `CACHE_STRATEGY` | `moderate` | 缓存策略：`conservative`/`moderate`/`aggressive`/`auto_aggressive` |
| `LOG_LEVEL` | `info` | 日志级别：`debug`/`info`/`warn`/`error` |
| `MAX_CACHE_BREAKPOINTS` | `4` | 最大缓存断点数 (1-4) |
| `TOKEN_MULTIPLIER` | `1.0` | Token 阈值乘数 |

### API 密钥配置

可以通过三种方式提供 Anthropic API 密钥（优先级从高到低）：

1. **请求头**（推荐用于多租户场景）：

   ```http
   Authorization: Bearer sk-ant-...
   # 或
   x-api-key: sk-ant-...
   ```
2. **环境变量**：

   ```bash
   ANTHROPIC_API_KEY=sk-ant-... ./autocache
   ```
3. **`.env` 文件**：

   ```bash
   ANTHROPIC_API_KEY=sk-ant-...
   ```

💡 **提示**：对于多用户环境（如具有多个 API 密钥的 n8n），请在请求头中传递密钥，并保持环境变量为空。

### 缓存策略

#### 🛡️ 保守型 (Conservative)

- **重点**：仅限系统提示词和工具。
- **断点数**：最多 2 个。
- **适用场景**：对成本极度敏感且内容可预测的应用。

#### ⚖️ 中等型 (Moderate) - 默认

- **重点**：系统提示词、工具和大型内容块。
- **断点数**：最多 3 个。
- **适用场景**：大多数平衡节省和效率的应用。

#### 🚀 激进型 (Aggressive)

- **重点**：最大的缓存覆盖范围。
- **断点数**：全部 4 个可用。
- **逻辑**：存在多个断点时自动将内容 TTL 升级为 1 小时，以确保协议合规。
- **适用场景**：具有重复内容的高吞吐量应用。

#### 🤖 自动激进型 (Auto-Aggressive)

- **重点**：智能多断点管理。
- **断点数**：全部 4 个可用。
- **逻辑**：自动检测现有的缓存断点，并智能地添加或合并新的断点。
- **适用场景**：具有迭代上下文的复杂 AI 代理（如 Claude Code）。

## ROI 分析

Autocache 通过响应头提供详细的 ROI 指标：

### 关键响应头

| 响应头 | 描述 |
| ------------------------------ | ------------------------------------------------ |
| `X-Autocache-Injected` | 是否应用了缓存 (`true`/`false`) |
| `X-Autocache-Cache-Ratio` | 缓存的 Token 百分比 (0.0-1.0) |
| `X-Autocache-ROI-Percent` | 大规模运行时的节省百分比 |
| `X-Autocache-ROI-BreakEven` | 达到盈亏平衡所需的请求数 |
| `X-Autocache-Savings-100req` | 发送 100 个请求后的总节省额 |

### 响应头示例

```http
X-Autocache-Injected: true
X-Autocache-Total-Tokens: 5120
X-Autocache-Cached-Tokens: 4096
X-Autocache-Cache-Ratio: 0.800
X-Autocache-ROI-FirstCost: $0.024
X-Autocache-ROI-Savings: $0.0184
X-Autocache-ROI-BreakEven: 2
X-Autocache-ROI-Percent: 85.2
X-Autocache-Breakpoints: system:2048:1h,tools:1024:1h,content:1024:5m
X-Autocache-Savings-100req: $1.75
```

## API 端点

### 主端点

```
POST /v1/messages
```

Anthropic `/v1/messages` 端点的无缝替代，具有自动缓存注入功能。

### 健康检查

```
GET /health
```

返回服务器健康状况和配置状态。

### 指标

```
GET /metrics
```

返回支持的模型、策略和缓存限制。

### 节省额分析

```
GET /savings
```

返回全面的 ROI 分析和缓存统计数据：

**响应包括：**

- **近期请求**：带有缓存元数据的近期请求完整历史。
- **聚合统计**：
  - 处理的总请求数。
  - 应用了缓存的请求数。
  - 处理的总 Token 数和缓存的 Token 数。
  - 平均缓存率。
  - 10 次和 100 次请求后的预计节省额。
- **调试信息**：
  - 按类型（系统、工具、内容）分类的断点。
  - 每种断点类型的平均 Token 数。
- **配置**：当前的缓存策略和历史记录大小。

**使用示例：**

```bash
curl http://localhost:8080/savings | jq '.aggregated_stats'
```

**响应示例：**

```json
{
  "aggregated_stats": {
    "total_requests": 25,
    "requests_with_cache": 20,
    "total_tokens_processed": 125000,
    "total_tokens_cached": 95000,
    "average_cache_ratio": 0.76,
    "total_savings_after_10_reqs": "$1.85",
    "total_savings_after_100_reqs": "$18.50"
  },
  "debug_info": {
    "breakpoints_by_type": {
      "system": 15,
      "tools": 12,
      "content": 8
    },
    "average_tokens_by_type": {
      "system": 2048,
      "tools": 1536,
      "content": 1200
    }
  }
}
```

**用例：**

- 📊 监控缓存随时间推移的有效性。
- 🔍 调试缓存注入决策。
- 💰 跟踪实际成本节省。
- 📈 分析哪些内容类型从缓存中受益最多。

## 高级用法

### 绕过缓存

添加以下 Header 以跳过缓存注入：

```http
X-Autocache-Bypass: true
# 或
X-Autocache-Disable: true
```

### 自定义配置

```bash
# 开启调试日志的激进缓存
CACHE_STRATEGY=aggressive LOG_LEVEL=debug ./autocache

# 具有更高阈值的保守缓存
CACHE_STRATEGY=conservative TOKEN_MULTIPLIER=1.5 ./autocache
```

### 生产部署示例

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  autocache:
    image: autocache:latest
    environment:
      - LOG_JSON=true
      - LOG_LEVEL=info
      - CACHE_STRATEGY=aggressive
    ports:
      - "8080:8080"
    restart: unless-stopped
```

## 工作原理

1. **请求分析**：分析传入的 Anthropic API 请求。
2. **Token 计数**：使用近似分词技术识别可缓存内容。
3. **智能注入**：在最佳断点处放置缓存控制字段：
   - 系统提示词 (1h TTL)
   - 工具定义 (1h TTL)
   - 大型内容块 (5m TTL)
4. **ROI 计算**：计算成本节省和盈亏平衡分析。
5. **请求转发**：将增强后的请求发送到 Anthropic API。
6. **响应增强**：将 ROI 元数据添加到响应头中。

## 缓存控制详情

### 支持的内容类型

- ✅ 系统消息 (System messages)
- ✅ 工具定义 (Tool definitions)
- ✅ 文本内容块 (Text content blocks)
- ✅ 消息内容 (Message content)
- ❌ 图像 (Images，根据 Anthropic 限制不可缓存)

### Token 要求

- **大多数模型**：最小 1024 tokens。
- **Haiku 模型**：最小 2048 tokens。
- **断点限制**：每个请求最多 4 个。

### TTL 选项

- **5 分钟**：动态内容，频繁更改。
- **1 小时**：稳定内容（系统提示词、工具）。

## 成本节省示例

### 示例 1：文档对话

```
请求：8,000 tokens (6,000 缓存的系统提示词 + 2,000 用户问题)
不使用缓存的成本：每次请求 $0.024
使用缓存的成本：
  - 首次请求：$0.027 (包含缓存写入)
  - 后续请求：$0.0066 (节省 90%)
  - 盈亏平衡点：2 次请求
  - 100 次请求后的节省：$1.62
```

### 示例 2：代码审查助手

```
请求：12,000 tokens (10,000 缓存的代码库 + 2,000 审查请求)
不使用缓存的成本：每次请求 $0.036
使用缓存的成本：
  - 首次请求：$0.045 (包含缓存写入)
  - 后续请求：$0.009 (节省 75%)
  - 盈亏平衡点：1 次请求
  - 100 次请求后的节省：$2.61
```

## 监控与调试

### 日志

```bash
# 详细记录缓存决策的调试模式
LOG_LEVEL=debug ./autocache

# 生产环境的 JSON 日志
LOG_JSON=true LOG_LEVEL=info ./autocache
```

### 关键日志字段

- `cache_injected`: 是否应用了缓存。
- `cache_ratio`: 缓存的 Token 百分比。
- `breakpoints`: 使用的缓存断点数量。
- `roi_percent`: 实现的节省百分比。

### 健康监控

```bash
# 检查代理健康状况
curl http://localhost:8080/health

# 获取指标和配置
curl http://localhost:8080/metrics

# 获取全面的节省额分析
curl http://localhost:8080/savings | jq .

# 监控聚合统计数据
curl http://localhost:8080/savings | jq '.aggregated_stats'

# 检查断点分布
curl http://localhost:8080/savings | jq '.debug_info.breakpoints_by_type'
```

## 故障排除

### 常见问题

**❌ 未应用任何缓存**

- 检查 Token 计数是否达到最小值 (1024/2048)。
- 确认内容是否可缓存（非图像）。
- 查看缓存策略配置。

**❌ 盈亏平衡点过高**

- 内容可能太小，无法进行有效缓存。
- 考虑使用更保守的策略。
- 检查 `TOKEN_MULTIPLIER` 设置。

**❌ API 密钥错误**

- 确保设置了 `ANTHROPIC_API_KEY` 或在 Header 中传递了密钥。
- 验证 API 密钥格式：`sk-ant-...`

### 调试模式

```bash
LOG_LEVEL=debug ./autocache
```

提供有关以下内容的详细信息：

- Token 计数决策。
- 缓存断点放置。
- ROI 计算。
- 请求/响应处理。

## 架构

Autocache 遵循 Go 最佳实践，采用清晰的模块化架构：

```
autocache/
├── cmd/
│   └── autocache/           # 应用入口点
├── internal/
│   ├── types/              # 共享数据模型
│   ├── config/             # 配置管理
│   ├── tokenizer/          # Token 计数（启发式、离线、基于 API）
│   ├── pricing/            # 成本计算和 ROI
│   ├── client/             # Anthropic API 客户端
│   ├── cache/              # 缓存注入逻辑
│   └── server/             # HTTP 处理程序和路由
└── test_fixtures.go        # 共享测试工具
```

### 核心组件

- **Server**: 支持流式传输的 HTTP 请求处理程序。
- **Cache Injector**: 具有 ROI 评分功能的智能缓存断点放置。
- **Tokenizer**: 多种实现（启发式、离线分词器、真实 API）。
- **Pricing Calculator**: ROI 和成本效益分析。
- **API Client**: 具有 Header 管理功能的 Anthropic API 通信。

有关详细的架构文档，请参阅 [CLAUDE.md](CLAUDE.md)。

## 贡献

1. Fork 本仓库。
2. 创建功能分支。
3. 为新功能添加测试。
4. 确保所有测试通过：`go test ./...`。
5. 提交 Pull Request。

## 许可证

MIT 许可证 - 详情见 LICENSE 文件。

## 支持

- 📧 邮箱: hi@montevive.ai
- 💬 问题: [GitHub Issues](https://github.com/montevive/autocache/issues)
- 📖 文档: [GitHub Wiki](https://github.com/montevive/autocache/wiki)

---

**Autocache** - 通过智能缓存最大化您的 Anthropic API 效率 🚀
