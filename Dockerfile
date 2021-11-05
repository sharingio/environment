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
  && yes '\n' | add-apt-repository ppa:git-core/ppa \
  && apt update \
  && DEBIAN_FRONTEND=noninteractive apt install --no-install-recommends -y \
    emacs-nox \
    tmate \
    bash-completion \
    less \
    xz-utils \
    sudo \
    curl \
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
