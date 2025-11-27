package agent

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"travelAgent/cmd/llm"
	"travelAgent/cmd/tool"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var (
	LlmClient           = llm.OpenAICompatibleClient{}
	AGENT_SYSTEM_PROMPT = `
	你是一个智能旅行助手。你的任务是分析用户的请求，并使用可用工具一步步地解决问题。
	
	# 可用工具:
	- get_weather(city string): 查询指定城市的实时天气 城市使用拼音。
	- get_attraction(city string, weather string): 根据城市和天气搜索推荐的旅游景点。
	
	# 行动格式:
	你的回答必须严格遵循以下格式。首先是你的思考过程，然后是你要执行的具体行动，每次回复只输出一对Thought-Action：
	Thought: [这里是你的思考过程和下一步计划]
	Action: [这里是你要调用的工具，格式为 function_name(arg_name="arg_value")]
	
	# 任务完成:
	当你收集到足够的信息，能够回答用户的最终问题时，你必须在 Action: 字段后使用 finish(answer="...") 来输出最终答案。
	
	请开始吧！
	`
)

func init() {
	_ = godotenv.Load()

	config := openai.DefaultConfig(os.Getenv("API_KEY"))
	config.BaseURL = os.Getenv("BASE_URL")

	LlmClient = llm.OpenAICompatibleClient{
		Model:  os.Getenv("MODEL"),
		Client: openai.NewClientWithConfig(config),
	}
}

func Run() {
	user_prompt := "你好，请帮我查询一下今天成都的天气，然后根据天气推荐一个合适的旅游景点。"
	prompt_history := []string{
		fmt.Sprintf("用户请求: %s", user_prompt),
	}
	fmt.Printf("用户输入: %s\n%s\n", user_prompt, strings.Repeat("=", 40))

	// 导入工具
	availableTools := getAvailableTools()

	// 主循环，最多运行10次
	for i := 0; i < 10; i++ {
		fmt.Printf("\n--- 循环 %d ---\n\n", i+1)

		// 3.1 构建完整的 prompt
		full_prompt := strings.Join(prompt_history, "\n")

		// 3.2 调用 LLM 进行思考
		llm_output := LlmClient.Generate(full_prompt, AGENT_SYSTEM_PROMPT)

		// 截断多余的 Thought-Action 对（可选）
		truncateRegex := regexp.MustCompile(`(?s)(Thought:.*?Action:.*?)(?:\n\s*(?:Thought:|Action:|Observation:)|\z)`)
		if matches := truncateRegex.FindStringSubmatch(llm_output); len(matches) > 1 {
			truncated := strings.TrimSpace(matches[1])
			if truncated != strings.TrimSpace(llm_output) {
				llm_output = truncated
				fmt.Println("已截断多余的 Thought-Action 对")
			}
		}

		fmt.Printf("模型输出:\n%s\n\n", llm_output)
		prompt_history = append(prompt_history, llm_output)

		// 3.3 解析并执行行动
		actionRegex := regexp.MustCompile(`(?s)Action:\s*(.*)`)
		actionMatch := actionRegex.FindStringSubmatch(llm_output)
		if len(actionMatch) < 2 {
			fmt.Println("解析错误: 模型输出中未找到 Action")
			break
		}

		actionStr := strings.TrimSpace(actionMatch[1])

		// 检查是否是 finish
		if strings.HasPrefix(actionStr, "finish") {
			finishRegex := regexp.MustCompile(`finish\(answer="(.*)"\)`)
			finishMatch := finishRegex.FindStringSubmatch(actionStr)
			if len(finishMatch) >= 2 {
				finalAnswer := finishMatch[1]
				fmt.Printf("任务完成，最终答案: %s\n", finalAnswer)
				break
			} else {
				fmt.Println("解析错误: finish 格式不正确")
				break
			}
		}

		// 解析工具名称
		toolNameRegex := regexp.MustCompile(`(\w+)\(`)
		toolNameMatch := toolNameRegex.FindStringSubmatch(actionStr)
		if len(toolNameMatch) < 2 {
			fmt.Printf("解析错误: 无法解析工具名称，actionStr=%s\n", actionStr)
			break
		}
		toolName := toolNameMatch[1]

		// 解析参数
		argsRegex := regexp.MustCompile(`\((.*)\)`)
		argsMatch := argsRegex.FindStringSubmatch(actionStr)
		if len(argsMatch) < 2 {
			fmt.Printf("解析错误: 无法解析参数，actionStr=%s\n", actionStr)
			break
		}
		argsStr := argsMatch[1]

		// 提取键值对参数
		kvRegex := regexp.MustCompile(`(\w+)="([^"]*)"`)
		kvMatches := kvRegex.FindAllStringSubmatch(argsStr, -1)
		params := make(map[string]string)
		for _, match := range kvMatches {
			if len(match) == 3 {
				params[match[1]] = match[2]
			}
		}

		// 调用工具
		var observation string
		if toolFunc, exists := availableTools[toolName]; exists {
			observation = toolFunc(params)
		} else {
			observation = fmt.Sprintf("错误: 未定义的工具 '%s'", toolName)
		}

		// 3.4 记录观察结果
		observationStr := fmt.Sprintf("Observation: %s", observation)
		fmt.Printf("%s\n%s\n", observationStr, strings.Repeat("=", 40))
		prompt_history = append(prompt_history, observationStr)
	}

	fmt.Println("\nAgent 执行完毕")
}

// getAvailableTools 返回可用的工具函数映射
func getAvailableTools() map[string]func(map[string]string) string {
	return map[string]func(map[string]string) string{
		"get_weather": func(params map[string]string) string {
			city, ok := params["city"]
			if !ok {
				return "错误: 缺少参数 city"
			}
			return tool.GetWeather(city)
		},
		"get_attraction": func(params map[string]string) string {
			city, ok := params["city"]
			if !ok {
				return "错误: 缺少参数 city"
			}
			weather, ok := params["weather"]
			if !ok {
				return "错误: 缺少参数 weather"
			}
			return tool.GetAttraction(city, weather)
		},
	}
}
