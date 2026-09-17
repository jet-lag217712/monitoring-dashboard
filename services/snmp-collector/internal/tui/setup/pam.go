package setup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const defaultApplianceGroup = "equate-appliance"

var pamUsernamePattern = regexp.MustCompile(`^[a-z][-a-z0-9_]*$`)

type hostCmd struct {
	Name  string
	Args  []string
	Stdin string
}

type hostRunner func(cmd hostCmd) (string, error)

var runHost hostRunner = defaultRunHost

func defaultRunHost(cmd hostCmd) (string, error) {
	c := exec.Command(cmd.Name, cmd.Args...)
	if cmd.Stdin != "" {
		c.Stdin = strings.NewReader(cmd.Stdin)
	}
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}
		if msg == "" {
			msg = err.Error()
		}
		return stdout.String(), fmt.Errorf("%s", msg)
	}
	return stdout.String(), nil
}

func applianceGroup() string {
	if v := strings.TrimSpace(os.Getenv("EQUATE_APPLIANCE_GROUP")); v != "" {
		return v
	}
	return defaultApplianceGroup
}

func validatePamUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) > 32 || !pamUsernamePattern.MatchString(username) {
		return fmt.Errorf("invalid username")
	}
	return nil
}

func pamUserCreate(username, password string) error {
	if err := validatePamUsername(username); err != nil {
		return err
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if err := ensureApplianceGroup(); err != nil {
		return err
	}
	if userExists(username) {
		return fmt.Errorf("user already exists")
	}
	if _, err := runHost(hostCmd{Name: "useradd", Args: []string{"-m", "-s", "/bin/bash", "-G", applianceGroup(), username}}); err != nil {
		return fmt.Errorf("pam helper create: %s", err)
	}
	if err := setPassword(username, password); err != nil {
		return fmt.Errorf("pam helper create: %s", err)
	}
	return nil
}

func pamUserList() (string, error) {
	if err := ensureApplianceGroup(); err != nil {
		return "", err
	}
	users, err := applianceUsers()
	if err != nil {
		return "", fmt.Errorf("pam helper list: %s", err)
	}
	if len(users) == 0 {
		return "No appliance users listed.\n", nil
	}
	var b strings.Builder
	for _, user := range users {
		status := "enabled"
		if userLocked(user) {
			status = "disabled"
		}
		fmt.Fprintf(&b, "%s (%s)\n", user, status)
	}
	return b.String(), nil
}

func pamHasExistingUsers() (bool, string, error) {
	body, err := pamUserList()
	if err != nil {
		return false, "", err
	}
	body = strings.TrimSpace(body)
	if body == "" || strings.Contains(body, "No appliance users") {
		return false, "", nil
	}
	return true, body, nil
}

func pamUserDisable(username string) error {
	if err := requireApplianceUser(username); err != nil {
		return err
	}
	if _, err := runHost(hostCmd{Name: "passwd", Args: []string{"-l", username}}); err != nil {
		return fmt.Errorf("pam helper disable: %s", err)
	}
	return nil
}

func pamUserEnable(username string) error {
	if err := requireApplianceUser(username); err != nil {
		return err
	}
	if _, err := runHost(hostCmd{Name: "passwd", Args: []string{"-u", username}}); err != nil {
		return fmt.Errorf("pam helper enable: %s", err)
	}
	return nil
}

func pamUserReset(username, password string) error {
	if err := requireApplianceUser(username); err != nil {
		return err
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if err := setPassword(username, password); err != nil {
		return fmt.Errorf("pam helper reset-password: %s", err)
	}
	return nil
}

func pamUserDelete(username string) error {
	if err := requireApplianceUser(username); err != nil {
		return err
	}
	if !userLocked(username) {
		enabled, err := countEnabledApplianceUsers()
		if err != nil {
			return fmt.Errorf("pam helper delete: %s", err)
		}
		if enabled <= 1 {
			return fmt.Errorf("cannot delete the last enabled appliance user")
		}
	}
	if _, err := runHost(hostCmd{Name: "userdel", Args: []string{"-r", username}}); err != nil {
		if _, err = runHost(hostCmd{Name: "userdel", Args: []string{username}}); err != nil {
			return fmt.Errorf("pam helper delete: %s", err)
		}
	}
	return nil
}

func ensureApplianceGroup() error {
	group := applianceGroup()
	if _, err := runHost(hostCmd{Name: "getent", Args: []string{"group", group}}); err == nil {
		return nil
	}
	if _, err := runHost(hostCmd{Name: "groupadd", Args: []string{"--system", group}}); err != nil {
		return fmt.Errorf("ensure group %s: %s", group, err)
	}
	return nil
}

func userExists(username string) bool {
	_, err := runHost(hostCmd{Name: "id", Args: []string{username}})
	return err == nil
}

func inApplianceGroup(username string) bool {
	out, err := runHost(hostCmd{Name: "id", Args: []string{"-nG", username}})
	if err != nil {
		return false
	}
	group := applianceGroup()
	for _, g := range strings.Fields(out) {
		if g == group {
			return true
		}
	}
	return false
}

func requireApplianceUser(username string) error {
	if err := validatePamUsername(username); err != nil {
		return err
	}
	if !userExists(username) {
		return fmt.Errorf("user not found")
	}
	if !inApplianceGroup(username) {
		return fmt.Errorf("user is not an appliance account")
	}
	return nil
}

func setPassword(username, password string) error {
	_, err := runHost(hostCmd{
		Name:  "chpasswd",
		Stdin: username + ":" + password + "\n",
	})
	return err
}

func userLocked(username string) bool {
	out, err := runHost(hostCmd{Name: "passwd", Args: []string{"-S", username}})
	if err != nil {
		return false
	}
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return false
	}
	status := fields[1]
	return status == "L" || status == "LK"
}

func applianceUsers() ([]string, error) {
	out, err := runHost(hostCmd{Name: "getent", Args: []string{"group", applianceGroup()}})
	if err != nil {
		return nil, err
	}
	parts := strings.Split(strings.TrimSpace(out), ":")
	if len(parts) < 4 {
		return nil, nil
	}
	raw := strings.TrimSpace(parts[3])
	if raw == "" {
		return nil, nil
	}
	var users []string
	for _, user := range strings.Split(raw, ",") {
		user = strings.TrimSpace(user)
		if user != "" {
			users = append(users, user)
		}
	}
	return users, nil
}

func countEnabledApplianceUsers() (int, error) {
	users, err := applianceUsers()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, user := range users {
		if !userLocked(user) {
			count++
		}
	}
	return count, nil
}

// CreateApplianceUser creates a PAM-backed appliance operator.
func CreateApplianceUser(username, password string) error {
	return pamUserCreate(username, password)
}

// ListApplianceUsers returns the formatted appliance user list.
func ListApplianceUsers() (string, error) {
	return pamUserList()
}

// DisableApplianceUser locks a PAM-backed appliance operator.
func DisableApplianceUser(username string) error {
	return pamUserDisable(username)
}

// EnableApplianceUser unlocks a PAM-backed appliance operator.
func EnableApplianceUser(username string) error {
	return pamUserEnable(username)
}

// ResetApplianceUserPassword resets a PAM-backed appliance operator password.
func ResetApplianceUserPassword(username, password string) error {
	return pamUserReset(username, password)
}

// DeleteApplianceUser removes a PAM-backed appliance operator.
func DeleteApplianceUser(username string) error {
	return pamUserDelete(username)
}
