#!/bin/sh

echo $(date)

sudo apt install unzip jq zip wget -y

wget https://github.com/dscharrer/innoextract/releases/download/1.9/innoextract-1.9-linux.tar.xz
tar xf innoextract-1.9-linux.tar.xz

getRelaseName() {
    info=$(curl https://api.github.com/repos/1121170088/qqwry/releases/latest)
    echo $info | jq -r '.name'
}

content=$(curl "https://mp.weixin.qq.com/mp/appmsgalbum?action=getalbum&album_id=2329805780276838401&f=json&count=10")
releasePage=$(echo $content | jq -r '.getalbum_resp|.article_list[0].url')


lastestName=$(echo $content | jq -r '.getalbum_resp|.article_list[0].title')
releaseName=$(getRelaseName)

if [ "$lastestName" != "$releaseName" ]
then
  echo "downloading lastest release."
  releasePage=$(echo $releasePage | sed 's/http/https/')
  page=$(curl "$releasePage")
  downloadUrl=$(echo $page  | grep -Po 'https://www.cz88.net/soft/.*?\.zip')
  curl "$downloadUrl" -o 1.zip
  unzip 1.zip
  ./innoextract-1.9-linux/innoextract setup.exe
  echo "RELEASE_NAME=$lastestName" >> $GITHUB_ENV
else
  echo "both names are same."
fi

