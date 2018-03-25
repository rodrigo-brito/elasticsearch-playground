run:
	docker run -p 9200:9200 -p 9300:9300 -e "http.host=0.0.0.0" -e "transport.host=127.0.0.1" -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:5.4.3
rm:
	curl -XDELETE -u elastic:changeme http://localhost:9200/playground