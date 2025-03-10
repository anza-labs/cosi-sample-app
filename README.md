# cosi-sample-app

## Installing the sample app

To install the `cosi-sample-app`, run the following commands. This will ensure you're always pulling the latest stable tag from the GitHub repository.

```sh
LATEST="$(curl -s 'https://api.github.com/repos/anza-labs/cosi-sample-app/tags' | jq -r '.[0].name')"
kubectl apply -k "https://github.com/anza-labs/cosi-sample-app//?ref=${LATEST}"
```

## Patching the Claims

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- https://github.com/anza-labs/cosi-sample-app//?ref=v0.1.3
patches:
- patch: |-
    - op: replace
      path: /spec/bucketClassName
      value: your-bucket-class
  target:
    kind: BucketClaim
    name: app
- patch: |-
    - op: replace
      path: /spec/bucketAccessClassName
      value: your-bucket-access-class
  target:
    kind: BucketAccess
    name: app
```
