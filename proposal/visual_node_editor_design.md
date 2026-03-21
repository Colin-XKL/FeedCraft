# 可视化节点编排系统设计 (Visual Node-based Orchestration Design)

## 1. 需求场景与核心目标 (Scenario & Goals)

随着 FeedCraft 架构向“流式数据处理管线”演进，系统内部将存在大量不同层级的组合：多个 Atom-Craft 组成 Flow-Craft，Source 和 Flow-Craft 组成 Recipe，多个 Recipe 汇聚成 Topic。
对于管理员而言，纯表单或 JSON 的管理方式会随着组合的复杂化而变得难以维护。

**核心目标：**
1. **拓扑可视化 (Topology Visualization)**：将复杂的嵌套关系（如 Flow-Craft 内部的 Atom-Craft 链路）展开为直观的拓扑图。
2. **可拖拽交互 (Drag & Drop UI)**：支持像搭建乐高积木一样，通过拖拽节点和连线来创建或修改 Recipe、Topic 和 Flow-Craft。
3. **参数配置所见即所得 (WYSIWYG Configuration)**：点击任意节点（如 AI 摘要节点），可在侧边栏直接配置其专有参数（如 Prompt、模型选择）。

## 2. 核心界面抽象映射 (UI Abstraction Mapping)

在可视化画布（Canvas）上，我们将底层架构概念映射为标准的图（Graph）元素：

### 2.1 节点 (Nodes)
* **触发器/入口节点 (Trigger/Source Node)**：
  * 对应后端的 `FeedProvider`。如：RSS 抓取节点、网页转换节点、搜索节点。
  * 视觉特征：只有输出端口（Output Port），没有输入端口。
* **加工节点 (Processor Node)**：
  * 对应后端的 `FeedProcessor` (Atom-Craft)。如：翻译节点、AI摘要节点、去广告节点。
  * 视觉特征：同时具有输入和输出端口（Input & Output Ports）。
* **聚合/结构节点 (Router/Struct Node)**：
  * 负责多流合并或条件分发。如：Topic 中的合并去重器 (Merge Aggregator)。
  * 视觉特征：具有多个输入端口，一个输出端口。
* **终点节点 (Output Node)**：
  * 业务上的最终输出端点。如：生成最终的 XML Feed 地址。

### 2.2 连线 (Edges)
* 代表了 `CraftFeed` (包含了多条 `CraftArticle`) 在各个节点之间的流转。
* 可以在连线上增加动画（如流动的小圆点）来表示数据流向。

## 3. 视图层级设计 (View Hierarchy Management)

FeedCraft 的业务模型具有明显的嵌套特征（Topic 包含 Recipe，Recipe 包含 Flow-Craft）。在 UI 设计上，我们提供两种查看和管理模式，以防止单张画布过于拥挤：

### 模式一：宏观/聚合画布 (Topic Workflow View)
用于管理 `Topic`。
* **展现形式**：多个独立的分支（各个 Recipe 的入口）汇总到一个聚合节点（Aggregator），再连接到全局加工节点（如全局排序、全局 AI 裂变），最终输出。
* **节点的折叠与展开**：在 Topic 画布中，一个 `Recipe` 默认显示为一个“复合节点 (Group Node)”。用户可以双击该节点，展开查看其内部的 `Flow-Craft` 细节，或者直接下钻（Drill-down）进入“微观画布”。

### 模式二：微观/流水线画布 (Recipe / Flow-Craft View)
用于管理单一 `Recipe` 或纯粹的 `Flow-Craft`（处理链条）。
* **展现形式**：线性的或简单的单路分支结构。例如：`HtmlSource` -> `AdFilter (Atom)` -> `Translate (Atom)` -> `Output`。
* 在这里，管理员可以精细调节每一个 `Atom-Craft` 的先后顺序，或者在中间插入新的数据处理节点。

## 4. 前后端交互数据结构 (Data Structure Design)

为了支持前端的可视化（如使用 React Flow 或 Vue Flow 等开源库），后端除了保存原有的业务模型配置外，还需要支持将其转换为 DAG（有向无环图）的标准化 JSON 格式。

**前端所需的核心数据结构示例：**
```json
{
  "nodes": [
    {
      "id": "source-1",
      "type": "rss_source",
      "position": { "x": 100, "y": 200 },
      "data": { "url": "https://example.com/rss", "label": "官方RSS" }
    },
    {
      "id": "craft-1",
      "type": "atom_translate",
      "position": { "x": 400, "y": 200 },
      "data": { "target_lang": "zh", "label": "中英翻译" }
    }
  ],
  "edges": [
    {
      "id": "edge-1",
      "source": "source-1",
      "target": "craft-1",
      "animated": true
    }
  ]
}
```
**后端的适配挑战**：后端的执行引擎需要能够动态解析这种基于 Node 和 Edge 的图结构，并将其还原为内存中的 `FeedProvider` 和 `FeedProcessor` 的调用链。

## 5. 渐进式演进路径 (Evolution Path)

开发一个完整的 n8n 级别可视化编辑器工程量巨大，建议分为三个阶段稳步推进：

### Phase 1: 只读拓扑可视化 (Read-only Topology Viewer)
* **目标**：方便查阅，不改变现有后端的配置管理方式。
* **实现**：后台依然通过表单增删改查 Topic、Recipe 和 Craft。前端额外增加一个“拓扑图视图”标签页，读取后端的配置后，自动计算布局（例如使用 dagre.js），在画布上渲染出当前配置的数据流向拓扑图。仅供查看和理清关系。

### Phase 2: 线性编排与参数配置 (Linear Drag & Drop Editor)
* **目标**：支持纯线性的拖拽编辑。
* **实现**：在编辑 `Recipe` (或 `Flow-Craft`) 时，引入可视化画布。用户可以从左侧组件库拖拽 `Atom-Craft` 到画布上，首尾相连形成直线。点击节点在右侧抽屉面板修改参数。保存时，前端将图结构序列化为后端可接受的顺序数组。

### Phase 3: 全自由度 DAG 编排 (Full DAG Orchestration)
* **目标**：打破 Topic、Recipe 和 Flow-Craft 之间僵硬的表结构边界，实现终极的自由编排。
* **实现**：整个系统只有一个“Workflow（工作流）”的概念。用户可以在同一张巨大画布上，任意拉取多个 Source，任意进行 Merge 和 Split，任意在各处安插 Craft。后端引擎升级为完全基于图数据结构的执行调度器。

---
*注：当前阶段（MVP）我们应专注于打牢底层的 `CraftFeed` 和两级接口抽象。本设计方案主要确立了系统最终的产品形态愿景，确保当下的底层重构不会偏离未来前端可视化的要求。*