dev-front:
	cd front && yarn && yarn start

build-front:
	cd front && yarn && yarn build

docker-front:
	docker build -f Dockerfile.front . -t bonitto-front

run-front: docker-front
	docker run --rm -p 8080:80 bonitto-front

docker-back:
	docker build -f Dockerfile.back . -t bonitto-back

deploy-front: docker-front
	@curl -X GET https://hub.docker.com/v2/repositories/ckcks12/bonitto-front/tags 2>/dev/null | jq -rM '.results[].name' \
	&& read -p "front tag: " TAG \
	&& docker tag bonitto-front:latest ckcks12/bonitto-front:$${TAG} \
	&& docker push ckcks12/bonitto-front:$${TAG} \
	&& echo "🚀 deployed: ckcks12/bonitto-front:$${TAG}"

deploy-back: docker-back
	@curl -X GET https://hub.docker.com/v2/repositories/ckcks12/bonitto-back/tags 2>/dev/null | jq -rM '.results[].name' \
	&& read -p "back tag: " TAG \
	&& docker tag bonitto-back:latest ckcks12/bonitto-back:$${TAG} \
	&& docker push ckcks12/bonitto-back:$${TAG} \
	&& echo "🚀 deployed: ckcks12/bonitto-back:$${TAG}"

deploy: deploy-front deploy-back
