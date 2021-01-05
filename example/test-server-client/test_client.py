from requests import get, post

ans = post(
    "http://localhost:8000/json/",
    json={"jsonrpc": "2.0", "method": "sum", "params": [1, 2]},
)
print(ans.text)


ans = post(
    "http://localhost:8000/json/",
    json={"jsonrpc": "2.0", "method": "prod", "params": [1, 2]},
)
print(ans.text)


ans = post(
    "http://localhost:8000/json/",
    json={"jsonrpc": "2.0", "method": "hello", "params": ["navid"]},
)
print(ans.text)