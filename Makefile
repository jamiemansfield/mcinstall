.DEFAULT: help

## Displays the help dialog
help:
	@echo 'mcinstall Build Tool (powered by Make)'
	@echo
	@echo 'Usage:'
	@echo '  make <target>'
	@echo
	@echo 'Targets:'
	@echo '  help                | Shows help'
	@echo '  buildForgeTool      | Builds the Minecraft Forge install tool'
	@echo '  buildLegacyLauncher | Builds our legacy launcher'
	@echo '  build               | Builds mcinstall'

## Build Minecraft Forge install tool
buildForgeTool:
	mkdir -p forge/tool/build/classes
	mkdir -p forge/tool/libs
	curl https://files.minecraftforge.net/maven/net/minecraftforge/installer/2.0.14/installer-2.0.14.jar -o forge/tool/libs/installer.jar
	javac -cp forge/tool/libs/installer.jar -d forge/tool/build/classes \
		forge/tool/src/ModernForgeClientTool.java
	jar cf forge/tool/build/forgetool.jar -C forge/tool/build/classes .

## Generates the Mule file for the Forge install tool
generateForgeTool: buildForgeTool
	go generate ./forge

## Builds Legacy Launcher
buildLegacyLauncher:
	mkdir -p minecraft/launcher/legacylaunch/build/classes
	javac -d minecraft/launcher/legacylaunch/build/classes \
		minecraft/launcher/legacylaunch/src/me/jamiemansfield/mcinstall/LegacyLauncher.java \
		minecraft/launcher/legacylaunch/src/net/minecraft/Launcher.java
	jar cf minecraft/launcher/legacylaunch/build/legacylaunch-1.0.0.jar -C minecraft/launcher/legacylaunch/build/classes .

## Generates the Mule file for LegacyLaunch
generateLegacyLauncher: buildLegacyLauncher
	go generate ./minecraft/launcher

## Builds mcinstall
build:
	go build ./cmd/ftbinstall
	go build ./cmd/technicinstall
