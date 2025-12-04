import getpass
import os


def _set_env(var: str):
    if not os.environ.get(var):
        os.environ[var] = getpass.getpass(f"{var}: ")

_set_env("OPENAI_API_KEY")


from langchain_openai import ChatOpenAI
from langgraph.graph import StateGraph, MessagesState, START, END
from IPython.display import Image, display

def chatbot(state: MessagesState):
    return {"messages": [ChatOpenAI(model="doubao-seed-1-6-vision-250815", base_url="https://ark.cn-beijing.volces.com/api/v3").invoke(state["messages"])]}

# 构建图
graph = StateGraph(MessagesState)
graph.add_node("chatbot", chatbot)
graph.add_edge(START, "chatbot")
graph.add_edge("chatbot", END)

# 编译并运行图 这个 app 就是我们构建的图，有时候也会命名为 graph
app = graph.compile()
res = app.invoke({"messages": [("user", "你好，请用一句话介绍 LangGraph")]})
print(res["messages"][-1].content)  # 输出模型的回复

# 🎨 可视化图结构

display(Image(app.get_graph().draw_mermaid_png()))