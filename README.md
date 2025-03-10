# cosi-sample-app

## Installing the sample app

To install the `cosi-sample-app`, run the following commands. This will ensure you're always pulling the latest stable tag from the GitHub repository.

```sh
LATEST="$(curl -s 'https://api.github.com/repos/anza-labs/cosi-sample-app/tags' | jq -r '.[0].name')"
kubectl apply -k "https://github.com/anza-labs/cosi-sample-app//?ref=${LATEST}"
```
