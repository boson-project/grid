# grid

[![Main Build Status](https://github.com/boson-project/grid/workflows/Main/badge.svg?branch=main)](https://github.com/boson-project/grid/actions?query=workflow%3AMain+branch%3Amain)
[![Develop Build Status](https://github.com/boson-project/grid/workflows/Develop/badge.svg?branch=develop&label=develop)](https://github.com/boson-project/grid/actions?query=workflow%3ADevelop+branch%3Adevelop)
[![Documentation](https://godoc.org/github.com/boson-project/grid?status.svg)](http://godoc.org/github.com/boson-project/grid)
[![GitHub Issues](https://img.shields.io/github/issues/boson-project/grid.svg)](https://github.com/boson-project/grid/issues)
[![License](https://img.shields.io/github/license/boson-project/grid)](https://github.com/boson-project/grid/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/boson-project/grid.svg?label=Release)](https://github.com/boson-project/grid/releases)


Grid adapter for converting faas platforms into a predictable pool of homogonized compute resources.

## Architecture

This service runs along side function instances and is implicitly available at `127.0.0.1:1111`.  The API presented is a lazily-populated intersection of feature sets provided by supported platforms.

Utilized during the build process of buildpacks, and transparent to a function developer insofar as implementation.

## Status

This is currently a work in progress, with merely stub implementations of some example endpoints.  These will be expanded as we settle on use cases we would like to be available to "Portable Functions" (Functions which will be able to run on multiple cloud platforms presuming they are built and packaged using Boson buildpacks).

