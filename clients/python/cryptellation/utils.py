import subprocess
import json

def cryptellation(args):
    res = subprocess.check_output(['cryptellation', '-j']+args).decode('utf-8')
    return json.loads(res)