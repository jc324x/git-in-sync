#!/bin/bash

project_directory="${HOME}/dev"
appIcon="${HOME}/dev/git-in-sync/media/icon.png"

rocket=$(echo $'\xF0\x9F\x9A\x80')
warning=$(echo $'\xF0\x9F\x9A\xA7')
repositories=($(find -L "$project_directory" -name .git -type d -prune | sed -e "s/\/[^\/]*$//" ))
clean_repositories=()
synced_repositories=()

for repository in "${repositories[@]}"
do
    cd -P "$repository" || return;
    short_repo=${repository##*/}
    git_status=$(git status --porcelain ) 
    if [ ! -z "$git_status" ]
    then 
        terminal-notifier -title "${short_repo}" -message "${warning}  Uncommitted changes" -appIcon "$appIcon"
    else 
        clean_repositories+=($repository)
    fi
done

for repository in "${clean_repositories[@]}"
do
    cd -P "$repository" || return;
    short_repo=${repository##*/}
    git remote update &> /dev/null
    native=$(git rev-parse @)
    remote=$(git rev-parse "@{u}")
    base=$(git merge-base @ "@{u}")
    upstream=$(git rev-parse --abbrev-ref --symbolic-full-name "@{u}")

    if [ $native = $remote ]; then
      synced_repositories+=($repository)
    elif [ $native = $base ]; then
      terminal-notifier -title "${short_repo}" -message "${warning}  ${short_repo} is behind ${upstream}" -appIcon "$appIcon"
    elif [ $remote = $base ]; then
      terminal-notifier -title "${short_repo}" -message "${warning}  ${short_repo} is ahead of ${upstream}" -appIcon "$appIcon"
    fi
done

if [ ${#repositories[@]} == ${#synced_repositories[@]} ]; then
  terminal-notifier -title "Git-In-Sync" -message "${rocket}   ${#repositories[@]} / ${#clean_repositories[@]} Up to date!" -appIcon "$appIcon"
fi 

exit
