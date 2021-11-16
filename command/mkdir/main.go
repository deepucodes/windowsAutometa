package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// PowerShell struct
type PowerShell struct {
	powerShell string
}

// New create new session
func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

var (
	// Below command will enable the HyperV module
	enableHyperVCmd    = `Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All`
	elevateProcessCmds = `
	$myWindowsID=[System.Security.Principal.WindowsIdentity]::GetCurrent()
	$myWindowsPrincipal=new-object System.Security.Principal.WindowsPrincipal($myWindowsID)
 
	# Get the security principal for the Administrator role
	$adminRole=[System.Security.Principal.WindowsBuiltInRole]::Administrator

	# Create a new process object that starts PowerShell
	$newProcess = new-object System.Diagnostics.ProcessStartInfo "PowerShell";
	# Specify the current script path and name as a parameter
	$newProcess.Arguments = $MyInvocation.MyCommand.Definition.Path;
	
	# Write-Host -NoNewLine $script:MyInvocation.MyCommand.Definition.Path

	# Indicate that the process should be elevated
	$newProcess.Verb = "runas";
	
	# Start the new process
	$process = [System.Diagnostics.Process]::Start($newProcess);
	
	# Exit from the current, unelevated, process
	exit	
`
)

func main() {
	posh := New()

	// Scenario 1
	// stdOut, stdErr, err := posh.execute(elevateProcessCmds)
	// fmt.Printf("ElevateProcessCmds:\nStdOut : '%s'\nStdErr: '%s'\nErr: %s", strings.TrimSpace(stdOut), stdErr, err)
	// ========= Above working and invoke a publisher permission dialog and Admin shell ================

	// Scenario 2
	// stdOut, stdErr, err := posh.execute(enableHyperVCmd)
	// fmt.Printf("\nEnableHyperV:\nStdOut : '%s'\nStdErr: '%s'\nErr: %s", strings.TrimSpace(stdOut), stdErr, err)
	// ========= Behavior(expected one): StdErr: 'Enable-WindowsOptionalFeature : The requested operation requires elevation.

	// Scenario 3 : Both scenario 1 and 2 combined
	enableHyperVScript := fmt.Sprintf("%s\n%s", elevateProcessCmds, enableHyperVCmd)
	stdOut, stdErr, err := posh.execute(enableHyperVScript)
	fmt.Printf("\nEnableHyperV:\nStdOut : '%s'\nStdErr: '%s'\nErr: %s", strings.TrimSpace(stdOut), stdErr, err)
	// ========= Above suppose to open a permission dialog, on click of "yes" should
	// ========= run the hyperv enable command and once done ask for restart operation
	// ========= Actual Behavior: Only invoking the Powershell in admin mode and not running the HyperV Enable command.

	stdOut, stdErr, err = posh.execute(`New-Item -Path . -Name "brightdrop" -ItemType "folder"`)
	fmt.Printf("\nEnableHyperV:\nStdOut : '%s'\nStdErr: '%s'\nErr: %s", strings.TrimSpace(stdOut), stdErr, err)

}
