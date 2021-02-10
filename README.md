# Flock: a high-level SSH library

Flock is a high-level library for executing commands on remote machines.

It's inspired by the Python [Fabric](http://www.fabfile.org) package, but
is written in Go.

**WARNING:** This library is still in development, has no test coverage,
and very likely contains some potentially serious vulnerabilities.  Do **NOT**
use in production settings or in sensitive contexts.

## Known Issues

- **Known command injection issues.**  Commands line arguments are currently not being escaped properly.  Do not use in production settings.
- No tests!!
- File transfer uses `cat` and `bash`.  It needs to be written to use something like `scp`
- Code quality is pretty dodgy.  It neads some refactoring and a good clean up (which is blocked by no tests).
- The use of `sudo` needs to be