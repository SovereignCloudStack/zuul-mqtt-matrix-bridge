# Docker images

The images for the "zuul-mqtt-matrix-bridge" project are stored in the [SCS registry](https://registry.scs.community/zuul). 

TODO: Automating the build and push of images for new releases has not been implemented yet.

If you want to build and publish an image for a new release, follow these steps:

```bash
# Build and tag image locally
docker build -t zuul-mqtt-matrix-bridge:<version> . -f docker/Dockerfile
docker tag  zuul-mqtt-matrix-bridge:<version> registry.scs.community/zuul/zuul-mqtt-matrix-bridge:<version>
# Push it to the SCS registry
docker push registry.scs.community/zuul/zuul-mqtt-matrix-bridge:<version>
```
