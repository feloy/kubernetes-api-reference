# kubernetes-api-reference

<!-- ![Bugs](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=bugs) -->
<!-- ![Code Smalls](https://sonarcloud.io/api/project_badges/measure? project=feloy_kubernetes-api-reference&metric=code_smells) -->
<!-- ![Duplicated lines](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=duplicated_lines_density) -->
<!-- ![Lines of code](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=ncloc) -->
<!-- ![technical debt](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=sqale_index)-->
<!-- ![vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=vulnerabilities) -->

![Maintainability](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=sqale_rating)
![reliability](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=reliability_rating)
![security](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=security_rating)
![quality gate](https://sonarcloud.io/api/project_badges/measure?project=feloy_kubernetes-api-reference&metric=alert_status)


Tool to create documentation of the Kubernetes API.

## Usage

```go
$ go build cmd/main.go
$ ./main hugo \
  --file api/v1.18/swagger.json \
  --config-dir config/v1.18/ \
  --output-dir website/content/en/docs/
$ cd website
$ hugo
$ cd public
$ python -m http.server
```

## OpenAPI specification

The source of truth for the Kubernetes API is an OpenAPI specification. A standard OpenAPI specification describes:

- a list of *Definitions*,
- a list of *Paths*, each describing a list of *Operations*.

## Kubernetes extensions

https://github.com/kubernetes/kubernetes/tree/master/api/openapi-spec

Kubernetes API extends OpenAPI using these extensions:

- `x-kubernetes-group-version-kind`:
  - Definitions associated with a Kubernetes *Resource* use this extension to declare the GVK to which the resource belongs.
  - Operations use this extension to declare on which Kubernetes resource they operate.
- `x-kubernetes-action`: OpenAPI Operations (get, post, etc) are mapped to Kubernetes API *actions* (get, list, watch, etc) with this extension.
- `x-kubernetes-patch-strategy`: a comma-separated list of strategic merge patch strategies supported by a field of a Kubernetes resource.
- `x-kubernetes-patch-merge-key`: when a field supports the `merge` strategy, this extension indicates the key used to identify the fields to merge.

## TODO

- additionalProperties: map[string]...
- merge extensions on properties
