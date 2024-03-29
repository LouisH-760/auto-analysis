import subprocess  # run commands
import sys  # get arguments
from base64 import b64decode  # decode arguments
import json  # format output

if len(sys.argv) != 2:
    print("Wrong / Malformed arguments", file=sys.stderr)
    sys.exit(1)
else:
    sample = b64decode(sys.argv[-1].encode('ascii')).decode('ascii')
    # Get imports information in a JSON format from rz-bin
    cmd = ["rz-bin", "-ij", sample]
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    imports = json.loads(output.stdout.decode())["imports"]
    print(json.dumps(imports), file=sys.stdout)
