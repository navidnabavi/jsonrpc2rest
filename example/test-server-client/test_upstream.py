from flask import Flask, request, jsonify

app = Flask(__name__)


@app.route("/sum/")
def sum():
    print("1")
    a = int(request.args.get("a", ""))
    b = int(request.args.get("b", ""))
    return jsonify({"result": {"sum": a + b, "diff": a - b}})


@app.route("/prod/")
def prod():
    obj = request.json
    a = obj["a"]
    b = obj["b"]
    return jsonify({"result": a + b})


@app.route("/hello/<name>/")
def hello(name):
    return jsonify({"result": "hello {}".format(name)})
