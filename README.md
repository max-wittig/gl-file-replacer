# gl-file-replacer

Tool to replace files in repositories in Gitlab via the merge request flow.

Can be used to automate the file replacements in many repositories.

## Usage

```
Usage of ./gl-file-replacer:
  -branch string
        The branch name to use. Optional. Defaults to refactor-change-file-<filename>
  -file string
        Specify the location of the local file
  -force
        Force branch creation
  -m string
        The commit message
  -repo string
        The location of the GitLab repo. E.g. 'gitlab-org/gitlab'
  -repo-file string
        The repository file to replace
```

### Example to replace the README.md in the repository with a local one

```sh
gl-file-replacer -repo max-wittig/gl-file-replacer -repo-file README.md -file ../other-repo/README.md -m "chore: update readme"
```

### Build

Install on your platform

```sh
make install
```

Build for all available platforms

```sh
make build-all
```

## Planned features

* Regex file replace mode
