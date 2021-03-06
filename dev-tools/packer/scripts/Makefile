DATE:=$(shell date +%y%m%d%H%M%S)
BUILDID?=$(DATE)

makefile_abspath:=$(abspath $(lastword $(MAKEFILE_LIST))/..)
packer_absdir=$(shell dirname ${makefile_abspath})


%/deb: % build/god-linux-386 build/god-linux-amd64 fpm-image
	ARCH=386 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/debian/build.sh
	ARCH=amd64 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/debian/build.sh

%/rpm: % build/god-linux-386 build/god-linux-amd64 fpm-image
	ARCH=386 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/centos/build.sh
	ARCH=amd64 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/centos/build.sh

%/darwin: %
	ARCH=amd64 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/darwin/build.sh

%/win: %
	ARCH=386 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/windows/build.sh
	ARCH=amd64 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/windows/build.sh

%/bin: %
	ARCH=386 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/binary/build.sh
	ARCH=amd64 BEAT=$(@D) BUILDID=$(BUILDID) $(packer_absdir)/platforms/binary/build.sh

.PHONY: deps
deps:
	go get -u github.com/tsg/xgo
	go get -u github.com/tsg/gotpl

.PHONY: xgo-image
xgo-image:
	cd $(packer_absdir)/docker/xgo-image/; ./build.sh
	cd $(packer_absdir)/docker/xgo-image-deb6/; ./build.sh

.PHONY: fpm-image
fpm-image:
	docker build --rm=true -t tudorg/fpm $(packer_absdir)/docker/fpm-image

.PHONY: go-daemon-image
go-daemon-image:
	docker build --rm=true -t tudorg/go-daemon $(packer_absdir)/docker/go-daemon/

build/god-linux-386 build/god-linux-amd64:
	docker run -v $(shell pwd)/build:/build tudorg/go-daemon

build/upload:
	mkdir -p build/upload

build/upload/build_id.txt:
	echo $(BUILDID) > build/upload/build_id.txt

.PHONY: s3-nightlies-upload
s3-nightlies-upload: all
	aws s3 cp --recursive --acl public-read build/upload s3://beats-nightlies

# Build the image required for package-upload.
.PHONY: deb-rpm-s3
deb-rpm-s3:
	docker/deb-rpm-s3/build.sh

# Run after building to sign packages and publish to APT and YUM repos.
.PHONY: package-upload
package-upload:
	# You must export AWS_ACCESS_KEY=<AWS access> and export AWS_SECRET_KEY=<secret>
	# before running this make target.
	docker/deb-rpm-s3/deb-rpm-s3.sh

.PHONY: release-upload
release-upload:
	aws s3 cp --recursive --acl public-read build/upload s3://download.elasticsearch.org/beats/

.PHONY: run-interactive
run-interactive:
	docker run -t -i -v $(shell pwd)/build:/build \
		-v $(shell pwd)/xgo-scripts/:/scripts \
		--entrypoint=bash tudorg/beats-builder-deb6

.PHONY: images
images: xgo-image fpm-image go-daemon-image

.PHONY: push-images
push-images:
	docker push tudorg/beats-builder
	docker push tudorg/beats-builder-deb6
	docker push tudorg/fpm
	docker push tudorg/go-daemon

.PHONY: pull-images
pull-images:
	docker pull tudorg/beats-builder
	docker pull tudorg/beats-builder-deb6
	docker pull tudorg/fpm
	docker pull tudorg/go-daemon


define rm-image =
@echo "Cleaning $(1) image..."
  @if [ $(shell docker ps  -n 1 -a -q  --filter="image=$(1)" )  ]; then \
      docker stop $(shell docker ps -a -q --filter="image=$(1)"); \
      docker rm   $(shell docker ps -a -q --filter="image=$(1)"); \
  fi; \
\
  if [ $(shell docker images -q $(1)) ]; then \
    docker rmi $(1); \
  fi
endef


.PHONY: clean-images
clean-images:
	@$(call rm-image,tudorg/beats-builder-deb6)
	@$(call rm-image,tudorg/beats-builder)

.PHONY: clean
clean:
	rm -rf build/ || true
	$(call rm-image,build-image)

# Creates a latest file for the most recent build
.PHONY: latest
latest:
	BUILDID=${BUILDID} \
	$(packer_absdir)/xgo-scripts/latest.sh
