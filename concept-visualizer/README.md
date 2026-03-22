# FeedCraft Concept Visualizer

This is a standalone Vue 3 application designed to visually explain the core concepts and data flow architecture of [FeedCraft](https://github.com/Colin-XKL/FeedCraft) to new users. It uses smooth, cross-page "shared element" transitions to demonstrate how data moves through the system.

## 🌟 The Core Concepts (The Animation Flow)

The visualizer is built as a 5-step interactive scrollytelling/onboarding experience. A central "magic block" (representing a piece of data) travels through the following stages:

### 1. AtomCraft (The Atomic Processor)

- **Concept**: The smallest unit of processing. If you already have a standard RSS feed (like Hacker News), you can pass it through a single `AtomCraft` to instantly apply an action (e.g., Extract Fulltext, AI Translate, AI Summary).
- **Visual Flow**: A standard RSS icon enters a single processing node. The magic block changes color/shape, showing the immediate transformation.

### 2. FlowCraft (The Pipeline)

- **Concept**: Why stop at one action? A `FlowCraft` is a sequence of multiple chained `AtomCraft`s acting like LEGO blocks.
- **Visual Flow**: The data block travels along a pipeline, sequentially passing through multiple modules (Extract -> AI Summary -> Translate), accumulating transformations along the way.

### 3. Source Generators & RawFeed (Everything is a Feed)

- **Concept**: What if a site doesn't have an RSS feed? FeedCraft can capture unstructured data from raw HTML pages, Search Engine results, or API responses (Curl). This normalized starting point is called a `RawFeed`.
- **Visual Flow**: Disparate icons (HTML, Search, API) converge into a standardized "RawFeed" data block, demonstrating the universal input capability.

### 4. RecipeFeed (The Binding)

- **Concept**: A `Recipe` is the permanent binding between a specific Input Source (`RawFeed`) and a processing pipeline (`FlowCraft`). Once bound, it yields a permanent, customized RSS URL that automatically updates.
- **Visual Flow**: Shows the formula `RawFeed + FlowCraft = RecipeFeed`. The data block enters a "Recipe Active" container, generating a shiny, subscription-ready RSS URL.

### 5. TopicFeed (The Aggregation Tree)

- **Concept**: To solve information overload, multiple feeds (`RecipeFeed`s) can be aggregated into a single `TopicFeed`. TopicFeeds can be nested infinitely (Sub-Topics) to build an auto-deduplicated knowledge tree.
- **Visual Flow**: A glowing network topology graph. Smaller nodes (Recipe Feeds and Sub-Topics) shoot data streams along connecting lines into a massive central `TopicFeed` aggregator.

## 🛠️ Technology Stack

- **Framework**: [Vue 3](https://vuejs.org/) (Composition API) + [Vite](https://vitejs.dev/) + TypeScript
- **Styling**: [Tailwind CSS](https://tailwindcss.com/) for rapid, responsive, glassmorphism UI design.
- **Animations**: [GSAP](https://gsap.com/) & **GSAP Flip Plugin**. The `Flip` plugin is the engine behind the "Shared Element Transition"—calculating the absolute position of the floating data block as it moves seamlessly between different anchor points across DOM changes.
- **Slider**: [Swiper.js](https://swiperjs.com/) for the horizontal paging framework.
- **Icons**: [Lucide Vue Next](https://lucide.dev/)

## 🚀 How to Run Locally

1.  Navigate to the visualizer directory:
    ```bash
    cd concept-visualizer
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Start the dev server:
    ```bash
    npm run dev
    ```

## 📦 Future Integration

This independent application is designed to be embedded into the official FeedCraft documentation site (`doc-site` driven by Astro) as an interactive landing page or conceptual guide component.
