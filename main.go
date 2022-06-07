package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Users struct {
	userAcc      []string
	newAdminName string
	newAdminPass string
}

func main() {
	users := Users{}
	// process flags
	users.processUserInput()

	// get local users
	if err := users.getUsers(); err != nil {
		log.Fatalln("could not get users list err=", err.Error())
	}

	// delete all local user accounts
	users.deleteAllLocalUsers()

	// create new admin user
	if err := users.createAdminUser(); err != nil {
		log.Fatalln("could not create new admin user err=", err.Error())
	}

	// disable builtin Administrator account
	if err := users.disableBuiltinAdmin(); err != nil {
		log.Fatalln("could not disable built-in administrator account err=", err.Error())
	}

	os.Exit(0)
}

func (u *Users) getUsers() error {
	// run win cmd to get users list
	out, err := runCmd("wmic USERACCOUNT WHERE LocalAccount=True GET Name")
	if err != nil {
		return fmt.Errorf(fmt.Sprint("could not run wmic command: ", out))
	}

	// remove system accounts and store users list
	for _, user := range strings.Split(strings.TrimSpace(out), " ")[1:] {
		if !strings.Contains(user, "Administrator") &&
			!strings.Contains(user, "ASPNET") &&
			!strings.Contains(user, "DefaultAccount") &&
			!strings.Contains(user, "Guest") &&
			!strings.Contains(user, "WDAGUtilityAccount") &&
			user != "" {
			u.userAcc = append(u.userAcc, strings.TrimSpace(user))
		}
	}

	return nil
}

func (u *Users) createAdminUser() error {
	// run net user command to create new user
	out, err := runCmd(fmt.Sprintf("net user %s %s /ADD", u.newAdminName, u.newAdminPass))
	if err != nil {
		return fmt.Errorf("could not add new user: %s", out)
	}

	// and add this user to a local administrators group
	out, err = runCmd(fmt.Sprintf("net localgroup Administrators %s /ADD", u.newAdminName))
	if err != nil {
		return fmt.Errorf("could not add new admin account to administrators group: %s", out)
	}

	// set password to not expire for this admin account
	out, err = runCmd(fmt.Sprintf("WMIC USERACCOUNT WHERE Name='%s' SET PasswordExpires=FALSE", u.newAdminName))
	if err != nil {
		return fmt.Errorf("could not set non-expiring password for admin account: %s", out)
	}

	return nil
}

func (u *Users) processUserInput() {
	flag.StringVar(&u.newAdminPass, "admin-pass", "P@ssw0rd", "the password to be set for new user account")
	flag.StringVar(&u.newAdminName, "admin-name", "SecureAdmin", "the username for new admin account")
	flag.Parse()
}

func (u *Users) disableBuiltinAdmin() error {
	// disable local admin account
	out, err := runCmd("net user Administrator /ACTIVE:NO")
	if err != nil {
		return fmt.Errorf("could not disable built-in administrator account: %s", out)
	}

	return nil
}

func (u *Users) deleteAllLocalUsers() {
	// delete all local user accounts
	for _, user := range u.userAcc {
		out, err := runCmd(fmt.Sprintf("net user %s /DELETE", user))
		if err != nil {
			log.Println(fmt.Sprintf("could not delete local account named %s: %s", user, out))
		}
	}
}

func runCmd(cmd string) (string, error) {
	args := strings.Split("/c "+cmd, " ")
	out, err := exec.Command("cmd.exe", args...).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
