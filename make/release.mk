# Production release, OVA, and update-channel targets. Invoked from the repository root Makefile.
.PHONY: appliance-bundle-all appliance-package appliance-ova-amd64-ci \
	appliance-prepare-ova appliance-publish-azure appliance-generate-keys \
	appliance-setup-azure-channel appliance-package-ova appliance-export-ova-arm64 \
	appliance-verify-ova

appliance-bundle-all:
	@$(MAKE) appliance-bundle ARCH=arm64 VERSION="$(VERSION)"
	@$(MAKE) appliance-bundle ARCH=amd64 VERSION="$(VERSION)"

appliance-package:
	@./$(RELEASE_SCRIPTS)/package-eqa.sh --arch "$(ARCH)" --version "$(VERSION)"

appliance-ova-amd64-ci:
	@./$(RELEASE_SCRIPTS)/build-ova-amd64-ci.sh --version "$(VERSION)"

appliance-publish-azure:
	@test -n "$(STORAGE_ACCOUNT)" || (echo "STORAGE_ACCOUNT is required" >&2; exit 1)
	@./$(RELEASE_SCRIPTS)/publish-update-channel-azure.sh \
		--storage-account "$(STORAGE_ACCOUNT)" \
		--container "$(or $(CONTAINER),updates)" \
		--channel "$(or $(CHANNEL),stable)" \
		--edition standard \
		--arch "$(ARCH)" \
		--version "$(VERSION)" \
		$(if $(CDN_BASE),--cdn-base "$(CDN_BASE)",)

appliance-prepare-ova:
	@sudo ./$(RELEASE_SCRIPTS)/prepare-ova.sh

appliance-generate-keys:
	@./$(RELEASE_SCRIPTS)/generate-update-keys.sh $(if $(KEYS_DIR),--out-dir "$(KEYS_DIR)",)

appliance-setup-azure-channel:
	@test -n "$(STORAGE_ACCOUNT)" || (echo "STORAGE_ACCOUNT is required" >&2; exit 1)
	@./$(RELEASE_SCRIPTS)/setup-update-channel-azure.sh --storage-account "$(STORAGE_ACCOUNT)"

appliance-package-ova:
	@test -n "$(VMDK)" || (echo "VMDK is required" >&2; exit 1)
	@./$(RELEASE_SCRIPTS)/package-ova.sh --arch "$(ARCH)" --version "$(VERSION)" --vmdk "$(VMDK)"

appliance-export-ova-arm64:
	@test -n "$(VMX)" || (echo "VMX is required" >&2; exit 1)
	@./$(RELEASE_SCRIPTS)/export-arm64-ova.sh --vmx "$(VMX)" --version "$(VERSION)"

appliance-verify-ova:
	@sudo ./$(RELEASE_SCRIPTS)/verify-ova-import.sh \
		$(if $(ARTIFACT),--artifact "$(ARTIFACT)",) \
		$(if $(CONFIGURED),--configured,)
