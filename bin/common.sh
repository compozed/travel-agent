target(){
  TARGET=$1

  if [  -z "$TARGET" ]
  then
    echo "${red}You must provide a concourse target${reset}"
    echo ''
    echo 'USAGE: ./travel-agent target CONCOURSE_TARGET'
    exit 1
  fi

  fly -t travel-agent login -c $TARGET
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
bootstrap   - Generates and upgrades travel agent project
target      - Sets concourse target. EG: https://1.2.3.4:9090
book  - test and deploy your pipeline template (manifest.ego) against 2 dummy assets(dev and prod)
EOM
}

book() {
  echo 'Booking...'
  TRAVEL_AGENT_CONFIG=$1
  SPRUCE_FILE=$2

  if [ -n "$TRAVEL_AGENT_CONFIG" ] ; then
    TRAVEL_AGENT_PROJECT=$(cat "$1" | grep -v -e "^#" | grep -e "^git_project:" |  awk -F"git_project:" '{print $2}' )
    clone_project "$TRAVEL_AGENT_PROJECT"
  fi

  if [ -d "ci/manifest" ]; then
    echo "Travel agent project found"
  else
    echo "This does not look like a travel agent project"
  fi

  pushd ci/manifest > /dev/null

  clean

  # Compile .ego into manifest.go
  ego -package main -o manifest.go

  # Run test suite to match template with assets/*
  ginkgo -r -failFast

  # Creates a manifest if a TRAVEL_AGENT_CONFIG was provided
  if [ -z "$TRAVEL_AGENT_CONFIG" ] ; then
    echo "${yellow}===> INFO: provide TAVEL_AGENT_CONFIG if you want to generate a manifest${reset}"
  else

    echo "${green}===> Generating manifest for enviroments provided by the travel-agent.yml config file ...${reset}"
    NAME=$(grep -E "^name:" $TRAVEL_AGENT_CONFIG | awk -F " " '{print $2}')
    pre_merged_manifest=.tmp/pre_merged_manifest.yml
    mkdir -p .tmp
    go run manifest.go main.go $TRAVEL_AGENT_CONFIG > $pre_merged_manifest
    echo "${green}===> done${reset}"

    if [ -n "$DEBUG" ] ; then
      echo "PREMERGED MANIFEST"
      cat $pre_merged_manifest
    fi
  fi


  # Merges file with spruce if SPRUCE_FILE was provided
  if [ -z "$SPRUCE_FILE" ] ; then
    manifest=$pre_merged_manifest
    echo "${yellow}===> INFO: provide SPRUCE_FILE if you want to spruce merge secrets to manifest${reset}"
  else
    echo "${green}===> Merging secrets from spruce-secret.yml into the generated manifest ...${reset}"
    manifest=.tmp/concourse_deployment_manifest.yml
    spruce merge --prune config $pre_merged_manifest $SPRUCE_FILE > $manifest
    echo "${green}===> done${reset}"
  fi

  # Outputs spruce merging errors, hides output
  if [ -n "$DEBUG" ] ; then
    spruce merge $manifest > /dev/null
  else
    spruce merge $manifest 
  fi

  # Deploys if TARGET is set
  if [ -z "$NAME" ] || [ ! -z "$DRY_RUN" ]; then
    echo "${yellow}===> INFO: set target and provie TRAVEL_AGENT_CONFIG if you want to deploy${reset}"
  else
    fly -t travel-agent set-pipeline -c $manifest -p $NAME
  fi
  popd > /dev/null

}
