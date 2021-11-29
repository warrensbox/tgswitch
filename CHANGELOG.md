# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.0] - 2021-11-14
### Context
- In the past, `tgswitch` uses github's API to get the list of releases.
- `tgswitch` uses client autorization key to access the terragrunt releases, however, github limits the number of api calls `tgswitch` can make.
- As a result, user cannot immediately download the version of terragrunt they want. They had to wait.

### Added
- `tgswitch` will now get th list of releases from [terragrunt list page maintained by warrensbox](https://warrensbox.github.io/terragunt-versions-list/)
- `tgswitch` will directly download from the [terragrunt release page](https://github.com/gruntwork-io/terragrunt/releases)

### Removed
- removed all functions that would make API github calls

## [0.5.0] - 2021-11-28
### Bug Fixes
- Fixed issue where symlink points to terraform
