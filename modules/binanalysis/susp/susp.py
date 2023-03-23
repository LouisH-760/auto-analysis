import sqlite3
import subprocess  # run commands
import sys  # get arguments
from base64 import b64decode  # decode arguments
import json  # format output

# Give an overview over what part of an exe might make it suspicious. Those indicators alone might be benign, but when combined can be telling.

if len(sys.argv) != 2:
    print("Wrong / Malformed arguments", file=sys.stderr)
    sys.exit(1)
else:
    out = {
        "sections": [],
        "imports": [],
        "libraries": [],
        "other": []
    }
    sample = b64decode(sys.argv[-1].encode('ascii')).decode('ascii')
    # GATHER DATA
    # Get linked libraries information in JSON format from rz-bin
    cmd = ["rz-bin", "-lj", sample]
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    libraries = json.loads(output.stdout.decode())["libs"]
    # Get imports information in a JSON format from rz-bin
    cmd = ["rz-bin", "-ij", sample]
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    imports = json.loads(output.stdout.decode())["imports"]
    # Get sections information in a JSON format from rz-bin
    cmd = ["rz-bin", "-Sj", sample]
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    sections = json.loads(output.stdout.decode())["sections"]
    # CHECK SECTIONS FOR RWX
    for section in sections: # more useful for exes out of a memdump
        perms = set([char for char in section["perm"]]) # hacky, get a set of all permissions, why can't you give me a number rizin???
        if 'r' in perms and 'w' in perms and 'x' in perms:
            out["sections"].append({
                "name": section["name"],
                "message": "This section has rwx perms!",
                "raw": section
            })
    with sqlite3.connect("/autoa/modules/binanalysis/susp/susp.db") as dbconn:
        cursor = dbconn.cursor()
        for im in imports: # check imports
            mc = cursor.execute("SELECT * FROM import WHERE LOWER( import.name ) = (?)", [im["name"].lower()])
            matches = mc.fetchall()
            if len(matches) > 0:
                out["imports"].append({
                    "name": im["name"],
                    "message": f"Suspicious import found: {matches[0]}",
                    "raw": im
                })
        for lib in libraries: # check linked libraries
            ml = cursor.execute("SELECT * FROM library WHERE LOWER( library.name ) = (?)", [lib.lower()])
            matches = ml.fetchall()
            if len(matches) > 0:
                out["libraries"].append({
                    "name": lib,
                    "message": f"Suspicious linked library found: {matches[0]}",
                    "raw": lib
                })
        mn = cursor.execute("SELECT * FROM process WHERE LOWER( process.name ) = (?)", [sample.lower()])
        matches = mn.fetchall()
        if len(matches) > 0:
            out["other"].append({
                "name": sample,
                "message": f"The sample name matches with a windows system process, this might be an attempt at impersonation!: {matches[0]}",
                "raw": matches
            })
    cmd = ["diec", "-j", sample] # use diec to detect a packer
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    buildchain = json.loads(output.stdout.decode())["detects"][0]
    for val in buildchain["values"]:
        if val["type"] == "Packer":
            out["other"].append({
                "name": sample,
                "message": f"The sample binary is packed using {val['name']}"
            })
    cmd = ["diec", "-je", sample] # use diec to get entropy
    output = subprocess.run(cmd, stdout=subprocess.PIPE)
    sectentropy = json.loads(output.stdout.decode())
    if sectentropy["total"] >= 6.5:
        out["other"].append({
            "name": sample,
            "message": f"High overall entropy: might be packed ({sectentropy['total']})"
        })
    for section in sectentropy["records"]:
        if section["entropy"] > 6.5:
            out["other"].append({
                "name": section["name"],
                "message": f"Section has unusually high entropy, might be packed: {section['entropy']}",
                "raw": section
            })
    print(json.dumps(out))