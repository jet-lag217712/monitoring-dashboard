# Develop, lab, and engineer helper targets. Invoked from the repository root Makefile.
.PHONY: test appliance-bundle appliance-stage appliance-configure appliance-upgrade \
	appliance-verify lab-validate lab-appliance-setup lab-appliance-upgrade \
	lab-gns3-bridge lab-smoke lab-mqtt-outage-drill \
	db-migrate db-bootstrap-roles mqtt-dev-certs

test:
	@./deployments/test.sh $(if $(QUICK),--quick,)

appliance-bundle:
	@./$(RELEASE_SCRIPTS)/build-release.sh --arch "$(ARCH)" --version "$(VERSION)"

appliance-stage:
	@test -n "$(HOST)" || (echo "HOST is required" >&2; exit 1)
	@./$(RELEASE_SCRIPTS)/stage-release.sh --host "$(HOST)" --user "$(USER)" --arch "$(ARCH)" --version "$(VERSION)"

appliance-configure:
	@test -n "$(BUNDLE)" || (echo "BUNDLE is required" >&2; exit 1)
	@sudo ./$(RUNTIME_SCRIPTS)/configure-vm.sh --bundle "$(BUNDLE)" --version "$(VERSION)"

appliance-upgrade:
	@test -n "$(BUNDLE)" || (echo "BUNDLE is required" >&2; exit 1)
	@sudo ./$(RUNTIME_SCRIPTS)/configure-vm.sh --upgrade --bundle "$(BUNDLE)" --version "$(VERSION)" $(if $(CANARY),--canary,)

appliance-verify:
	@sudo ./$(RELEASE_SCRIPTS)/verify-appliance.sh

lab-validate:
	@./remote-server/validate-lab.sh

lab-appliance-setup:
	@sudo ./remote-server/lab-appliance-setup.sh

lab-appliance-upgrade:
	@test -n "$(OLD_VERSION)" || (echo "OLD_VERSION is required" >&2; exit 1)
	@test -n "$(NEW_VERSION)" || (echo "NEW_VERSION is required" >&2; exit 1)
	@test -n "$(BUNDLE)" || (echo "BUNDLE is required" >&2; exit 1)
	@sudo OLD_VERSION="$(OLD_VERSION)" NEW_VERSION="$(NEW_VERSION)" BUNDLE="$(BUNDLE)" \
		./remote-server/lab-appliance-upgrade.sh

lab-gns3-bridge:
	@sudo ./remote-server/setup-gns3-bridge.sh

lab-smoke:
	@./remote-server/smoke_mqtt_v2_to_api.sh

lab-mqtt-outage-drill:
	@sudo ./remote-server/mqtt_outage_drill.sh

db-migrate:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required" >&2; exit 1)
	@./infrastructure/script/migrate.sh up

db-bootstrap-roles:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required" >&2; exit 1)
	@./infrastructure/script/bootstrap-db-roles.sh

mqtt-dev-certs:
	@./infrastructure/docker/mqtt-broker/scripts/gen-dev-certs.sh
	@./infrastructure/docker/mqtt-broker/scripts/gen-passwords.sh
