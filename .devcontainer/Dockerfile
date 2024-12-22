FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Install k9s 
WORKDIR /tmp
RUN wget https://github.com/derailed/k9s/releases/download/v0.32.5/k9s_linux_amd64.deb && \
  dpkg -i k9s_linux_amd64.deb && \
  rm k9s_linux_amd64.deb

WORKDIR /workspaces/gopherproxy