package setup

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidatePamUsername(t *testing.T) {
	if err := validatePamUsername("alice"); err != nil {
		t.Fatal(err)
	}
	if err := validatePamUsername("Alice"); err == nil {
		t.Fatal("expected invalid username")
	}
	if err := validatePamUsername(""); err == nil {
		t.Fatal("expected required username")
	}
}

func TestPamHasExistingUsersFromList(t *testing.T) {
	orig := runHost
	t.Cleanup(func() { runHost = orig })
	runHost = func(cmd hostCmd) (string, error) {
		if cmd.Name == "getent" && len(cmd.Args) >= 2 && cmd.Args[0] == "group" {
			return "equate-appliance:x:900:admin\n", nil
		}
		if cmd.Name == "passwd" && len(cmd.Args) >= 2 && cmd.Args[0] == "-S" {
			return cmd.Args[1] + " P 2026-01-01 0 99999 7 -1\n", nil
		}
		return "", fmt.Errorf("unexpected %#v", cmd)
	}
	has, body, err := pamHasExistingUsers()
	if err != nil {
		t.Fatal(err)
	}
	if !has || !strings.Contains(body, "admin") {
		t.Fatalf("has=%v body=%q", has, body)
	}
}

func TestPamHasExistingUsersEmpty(t *testing.T) {
	orig := runHost
	t.Cleanup(func() { runHost = orig })
	runHost = func(cmd hostCmd) (string, error) {
		if cmd.Name == "getent" {
			return "equate-appliance:x:900:\n", nil
		}
		return "", fmt.Errorf("unexpected %#v", cmd)
	}
	has, body, err := pamHasExistingUsers()
	if err != nil {
		t.Fatal(err)
	}
	if has || body != "" {
		t.Fatalf("has=%v body=%q", has, body)
	}
}

func TestPamUserCreateInvokesUseradd(t *testing.T) {
	orig := runHost
	t.Cleanup(func() { runHost = orig })
	var names []string
	runHost = func(cmd hostCmd) (string, error) {
		names = append(names, cmd.Name)
		switch cmd.Name {
		case "getent":
			return "", fmt.Errorf("not found")
		case "groupadd", "useradd", "chpasswd":
			return "", nil
		case "id":
			return "", fmt.Errorf("no such user")
		default:
			return "", fmt.Errorf("unexpected %s", cmd.Name)
		}
	}
	if err := pamUserCreate("alice", "secret"); err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(names, ",")
	if !strings.Contains(joined, "useradd") || !strings.Contains(joined, "chpasswd") {
		t.Fatalf("commands=%v", names)
	}
}

func TestPamUserListFormatsStatus(t *testing.T) {
	orig := runHost
	t.Cleanup(func() { runHost = orig })
	runHost = func(cmd hostCmd) (string, error) {
		if cmd.Name == "getent" {
			return "equate-appliance:x:900:alice,bob\n", nil
		}
		if cmd.Name == "passwd" && cmd.Args[0] == "-S" {
			if cmd.Args[1] == "bob" {
				return "bob L 2026-01-01 0 99999 7 -1\n", nil
			}
			return "alice P 2026-01-01 0 99999 7 -1\n", nil
		}
		return "", fmt.Errorf("unexpected %#v", cmd)
	}
	out, err := pamUserList()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "alice (enabled)") || !strings.Contains(out, "bob (disabled)") {
		t.Fatalf("output=%q", out)
	}
}
