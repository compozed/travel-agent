assert_dir_exists() {
  if [ -d "$1" ]; then
    return 0
  else 
    flunk "Dir does not exists \`$1'"
  fi
}

assert_dir_not_exists() {
  if [ -d "$1" ]; then
    flunk "Dir exists \`$1'"
  else 
    return 0
  fi
}

assert_dir_not_empty(){
  if [ -z "$(ls -A $1)" ]; then
    flunk "Dir is empty \`$1'"
  else 
    return 0
  fi
}

flunk() {
  { if [ "$#" -eq 0 ]; then cat -
    else echo -e "$@"
    fi
  }
  return 1
}

assert_line() {
  if [ "$1" -ge 0 ] 2>/dev/null; then
    assert_equal "$2" "${lines[$1]}"
  else
    local line
    for line in "${lines[@]}"; do
      if [ "$line" = "$1" ]; then return 0; fi
    done
    flunk "expected line \`$1'"
  fi
}

assert_equal() {
  if [ "$1" != "$2" ]; then
    { echo "expected: $1"
      echo "actual:   $2"
    } | flunk
  fi
}

assert_match() {
  out="${2:-${output}}"
  if [ ! $(echo "${out}" | grep -- "${1}") ]; then
    { echo "expected match: $1"
      echo "actual: $out"
    } | flunk
  fi
}

assert_not_match() {
  out="${2:-${output}}"
  if [ $(echo "${out}" | grep -- "${1}") ]; then
    flunk "expected not to match: $1"
  fi
}

assert_success() {
  if [ "$status" -ne 0 ]; then
    flunk "command failed with exit status ${status}\\noutput: ${output}"
  elif [ "$#" -gt 0 ]; then
    assert_output "$1"
  fi
}

assert_failure() {
  if [ "$status" -eq 0 ]; then
    flunk "expected failed exit status"
  elif [ "$#" -gt 0 ]; then
    assert_output "$1"
  fi
}
