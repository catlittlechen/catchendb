#!/bin/bash

rm -f main/version.go
git rev-list HEAD | sort > config.git-hash
LOCALVER=`wc -l config.git-hash | awk '{print $1}'`
if [ $LOCALVER \> 1 ] ; then
    VER=`git rev-list origin/master | sort | join config.git-hash - | wc -l | awk '{print $1}'`
    if [ $VER != $LOCALVER ] ; then
        VER="$VER+$(($LOCALVER-$VER))"
    fi
    if git status -s --porcelain | grep "[ADM]2*\ " ; then
        VER="${VER}M" 
    fi
    VER="$VER $(git rev-list HEAD -n 1 )"
    GIT_VERSION=r$VER
else    
    GIT_VERSION=
    VER="x"
fi
rm -f config.git-hash
 
cat main/version.go.template | sed "s/\$GIT_VERSION/$GIT_VERSION/g" > main/version.go
 
echo "Generated version.go"
