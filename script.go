import subprocess
from concurrent.futures import ThreadPoolExecutor

// # Commands to run in two different terminals

command1 := "python ping.py"
command2 := "go run ./cmd/ndt7-client/main.go"

// # Function to run a command in a new terminal
func run_in_terminal(command):
	subprocess.run(["gnome-terminal", "-e", command])

// # Run the commands simultaneously using ThreadPoolExecutor
with ThreadPoolExecutor(max_workers=2) as executor:
	future1 = executor.submit(run_in_terminal, command1)
	future2 = executor.submit(run_in_terminal, command2)


// # Wait for both commands to complete (optional)
future1.result()
future2.result()