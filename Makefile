.PHONY: \
	all \
	build build-release build-debug \
	install install-infocenter \
	uninstall uninstall-infocenter

--build:
	@echo "Build infocenter"
	@echo $(BUILD_ARGS)
	go build $(BUILD_ARGS) -o ".bin/infocenter" cmd/infocenter/main.go

build: build-debug

build-release: BUILD_TYPE = release
build-release: --build

build-debug: BUILD_TYPE = debug
build-debug: BUILD_ARGS += -tags pprof
build-debug: --build

build-profile: BUILD_TYPE = profile
build-profile: BUILD_ARGS += -tags pprof
build-profile: --build

install-infocenter:
	@echo "Install infocenter service"
	mkdir -p /etc/infocenter
	cp -rf http /etc/infocenter/
	cp -f .bin/infocenter /usr/sbin/
	cp -f debian/systemd/infocenter.service /etc/systemd/system/
	systemctl enable infocenter
	systemctl start infocenter

uninstall-infocenter:
	@echo "Uninstall infocenter service"
	systemctl stop infocenter
	systemctl disable infocenter
	rm -f /etc/systemd/system/infocenter.service
	rm -f /usr/sbin/infocenter
	rm -rf /etc/infocenter

install: build-release install-infocenter
uninstall: uninstall-infocenter


