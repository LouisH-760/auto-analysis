import subprocess  # run commands
import sys  # get arguments
from base64 import b64decode  # decode arguments
import json  # format output

if len(sys.argv) != 3:
    print("Wrong / Malformed arguments", file=sys.stderr)
    sys.exit(1)
else:
    memorydump = b64decode(sys.argv[-2].encode('ascii')).decode('ascii')
    profile = b64decode(sys.argv[-1].encode('ascii')).decode('ascii')
    # Get process tree in a JSON format from volatility
    cmd = ["volatility2", "-f", memorydump, "--output=json", f"--profile={profile}", "pstree"]
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    imginfo = json.loads(output.stdout.decode())
    print(json.dumps(imginfo), file=sys.stdout)
