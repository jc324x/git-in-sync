# Git-In-Sync

**Summary** A collection of bash scripts to keep Git repositories in sync between multiple machines. 
The project is being re-written in Go now...

## Prerequisites
* [Xcode Command Line Tools](http://railsapps.github.io/xcode-command-line-tools.html)
* [Homebrew](https://brew.sh)

## Installation

### 1.) Clone this repository

```
git clone https://jychri.com/git-in-sync
```

### 2.) Set `project_directory` in all scripts

```
project_directory="${HOME}/your-git-projects-live-here"
```
or

```
project_directory="/usr/local/your-git-projects-live-here"
```
etc.

### 3.) Set a custom icon for notifications in `agent-git-status.sh`

```
appIcon="${HOME}/dev/git-in-sync/icon.png"
```

### 4.) Install `terminal-notifier`

```
brew install terminal-notifier
```

### 5.) Setup a Launch Agent for `agent-git-status.sh`

* [launchd.info](http://www.launchd.info/)
* [Lingon X](https://www.peterborgapps.com/lingon/)

## Demos

![Demo#1](https://media.giphy.com/media/l3mZbV8aFlhGSaEjS/giphy.gif)

![Demo#2](https://media.giphy.com/media/3oxHQtInRn0WbyFICA/giphy.gif)
