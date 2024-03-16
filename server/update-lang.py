from os import path
import json
import sys

if len(sys.argv) != 2:
    print(f"usage: {sys.argv[0]} <lang>")
    exit(1)

lang = sys.argv[1]
if not lang.endswith(".json"):
    lang = lang+".json"
lang = path.join("lang", lang)

f = open("lang/en.json", "r", encoding="utf-8")
src = json.loads(f.read())
f.close()

try:
    f = open(lang, "r", encoding="utf-8")
    dst = json.loads(f.read())
    f.close()
except:
    dst = {}

for k in src.keys():
    if k in dst.keys():
        continue
    dst[k] = "[EDIT THIS] "+src[k]

dmp = json.dumps(dst, indent=2, ensure_ascii=False).encode("utf-8").decode()
f = open(lang, "w", encoding="utf-8")
f.write(dmp)
f.close()
