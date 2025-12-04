import getpass
import os


def _set_env(var: str):
    if not os.environ.get(var):
        os.environ[var] = getpass.getpass(f"{var}: ")

_set_env("OPENAI_API_KEY")

from langchain_openai import ChatOpenAI
from langgraph.graph import StateGraph, MessagesState, START, END
from langchain_core.tools import tool
from langgraph.prebuilt import ToolNode

# 定义工具

@tool
def calculator_tool(expression: str) -> str:
    """执行数学计算。输入格式：'数字1 运算符 数字2'（如 '25 * 4'）"""
    try:
        parts = expression.strip().split()
        if len(parts) != 3:
            return "输入格式错误，请使用 '数字1 运算符 数字2' 格式。"
        num1, operator, num2 = parts
        num1, num2 = float(num1), float(num2)
        if operator == '+':
            return str(num1 + num2)
        elif operator == '-':
            return str(num1 - num2)
        elif operator == '*':
            return str(num1 * num2)
        elif operator == '/':
            return str(num1 / num2) if num2 != 0 else "除数不能为零。"
        else:
            return f"不支持的运算符:{operator}"
    except:
        return "计算错误"

@tool
def get_weather(city: str) -> str:
    """模拟获取天气信息的工具函数"""
    weather_db = {
        "北京": "晴，25°C",
        "上海": "多云，22°C",
        "广州": "小雨，28°C",
    }
    return weather_db.get(city, "无法获取该城市的天气信息。")

tools = [calculator_tool, get_weather]

# 定义节点
llm = ChatOpenAI(model="doubao-seed-1-6-vision-250815", base_url="https://ark.cn-beijing.volces.com/api/v3")
llm_with_tools = llm.bind_tools(tools)

def chatbot(state: MessagesState):
    return {"messages": [llm_with_tools.invoke(state["messages"])]}

# 条件边: 判断是否需要调用工具
def should_continue(state: MessagesState):
    last_message = state["messages"][-1]
    if hasattr(last_message, "tool_calls") and last_message.tool_calls:
        return "tools"
    return END
# 构建图
graph = StateGraph(MessagesState)
graph.add_node("chatbot", chatbot)
graph.add_node("tools", ToolNode(tools))

graph.add_edge(START, "chatbot")
graph.add_conditional_edges("chatbot", should_continue, ["tools", END])
graph.add_edge("tools", "chatbot")

app = graph.compile()

# 测试
response = app.invoke({"messages":[("user", "请帮我计算25乘以4的结果，并告诉我北京的天气如何？")]})
print(response["messages"][-1].content)  # 输出模型的回复
