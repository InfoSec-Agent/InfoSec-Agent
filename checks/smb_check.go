package checks

import (
	"fmt"
	"strings"

	"github.com/InfoSec-Agent/InfoSec-Agent/commandmock"
)

// SmbCheck is a function that checks the status of SMB1 (Server Message Block) and SMB2 protocols on the system.
//
// Parameters:
//   - smb1executor commandmock.CommandExecutor: An executor to run the command for checking the status of SMB1.
//   - smb2executor commandmock.CommandExecutor: An executor to run the command for checking the status of SMB2.
//
// Returns:
//   - Check: A struct containing the results of the checks. The result indicates whether SMB1 and SMB2 protocols are enabled or not.
//
// The function works by executing the commands to check the status of SMB1 and SMB2 protocols using the provided executors. It then parses the output of the commands to determine whether the protocols are enabled or not. The function returns a Check instance containing the results of the checks.
func SmbCheck(smb1executor commandmock.CommandExecutor, smb2executor commandmock.CommandExecutor) Check {
	smb1, err := SmbEnabled("SMB1", smb1executor)

	if err != nil {
		return NewCheckError("smb", err)
	}
	smb2, err := SmbEnabled("SMB2", smb2executor)

	if err != nil {
		return NewCheckError("smb", err)
	}

	return NewCheckResult("smb", smb1, smb2)
}

// SmbEnabled is a function that determines the status of a specified SMB (Server Message Block) protocol on the system.
//
// Parameters:
//   - smb string: The SMB protocol to check. This should be either "SMB1" or "SMB2".
//   - executor commandmock.CommandExecutor: An executor to run the command for checking the status of the specified SMB protocol.
//
// Returns:
//   - string: A string indicating the status of the specified SMB protocol. The string is in the format "<SMB>: enabled" if the protocol is enabled, and "<SMB>: not enabled" if the protocol is not enabled.
//   - error: An error object that describes the error, if any occurred during the execution of the command.
//
// The function works by executing a PowerShell command to get the server configuration of the specified SMB protocol. It then parses the output of the command to determine whether the protocol is enabled or not. The function returns a string indicating the status of the protocol and an error object if an error occurred during the execution of the command.
func SmbEnabled(smb string, executor commandmock.CommandExecutor) (string, error) {
	// Get the status of the specified SMB protocol
	command := fmt.Sprintf("Get-SmbServerConfiguration | Select-Object Enable%sProtocol", smb)
	output, err := executor.Execute("powershell", command)

	if err != nil {
		return "", err
	}

	outputString := strings.Split(string(output), "\r\n")
	line := strings.TrimSpace(outputString[3])
	if line == "True" {
		return smb + ": enabled", nil
	}

	return smb + ": not enabled", nil
}
