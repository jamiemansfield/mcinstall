.DEFAULT: help

## Displays the help dialog
help:
	@echo 'ftbinstall Build Tool (powered by Make)'
	@echo
	@echo 'Usage:'
	@echo '  make <target>'
	@echo
	@echo 'Targets:'
	@echo '  help             | Shows help'
	@echo '  compileForgeTool | Compiles the Minecraft Forge install tool'
	@echo '  build            | Builds ftbinstall'

## Compiles the Minecraft Forge install tool
compileForgeTool:
	mkdir -p forgeinstall/tool/build
	mkdir -p forgeinstall/tool/libs
	curl https://files.minecraftforge.net/maven/net/minecraftforge/installer/2.0.14/installer-2.0.14.jar -o forgeinstall/tool/libs/installer.jar
	javac -cp forgeinstall/tool/libs/installer.jar -d forgeinstall/tool/build forgeinstall/tool/src/ModernForgeClientTool.java

## Builds ftbinstall
build: compileForgeTool
	go generate ./forgeinstall
	go build ./cmd/ftbinstall
