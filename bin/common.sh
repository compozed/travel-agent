#!/usr/bin/env bash

exec >&2
set -e

if [[ "${DEBUG,,}" == "true" ]]; then
  set -x
  echo "Environment Variables:"
  env
fi

OUTPUT=.tmp
VENDORED_EGO=$GOPATH/src/github.com/compozed/travel-agent/vendor/github.com/benbjohnson/ego/cmd/ego/main.go

red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
magenta=`tput setaf 5`
reset=`tput sgr0`

function log() {
  printf "${green}$1${reset}"
}

function log_error() {
  printf "${red}$1${reset}"
}

function target(){
  TARGET=$1

  if [  -z "$TARGET" ]; then
    log_error "You must provide a concourse target"
    echo ''
    echo 'USAGE: ./travel-agent target CONCOURSE_TARGET. eg: http://ATC_IP:8080'
    exit 1
  fi

    fly -t concourse login -k -c $TARGET
}

function clean(){
  rm -f manifest.go manifest
}

function clone_project(){
  local project_path=$1
  echo $project_path
  if [ ! -z "$project_path" ] && [[ $1 == *".git" ]]  ; then
    project_dir=`basename $project_path .git`

    echo "Cloning $project_path ..."

    pushd /tmp > /dev/null
    rm -rf $project_dir
    git clone $project_path
    popd > /dev/null

    pushd /tmp/$project_dir > /dev/null
  else
    pushd $project_path > /dev/null
  fi
}

help() {
  cat << EOM

Travel agent helps you write concourse pipelines without repeating yourself.
TDD your pipeline templates and create jobs and resources for 1..N environments.

Running travel agent:

travel-agent SUBCOMMAND

Subcommands:

help
intit       - Generates and upgrades travel agent project
target      - Sets concourse target. EG: https://1.2.3.4:9090
book        - compiles and deploys manifest.ego (manifest.ego)
EOM
}

function book() {
  TRAVEL_AGENT_CONFIG=$1
  shift
  FILES_TO_MERGE=$@

  local manifest=$OUTPUT/concourse_deployment_manifest.yml
  local pre_merged_manifest=$OUTPUT/pre_merged_manifest.yml

  if [ -z "$TRAVEL_AGENT_CONFIG" ] ; then
    log_error "===> provide TAVEL_AGENT_CONFIG if you want to generate a manifest"
    exit 1
  fi

  if [ -z "$FILES_TO_MERGE" ] ; then
    log_error "$===> provide FILES_TO_MERGE if you want to spruce merge secrets to manifest"
    exit 1
  fi

  if [ -n "$TRAVEL_AGENT_CONFIG" ] ; then
    TRAVEL_AGENT_PROJECT=$( \
      cat "$TRAVEL_AGENT_CONFIG" | \
      grep -v -e "^#" | \
      grep -e "^git_project:" |  \
      awk -F"git_project:" '{print $2}' )

    clone_project "$TRAVEL_AGENT_PROJECT"
  fi

  if ! [ -d "ci/manifest" ]; then
    log_error "This does not look like a travel agent project"
    exit 1
  fi

  pushd ci/manifest > /dev/null

  clean

  log "===> Rendering the ${magenta}.ego${green} manifest..."
  go run $VENDORED_EGO -package main -o manifest.go
  NAME=$(grep -E "^name:" $TRAVEL_AGENT_CONFIG | awk -F " " '{print $2}')
  mkdir -p .tmp
  go run manifest.go main.go $TRAVEL_AGENT_CONFIG > $pre_merged_manifest
  echo "${green}done${reset}"

  log "===> Authenticating to Vault via ${magenta}safe..."
  export VAULT_ADDR="$(cat ~/.svtoken | grep '^vault:' | cut -d " " -f2)"
  export VAULT_TOKEN="$(safe vault token-renew --format=json | jq -r '.auth.client_token')"
  log "done\n"

  log "===> ${magenta}spruce${green} merging secrets from settings.yml into the generated manifest ..."
  spruce merge --prune meta --prune envs --prune git_project --prune name \
    $TRAVEL_AGENT_CONFIG $pre_merged_manifest $FILES_TO_MERGE > $manifest
  log "done\n"

  log "===> Updating concourse pipeline configuration via ${magenta}fly${green}..."
  fly -t concourse set-pipeline --check-creds -c $manifest -p $NAME $fly_opts

  popd > /dev/null
}
