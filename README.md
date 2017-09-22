This is a JsonRpc to REST and HTTP Apis translator which is even capable to act as a gateway.

# How it works?

## Simple Example
To run this JsonRpc2Rest you need a config file contains endpoints and apis. As an Example:
```
{
    "bind": "0.0.0.0:8080"
    "upstreams": {
        "sum": {
            "url": "http://192.168.1.1/sum/",
            "params": ["a","b"],
            "method": "GET"
        }
    }
}
```
As you see in the example config.json file above, `bind` key shows the bind address of the gateway. `upstreams` tells the methods of jsonrpc mapped to an endpoint. Every method has to have `url` to the endpoint and list of its parameters in `params`. Finally the http method comes in `method`. these three keys in config are mandatory.
For Example if the json rpc is a function like sum(1,2) it will be translated to http://192.168.1.1/sum/?a=1&b=2.

## Payload Params
```
...
"my_awesome_method": {
    "url": "http://192.168.1.1/my_awesome_endpoint/",
    "params": ["a","b","c"],
    "method": "POST",
    "payload": ["a","b"]
}
...
```

It is even possible to send some parameters (or all of them) to payload as json. To do this it is only needed to add `payload` key and add payload parameters to it. The above example translates a `my_awesome_method(1,2,3)` to http://192.168.1.1/my_awsome_endpoint/?c=3 and payload of `{"a":1,"b":2}.

## Url Paramteres

To have paramaters in url address in rest endpoint the below example is able to handle:
```
...
"hello": {
    "url": "http://192.168.1.1/hello/:name/", <---- name in url as parameter
    "params": ["name"],
    "method": "GET"
}
...
```
This config tells the translator to replace `:name` with paramter of `name`. If we have a rpc like `hello("Yadollah")` then it will be translated to http://192.168.1.1/hello/Yadollah/.

## Load Balancing

JsonRpc2Rest is able to load balance. To perform this you need to define `hosts` in config json:
```
{
    "hosts" : {
        "myhost" : {
            "127.0.0.1:5000" : 2,
            "127.0.0.2:5000" : 3
        },
        "myhost2" : {
            "127.0.0.1:5000" : 2,
            "127.0.0.2:5000" : 1
        }
    },
    "upstreams": {
        "sum": {
            "url": "http://myhost/sum/",
            "params": ["a","b"],
            "method": "GET"
        },
       "prod": {
            "url": "http://myhost2/prod/",
            "mehod":"POST",
            "params": ["a","b"],
            "payload_params": ["a","b"]
        }
    }
}
```
Hosts contains alias names and their addresses. Every address has to have an unsigned integer as its weight in the value part. As an example in the config above, `myhost` has to address with weights of 2 and 3. It means that the gateway randomly chooses one of the with 40% chance for first endpoint and 60% for second one.

## Running
After you compiled the project, you can run it like `./jsonrpc2rest config.json`.