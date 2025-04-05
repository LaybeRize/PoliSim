import glob
import os
import subprocess
import re
public_folder_regex = re.compile(r"(\\welcome.\w\w.html|\\sim)$", re.MULTILINE)

# Constants:
repo_name = "layberize/polisim"

# Function get the new version
def upgrade_version(old_version_str: str, upgrade: str) -> str:
    if upgrade.upper() == "MAJOR":
        return str(int(old_version_str.split(".",1)[0]) + 1) + ".0.0"

    if upgrade.upper() == "FIX":
        split_res = old_version_str.rsplit(".", 1)
        return split_res[0]+"."+str(int(split_res[1])+1)

    split_res = old_version_str.split(".")
    return split_res[0]+"."+str(int(split_res[1])+1)+".0"


# Function to update the dockerfiles for all new/removed files in the public folder
def update_docker_files(docker_file_paths: list[str], docker_file_languages: list[str],
                        public_folder_docker_paths: list[str]):
    for pos, path in enumerate(docker_file_paths):
        with open("." + path.removeprefix(".\\resources"), "r", encoding="UTF-8") as docker_file:
            content = docker_file.readlines()

        start = 0
        end = 0

        for line_pos, line in enumerate(content):
            if line.startswith("# Public Folder"):
                start = line_pos
            if start != 0 and end == 0 and line.strip() == "":
                end = line_pos
                break

        content = (content[:start+1] +
                   ["COPY " + s + " " + s + "\n" for s in public_folder_docker_paths] +
                   ["COPY ./public/welcome." + docker_file_languages[pos].lower() + ".html ./public/welcome."
                    + docker_file_languages[pos].lower() + ".html\n"] +
                   content[end:])

        with open("." + path.removeprefix(".\\resources"), "w", encoding="UTF-8") as docker_file:
            docker_file.writelines(content)



def main_function():
    # Get the type of version upgrade
    upgrade_type = input("New version is major, minor or fix?\n>")
    if upgrade_type.upper() not in ["MAJOR", "MINOR", "FIX"]:
        print("Could not find the upgrade type. ending program")
        exit(1)

    # Define the basic information that can be gleaned from just the file system alone
    sim_dir = os.getcwd().removesuffix("resources")
    docker_files = [".\\resources\\" + s.rsplit("\\", 1)[1] for s in sorted(glob.glob(os.getcwd()+"\\Dockerfile*"))]
    public_files = ["./public/" + s.rsplit("\\", 1)[1] for s in sorted(glob.glob(sim_dir+"public\\*")) if not public_folder_regex.search(s)]
    language_abbreviations = [s.split("-")[1] for s in docker_files]
    update_docker_files(docker_files, language_abbreviations, public_files)
    commands = []
    old_version = ""
    build_version = ""

    language_link_string = ("![Supported Languages are " + ", ".join(language_abbreviations) +
                            "](https://img.shields.io/badge/languages-" + ",_".join(language_abbreviations) +
                            "-yellow)\n")

    # Open the README.md to read it
    with open("../README.md", "r", encoding="UTF-8") as readme:
        full_readme_text = readme.readlines()

    # Open log file, because now it is getting interesting
    file = open("./run.log", "w", encoding="UTF-8")

    # Update README.md and extract new version
    new_readme_text = []
    for line in full_readme_text:
        if line.startswith("![Version is"):
            old_version = line.removeprefix("![Version is ").split("]",1)[0]
            file.write("Current Version: "+old_version+"\n")
            build_version = upgrade_version(old_version, upgrade_type)
            new_readme_text.append(f"![Version is {build_version}]"
                                   f"(https://img.shields.io/badge/version-{build_version}-blue)\n")
            file.write("New Version: "+build_version+"\n")
            continue
        if line.startswith("![Supported Languages are"):
            file.write("Supported Languages: "+",".join(language_abbreviations)+"\n")
            new_readme_text.append(language_link_string)
            continue
        new_readme_text.append(line)

    # Write updated README.md back
    with open("../README.md", "w", encoding="UTF-8") as readme:
        readme.writelines(new_readme_text)

    # Force a new version load of anything that is versioned in the html header
    with open("../handler/_templates/base.gohtml", "r+", encoding="UTF-8") as base:
        data = base.read()
        base.seek(0)
        base.write(data.replace("?v="+old_version, "?v="+build_version))
        base.truncate()

    # Generate the docker commands for the new version
    for pos, docker_file in enumerate(docker_files):
        commands.append("cd "+sim_dir+"; docker build -t "+
                        repo_name+":v"+build_version+"-"+language_abbreviations[pos]+
                        " -f "+docker_file+" .")
        if language_abbreviations[pos] == "EN":
            commands.append("cd "+sim_dir+"; docker build -t "+
                            repo_name+
                            " -f "+docker_file+" .")

    file.write("--> Creating containers for new version "+build_version+".\n")

    # Generate te new docker version
    for command in commands:
        file.write("--> Executing command: '"+command+"':\n")
        p = subprocess.run(["powershell.exe", command],
                           stdout = subprocess.PIPE,
                           stderr = subprocess.PIPE)
        file.write(str(p.stdout.decode('utf-8')))
        if p.stderr is not None:
            file.write("-> Error:\n")
            file.write(str(p.stderr.decode('utf-8')))

    file.close()

if __name__ == '__main__':
    main_function()
