#!/bin/bash
# pretty=$(git log -1 --date=short --pretty=format:%cd)

cd "$1" || exit

is_repo=$(git rev-parse --is-inside-work-tree)

if [ ! "$is_repo" == "true" ];then
  exit
fi

bold=$(tput bold)
normal=$(tput sgr0)
short_hash=$(git log -1 --pretty=format:%h)
ago=$(git log -1 --pretty=format:%ar)
person=$(git log -1 --pretty=format:%cn)
message=$(git log -1 --pretty=%B)
local_repo=$(git rev-parse @)
remote=$(git rev-parse "@{u}")
base=$(git merge-base @ "@{u}")
upstream=$(git rev-parse --abbrev-ref --symbolic-full-name "@{u}")
git_status=$(git status --porcelain ) 
current=$(pwd)
short_repo=${current##*/}

git remote update &> /dev/null

echo "${bold}Details:${normal} $person $ago ($short_hash)" 
echo "${bold}Message:${normal} $message"

if [ ! -z "$git_status" ]; then
  echo "${bold}Status :${normal} Uncommitted changes in $short_repo" 
  exit
fi

if [ $local_repo = $remote ]; then
    echo "${bold}Status :${normal} $short_repo is up to date with ${upstream}" 
  elif [ $local_repo = $base ]; then
    echo "${bold}Status :${normal} $short_repo is behind ${upstream}" 
  elif [ $remote = $base ]; then
    echo "${bold}Attention:${normal} $short_repo is behind ${upstream}" 
  else
    echo "${bold}Attention:${normal} ???" 
fi
