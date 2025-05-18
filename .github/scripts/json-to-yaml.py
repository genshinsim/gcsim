
import os
import json
import yaml

def json_to_yaml(subdir, file):
    obj = None

    json_file = os.path.join(subdir, file)
    with open(json_file) as f:
        obj = json.load(f)

    yaml_file = os.path.join(subdir, "metadata.yaml")
    with open(yaml_file, "w") as f:
        yaml.dump(obj, f)

    os.remove(json_file)


if __name__ == "__main__":

    for subdir, dirs, files in os.walk("./containers"):
        for f in files:
            if f != "metadata.json":
                continue
            json_to_yaml(subdir, f)