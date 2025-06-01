# LiberGo
This is a Liber Primus Analysis Toolkit.  I have been using it for investigating the Liber Primus.  It is **very early** in its development and I am working on it as I have the time to do it.  I have also used it on some other puzzles so its application is not limited to the Liber Primus.

## Building release binaries
You will need to have go installed on you system.  You can get it from https://golang.org.

You will also need fyne to build the GUI tools.  You can get it from https://fyne.io.

This script is for Linux and not Windows or Mac.

```
chmod 777 build.sh
./build.sh
```

**This will not destroy the database!**

## Installation
If you on Linux, then you can use install.sh.  It will copy the binaries to /opt.  Then it will create the symlinks to the bin directory.

There are no installers for Windows at this time.  You will need to copy the files to a directory of your choosing.

After you get it installed, run the following command.  It will set up everything for you.

```
libergo -init
```

**Note: If you are on Windows, you will need to copy in the text files into the .libergo directory in your user directory.**

I don't use Windows so if anyone wants to help, it would be much appreciated.

## Database Required!!!
Some tools require a Postgres database to be present.

I prefer to use podman to host the database.  You can use docker if you want.  You can use the following command to start the database for podman.

### Powershell
You will need to install WSL2, Python 3.9+, and Podman.  You can get Podman from https://podman.io/docs/installation.  You can get Python from https://www.python.org/downloads/.
```
pip3 install podman-compose
.\create_podman_db.ps1
```

### Bash
```
chmod 777 create_podman_db.sh
./create_podman_db.sh
```

### libergo
```
libergo -initdbserver
```
or use this if you are using alternate credentials or your own database server.
```
libergo -initTables
```

If you use your own server, you will need to set the variables in the appsettings.json file in the .libergo directory in you user directory.

## Upcoming Tools available in libergo
- circular shift
- least significant bits in message.
- least significant bit in number strings.
- skip and take
- dictionary checker
- clock angle claculator
- letter frequency analysis stuff
- Scytale
- Text spiraler