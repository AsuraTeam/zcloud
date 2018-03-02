echo $@
echo $#
if [ $# < 1 ] ; then
exit
fi
svn up
rsync -avr ../cloud/ ../zcloud/ --exclude=key/* --exclude=conf/*  --exclude=.svn --exclude=.git --delete --exclude=zcloud.iml --exclude=make.go
cd ../zcloud/
git add .
git commit -m $1
git push origin master
