cd
TESTDIR="testarcbranch"
mkdir "$TESTDIR"
cd "$TESTDIR"
git init
cp -r ~/arcbranch/examples/* .
git add .
git commit -m "Initial commit with example files"
arcbranch 4
echo "Test repo created and arcbranch 4 run in $PWD"