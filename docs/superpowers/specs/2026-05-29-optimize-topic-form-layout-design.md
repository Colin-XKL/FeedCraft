# 优化 Topic 创建页面表单布局设计方案

## 1. 痛点分析
当前 Topic（主题）创建/编辑弹窗（Modal）中的“输入源（Input Sources）”与“聚合规则（Aggregator Steps）”表单布局不合理，导致以下严重体验问题：
* **横向挤压与折行混乱**：由于 Arco Design 的 `a-form-item` 默认对多子节点采用弹性盒横向排列布局（`.arco-form-item-content` 为 `display: flex; flex-direction: row` 样式），多个并列的输入源组件 `.input-source-row` 和“添加输入源”按钮被强行挤在同一行展示。这导致 Radio Group 中的选项（“外部 RSS”, “Recipe”, “Topic”）文字因空间严重受限而发生垂直折行，输入框和下拉框几乎被压缩至不可见，整体视觉极其凌乱。
* **规则配置缺乏清晰的物理区块感**：动态表单项（可自由增删的行）如果直接平铺在白色背景上，加上上述的挤压问题，用户难以一眼分辨每行输入项的起止界限，也无法感知各个输入项之间的独立层级。
* **移动与操作控件的对齐及微调不精细**：原有的 grid 或 flex 布局没有实现纵向完美居中，导致 Radio Group 与 Input/Select 高度参差不齐。

---

## 2. 优化方案：区块卡片化与垂直流式容器 (Block Card & Vertical Container)
为了彻底解决上述痛点，同时保证在 860px 宽度的 Modal 中完美展示，我们采用**方案 2（区块卡片化 + 垂直 Flex 布局）**。

### 2.1 结构与布局调整 (HTML/CSS)
我们将所有动态表单项列表和操作按钮，使用专门的垂直容器包裹，从而让它们在 `a-form-item` 中作为**单个子元素**被渲染，从而独占 100% 宽度，并垂直向下排布。

```html
<!-- 输入源部分包裹容器 -->
<div class="input-sources-container">
  <div v-for="(source, idx) in formData.inputSources" ... class="input-source-card">
    <!-- 单个输入源卡片内部布局 -->
  </div>
  <a-button type="dashed" long @click="addSource" class="add-button">
    <icon-plus /> 添加输入源
  </a-button>
</div>

<!-- 聚合规则部分包裹容器 -->
<div class="steps-container">
  <div v-for="(step, idx) in formData.aggregator_config" ... class="step-card">
    <!-- 单个步骤卡片内部布局 -->
  </div>
  <a-button type="dashed" long @click="addStep" class="add-button">
    <icon-plus /> 添加规则
  </a-button>
</div>
```

### 2.2 视觉美化：卡片化区块 (Block Card)
每一项动态条目（即每个输入源、每个聚合步骤）都将被重构为一个微型的“卡片/色块”：
* **背景色**：`background-color: var(--color-fill-1)`（采用 Arco Design 内置的轻量填充色，护眼且有柔和的对比度）
* **边框**：`border: 1px solid var(--color-border-1)` 或 `var(--color-fill-2)`（淡雅的灰色边框，进一步勾勒轮廓）
* **圆角**：`border-radius: 6px`（使界面更加精致现代）
* **内边距**：`padding: 12px 16px`（提供充足的呼吸感）
* **底部边距**：`margin-bottom: 12px`（确保项与项之间有良好的物理间距）
* **悬停效果**：在鼠标悬浮时微调边框色或阴影，提升交互回馈感（`border-color: var(--color-border-3)`）。

### 2.3 极致行内对齐 (Grid & Flex Alignment)
在每个卡片内部：
* **输入源行内部组件**：
  * 使用 Grid 布局：`grid-template-columns: auto 1fr auto;`
  * 配合 `align-items: center;` 和较大的 `gap: 12px`
  * 移除按钮（Remove Button）进行视觉降噪与突出：采用 `type="text" status="danger"` 并搭配删除图标 `icon-delete`。
  * `a-radio-group` 容器使用绝对不折行设置 `white-space: nowrap; flex-shrink: 0;`。
* **聚合步骤行内部组件**：
  * 精细调节每一步的编辑行 `.editor-row`
  * 类型选择 `a-select`（如去重/排序/限制数量）拥有恰当的固定宽度（如 160px）
  * 具体选项选择 `a-select`（如去重策略/排序方式）和阈值输入 `a-input-number` 自适应排开。
  * 卡片内部的提示信息 `.step-hint` 紧贴在输入控件下方，具有适当的左边距和更弱的灰色字体，优雅呈现对 SimHash / Embedding 算法的通俗解释。

---

## 3. 受影响文件
* `web/admin/src/views/dashboard/topic_feed/topic_feed.vue`：主要的 Vue 模板和样式修改。
* `web/admin/src/locale/zh-CN/topic.ts`（非必须，如果有必要，检查是否有翻译缺失）：无。

---

## 4. 自检清单 (Spec Self-Review)
* **占位符检查**：不含 TBD 或 TODO 等含糊字段。
* **一致性检查**：卡片样式与 Arco Design 内置风格（如 Arco Vue 的其它表单卡片、详情页）高度契合，保持视觉一致性。
* **边界与极限情况**：
  * **0 个或 1 个输入源**：通过 Vue `removeSource` 的逻辑，若删光了输入源会自动保留/追加 1 个默认外部 RSS，体验正常。
  * **宽屏与窄屏自适应**：在 860px 宽度的 Modal 中，由于使用了 `1fr` 弹性列，输入框会自适应填满剩余空间，在普通分辨率下均不会发生重叠。
* **双语支持**：完全复用已有的 `$t(...)` 多语言函数，不破坏已有的中英文国际化，若新增极少数文案则同步配置对应语言包。

---

## 5. 实现与测试计划

### 5.1 实现步骤
1. **重构模板代码**：
   * 将输入源循环和添加按钮包裹进 `.input-sources-container` 中，把 `.input-source-row` 升级为 `.input-source-card`。
   * 将聚合规则循环和添加按钮包裹进 `.steps-container` 中，把 `.step-wrapper` 升级为 `.step-card`。
   * 优化删除按钮，加入 `icon-delete`，增强辨识度。
2. **优化 CSS 样式**：
   * 定义容器类 `input-sources-container` 和 `steps-container`，使用 flex-column 排布。
   * 完善 `.input-source-card` 和 `.step-card` 视觉样式（轻背景、浅圆角、轻边框、适当 gap、悬浮状态）。
   * 精细微调行内元素垂直对齐、间距等。

### 5.2 验证测试
1. 启动前端开发服务器：在 `web/admin` 目录下运行 `pnpm run dev`（根据 `DEVELOPMENT.md`，可结合后端服务一起查看效果）。
2. 使用 `computerUse` 子代理或手动截图进行视觉验证。
3. 录制最终效果视频，确保：
   * Modal 弹出时没有任何挤压或折行。
   * 添加/删除输入源功能正常。
   * 添加/删除聚合步骤功能正常，SimHash / Embedding 下拉及阈值输入正常。
   * 校验配置、保存、取消均行为正常。
4. 运行前端 lint 和 build 检查，确保代码符合规范且可以成功编译。
