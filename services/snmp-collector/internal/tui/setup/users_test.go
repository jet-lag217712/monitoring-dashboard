package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPamUserCreateRejectsInvalidUsername(t *testing.T) {
	if err := pamUserCreate("Alice", "secret"); err == nil {
		t.Fatal("expected invalid username")
	}
}

func TestPamUserCreateWithStubHost(t *testing.T) {
	orig := runHost
	t.Cleanup(func() { runHost = orig })
	runHost = func(cmd hostCmd) (string, error) {
		switch cmd.Name {
		case "getent":
			return "", os.ErrNotExist
		case "groupadd", "useradd", "chpasswd":
			return "", nil
		case "id":
			return "", os.ErrNotExist
		default:
			return "", os.ErrNotExist
		}
	}
	if err := pamUserCreate("alice", "secret"); err != nil {
		t.Fatal(err)
	}
}

func TestApplianceUsersParsesGetent(t *testing.T) {
	orig := runHost
	t.Cleanup(func() { runHost = orig })
	runHost = func(cmd hostCmd) (string, error) {
		if cmd.Name == "getent" {
			return "equate-appliance:x:900:alice\n", nil
		}
		if cmd.Name == "passwd" {
			return "alice P\n", nil
		}
		return "", os.ErrNotExist
	}
	out, err := pamUserList()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "alice") {
		t.Fatalf("output=%q", out)
	}
}

func TestSiteArtifactDirStillCreated(t *testing.T) {
	deployDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(deployDir, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
}
