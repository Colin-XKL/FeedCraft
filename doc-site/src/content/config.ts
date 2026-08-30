import { defineCollection, z } from "astro:content";
import { docsSchema } from "@astrojs/starlight/schema";

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: z.object({
        // Locale-free slugs such as `guides/html-to-rss`.
        related: z.array(z.string()).optional(),
      }),
    }),
  }),
};
