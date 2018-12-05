# git-in-sync

## Synopsis

`git-in-sync` parses flags, reads `~/.gisrc.json` then configures directories and
repositories as configured.

## Overview

1. `~/.gisrc.json` outlines directories and their repos as `bundles` and `zones`:

{json block showing configuration}

2. Modes are set with `-m`:

{Here is where the modes go}

3. Additional options can also be set:

{Here is where the additional options go}

4. .gisrc.json unmarshalled to `Config` struct

```go
type Config struct {
	Zones []struct {
		Path    string `json:"path"`
		Bundles []struct {
			User     string   `json:"user"`
			Remote   string   `json:"remote"`
			Division string   `json:"division"`
			Repos    []string `json:"repositories"`
		} `json:"bundles"`
	} `json:"zones"`
}
```

5. The default mode is `verify`. In this mode:

* Verifies directory structure
* Verifies repositories existence [highlight that it's async]
* Gets Git information for each repository

6. If the repository is complete, it's marked as such. Otherwise the user is prompted
   to confirm the next step

   {Steps go here}


7. Once all actions have been cued...final approval...actions dispatched 

## Getting Started
