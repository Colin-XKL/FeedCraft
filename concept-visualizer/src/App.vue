<script setup lang="ts">
import { onMounted, ref, nextTick } from "vue";
import { Swiper, SwiperSlide } from "swiper/vue";
import { Navigation, Pagination, EffectFade } from "swiper/modules";
import gsap from "gsap";
import { Flip } from "gsap/Flip";
import {
  Zap,
  ArrowRight,
  Layers,
  Rss,
  Combine,
  Filter,
  Languages,
  Search,
  ChevronRight,
  ChevronLeft,
  Sparkles,
  Database,
  Cpu,
  FileJson,
} from "lucide-vue-next";

import "swiper/css";
import "swiper/css/navigation";
import "swiper/css/pagination";
import "swiper/css/effect-fade";

gsap.registerPlugin(Flip);

const modules = [Navigation, Pagination, EffectFade];
const swiperInstance = ref<any>(null);
const activeIndex = ref(0);
const floatingItemRef = ref<HTMLElement | null>(null);

const onSwiper = (swiper: any) => {
  swiperInstance.value = swiper;
};

const updateFlip = async () => {
  await nextTick();
  const state = Flip.getState(floatingItemRef.value);

  // Find the anchor in the current active slide
  const activeAnchor = document.querySelector(
    `.swiper-slide-active .item-anchor`,
  ) as HTMLElement;

  if (activeAnchor && floatingItemRef.value) {
    activeAnchor.appendChild(floatingItemRef.value);

    // Change the floating item's color based on the step to show progression
    const colors = [
      "from-cyan-400 to-blue-600 shadow-blue-500/50 text-white", // P1: AtomCraft
      "from-indigo-400 to-purple-600 shadow-indigo-500/50 text-white", // P2: FlowCraft
      "from-slate-400 to-slate-600 shadow-slate-500/50 text-slate-100", // P3: Source Generators
      "from-emerald-400 to-teal-600 shadow-emerald-500/50 text-white", // P4: RecipeFeed
      "from-orange-400 to-red-600 shadow-orange-500/50 text-white", // P5: TopicFeed
    ];
    
    floatingItemRef.value.className = `w-16 h-16 bg-gradient-to-br rounded-2xl flex items-center justify-center z-50 transition-[background-color,box-shadow,color] duration-1000 shadow-[0_0_30px_currentColor] ${colors[activeIndex.value]}`;

    Flip.from(state, {
      duration: 0.8,
      ease: "power3.inOut",
      absolute: true,
      scale: true,
      onComplete: () => {
        animatePageElements(activeIndex.value);
      },
    });
  }
};

const onSlideChange = (swiper: any) => {
  activeIndex.value = swiper.activeIndex;
  updateFlip();
};

const animatePageElements = (index: number) => {
  if (index === 0) {
    // Page 1: AtomCraft
    gsap.fromTo(
      ".atom-craft-node",
      { scale: 0.8, opacity: 0 },
      { scale: 1, opacity: 1, duration: 0.6, ease: "back.out(1.5)" },
    );
  } else if (index === 1) {
    // Page 2: FlowCraft
    gsap.fromTo(
      ".craft-module",
      { y: 30, opacity: 0 },
      { y: 0, opacity: 1, duration: 0.6, stagger: 0.15, ease: "back.out(1.2)" },
    );
    gsap.fromTo(
      ".flow-line",
      { scaleX: 0, transformOrigin: "left center" },
      { scaleX: 1, duration: 0.8, ease: "power2.inOut" },
    );
  } else if (index === 4) {
    // Page 5: TopicFeed
    gsap.fromTo(
      ".topic-node",
      { scale: 0, opacity: 0 },
      {
        scale: 1,
        opacity: 1,
        duration: 0.6,
        stagger: 0.1,
        ease: "back.out(1.5)",
      },
    );
    gsap.fromTo(
      ".topic-line",
      { strokeDashoffset: 100 },
      { strokeDashoffset: 0, duration: 1, stagger: 0.2, ease: "power1.inOut" },
    );
  }
};

