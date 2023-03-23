import sys  # get arguments
import json # format output

# Small debugging modules, just passes sys.argv back to you.

res = {
    "argc": len(sys.argv),
    "argv": sys.argv
}
print(json.dumps(res), file=sys.stdout)
