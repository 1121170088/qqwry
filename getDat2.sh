#!/bin/sh

echo $(date)

sudo apt install unzip jq zip wget -y

ninfo=$(curl "https://api.github.com/repos/FW27623/qqwry/releases/latest" | sed 's/\\n//g')
echo $ninfo
nrelease=$(echo $ninfo | jq -r '.name')

getRelaseName() {
    info=$(curl "https://api.github.com/repos/1121170088/qqwry/releases/latest" | sed 's/\\n//g')
    echo $info | jq -r '.name'
}

releaseName=$(getRelaseName)

if [ "$nrelease" != "$releaseName" ]
then
  echo "downloading lastest release."
  downloadUrl=$(echo $ninfo | jq -r '.assets[0].browser_download_url')
  mkdir app
  cd app
  echo $downloadUrl
  pwd
  curl -LO "$downloadUrl"
  echo "RELEASE_NAME=$nrelease" >> $GITHUB_ENV
else
  echo "both names are same."
fi

