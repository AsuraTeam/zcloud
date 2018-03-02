# Miscellaneous Go utilities

This repository contains miscellaneous Go utility packages. Mature packages get
migrated to here from my holditall,
[degoutils](https://github.com/hlandau/degoutils), or to their own package.

Currently, all packages in here depend only on the standard library. In the
future, non-standard library dependencies in this repository should be kept
minimized, to mitigate the impact on tools, such as vendorization tools, which
work on a repository and not a package level.
