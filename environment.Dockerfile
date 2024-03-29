# Dockerfile
#   definition to the environment container image

ARG BASE_IMAGE=ubuntu:22.04
FROM $BASE_IMAGE

ARG ENV_DOOM_REPO=https://github.com/doomemacs/doomemacs
ARG ENV_DOOM_CONFIG_REPO=https://github.com/humacs/.doom.d
ARG ENV_DOOM_REF=b06fd63dcb686045d0c105f93e07f80cb8de6800
ARG ENV_DOOM_CONFIG_REF=13d47156d1eeda72fff473dacbef3cfd499bae3c
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
  && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
  && apt-key fingerprint 0EBFCD88 \
  && add-apt-repository "deb https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
  && curl -fsSL https://deb.nodesource.com/setup_16.x | bash -x - \
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
    docker-ce \
    tree \
    iproute2 \
    net-tools \
    tcpdump \
    htop \
    iftop \
    tmux \
    language-pack-en \
    openjdk-17-jdk \
    rlwrap \
    fonts-powerline \
    dnsutils \
    python3-pip \
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
    ssh-import-id \
    bsdmainutils \
    netcat \
    asciinema \
  && ln -s /usr/bin/fdfind /usr/local/bin/fd \
  && rm -rf /var/lib/apt/lists/*

COPY etc/ /etc/

RUN mkdir -p /etc/sudoers.d && \
  echo "%sudo    ALL=(ALL:ALL) NOPASSWD: ALL" > /etc/sudoers.d/sudo && \
  useradd -m -G users,sudo -u 1000 -s /bin/bash ii && \
  chmod 0775 /usr/local/lib && chgrp users /usr/local/lib && \
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

USER root
COPY ./usr/lib /usr/lib
ENV DOCKER_VERSION=20.10.10 \
  KIND_VERSION=0.11.1 \
  KUBECTL_VERSION=1.24.2 \
  GO_VERSION=1.19.1 \
  TILT_VERSION=0.22.15 \
  TMATE_VERSION=2.4.0 \
  HELM_VERSION=3.7.1 \
  GH_VERSION=2.13.0 \
  LEIN_VERSION=stable \
  CLOJURE_VERSION=1.10.1.697 \
  CLUSTERCTL_VERSION=1.2.2 \
  TALOSCTL_VERSION=1.1.0 \
  TERRAFORM_VERSION=1.2.4 \
  DIVE_VERSION=0.10.0 \
  CRICTL_VERSION=1.22.0 \
  KUBECTX_VERSION=0.9.4 \
  FZF_VERSION=0.26.0 \
  NERDCTL_VERSION=0.18.0 \
  METALCLI_VERSION=0.6.0 \
  KO_VERSION=0.11.2 \
  KN_VERSION=1.6.0 \
  UPTERM_VERSION=0.7.6 \
  GOPLS_VERSION=0.8.3 \
# GOLANG, path vars
  GOROOT=/usr/local/go \
  PATH="$PATH:/usr/local/go/bin:/usr/libexec/flatpak-xdg-utils:/home/ii/go/bin" \
  CONTAINERD_NAMESPACE=k8s.io
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://dl.google.com/go/go${GO_VERSION}.linux-${ARCH_TYPE_2}.tar.gz \
    | tar --directory /usr/local --extract --ungzip
# kind binary
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -Lo /usr/local/bin/kind \
    https://github.com/kubernetes-sigs/kind/releases/download/v${KIND_VERSION}/kind-linux-${ARCH_TYPE_2} \
    && chmod +x /usr/local/bin/kind
# kubectl binary
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/${ARCH_TYPE_2}/kubectl -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl
# tilt binary
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -fsSL \
    https://github.com/tilt-dev/tilt/releases/download/v${TILT_VERSION}/tilt.${TILT_VERSION}.linux.${ARCH_TYPE_1}.tar.gz \
    | tar --directory /usr/local/bin --extract --ungzip tilt
# gh cli
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -sSL https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_${ARCH_TYPE_2}.tar.gz \
    | tar --directory /usr/local --extract --ungzip \
     --strip-components 1 gh_${GH_VERSION}_linux_${ARCH_TYPE_2}/bin/gh \
    && chmod +x /usr/local/bin/gh
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L \
    https://github.com/tmate-io/tmate/releases/download/${TMATE_VERSION}/tmate-${TMATE_VERSION}-static-linux-${ARCH_TYPE_3}.tar.xz \
    | tar --directory /usr/local/bin --extract --xz \
  --strip-components 1 tmate-${TMATE_VERSION}-static-linux-${ARCH_TYPE_3}/tmate
# helm binary
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://get.helm.sh/helm-v${HELM_VERSION}-linux-${ARCH_TYPE_2}.tar.gz | tar --directory /usr/local/bin --extract -xz --strip-components 1 linux-${ARCH_TYPE_2}/helm
# talosctl
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L -o /usr/local/bin/talosctl https://github.com/talos-systems/talos/releases/download/v${TALOSCTL_VERSION}/talosctl-linux-${ARCH_TYPE_2} && \
  chmod +x /usr/local/bin/talosctl
# terraform
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${ARCH_TYPE_2}.zip \
  | gunzip -c - > /usr/local/bin/terraform && \
  chmod +x /usr/local/bin/terraform
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v${CRICTL_VERSION}/crictl-v${CRICTL_VERSION}-linux-${ARCH_TYPE_2}.tar.gz \
  | tar --directory /usr/local/bin --extract --gunzip crictl
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://github.com/ahmetb/kubectx/releases/download/v${KUBECTX_VERSION}/kubectx_v${KUBECTX_VERSION}_linux_x86_64.tar.gz | tar --directory /usr/local/bin --extract --ungzip kubectx
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://github.com/junegunn/fzf/releases/download/${FZF_VERSION}/fzf-${FZF_VERSION}-linux_${ARCH_TYPE_2}.tar.gz | tar --directory /usr/local/bin --extract --ungzip fzf
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://github.com/containerd/nerdctl/releases/download/v${NERDCTL_VERSION}/nerdctl-${NERDCTL_VERSION}-linux-${ARCH_TYPE_2}.tar.gz | tar --directory /usr/local/bin --extract --ungzip nerdctl
# Leiningen for clojure
RUN curl -fsSL https://raw.githubusercontent.com/technomancy/leiningen/${LEIN_VERSION}/bin/lein \
    -o /usr/local/bin/lein \
    && chmod +x /usr/local/bin/lein \
    && lein version
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://github.com/knative/client/releases/download/knative-v${KN_VERSION}/kn-linux-${ARCH_TYPE_2} -o /usr/local/bin/kn \
  && chmod +x /usr/local/bin/kn
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://github.com/wagoodman/dive/releases/download/v${DIVE_VERSION}/dive_${DIVE_VERSION}_linux_${ARCH_TYPE_2}.tar.gz | tar --directory /usr/local/bin --extract --ungzip dive
RUN . /usr/lib/sharingio/environment/helper.sh \
  && curl -L https://github.com/owenthereal/upterm/releases/download/v${UPTERM_VERSION}/upterm_linux_${ARCH_TYPE_2}.tar.gz | tar --directory /usr/local/bin --extract --ungzip upterm
RUN set -x \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install golang.org/x/tools/gopls@v$GOPLS_VERSION \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/mikefarah/yq/v4@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/stamblerre/gocode@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/go-delve/delve/cmd/dlv@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/fatih/gomodifytags@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/cweill/gotests/...@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/motemen/gore/cmd/gore@v0.5.2 \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install golang.org/x/tools/cmd/guru@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/minio/mc@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/jessfraz/dockfmt@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install gitlab.com/safesurfer/go-http-server@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/google/go-containerregistry/cmd/crane@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/google/go-containerregistry/cmd/gcrane@latest \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/equinix/metal-cli/cmd/metal@v$METALCLI_VERSION \
  && /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go install github.com/google/ko@v$KO_VERSION

RUN git clone https://github.com/kubernetes-sigs/cluster-api /tmp/cluster-api && \
  cd /tmp/cluster-api && \
  git checkout v${CLUSTERCTL_VERSION} && \
  /bin/env GO111MODULE=on GOPATH=/usr/local/go /usr/local/go/bin/go build -a -trimpath -ldflags "$(bash ./hack/version.sh) -extldflags '-static'" -o /usr/local/bin/clusterctl ./cmd/clusterctl && \
  rm -rf /tmp/cluster-api

# Install Clojure
RUN curl -OL https://download.clojure.org/install/linux-install-${CLOJURE_VERSION}.sh \
    && bash linux-install-${CLOJURE_VERSION}.sh \
    && rm ./linux-install-${CLOJURE_VERSION}.sh
RUN npm install --global prettier @prettier/plugin-php prettier-plugin-solidity

# Test dependencies
RUN set -x && \
  go version && \
  kind version && \
  kubectl version --client && \
  tilt version && \
  gh version && \
  tmate -V && \
  helm version && \
  clusterctl && \
  talosctl version --client && \
  terraform version && \
  dive version && \
  crictl && \
  fzf --version && \
  lein version && \
  clojure --help | head -1 && \
  gopls version && \
  yq --version && \
  dlv version && \
  gore --version && \
  mc --version && \
  dockfmt version && \
  crane version && \
  gcrane version && \
  nerdctl version 2> /dev/null | grep Version && \
  metal --version && \
  docker version -f '{{ . }}' 2> /dev/null | grep 'Docker Engine - Community' && \
  ko version | grep "${KO_VERSION}"

RUN localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8 \
  && touch /etc/localtime
COPY ./usr /usr
ENV LANG=en_US.utf8 \
  DOCKER_CLI_EXPERIMENTAL=enabled \
  USER=ii
USER ii
