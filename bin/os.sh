#!/bin/bash
#
# Shell script determining operating system.
#
# ainsley.dev - 16/09/2021

# Convert to WebP
echo '--------------------------------------------'
echo 'System Checks'
echo '--------------------------------------------'

case "$OSTYPE" in
  solaris*) echo "SOLARIS" ;;
  darwin*)  echo "OSX" ;;
  linux*)   echo "LINUX" ;;
  bsd*)     echo "BSD" ;;
  msys*)    echo "WINDOWS" ;;
  cygwin*)  echo "ALSO WINDOWS" ;;
  *)        echo "unknown: $OSTYPE" ;;
esac
