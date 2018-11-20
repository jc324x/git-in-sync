#!/bin/bash

project_directory="${HOME}/dev"

rocket=$(echo $'\xF0\x9F\x9A\x80')
warning=$(echo $'\xF0\x9F\x9A\xA7')
prohibited=$(echo $'\xF0\x9F\x9A\xAB')
bold=$(tput bold)
normal=$(tput sgr0)
repositories=($(find -L "$project_directory" -name .git -type d -prune | sed -e "s/\/[^\/]*$//" ))

for repository in "${repositories[@]}"
do
  cd -P "$repository" || return;
  git remote update &> /dev/null
  short_repo=${repository##*/}
  git_status=$(git status --porcelain ) 
  if [ ! -z "$git_status" ]; then 
      echo "${warning}   ${bold}Attention:${normal} Uncommitted changes in ${bold}$short_repo${normal}" 
    else 
      local_repo=$(git rev-parse @)
      remote=$(git rev-parse "@{u}")
      base=$(git merge-base @ "@{u}")
      upstream=$(git rev-parse --abbrev-ref --symbolic-full-name "@{u}")

      if [ "$local_repo" == "$remote" ]; then
          echo "${rocket}   ${bold}$short_repo${normal} Up to date!" 
        elif [ "$local_repo" == "$base" ]; then
          echo "${warning}   ${bold}Attention:${normal} ${bold}$short_repo${normal} is behind ${bold}${upstream}${normal}" 
          read -p "${warning}   ${bold}Attention:${normal} ${bold}Pull${normal} changes? " input
          case $input in yes | pull | Pull ) git pull;; * ) continue;; esac
        elif [ "$remote" == "$base" ]; then
          echo "${warning}   ${bold}Attention:${normal} ${bold}$short_repo${normal} is ahead of ${bold}${upstream}${normal}" 
          read -p "${warning}   ${bold}Attention:${normal} ${bold}Push${normal} changes? " input
          case $input in yes | push | Push ) git push;; * ) continue;; esac
        else
          echo "${prohibited}   ${bold}Attention:${normal} Check on ${bold}$short_repo${normal}." 
      fi
  fi
done
