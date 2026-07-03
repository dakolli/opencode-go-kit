FROM ubuntu:24.04

# Avoid prompt interruptions during package installation
ENV DEBIAN_FRONTEND=noninteractive

# 1. Install system utilities and development dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    wget \
    unzip \
    ripgrep \
    build-essential \
    sqlite3 \
    python3 \
    python3-pip \
    python3-venv \
    golang-go \
    && rm -rf /var/lib/apt/lists/*

# 2. Install Node.js LTS (v20) system-wide
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    rm -rf /var/lib/apt/lists/*

# 3. Securely install OpenCode globally via NPM (Standard method)
RUN npm install -g opencode-ai

# 4. Set up non-root user parameters for host alignment
ARG USER_ID=1000
ARG GROUP_ID=1000
ARG USER_NAME=opencode
ARG GROUP_NAME=opencode

# FIX: Forcefully evict the default 'ubuntu' user/group if they exist to free up UID/GID 1000
RUN if getent passwd ubuntu; then userdel -f ubuntu; fi && \
    if getent group ubuntu; then groupdel ubuntu; fi

# 5. Create the customizable user space
RUN groupadd --gid ${GROUP_ID} ${GROUP_NAME} && \
    useradd --uid ${USER_ID} --gid ${GROUP_ID} --shell /bin/bash --create-home ${USER_NAME}

# 6. Initialize the workspace and pass ownership to our user
RUN mkdir -p /workspace && \
    chown -R ${USER_NAME}:${GROUP_NAME} /workspace && \
    chmod 755 /workspace

WORKDIR /workspace
USER ${USER_NAME}

# ARG CONTAINER_PORT=4096
# ENV CONTAINER_PORT=${CONTAINER_PORT}
# EXPOSE ${CONTAINER_PORT}

CMD ["sh", "-c", "exec opencode serve --hostname 0.0.0.0 --port 4096"]