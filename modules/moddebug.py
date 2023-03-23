import subprocess  # run commands
import sys  # get arguments
import base64
import traceback
import json
try:
    moduleandargs = base64.b64decode(sys.argv[-1].encode("ascii")).decode("ascii")
    # Get binary information in a JSON format from rz-bin
    cmd = ["python3"]
    cmd.extend(moduleandargs.split(" "))
    output = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    res = {
        "stdout": output.stdout.decode(),
        "stderr": output.stderr.decode()
    }
    print(json.dumps(res), file=sys.stdout)
except:
    print(f"could not run module: {base64.b64encode(traceback.format_exc().encode('ascii')).decode('ascii')}")