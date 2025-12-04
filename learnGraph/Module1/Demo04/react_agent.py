import getpass
import os


def _set_env(var: str):
    if not os.environ.get(var):
        os.environ[var] = getpass.getpass(f"{var}: ")

_set_env("OPENAI_API_KEY")
_set_env("TAVILY_API_KEY")

from langchain_openai import ChatOpenAI
from langchain_community.tools.tavily_search import TavilySearchResults
from langgraph.graph import StateGraph, MessagesState, START, END
from langgraph.prebuilt import ToolNode

# 初始化工具
search = TavilySearchResults(max_results=2)
tools = [search]

# 定义agent (添加 ReAct 系统提示词)
llm = ChatOpenAI(model="doubao-seed-1-6-vision-250815", base_url="https://ark.cn-beijing.volces.com/api/v3")
llm_with_tools = llm.bind_tools(tools)

def agent(state: MessagesState):
    # 添加 ReAct 提示词
    system_message = """你是一个 ReAct (Reasoning + Acting) Agent。
    处理用户问题时，请遵循以下步骤：
    1. Thought（思考）：分析问题需要什么信息
    2. Action（行动）：决定调用哪个工具
    3. Observation（观察）：分析工具返回的结果
    4. Answer（回答）：基于观察给出最终答案
    你可以多次搜索并且将信息进行综合，以确保答案的准确性。
    始终展示你的推理过程。"""
    messages = [{"role":"system", "content": system_message}] + state["messages"]
    return {"messages": [llm_with_tools.invoke(messages)]}

def should_continue(state: MessagesState):
    last_message = state["messages"][-1]
    if hasattr(last_message, "tool_calls") and last_message.tool_calls:
        return "tools"
    return END

# 构建图
graph = StateGraph(MessagesState)
graph.add_node("agent", agent)
graph.add_node("tools", ToolNode(tools=tools))

graph.add_edge(START, "agent")
graph.add_conditional_edges("agent", should_continue, ["tools", END])
graph.add_edge("tools", "agent")

app = graph.compile()

response = app.invoke({
    "messages":[("user", "贝壳找房有什么业务")]
})

print("Agent Response:", response["messages"][-1].content)




