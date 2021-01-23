# Status

This project has been in production use since 2017, running behind a WAF and load balancer.

# Objective

To serve as a low-fuss, minimal configuration reverse proxy for APIs. No complicated configurations, user portals, authentication, load balancing or anything. The point is that this can run behind something like Cloudflare that can take care of many of these concerns. Indeed, many API gateway solutions run behind load balancers or CDNs.

### Minimal configuration:

```
{
  "port": "8080",
  "versions": {
    "v1": {
        "/api": "localhost:8089"
    }
  }
}
```

### Full configuration:

```
{
  "port": "8081",
  "versions": {
    "v1": {
      "/google": "google.com:443",
      "/history": "127.0.0.1:8089"
    },
    "v2": {
      "/microsoft": "microsoft.com:443"
    }
  },
  "not_found_error": {
    "code": 404,
    "domain": "Route not found",
    "message": "The requested route was not found"
  },
  "fallback_rule": "www.openenergi.com",
  "scheme": "https"
}
```

# Getting started

You'll need to build the binary with `go build`. It's a static, standalone binary with no dependencies. The code can also be used as a Go library, and it has no dependencies outside the standard library.

Next, create a configuration file. The options are:

* `versions` : This is an object that gets prefixed to URLs and ignored by backends e.g. a version string of `v1` would lead the gateway to route `/v1/users/asdf` to `/users/asdf`. The object contains a map of route prefix to backend. The backend will see the entire URL except for the version string, if that was specified. Rules are not applied in any particular order.
* `port` (Optional): What port to run the gateway on. If this is not specified the `HTTP_PLATFORM_PORT` environment variable will be used.
* `not_found_error` (Optional): Custom error object to return in case the request URL does not match any rules.
* `scheme` (Optional): The URL scheme to use (http or https).

When you have created the configuration and built the gateway, just run it as

```
gateway path-to-config.json
```

Alternatively, you can set the path of the config file in an environment variable `GATEWAY_CONFIG_FILE` and run `gateway` without any additional arguments.

# Library use 

You can also use the code as a library, importing `github.com/michaelbironneau/gateway/lib`. For this purpose you will need to create a new gateway using `lib.New`.

The benefit of using the gateway as a library is that you can intercept all requests both before they are handled (to filter them), and afterwards, before they are returned to the client and log (possibly modify) the response. At this time, you *cannot* modify responses from the gateway itself (that is, when a backend is not found so your configurable "not found response" is returned, or when a Filter you have implemented denies access) - but you can modify any response that one of your backends creates. If you read the response, remember to re-set it, eg. 

```go
bodyBytes, _ := ioutil.ReadAll(resp.Body)
resp.Body.Close()  //  must close
resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
```

Example that intercepts and prints out response body):

```go
func main(){            
    //...read config from config file into c (see main.go for example)
    c.Interceptor = func(r *http.Request, resp *http.Response){
        log.Println("Response:/n--------/n")
        bodyBytes, _ := ioutil.ReadAll(resp.Body)
        resp.Body.Close()  //  must close
        resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // reset it now we've read it all
        log.Println(string(bodyBytes) + "\n")
    }
    http.HandleFunc("/", lib.New(c))
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

```