import subprocess  # run commands
import sys  # get arguments
from base64 import b64decode  # decode arguments
import json  # format output

if len(sys.argv) != 2:
    print("Wrong / Malformed arguments", file=sys.stderr)
    sys.exit(1)
else:
    memorydump = b64decode(sys.argv[-1].encode('ascii')).decode('ascii')
    # Get image information in a JSON format from volatility
    cmd = ["volatility2", "-f", "--output=json", memorydump]
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    imginfo = json.loads(output.stdout.decode())
    print(json.dumps(imginfo), file=sys.stdout)
