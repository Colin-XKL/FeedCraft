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
语言要求：使用简体中文
字数限制：中文不超过120字
语言风格：简洁、生动、引人入胜。口语化但专业，避免术语堆砌。

请根据以下文章内容生成导读：

`,
	ProcessorTypeSummary: `
Act as a professional summarizer. Create a concise and comprehensive summary of the text. while adhering to the guidelines enclosed below. 

Guidelines:  
- Create a summary that is detailed, thorough, in-depth, and complex, while maintaining clarity and conciseness. 
The summary must cover all the key points and main ideas presented in the original text, while also condensing the information into a concise and easy-to-understand format. 
- Ensure that the summary includes relevant details and examples that support the main ideas, while avoiding any unnecessary information or repetition. 
- Rely strictly on the provided text, without including external information. 
- The length of the summary must be appropriate for the length and complexity of the original text. The length must allow to capture the main points and key details, without being overly long.  
- Ensure that the summary is well-organized and easy to read, with clear headings and subheadings to guide the reader through each section. Format each section in paragraph form. you can use markdown
- Output result in simplified Chinese

Input Article Content:
`,
}
