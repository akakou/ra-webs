NAME=$1
REPO=$2
COMMITID=$3
BASE_PATH=$4
BASE_PROGRAM_PATH=$5
EXECUTABLE=$6

# echo "NAME: $NAME"
# echo "REPO: $REPO"
# echo "BASE_PATH: $BASE_PATH"
# echo "BRANCH: $BRANCH"
# echo "BASE_PROGRAM_PATH: $BASE_PROGRAM_PATH"
# echo "EXECUTABLE: $EXECUTABLE"

LOG_PATH=`pwd`/$BASE_PATH/$NAME.txt
REPO_PATH=`pwd`/$BASE_PATH/$NAME

### Clone
echo -en "yes\n" | git clone $REPO $REPO_PATH --depth 1 --branch $COMMITID --single-branch --quiet &> $LOG_PATH
cd $REPO_PATH/$BASE_PROGRAM_PATH &> $LOG_PATH

### Show Git Log
git log -1 --pretty=format:%H
echo ""

### Build
ego-go build -buildvcs=false -trimpath=true &> $LOG_PATH
ego sign $EXECUTABLE &> $LOG_PATH

### Show UniqueID
ego uniqueid $EXECUTABLE
