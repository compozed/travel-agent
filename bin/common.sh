manifest=.tmp/concourse_deployment_manifest.yml
pre_merged_manifest=.tmp/pre_merged_manifest.yml
vendored_ego=$GOPATH/src/github.com/compozed/travel-agent/vendor/github.com/benbjohnson/ego/cmd/ego/main.go

target(){
  TARGET=$1

  if [  -z "$TARGET" ]
  then
    echo "${red}You must provide a concourse target${reset}"
    echo ''
    echo 'USAGE: ./travel-agent target CONCOURSE_TARGET. eg: http://ATC_IP:8080'
    exit 1
  fi

  fly -t concourse login -k -c $TARGET
}

# Clean old atifacts
clean(){
  rm -f manifest.go manifest 
}

clone_project(){
  project_path=$1
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

book() {
  TRAVEL_AGENT_CONFIG=$1
  shift
  FILES_TO_MERGE=$@

  if [ -z "$TRAVEL_AGENT_CONFIG" ] ; then
    echo "${red}===> provide TAVEL_AGENT_CONFIG if you want to generate a manifest${reset}"
    exit 1
  fi

  if [ -z "$FILES_TO_MERGE" ] ; then
    echo "${red}===> provide FILES_TO_MERGE if you want to spruce merge secrets to manifest${reset}"
    exit 1
  fi

  if [ -n "$TRAVEL_AGENT_CONFIG" ] ; then
    TRAVEL_AGENT_PROJECT=$(cat "$TRAVEL_AGENT_CONFIG" | grep -v -e "^#" | grep -e "^git_project:" |  awk -F"git_project:" '{print $2}' )
    clone_project "$TRAVEL_AGENT_PROJECT"
  fi

  if ! [ -d "ci/manifest" ]; then
    echo "This does not look like a travel agent project"
    exit 1
  fi

  pushd ci/manifest > /dev/null

  clean

  printf "${green}===> Rendering the ${magenta}.ego${green} manifest...${reset}"
  go run $vendored_ego -package main -o manifest.go
  NAME=$(grep -E "^name:" $TRAVEL_AGENT_CONFIG | awk -F " " '{print $2}')
  mkdir -p .tmp
  go run manifest.go main.go $TRAVEL_AGENT_CONFIG > $pre_merged_manifest
  echo "${green}done${reset}"

  printf "${green}===> Authenticating to Vault via ${magenta}safe${green}...${reset}"
  export VAULT_ADDR="$(safe target 2>&1 | sed 's/.* http\([^ ]*\) .*/http\1/')"
  export VAULT_TOKEN="$(safe vault token renew --format=json | jq -r '.auth.client_token')"
  echo "${green}done${reset}"

  printf "${green}===> ${magenta}spruce${green} merging secrets from settings.yml into the generated manifest ...${reset}"
  spruce merge --prune meta $pre_merged_manifest $FILES_TO_MERGE > $manifest
  echo "${green}done${reset}"

  echo "${green}===> Updating concourse pipeline configuration via ${magenta}fly${green}...${reset}"
  fly -t concourse set-pipeline -c $manifest -p $NAME $fly_opts

  popd > /dev/null
}
