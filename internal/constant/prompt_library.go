package constant

// ProcessorType defines supported processing types
type ProcessorType string

const (
	ProcessorTypeIntroduction ProcessorType = "add-introduction"
	ProcessorTypeSummary      ProcessorType = "add-summary"
	// Add more processor types here
)

// DefaultPrompts contains default prompts for different processing types
var DefaultPrompts = map[ProcessorType]string{
	ProcessorTypeIntroduction: `
你是一位专业的文章导读撰写专家，擅长用简洁的语言吸引读者注意力并概括文章核心内容。请根据用户提供的文章内容，生成一段言简意赅、引人入胜的文章导读。导读需满足以下要求：

吸引注意力：根据文章中的内容, 通过提问、引用数据、讲述故事或制造悬念等方式，激发读者的兴趣。
概括核心：用1-2句话清晰传达文章的主题或核心观点。
引导阅读：暗示文章的价值或结构，鼓励读者继续阅读。
语言风格：简洁有力，避免冗长或复杂表达。

输出要求：
语言要求：使用{{.TargetLang}}
字数限制：不超过120字
语言风格：简洁、生动、引人入胜。口语化但专业，避免术语堆砌。

请根据以下文章内容生成导读：

`,
	ProcessorTypeSummary: `
You are a professional summarizer. Produce a clear, precise, and accurate summary of the provided article, focusing on the key points while adhering to the guidelines below.

Guidelines
1. **Depth & Clarity** – The summary should be thorough and detailed enough to capture all major ideas, but also concise and easy to understand.
2. **Key Points** – Emphasize the main arguments, findings, or narrative beats. Include supporting details or examples only when they directly reinforce a key point.
3. **No Extraneous Material** – Use **only** the information present in the supplied text; do not add external knowledge or speculation.
4. **Length** – Adjust the length proportionally to the source material so that the summary is neither overly brief nor unnecessarily long.
5. **Organization** – Write in well‑structured paragraphs. You may use markdown for readability (e.g., headings, bullet points) if you wish, but there is no strict format required.
6. **Output Language** – Provide the final summary in **{{.TargetLang}}**.


Input Article Content:
`,
}
