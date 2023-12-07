# Docker images

The images for the "zuul-mqtt-matrix-bridge" project are stored in the [SCS registry](https://registry.scs.community/zuul). 

TODO: Automating the build and push of images for new releases has not been implemented yet.

If you want to build, tag and publish multiplatform (darwin/arm64,linux/amd64) images for a new release, follow these steps:

```bash
# Build, tag and push images to registry.scs.community
docker buildx create --use --name buildx_instance
docker buildx build \
  --file docker/Dockerfile \
  --platform=darwin/arm64,linux/amd64 \
  --tag registry.scs.community/zuul/zuul-mqtt-matrix-bridge:0.3.0 \
  --tag registry.scs.community/zuul/zuul-mqtt-matrix-bridge:latest \
  --provenance=false \
  --push .
```
