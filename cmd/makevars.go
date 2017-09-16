package cmd

// GitSha is populated at compile time via ld flags in the makefile and represents the git sha of the most recent commit used to build the binary.
var GitSha string

// Version is populated at compile time via ld flags in the makefile and represents the version of the program and is read from the VERSION file in the top level directory of the repository.
var Version string