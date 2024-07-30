$image="traefikcertextractor"
$version="latest"


echo "$image":"$version"

docker build -t jmgitde/"$image":"$version" .
docker push jmgitde/"$image":"$version"