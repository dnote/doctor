# Dnote Doctor

A utility to automatically fix any issues with local version of [Dnote](https://github.com/dnote/cli).

## Installation

On Linux or macOS, you can use the installation script:

    curl -s https://raw.githubusercontent.com/dnote/doctor/master/install.sh | sh

In some cases, you might need an elevated permission:

    curl -s https://raw.githubusercontent.com/dnote/doctor/master/install.sh | sudo sh

Otherwise, you can download the binary for your platform manually from the [releases page](https://github.com/dnote/doctor/releases).

## Usage

```bash
dnote-doctor
```

The program will:

* detect the version of the Dnote installed on the system
* diagnose any possible issues applicable to the version
* automatically resolve issues, if any.

## LICENSE

GPL

[![Build Status](https://travis-ci.org/dnote/dnote-doctor.svg?branch=master)](https://travis-ci.org/dnote/dnote-doctor)
