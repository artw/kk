# kk: Kubernetes Klient

A simple microservice that tries to feed your input payload to the kubernetes server it runs on in a way similar to `kubectl` with `-f` flag.
|method|kubectl|
|------|-------|
|PUT|apply|
|POST|create|
|DELETE|delete|

supports `application/json` or `application/yaml`, no defaults
## Prerequisites

- Must be run inside a Kubernetes cluster.
- The associated pod should have a service account with the appropriate roles and permissions to perform operations on the intended resources.

## Usage

### Create a Kubernetes resource

Use the `POST` method:

```bash
curl -X POST -H "Content-Type: application/yaml" --data-binary @your-resource-file.yaml http://server-address:8080/
```

Or with JSON:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"apiVersion":"v1", ... }' http://server-address:8080/
```

### Apply (create/update) a Kubernetes resource

Use the `PUT` method:

```bash
curl -X PUT -H "Content-Type: application/yaml" --data-binary @your-resource-file.yaml http://server-address:8080/
```

Or with JSON:

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"apiVersion":"v1", ... }' http://server-address:8080/
```

### Delete a Kubernetes resource

The `DELETE` method requires a minimal resource definition (kind, apiVersion, and metadata with name, optionally namespace):

```bash
curl -X DELETE -H "Content-Type: application/yaml" --data-binary @minimal-resource-definition.yaml http://server-address:8080/
```

Or with JSON:

```bash
curl -X DELETE -H "Content-Type: application/json" -d '{"apiVersion":"v1", ... }' http://server-address:8080/
```

Ensure the service account tied to this application's pod is granted the necessary roles to act on the target Kubernetes resources.

