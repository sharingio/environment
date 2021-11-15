#!/usr/bin/env bash
set -euo pipefail

if [ "${PAIR_ENVIRONMENT_DEBUG:-}" = "true" ]; then
    set -x
fi

cd "$HOME"

if [ -z "${GIT_AUTHOR_EMAIL:-}" ]; then
    echo "ERROR: GIT_AUTHOR_EMAIL env Must be set"
    exit 1
fi
if [ -z "${GIT_AUTHOR_NAME:-}" ]; then
    echo "ERROR: GIT_AUTHOR_NAME env Must be set"
    exit 1
fi

export TMATE_SOCKET="${TMATE_SOCKET:-/tmp/ii.default.target.iisocket}"
export TMATE_SOCKET_NAME="$(basename "${TMATE_SOCKET}")"
if tmate -S "${TMATE_SOCKET}" wait-for tmate-ready 2> /dev/null; then
    set +x
    echo "Already initialised with tmate ready."
    echo "Use: attach"
    exit 0
fi

# Generate an ssh-key if one doesn't exist
# Ensure that tmate can launch a session
if [ ! -f "$HOME/.ssh/id_rsa" ]
then
    ssh-keygen -b 4096 -t rsa -f ~/.ssh/id_rsa -q -N ""
fi

export ALTERNATE_EDITOR=""
export INIT_ORG_FILE="${INIT_ORG_FILE:-$HOME}"
export INIT_DEFAULT_DIR="${INIT_DEFAULT_DIR:-$HOME}"
export INIT_DEFAULT_REPOS="${INIT_DEFAULT_REPOS:-}"
export INIT_DEFAULT_REPOS_FOLDER="${INIT_DEFAULT_REPOS_FOLDER:-$INIT_DEFAULT_DIR}"
export INIT_PREFINISH_BLOCK="${INIT_PREFINISH_BLOCK:-}"
export PROJECT_CLONE_STRUCTURE="${PROJECT_STRUCTURE:-structured}"

# Load SSH_AUTH_SOCK
. /usr/local/bin/ssh-find-agent.sh

# Ensure that the home folder has been repopulated after a PVC recreates it
if [ "${REINIT_HOME_FOLDER:-false}" = "true" ]; then
    (
        if [ ! -f "${HOME}"/.pair-environment-has-reinit ]; then
            cd /etc/skel
            cp -r . /home/ii
            [ -f "${HOME}"/.kube/config ] \
                && chmod 0600 "${HOME}"/.kube/config
            touch "${HOME}"/.pair-environment-has-reinit
        fi
    )
fi

# Clone all projects
(
    if [ -n "$INIT_DEFAULT_REPOS" ]; then
        mkdir -p "${INIT_DEFAULT_REPOS_FOLDER}"
        for repo in ${INIT_DEFAULT_REPOS}; do
            if [ "$PROJECT_CLONE_STRUCTURE" = "structured" ]; then
                git-clone-structured "$repo" || true
            elif [ "$PROJECT_CLONE_STRUCTURE" = "plain" ]; then
                git clone -v --recursive "$repo" || true
            fi
        done
    fi
    cd
    eval "$INIT_PREFINISH_BLOCK"
)

if [ ! -d "/home/ii/.sharing.io" ]; then
    git clone "https://github.com/${SHARINGIO_PAIR_USER:-$USER}/.sharing.io" || \
        git clone https://github.com/sharingio/.sharing.io
fi

. /home/ii/.sharing.io/sharingio-pair-preinit-script.sh

# This background process will ensure tmate attach commands
# call osc52-tmate.sh to set the ssh/web uri for this session via osc52
# We need to wait's until the socket exists, and tmate is ready for commands
# before doing so. (Would be easier if this were a config option for .tmate.conf)
cd "${INIT_DEFAULT_DIR}"
(
    /usr/local/bin/tmate-wait-for-socket.sh
    tmate -S "${TMATE_SOCKET}" set-hook -ug client-attached # unset
    tmate -S "${TMATE_SOCKET}" set-hook -g client-attached 'run-shell "tmate new-window osc52-tmate.sh"'
)&

# This is our primary background process for humacs
# a tmate session in foreground mode, respawning if it dies
# A default directory and org file are used to start emacsclient as the main window
tmate -F -v -S "${TMATE_SOCKET}" new-session -d -c "${INIT_DEFAULT_DIR}" emacsclient --tty "${INIT_ORG_FILE}"