onMounted(() => {
  updateFlip();
});
</script>

<template>
  <div
    class="min-h-screen w-full bg-slate-950 flex items-center justify-center p-4 lg:p-12 font-sans select-none overflow-hidden relative"
  >
    <!-- Background glow -->
    <div
      class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-blue-900/10 rounded-full blur-[120px] pointer-events-none"
    ></div>

    <!-- The Aspect-Ratio Stage Container -->
    <div class="relative w-full max-w-6xl min-h-[600px] lg:aspect-video bg-slate-900/80 backdrop-blur-xl rounded-[2rem] border border-slate-700/50 shadow-2xl flex flex-col overflow-hidden ring-1 ring-white/10">
      
      <!-- The Shared Element (Floating) -->
      <div ref="floatingItemRef" id="floating-craft-item" class="w-16 h-16 bg-gradient-to-br from-cyan-400 to-blue-600 rounded-2xl flex items-center justify-center z-50 transition-[background-color,box-shadow,color] duration-1000 shadow-[0_0_30px_currentColor] text-white">
        <Cpu class="w-8 h-8 drop-shadow-md" v-if="activeIndex === 0" />
        <Sparkles
          class="w-8 h-8 drop-shadow-md"
          v-else-if="activeIndex === 1"
        />
        <Database
          class="w-8 h-8 drop-shadow-md"
          v-else-if="activeIndex === 2"
        />
        <Rss class="w-8 h-8 drop-shadow-md" v-else-if="activeIndex === 3" />
        <Combine class="w-8 h-8 drop-shadow-md" v-else />
      </div>

      <swiper
        :modules="modules"
        :slides-per-view="1"
        :space-between="0"
        :pagination="{ clickable: true }"
        @swiper="onSwiper"
        @slideChange="onSlideChange"
        class="flex-1 w-full h-full"
      >
        <!-- Page 1: AtomCraft -->
        <swiper-slide
          class="w-full h-full flex flex-col items-center justify-center pt-8 pb-12 relative px-4"
        >
          <div class="text-center z-10 px-4 md:px-8 mb-8 md:mb-12">
            <h2
              class="text-3xl md:text-5xl font-black mb-3 md:mb-4 tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-cyan-300 to-blue-500"
            >
              1. AtomCraft
            </h2>
            <p
              class="text-sm md:text-lg lg:text-xl text-slate-400 max-w-2xl mx-auto leading-relaxed"
            >
              Start simple. Use an <strong>AtomCraft</strong> to instantly
              process any existing RSS feed. Extract full text, translate, or
              summarize with a single operation.
            </p>
          </div>

          <div class="flex-1 w-full flex items-center justify-center relative">
            <div
              class="flex flex-col md:flex-row items-center space-y-6 md:space-y-0 md:space-x-12"
            >
              <div class="flex flex-col items-center">
                <Rss
                  class="w-12 h-12 md:w-16 md:h-16 text-orange-400 mb-2 opacity-50"
                />
                <span
                  class="text-[10px] md:text-xs font-bold text-slate-500 uppercase tracking-wider"
                  >Standard RSS</span
                >
              </div>

              <ArrowRight class="hidden md:block text-slate-600 w-8 h-8" />
              <!-- Mobile arrow down -->
              <ArrowRight class="md:hidden text-slate-600 w-6 h-6 rotate-90" />

              <!-- AtomCraft Processing Node -->
              <div
                class="atom-craft-node w-28 h-28 md:w-32 md:h-32 bg-slate-900 rounded-[2rem] border-2 border-cyan-500/50 shadow-[0_0_30px_rgba(6,182,212,0.2)] flex flex-col items-center justify-center relative z-10"
              >
                <div
                  class="absolute -top-3 bg-cyan-950 text-cyan-400 text-[9px] md:text-[10px] font-black px-3 py-1 rounded border border-cyan-500/50 uppercase tracking-widest whitespace-nowrap"
                >
                  AtomCraft
                </div>
                <!-- Anchor inside the single processor -->
                <div
                  class="item-anchor w-20 h-20 flex items-center justify-center"
                ></div>
              </div>

              <ArrowRight class="hidden md:block text-slate-600 w-8 h-8" />
              <!-- Mobile arrow down -->
              <ArrowRight class="md:hidden text-slate-600 w-6 h-6 rotate-90" />

              <div
                class="w-20 h-20 md:w-24 md:h-24 bg-slate-900/50 border border-slate-700 border-dashed rounded-2xl flex items-center justify-center"
              >
                <span
                  class="text-[9px] md:text-[10px] font-black text-slate-500 uppercase text-center"
                  >Processed<br />Output</span
                >
              </div>
            </div>
          </div>
        </swiper-slide>

        <!-- Page 2: FlowCraft (The Pipeline) -->
        <swiper-slide
          class="w-full h-full flex flex-col items-center justify-center pt-8 pb-12 relative px-4"
        >
          <div class="text-center z-10 px-4 md:px-8 mb-8 md:mb-12">
            <h2
              class="text-3xl md:text-5xl font-black mb-3 md:mb-4 tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-indigo-300 to-purple-500"
            >
              2. FlowCraft
            </h2>
            <p
              class="text-sm md:text-lg lg:text-xl text-slate-400 max-w-2xl mx-auto leading-relaxed"
            >
              Why stop at one? Combine multiple AtomCrafts into a
              <strong>FlowCraft</strong>. Build a pipeline that extracts text,
              generates an AI summary, and translates it all at once.
            </p>
          </div>

          <div
            class="flex-1 w-full flex items-center justify-center relative px-2 md:px-12"
          >
            <!-- Pipeline Container -->
            <div
              class="w-full max-w-4xl relative flex items-center justify-between"
            >
              <!-- Connecting Line Background -->
              <div
                class="absolute left-6 right-6 md:left-10 md:right-10 h-1 bg-slate-800 top-1/2 -translate-y-1/2 rounded-full z-0 overflow-hidden"
              >
                <div
                  class="flow-line w-full h-full bg-gradient-to-r from-slate-600 via-indigo-500 to-purple-600"
                ></div>
              </div>

              <!-- Starting Anchor -->
              <div
                class="item-anchor w-16 h-16 md:w-24 md:h-24 z-10 bg-slate-950 rounded-2xl flex items-center justify-center ring-2 md:ring-4 ring-slate-900 shadow-xl"
              ></div>

              <!-- FlowCraft Label Container -->
              <div
                class="absolute -top-8 md:-top-12 left-1/2 -translate-x-1/2 bg-indigo-500/10 border border-indigo-500/30 px-4 md:px-6 py-1.5 md:py-2 rounded-full text-indigo-300 text-[9px] md:text-xs font-black tracking-widest z-20 whitespace-nowrap"
              >
                FLOWCRAFT PIPELINE
              </div>

              <!-- AtomCraft Modules -->
              <div class="flex space-x-2 md:space-x-8 z-10 pt-4">
                <div
                  class="craft-module w-16 md:w-28 bg-slate-900 rounded-xl md:rounded-2xl border-2 border-indigo-500/50 flex flex-col items-center py-2 md:py-4 shadow-[0_0_20px_rgba(99,102,241,0.2)] transform transition-transform hover:-translate-y-2 relative"
                >
                  <div
                    class="hidden md:block absolute -top-3 text-[8px] bg-indigo-950 text-indigo-400 px-2 py-0.5 rounded border border-indigo-500/50"
                  >
                    AtomCraft
                  </div>
                  <div
                    class="w-8 h-8 md:w-10 md:h-10 bg-indigo-500/20 rounded-full flex items-center justify-center mb-1 md:mb-2 md:mt-2"
                  >
                    <Filter class="w-4 h-4 md:w-5 md:h-5 text-indigo-400" />
                  </div>
                  <span
                    class="text-[8px] md:text-[10px] font-black uppercase text-indigo-300 text-center leading-tight"
                    >Fulltext</span
                  >
                </div>

                <div
                  class="craft-module w-16 md:w-28 bg-slate-900 rounded-xl md:rounded-2xl border-2 border-purple-500/50 flex flex-col items-center py-2 md:py-4 shadow-[0_0_20px_rgba(168,85,247,0.2)] transform transition-transform hover:-translate-y-2 relative"
                >
                  <div
                    class="hidden md:block absolute -top-3 text-[8px] bg-purple-950 text-purple-400 px-2 py-0.5 rounded border border-purple-500/50"
                  >
                    AtomCraft
                  </div>
                  <div
                    class="w-8 h-8 md:w-10 md:h-10 bg-purple-500/20 rounded-full flex items-center justify-center mb-1 md:mb-2 md:mt-2"
                  >
                    <Sparkles class="w-4 h-4 md:w-5 md:h-5 text-purple-400" />
                  </div>
                  <span
                    class="text-[8px] md:text-[10px] font-black uppercase text-purple-300 text-center leading-tight"
                    >AI<br class="md:hidden" />
                    Summary</span
                  >
                </div>

                <div
                  class="craft-module w-16 md:w-28 bg-slate-900 rounded-xl md:rounded-2xl border-2 border-fuchsia-500/50 flex flex-col items-center py-2 md:py-4 shadow-[0_0_20px_rgba(217,70,239,0.2)] transform transition-transform hover:-translate-y-2 relative"
                >
                  <div
                    class="hidden md:block absolute -top-3 text-[8px] bg-fuchsia-950 text-fuchsia-400 px-2 py-0.5 rounded border border-fuchsia-500/50"
                  >
                    AtomCraft
                  </div>
                  <div
                    class="w-8 h-8 md:w-10 md:h-10 bg-fuchsia-500/20 rounded-full flex items-center justify-center mb-1 md:mb-2 md:mt-2"
                  >
                    <Languages class="w-4 h-4 md:w-5 md:h-5 text-fuchsia-400" />
                  </div>
                  <span
                    class="text-[8px] md:text-[10px] font-black uppercase text-fuchsia-300 text-center leading-tight"
                    >Translate</span
                  >
                </div>
              </div>

              <!-- Output Dummy -->
              <div
                class="w-16 h-16 md:w-24 md:h-24 z-10 bg-slate-900/50 border border-slate-700 border-dashed rounded-xl md:rounded-2xl flex items-center justify-center mt-4"
              >
                <span
                  class="text-[8px] md:text-[10px] font-black text-slate-600 uppercase text-center leading-tight"
                  >Processed</span
                >
              </div>
            </div>
          </div>
        </swiper-slide>

        <!-- Page 3: Source Generators -->
        <swiper-slide
          class="w-full h-full flex flex-col items-center justify-center pt-8 pb-12 relative px-4"
        >
          <div class="text-center z-10 px-4 md:px-8 mb-8 md:mb-12">
            <h2
              class="text-3xl md:text-5xl font-black mb-3 md:mb-4 tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-slate-100 to-slate-400"
            >
              3. Source Generators
            </h2>
            <p
              class="text-sm md:text-lg lg:text-xl text-slate-400 max-w-2xl mx-auto leading-relaxed"
            >
              No RSS? No problem. FeedCraft's built-in generators convert
              regular <strong>HTML</strong> pages,
              <strong>Search</strong> engine results, or API responses
              (<strong>Curl</strong>) into standard RSS feeds.
            </p>
          </div>

          <div class="flex-1 w-full flex items-center justify-center relative">
            <div
              class="flex flex-col md:flex-row items-center space-y-6 md:space-y-0 md:space-x-8 lg:space-x-16"
            >
              <!-- Source Icons -->
              <div
                class="grid grid-cols-2 gap-3 md:gap-6 relative w-32 md:w-48 mx-auto md:mx-0"
              >
                <div
                  class="w-14 h-14 md:w-20 md:h-20 bg-slate-800/80 rounded-2xl border border-slate-700 flex flex-col items-center justify-center shadow-lg transform -rotate-6"
                >
                  <Layers class="text-emerald-400 w-6 h-6 md:w-8 md:h-8 mb-1" />
                  <span
                    class="text-[8px] md:text-[10px] font-bold text-slate-400"
                    >HTML</span
                  >
                </div>
                <div
                  class="w-14 h-14 md:w-20 md:h-20 bg-slate-800/80 rounded-2xl border border-slate-700 flex flex-col items-center justify-center shadow-lg transform rotate-3 translate-y-2 md:translate-y-4"
                >
                  <Search class="text-blue-400 w-6 h-6 md:w-8 md:h-8 mb-1" />
                  <span
                    class="text-[8px] md:text-[10px] font-bold text-slate-400"
                    >SEARCH</span
                  >
                </div>
                <div
                  class="w-14 h-14 md:w-20 md:h-20 bg-slate-800/80 rounded-2xl border border-slate-700 flex flex-col items-center justify-center shadow-lg transform rotate-12 translate-x-2 md:translate-x-4"
                >
                  <FileJson class="text-amber-400 w-6 h-6 md:w-8 md:h-8 mb-1" />
                  <span
                    class="text-[8px] md:text-[10px] font-bold text-slate-400"
                    >API</span
                  >
                </div>

                <!-- Converge Lines (Decorative) -->
                <svg
                  class="hidden md:block absolute inset-0 w-[150%] h-[150%] -right-[50%] -top-[25%] pointer-events-none -z-10"
                  viewBox="0 0 100 100"
                  preserveAspectRatio="none"
                >
                  <path
                    d="M 80 20 L 100 50"
                    stroke="#334155"
                    stroke-width="1"
                    stroke-dasharray="4 4"
                    fill="none"
                  />
                  <path
                    d="M 80 80 L 100 50"
                    stroke="#334155"
                    stroke-width="1"
                    stroke-dasharray="4 4"
                    fill="none"
                  />
                </svg>
              </div>

              <ArrowRight
                class="hidden md:block text-slate-600 w-8 h-8 md:w-12 md:h-12"
              />
              <!-- Mobile arrow down -->
              <ArrowRight
                class="md:hidden text-slate-600 w-6 h-6 rotate-90 my-2"
              />

              <!-- Anchor for Shared Element -->
              <div class="relative">
                <div
                  class="absolute -inset-4 bg-slate-800/50 rounded-full blur-xl -z-10"
                ></div>
                <div
                  class="item-anchor w-28 h-28 md:w-40 md:h-40 bg-slate-900/90 rounded-[2rem] md:rounded-[2.5rem] flex items-center justify-center border-2 border-dashed border-slate-600 relative ring-2 md:ring-4 ring-slate-950 shadow-2xl"
                >
                  <div
                    class="absolute -top-6 md:-top-10 text-[9px] md:text-xs font-black text-slate-500 uppercase tracking-[0.2em] bg-slate-950 px-3 md:px-4 py-1 rounded-full border border-slate-800 whitespace-nowrap"
                  >
                    New RawFeed
                  </div>
                </div>
              </div>
            </div>
          </div>
        </swiper-slide>

        <!-- Page 4: RecipeFeed -->
        <swiper-slide
          class="w-full h-full flex flex-col items-center justify-center pt-8 pb-12 relative px-4"
        >
          <div class="text-center z-10 px-4 md:px-8 mb-4 md:mb-8">
            <h2
              class="text-3xl md:text-5xl font-black mb-3 md:mb-4 tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-emerald-300 to-teal-500"
            >
              4. RecipeFeed
            </h2>
            <p
              class="text-sm md:text-lg lg:text-xl text-slate-400 max-w-2xl mx-auto leading-relaxed"
            >
              Combine your Source Generator (RawFeed) with a FlowCraft pipeline
              to create a <strong>Recipe</strong>. This yields a permanent,
              customized RSS URL that automatically updates.
            </p>
          </div>

          <div class="flex-1 w-full flex items-center justify-center relative">
            <div class="flex flex-col items-center">
              <!-- Visual Formula -->
              <div
                class="flex flex-wrap justify-center items-center gap-2 mb-6 md:mb-8 text-slate-500 font-mono text-xs md:text-base font-bold bg-slate-900/50 px-4 md:px-6 py-2 rounded-full border border-slate-800"
              >
                <span class="text-slate-300">RawFeed</span>
                <span class="text-emerald-500">+</span>
                <span class="text-indigo-400">FlowCraft</span>
                <span class="text-emerald-500">=</span>
                <span class="text-emerald-400">RecipeFeed</span>
              </div>

              <div class="relative group mt-2">
                <!-- Glow Effect -->
                <div
                  class="absolute -inset-1 bg-gradient-to-r from-emerald-500 to-teal-500 rounded-[2rem] md:rounded-[3rem] blur opacity-25 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"
                ></div>

                <!-- Card Content -->
                <div
                  class="relative bg-slate-950 p-8 md:p-12 rounded-[2rem] md:rounded-[3rem] border border-slate-700 shadow-2xl flex flex-col items-center w-[260px] md:min-w-[320px]"
                >
                  <div
                    class="absolute -top-4 md:-top-5 bg-gradient-to-r from-emerald-500 to-teal-500 text-slate-950 px-4 md:px-6 py-1.5 md:py-2 rounded-full text-[10px] md:text-xs font-black tracking-widest uppercase shadow-lg"
                  >
                    Recipe Active
                  </div>

                  <div
                    class="item-anchor w-24 h-24 md:w-32 md:h-32 flex items-center justify-center mb-6 md:mb-8"
                  ></div>

                  <div class="w-full flex flex-col items-center">
                    <div
                      class="w-full h-1.5 md:h-2 bg-slate-800 rounded-full overflow-hidden mb-4 md:mb-6"
                    >
                      <div
                        class="h-full bg-gradient-to-r from-emerald-500 to-teal-400 w-full animate-pulse"
                      ></div>
                    </div>
                    <div
                      class="bg-slate-900 border border-slate-800 px-3 md:px-4 py-2 md:py-3 rounded-lg md:rounded-xl w-full flex items-center justify-between"
                    >
                      <span
                        class="text-[10px] md:text-xs font-mono text-slate-400 truncate mr-2 md:mr-4"
                        >feedcraft.io/rss/my-recipe</span
                      >
                      <Rss
                        class="w-3 h-3 md:w-4 md:h-4 text-emerald-500 shrink-0"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </swiper-slide>

        <!-- Page 5: TopicFeed (The Nesting) -->
        <swiper-slide
          class="w-full h-full flex flex-col items-center justify-center pt-8 pb-12 relative px-4"
        >
          <div class="text-center z-20 px-4 md:px-8 relative mb-4">
            <h2
              class="text-3xl md:text-5xl font-black mb-3 md:mb-2 tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-orange-400 to-red-500"
            >
              5. TopicFeed
            </h2>
            <p
              class="text-sm md:text-lg text-slate-400 max-w-2xl mx-auto leading-relaxed"
            >
              Information overload? Combine multiple
              <strong>RecipeFeeds</strong> into a single
              <strong>TopicFeed</strong>. Nest them infinitely to build an
              auto-deduplicated knowledge tree.
            </p>
          </div>

          <div
            class="flex-1 w-full relative overflow-hidden flex items-center justify-center mt-2"
          >
            <!-- SVG Connecting Lines (Absolute to container) -->
            <svg
              class="absolute inset-0 w-full h-full pointer-events-none z-0"
              preserveAspectRatio="xMidYMid slice"
            >
              <defs>
                <linearGradient
                  id="lineGrad"
                  x1="0%"
                  y1="0%"
                  x2="100%"
                  y2="100%"
                >
                  <stop offset="0%" stop-color="#10b981" stop-opacity="0.5" />
                  <stop offset="100%" stop-color="#f97316" stop-opacity="0.8" />
                </linearGradient>
                <linearGradient
                  id="lineGrad2"
                  x1="100%"
                  y1="0%"
                  x2="0%"
                  y2="100%"
                >
                  <stop offset="0%" stop-color="#10b981" stop-opacity="0.5" />
                  <stop offset="100%" stop-color="#f97316" stop-opacity="0.8" />
                </linearGradient>
                <linearGradient
                  id="lineGrad3"
                  x1="50%"
                  y1="100%"
                  x2="50%"
                  y2="0%"
                >
                  <stop offset="0%" stop-color="#64748b" stop-opacity="0.5" />
                  <stop offset="100%" stop-color="#f97316" stop-opacity="0.8" />
                </linearGradient>
              </defs>
              <!-- Lines to center (50%, 50%) -->
              <path class="topic-line md:stroke-[3px]" d="M 25% 30% L 50% 50%" stroke="url(#lineGrad)" stroke-width="2" stroke-dasharray="100" stroke-dashoffset="100" fill="none" />
              <path class="topic-line md:stroke-[3px]" d="M 75% 30% L 50% 50%" stroke="url(#lineGrad2)" stroke-width="2" stroke-dasharray="100" stroke-dashoffset="100" fill="none" />
              <path class="topic-line md:stroke-[3px]" d="M 50% 80% L 50% 50%" stroke="url(#lineGrad3)" stroke-width="2" stroke-dasharray="100" stroke-dashoffset="100" fill="none" />
            </svg>

            <!-- Network Nodes Container -->
            <div
              class="relative w-full max-w-lg md:max-w-2xl h-64 md:h-96 z-10"
            >
              <!-- Center Node (Parent Topic) -->
              <div
                class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 topic-node flex flex-col items-center"
              >
                <div
                  class="absolute -inset-4 md:-inset-8 bg-orange-500/20 rounded-full blur-xl md:blur-2xl -z-10"
                ></div>
                <div
                  class="w-28 h-28 md:w-40 md:h-40 rounded-full border-2 md:border-4 border-orange-500/80 bg-slate-950 flex flex-col items-center justify-center shadow-[0_0_20px_rgba(249,115,22,0.4)] relative"
                >
                  <div
                    class="absolute -top-3 md:-top-4 text-[8px] md:text-[10px] font-black text-orange-400 bg-orange-950 px-2 md:px-3 py-0.5 md:py-1 rounded-full border border-orange-500/50 whitespace-nowrap"
                  >
                    TOPIC FEED
                  </div>
                  <!-- The Final Destination Anchor -->
                  <div
                    class="item-anchor w-16 h-16 md:w-20 md:h-20 flex items-center justify-center z-20"
                  ></div>
                </div>
              </div>

              <!-- Top Left Node (Sub Topic A) -->
              <div
                class="absolute left-[20%] md:left-[25%] top-[25%] md:top-[30%] -translate-x-1/2 -translate-y-1/2 topic-node flex flex-col items-center"
              >
                <div
                  class="w-16 h-16 md:w-24 md:h-24 rounded-full border-2 border-emerald-500/80 bg-slate-900 flex flex-col items-center justify-center shadow-lg shadow-emerald-500/20 relative"
                >
                  <div
                    class="absolute -top-3 text-[7px] md:text-[9px] font-black text-emerald-400 bg-emerald-950 px-1.5 md:px-2 py-0.5 rounded-full border border-emerald-500/50 whitespace-nowrap"
                  >
                    RECIPE FEED A
                  </div>
                  <Rss class="w-4 h-4 md:w-6 md:h-6 text-emerald-400" />
                </div>
              </div>

              <!-- Top Right Node (Sub Topic B) -->
              <div
                class="absolute left-[80%] md:left-[75%] top-[25%] md:top-[30%] -translate-x-1/2 -translate-y-1/2 topic-node flex flex-col items-center"
              >
                <div
                  class="w-16 h-16 md:w-24 md:h-24 rounded-full border-2 border-emerald-500/80 bg-slate-900 flex flex-col items-center justify-center shadow-lg shadow-emerald-500/20 relative"
                >
                  <div
                    class="absolute -top-3 text-[7px] md:text-[9px] font-black text-emerald-400 bg-emerald-950 px-1.5 md:px-2 py-0.5 rounded-full border border-emerald-500/50 whitespace-nowrap"
                  >
                    RECIPE FEED B
                  </div>
                  <Rss class="w-4 h-4 md:w-6 md:h-6 text-emerald-400" />
                </div>
              </div>

              <!-- Bottom Node (Raw Recipe Feed) -->
              <div
                class="absolute left-[50%] top-[85%] md:top-[80%] -translate-x-1/2 -translate-y-1/2 topic-node flex flex-col items-center"
              >
                <div
                  class="w-14 h-14 md:w-20 md:h-20 rounded-xl md:rounded-2xl border-2 border-slate-600 bg-slate-800 flex flex-col items-center justify-center shadow-lg relative"
                >
                  <div
                    class="absolute -bottom-3 text-[7px] md:text-[9px] font-black text-slate-400 bg-slate-950 px-1.5 md:px-2 py-0.5 rounded-full border border-slate-600 whitespace-nowrap"
                  >
                    SUB TOPIC
                  </div>
                  <Combine class="w-4 h-4 md:w-6 md:h-6 text-slate-400" />
                </div>
              </div>
            </div>
          </div>
        </swiper-slide>
      </swiper>

      <!-- Navigation Overlay (Fixed within the aspect-ratio container) -->
      <div
        class="absolute bottom-8 left-0 right-0 flex justify-between px-8 md:px-12 z-40 pointer-events-none"
      >
        <button
          @click="swiperInstance?.slidePrev()"
          class="p-4 bg-slate-800/80 backdrop-blur rounded-full border border-slate-700 pointer-events-auto hover:bg-slate-700 transition-all shadow-lg"
          :class="{
            'opacity-0 scale-90 pointer-events-none': activeIndex === 0,
          }"
        >
          <ChevronLeft class="w-6 h-6 text-slate-300" />
        </button>
        <button
          @click="swiperInstance?.slideNext()"
          class="p-4 bg-slate-800/80 backdrop-blur rounded-full border border-slate-700 pointer-events-auto hover:bg-slate-700 transition-all shadow-lg"
          :class="{
            'opacity-0 scale-90 pointer-events-none': activeIndex === 4,
          }"
        >
          <ChevronRight class="w-6 h-6 text-slate-300" />
        </button>
      </div>

      <!-- Pagination Custom Styling override via Swiper inject -->
      <div class="absolute top-8 left-0 right-0 flex justify-center z-40">
         <div class="flex space-x-2">
            <div v-for="i in 5" :key="i" 
                 @click="swiperInstance?.slideTo(i - 1)" class="h-1.5 rounded-full transition-all duration-300 cursor-pointer hover:bg-slate-400"
                 :class="i - 1 === activeIndex ? 'w-8 bg-blue-500' : 'w-2 bg-slate-700'">
            </div>
         </div>
      </div>
    </div>
  </div>
</template>

<style>
/* Override default swiper pagination since we built a custom one visually */
.swiper-pagination {
  display: none !important;
}

#floating-craft-item {
  pointer-events: none;
}
/* Ensure the item stays inside the anchor even if it's rounded */
.item-anchor {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
