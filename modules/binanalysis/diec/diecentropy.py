import subprocess  # run commands
import sys  # get arguments
from base64 import b64decode  # decode arguments
import json  # format output

if len(sys.argv) != 2:
    print("Wrong / Malformed arguments", file=sys.stderr)
    sys.exit(1)
else:
    sample = b64decode(sys.argv[-1].encode('ascii')).decode('ascii')
    # Get binary information in a JSON format from rz-bin
    cmd = ["diec", "-je", sample]
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    infobj = json.loads(output.stdout.decode())
    print(json.dumps(infobj), file=sys.stdout)
