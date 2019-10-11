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