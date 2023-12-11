# Kubewarden Policies for dot-sql Crossplane Compositions

## Test

```bash
make test
```

## Publish

```bash
docker run --rm -e GOFLAGS="-buildvcs=false" \
    -v /Users/viktorfarcic/code/kubewarden-deployment:/src \
    -w /src tinygo/tinygo:0.30.0 \
    tinygo build -o policy.wasm -target=wasi -no-debug .

kwctl annotate policy.wasm \
    --metadata-path metadata.yml \
    --output-path annotated-policy.wasm

kwctl push annotated-policy.wasm \
    c8n.io/vfarcic/kubewarden-crossplane-sql:v0.0.2
```