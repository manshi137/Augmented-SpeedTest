import subprocess
from concurrent.futures import ThreadPoolExecutor

# Commands to run in two different terminals
command1 = "python ping.py"
command2 = "ndt7-client"

# Function to run a command in a new terminal
def run_in_terminal(command):
    subprocess.Popen(['start', 'cmd', '/k', command], shell=True)
    subprocess.Popen(['gnome-terminal', '--', 'bash', '-c', command])

# Run the commands simultaneously using ThreadPoolExecutor
with ThreadPoolExecutor(max_workers=2) as executor:
    future1 = executor.submit(run_in_terminal, command1)
    future2 = executor.submit(run_in_terminal, command2)

# Wait for both commands to complete (optional)
# You can remove this part if you don't want to wait for the commands to finish.
future1.result()
future2.result()
