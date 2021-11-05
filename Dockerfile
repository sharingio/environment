# Dockerfile
#   definition to the environment container image

ARG BASE_IMAGE=ubuntu:21.10
FROM $BASE_IMAGE

ARG ENV_DOOM_REPO=https://github.com/hlissner/doom-emacs \
  ENV_DOOM_CONFIG_REPO=https://github.com/humacs/.doom.d \
  ENV_DOOM_REF=0adaf03088ee7f39b3b2bca76e24fb1828721722 \
  ENV_DOOM_CONFIG_REF=50b4e68504debfc3c59fe811770328ec769267d6
ENV TERM=screen-256color \
  DOOMDIR=/home/ii/.doom.d \
  PATH=$PATH:/var/local/doom-emacs/bin

RUN DEBIAN_FRONTEND=noninteractive \
  apt update \
  && apt upgrade -y \
  && DEBIAN_FRONTEND=noninteractive apt install --no-install-recommends -y \
    software-properties-common \
    gpg-agent \
    curl \
  && yes '\n' | add-apt-repository ppa:git-core/ppa \
  && echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" \
    | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list \
  && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg  add - \
  && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
  && echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
  && apt update \
  && DEBIAN_FRONTEND=noninteractive apt install --no-install-recommends -y \
    emacs-nox \
    tmate \
    bash-completion \
    less \
    xz-utils \
    sudo \
    ca-certificates \
    libcap2-bin \
    git \
    kitty \
    openssh-client \
    postgresql-client \
    jq \
    inotify-tools \
    xtermcontrol \
    nodejs \
    gnupg2 \
    tzdata \
    wget \
    python3-dev \
    xz-utils \
    apache2-utils \
    sqlite3 \
    silversearcher-ag \
    build-essential \
    vim \
    rsync \
    unzip \
    direnv \
    iputils-ping \
    file \
    psmisc \
    docker-ce-cli \
    tree \
    iproute2 \
    net-tools \
    tcpdump \
    htop \
    iftop \
    tmux \
    language-pack-en \
    openjdk-16-jdk \
    rlwrap \
    fonts-powerline \
    dnsutils \
    python3-pip \
    npm \
    ripgrep \
    python-is-python3 \
    shellcheck \
    pipenv \
    fd-find \
    gettext-base \
    libcap2-bin \
    locate \
    flatpak-xdg-utils \
    google-cloud-sdk \
    awscli \
    expect \
    graphviz \
    runit \
    ssh-import-id \
    bsdmainutils \
    netcat \
  && ln -s /usr/bin/fdfind /usr/local/bin/fd \
  && rm -rf /var/lib/apt/lists/*

COPY etc/ /etc/

RUN mkdir -p /etc/sudoers.d && \
  echo "%sudo    ALL=(ALL:ALL) NOPASSWD: ALL" > /etc/sudoers.d/sudo && \
  useradd -m -G users,sudo -u 1000 -s /bin/bash ii && \
  chmod 0775 /usr/local/lib && chgrp users /usr/local/lib && \
  chmod 0770 -R /etc/service/ && \
  chgrp -R users /etc/service/ && \
  mkdir /usr/local/lib/node_modules && \
  chown -R ii:ii /usr/local/lib/node_modules /var/local

USER ii
WORKDIR /home/ii
RUN git clone $ENV_DOOM_REPO /var/local/doom-emacs && \
  cd /var/local/doom-emacs && \
  git checkout $ENV_DOOM_REF
RUN git clone $ENV_DOOM_CONFIG_REPO && \
  cd .doom.d && \
  git checkout $ENV_DOOM_CONFIG_REF
RUN /var/local/doom-emacs/bin/org-tangle .doom.d/ii.org \
  && yes | /var/local/doom-emacs/bin/doom install --no-env \
  && yes | /var/local/doom-emacs/bin/doom sync -e
