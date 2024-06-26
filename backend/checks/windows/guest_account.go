package windows

import (
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/checks"

	"github.com/InfoSec-Agent/InfoSec-Agent/backend/mocking"
)

// GuestAccount checks the status of the Windows guest account.
//
// Parameters:
//   - executorLocalGroup (commandmock.CommandExecutor): An instance of CommandExecutor used to execute the Get-WmiObject command to retrieve local group information.
//   - executorLocalGroupMembers (commandmock.CommandExecutor): An instance of CommandExecutor used to execute the 'net localgroup' command to retrieve local group members.
//   - executorYesWord (commandmock.CommandExecutor): An instance of CommandExecutor used to execute the 'net user' command to retrieve the word for 'yes' in the current user's language.
//   - executorNetUser (commandmock.CommandExecutor): An instance of CommandExecutor used to execute the 'net user' command to retrieve all users.
//   - usernameRetriever (mocking.UsernameRetriever): An instance of UsernameRetriever used to retrieve the current username.
//
// Returns:
//   - Check: A Check instance encapsulating the results of the guest account check. If the guest account is active, the Result field of the Check instance will contain the message "Guest account is active". If the guest account is not active, the Result field will contain the message "Guest account is not active". If an error occurs during the check, it is encapsulated in the Error and ErrorMSG fields of the Check instance.
//
// This function is primarily used to identify potential security risks associated with an active guest account on the Windows system.
func GuestAccount(
	executorLocalGroup mocking.CommandExecutor,
	executorLocalGroupMembers mocking.CommandExecutor,
	executorYesWord mocking.CommandExecutor,
	executorNetUser mocking.CommandExecutor,
	usernameRetriever mocking.UsernameRetriever,
) checks.Check {
	// Get localgroup name using GetWmiObject
	// output, err := GuestAccountLocalGroup(executorLocalGroup)

	command := "Get-WmiObject Win32_Group | Select-Object SID,Name"
	output, err := executorLocalGroup.Execute("powershell", command)
	if err != nil {
		return checks.NewCheckErrorf(checks.GuestAccountID, "error executing command Get-WmiObject", err)
	}
	outputString := strings.Split(string(output), "\r\n")
	found := false
	guestGroup := ""
	for _, line := range outputString {
		// Check for the guest account SID
		if strings.Contains(line, "S-1-5-32-546") {
			line = line[13 : len(line)-1]
			line = strings.TrimSpace(line)
			found = true
			guestGroup = line
		}
	}
	if !found {
		return checks.NewCheckResult(checks.GuestAccountID, 0)
	}

	// Get local group members using net localgroup command
	output, err = executorLocalGroupMembers.Execute("net", "localgroup", guestGroup)
	if err != nil {
		return checks.NewCheckErrorf(checks.GuestAccountID, "error executing command net localgroup", err)
	}
	outputString = strings.Split(string(output), "\r\n")
	guestUser := ""
	for i := range outputString {
		// Find the line containing the guest account
		if strings.Contains(outputString[i], "-----") {
			guestUser = outputString[i+1]
		}
	}
	if guestUser == "" {
		return checks.NewCheckResult(checks.GuestAccountID, 0)
	}

	// Retrieve current username
	currentUser, err := usernameRetriever.CurrentUsername()
	if err != nil {
		return checks.NewCheckErrorf(checks.GuestAccountID, "error retrieving current username", err)
	}

	// Retrieve the word for 'yes' from the currentUser language
	output, err = executorYesWord.Execute("net", "user", currentUser)
	if err != nil {
		return checks.NewCheckErrorf(checks.GuestAccountID, "error executing command net user", err)
	}
	outputString = strings.Split(string(output), "\r\n")
	line := strings.Split(outputString[5], " ")
	yesWord := line[len(line)-1]

	// Get all users using net user command
	output, err = executorNetUser.Execute("net", "user", guestUser)
	if err != nil {
		return checks.NewCheckErrorf(checks.GuestAccountID, "error executing command net user", err)
	}
	outputString = strings.Split(string(output), "\r\n")
	// Check if the guest account is active based on the presence of the word 'yes' in the user's language
	if strings.Contains(outputString[5], yesWord) {
		return checks.NewCheckResult(checks.GuestAccountID, 1)
	}
	return checks.NewCheckResult(checks.GuestAccountID, 2)
}
