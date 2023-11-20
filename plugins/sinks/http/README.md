# Http

Compass is a search and discovery engine built for querying application deployments, datasets and meta resources. It can also optionally track data flow relationships between these resources and allow the user to view a representation of the data flow graph.

## Usage

```yaml
sinks:
  name: http
  config:
    method: POST
    success_code: 200
    url: https://compass.requestcatcher.com/v1beta2/asset/{{ .Type }}/{{ .Urn }}
    headers:
      Header-1: value11,value12
```

## Config Defination

| Key | Value | Example | Description |  |
| :-- | :---- | :------ | :---------- | :-- |
|`url` | `string` | `https://compass.requestcatcher.com/v1beta2/asset/{{ .Type }}/{{ .Urn }}` | URL to the http server, contains all the info needed to make the request, like port and route, support go [text/template](https://pkg.go.dev/text/template) (see the properties in [v1beta2.Asset](https://github.com/goto/meteor/blob/main/models/gotocompany/assets/v1beta2/asset.pb.go#L25-L68))  | *required*|
| `method` | `string` | `POST` | the method string of by which the request is to be made, e.g. POST/PATCH/GET | *required* |
| `success_code` | `integer` | `200` |  to identify the expected success code the http server returns, defult is `200` | *optional* |
| `headers` | `map` | `"Content-Type": "application/json"` | to add any header/headers that may be required for making the request | *optional* |

## Contributing

Refer to the [contribution guidelines](../../../docs/docs/contribute/guide.md#adding-a-new-sink) for information on contributing to this module.
