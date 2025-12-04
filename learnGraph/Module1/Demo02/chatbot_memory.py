import getpass
import os

def _set_env(var: str):
    if not os.environ.get(var):
        os.environ[var] = getpass.getpass(f"{var}: ")

_set_env("OPENAI_API_KEY")

from langchain_openai import ChatOpenAI
from langgraph.graph import StateGraph, MessagesState, START, END
from langgraph.checkpoint.memory import MemorySaver

def chatbot(state: MessagesState):
    return {"messages": [ChatOpenAI(model="doubao-seed-1-6-vision-250815", base_url="https://ark.cn-beijing.volces.com/api/v3", api_key="97e40081-f579-4ae6-91c9-56c29f4abb14").invoke(state["messages"])]}

# 构建图
graph = StateGraph(MessagesState)
graph.add_node("chatbot", chatbot)
graph.add_edge(START, "chatbot")
graph.add_edge("chatbot", END)

# 使用 memorySaver 保存对话历史
memory = MemorySaver()

# 编译并运行图 这个 app 就是我们构建的图，有时候也会命名为 graph
app = graph.compile(checkpointer=memory)

# 多轮对话
config = {"configurable":{"thread_id":"user_001"}}

# 第一轮
response1 = app.invoke(
    {"messages":[("user", "我的名字是小明")]},
    config=config
)
print("Round 1:", response1["messages"][-1].content)

# 第二轮
response2 = app.invoke(
    {"messages":[("user", "我的名字是什么？")]},
    config=config
)
print("Round 2:", response2["messages"][-1].content)

