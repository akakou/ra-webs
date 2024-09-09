NAME=$1
REPO=$2
BASE_PATH=$3
BRANCH=$4
BASE_PROGRAM_PATH=$5
EXECUTABLE=$6

LOG_PATH=`pwd`/$BASE_PATH/$NAME.txt
REPO_PATH=`pwd`/$BASE_PATH/$NAME

### Clone
echo -en "yes\n" | git clone $REPO $REPO_PATH --depth 1 --branch $BRANCH --single-branch --quiet &> $LOG_PATH
cd $REPO_PATH/$BASE_PROGRAM_PATH &> $LOG_PATH

### Show Git Log
git log -1 --pretty=format:%H
echo ""

### Build
ego-go build -buildvcs=false -trimpath=true &> $LOG_PATH
ego sign $EXECUTABLE &> $LOG_PATH

### Show UniqueID
ego uniqueid $EXECUTABLE
