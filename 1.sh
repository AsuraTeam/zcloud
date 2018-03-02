if [ $# -lt 1 ] ; then
echo "1.sh 描述信息"
exit
fi
svn up
rsync -avr ../cloud/ ../zcloud/ --exclude=key/* --exclude=conf/*  --exclude=.svn --exclude=.git --delete --exclude=zcloud.iml --exclude=make.go --exclude=.idea
cd ../zcloud/
git add .
git commit -m $1
git status
git push origin master
