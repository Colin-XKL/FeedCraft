# 命名重构实施规划 (Naming Refactoring Plan)

## 1. 目标 (Objective)

将 FeedCraft 项目的核心概念命名从 "Craft 体系" 迁移到 "工程/频道体系"，以提高概念清晰度和用户友好度。

| 旧名称 (Old)  | 新名称 (New)  | 中文名 (CN) | 定义 (Definition)             |
| :------------ | :------------ | :---------- | :---------------------------- |
| **CraftAtom** | **Tool**      | **工具**    | 基础处理单元 (如翻译、摘要)。 |
| **CraftFlow** | **Blueprint** | **蓝图**    | 处理逻辑的编排配置。          |
| **Recipe**    | **Channel**   | **频道**    | 持久化的订阅源配置。          |

---

## 2. 实施策略 (Strategy)

采用 **"后端优先 -> 接口迁移 -> 前端适配 -> 文档更新"** 的分阶段实施策略。
_建议在单独的 `refactor/naming` 分支进行开发。_

### 阶段一：后端核心重构 (Phase 1: Backend Core)

**目标**: 完成 Go 结构体、DAO 层和数据库迁移逻辑的更名。

1.  **数据模型 (Models & DAO)**:
    - 重命名 `internal/dao/craft_atom.go` -> `tool.go`
      - Struct: `CraftAtom` -> `Tool`
      - Table: `craft_atoms` -> `tools` (需编写 GORM 迁移脚本)
    - 重命名 `internal/dao/craft_flow.go` -> `blueprint.go`
      - Struct: `CraftFlow` -> `Blueprint`
      - Table: `craft_flows` -> `blueprints`
    - 重命名 `internal/dao/recipe.go` -> `channel.go`
      - Struct: `CustomRecipeV2` -> `Channel`
      - Table: `custom_recipes_v2` -> `channels`

2.  **业务逻辑 (Logic)**:
    - 更新 `internal/craft` 包中的引用。建议将 `craft` 包重命名为 `engine` 或 `processor` 以避免混淆，或者暂时保留包名但重命名内部变量。
    - 变量重命名: `craftName` -> `processorName` (泛指 Tool 或 Blueprint)。

3.  **数据库迁移 (Migration)**:
    - 在 `internal/dao/migrate.go` 中添加 `RenameTable` 逻辑，确保旧数据无损迁移到新表。

### 阶段二：API 与路由重构 (Phase 2: API & Routing)

**目标**: 更新 RESTful API 路径和 Controller 命名。

1.  **管理端 API (Admin API)**:
    - `GET /api/admin/craft-atoms` -> `/api/admin/tools`
    - `GET /api/admin/craft-flows` -> `/api/admin/blueprints`
    - `GET /api/admin/recipes` -> `/api/admin/channels`
    - _Action_: 更新 `internal/router/registry.go` 和对应的 `internal/controller` 文件。

2.  **公共服务 API (Public Serving API)**:
    - **持久化订阅**: 默认使用 `/channel/:id` 代替 `/recipe/:id`。
      - _兼容性_: **必须保留** `/recipe/:id` 入口，其逻辑与 `/channel/:id` 完全一致，确保老用户订阅地址不失效。
    - **即时处理 (Portable Mode)**: **保留** `/craft/:name` 路径。
      - _理由_: `craft` 作为一个动词（加工/处理）在此处依然非常贴切且符合品牌名。

### 阶段三：前端重构 (Phase 3: Frontend Refactor)

**目标**: 更新管理后台的 UI、变量和多语言文案。

1.  **API 客户端**:
    - `web/admin/src/api/craft_atom.ts` -> `tool.ts`
    - `web/admin/src/api/craft_flow.ts` -> `blueprint.ts`
    - `web/admin/src/api/custom_recipe.ts` -> `channel.ts`

2.  **视图组件 (Views)**:
    - 重命名目录:
      - `views/dashboard/craft_atom` -> `views/dashboard/tool`
      - `views/dashboard/craft_flow` -> `views/dashboard/blueprint`
      - `views/dashboard/custom_recipe` -> `views/dashboard/channel`
    - 更新组件内的变量名和引用。

3.  **多语言 (Locales)**:
    - 全面更新 `en-US` 和 `zh-CN` 中的文案。
    - Key 重命名: `craftAtom.*` -> `tool.*`, `customRecipe.*` -> `channel.*` 等。

### 阶段四：文档与清理 (Phase 4: Documentation & Cleanup)

**目标**: 确保所有外部文档与代码一致。

1.  **项目文档**:
    - 更新 `README.md`, `INTRODUCTION.md`, `AGENTS.md`。
    - 删除不再需要的临时规划文档 (`NAMING_*.md`)。

2.  **文档站点 (Doc Site)**:
    - 批量搜索替换 `doc-site/` 目录下的 Markdown 文件。
    - 更新截图 (ScreenShots) —— _这是一个较大的工作量，需重新截取管理后台图片_。

---

## 3. 风险控制 (Risk Management)

1.  **数据库兼容性**:
    - 使用 GORM 的 `Migrator().RenameTable()` 进行表名变更。
    - 保留旧表备份，直到确认新表工作正常。

2.  **URL 兼容性**:
    - 对于已分发的 RSS 链接 (如 `/recipe/xyz`)，需要在路由层保留 `Alias` (别名) 支持，确保老用户订阅不失效。
    - 建议保留 `/recipe/:id` 路由，但在 Log 中输出 Deprecation Warning。

3.  **外部依赖**:
    - 检查 Docker 环境变量命名 (`FC_...`) 是否有涉及旧术语的配置项，如有则需增加兼容逻辑。

## 4. 检查清单 (Checklist)

- [x] **Backend**: Structs renamed (Tool, Blueprint, Channel).
- [x] **Backend**: DB Tables renamed (tools, blueprints, channels).
- [x] **Backend**: Controllers & Routers updated.
- [x] **Backend**: Route `/channel/:id` added and `/recipe/:id` preserved for compatibility.
- [x] **Backend**: Route `/craft/:name` preserved.
- [ ] **Frontend**: API clients updated.
- [ ] **Frontend**: UI components renamed and updated.
- [ ] **Frontend**: Locales (Translation keys) updated.
- [ ] **Docs**: README & Introduction updated.
- [ ] **Docs**: Doc-site content updated.
