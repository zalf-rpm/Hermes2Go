# This script is used to tag a new release version of the project.

# this script can be used on windows with git bash
# or on linux

# when to increment which version number?
# major: when you make incompatible API changes
# minor: when you add functionality in a backwards-compatible manner
# patch: when you make backwards-compatible bug fixes


# input parameters:
# major=(true/false), minor=(true/false), patch=(true/false): to increase the version number to be released
# extra=(string): to add an extra string to the version number to be released (e.g. "beta1", "test", "pre-release", etc.)
MAJOR=$1 
MINOR=$2
PATCH=$3
EXTRA=$4

git describe --always --tags --long > version.txt
VERSION=$(cat version.txt) 
echo "VERSION: ${VERSION}"

# split version number into parts v0.3.2
MAJOR_VERSION=$(echo $VERSION | cut -d'-' -f1 | cut -d'.' -f1 | cut -d'v' -f2)
MINOR_VERSION=$(echo $VERSION | cut -d'-' -f1 | cut -d'.' -f2)
PATCH_VERSION=$(echo $VERSION | cut -d'-' -f1 | cut -d'.' -f3)

# check if version number is valid
if [ -z "$MAJOR_VERSION" ] || [ -z "$MINOR_VERSION" ] || [ -z "$PATCH_VERSION" ] ; then
    echo "ERROR: version number is not valid"
    exit 1
fi

# if a version number change is requested
if [ "$MAJOR" = true ] || [ "$MINOR" = true ] || [ "$PATCH" = true ] ; then
    # check if the version number is valid
    if [ "$MAJOR" = true ] && [ "$MINOR" = true ] ; then
        echo "ERROR: major and minor version number cannot be increased at the same time"
        exit 1
    fi
    if [ "$MAJOR" = true ] && [ "$PATCH" = true ] ; then
        echo "ERROR: major and patch version number cannot be increased at the same time"
        exit 1
    fi
    if [ "$MINOR" = true ] && [ "$PATCH" = true ] ; then
        echo "ERROR: minor and patch version number cannot be increased at the same time"
        exit 1
    fi
    if [ "$MAJOR" = true ] && [ "$MINOR" = true ] && [ "$PATCH" = true ] ; then
        echo "ERROR: major, minor and patch version number cannot be increased at the same time"
        exit 1
    fi

    # increase version number
    if [ "$MAJOR" = true ] ; then
        MAJOR_VERSION=$((MAJOR_VERSION+1))
        MINOR_VERSION=0
        PATCH_VERSION=0
    elif [ "$MINOR" = true ] ; then
        MINOR_VERSION=$((MINOR_VERSION+1))
        PATCH_VERSION=0
    elif [ "$PATCH" = true ] ; then
        PATCH_VERSION=$((PATCH_VERSION+1))
    fi
fi

# create new version number
NEW_VERSION="v${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}"

# add extra string to version number
if [ ! -z "$EXTRA" ] ; then
    NEW_VERSION="${NEW_VERSION}.${EXTRA}"
fi

# create a git tag with the version number
echo "TAG: ${NEW_VERSION}"
git tag -a ${NEW_VERSION} -m "version ${NEW_VERSION}"
push the tag to the remote repository
git push origin ${NEW_VERSION}
