#!/bin/bash

user=ffonord
cameraIp=192.168.1.133
ftpWorkTime=1800
baseDir=/home/$user/Videos/
ftpPath=$cameraIp/DCIM/101MEDIA/

echo "Checking if the camera is on the home network"
cameraNotExists=$(ping $cameraIp -c 1 | grep '0 received')
[ "$cameraNotExists" ] && echo 'Camera not found' && exit

echo "Start ftp server, ip: $cameraIp"
setsid "$( (
  sleep 1
  echo root
  sleep 1
  echo '/tmp/fuse_d/ftp.sh'
  sleep $ftpWorkTime
  echo exit
  sleep 1
) | telnet $cameraIp)" >/dev/null 2>&1 </dev/null &

#start spinner
i=1
sp="/-\|"
while [ -z "$(ps -fade | grep "sleep $ftpWorkTime" | grep -v grep | awk '{print $2}')" ]; do
  printf "\b${sp:i++%${#sp}:1}"
done
sleep 1

echo "Start download videos, ip: $cameraIp"
dirPath="$baseDir"$(date +%F)
[ ! -d "$dirPath" ] && mkdir -p "$dirPath"
cd "$dirPath" || exit

#fetch files list
wget --no-remove-listing --user=root ftp://$ftpPath >/dev/null 2>&1
#only filenames *.MP4
fileNames="$(cat .listing | sed -n '/.MP4/p' | awk '{print $NF}')"

for fileName in $fileNames; do
  mp4File="${fileName%.*}".MP4
  secFile="${fileName%.*}".SEC

  filePath=$dirPath/$mp4File

  wget --no-remove-listing \
    --no-parent --recursive \
    --level=1 --no-directories \
    --no-host-directories \
    --user=root \
    --no-clobber \
    --continue \
    -O "$filePath" \
    ftp://"$ftpPath$mp4File"

  localFileSize=$(stat -c%s "$filePath")
  externalFileSize=$(cat .listing | grep "$mp4File" | awk '{print $5}')

  if [ "$localFileSize" = "$externalFileSize" ]; then
    curl ftp://$ftpPath -X "DELE $mp4File" --user root:
    curl ftp://$ftpPath -X "DELE $secFile" --user root:
  else
    echo "Files not equals $mp4File, local: $localFileSize, external: $externalFileSize"
    rm "$filePath"
  fi
done

echo "Stop ftp server, ip: $cameraIp"
kill -9 $(ps -fade | grep "sleep $ftpWorkTime" | grep -v grep | awk '{print $2}')
rm ".listing" index.html*
