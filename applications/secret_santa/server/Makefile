.PHONY: packages
packages:
	sudo apt-get install software-properties-common
	sudo apt-get install python-software-repositories
	sudo apt update
	sudo apt-get install openjdk-16-jdk-headless

minecraft_server.jar:
	curl https://launcher.mojang.com/v1/objects/a16d67e5807f57fc4e550299cf20226194497dc2/server.jar --output minecraft_server.jar

.PHONY: server
server:
	java -Xmx4G minecraft_server.jar

