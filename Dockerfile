FROM alpine:latest

# Install curl and required dependencies
RUN apk add --no-cache curl

# Set the kubectl version to install
ENV KUBECTL_VERSION=v1.27.8

# Download and install kubectl
RUN curl -LO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" && \
  chmod +x kubectl && \
  mv kubectl /usr/local/bin/

# Verify the installation
RUN kubectl version --client

# Install Node.js and npm
RUN apk add --no-cache nodejs npm

# Verify installation
RUN node --version && npm --version
