# goms-tool
CMS V138 [Client](https://mega.nz/file/ml01DAjK#ARUbHJr1mKdgoQwIeW5P4qmEWAYgx9xj0mCskcMcTlU) Tool
## Dllhook
The ijl15.dll patch forward client connect request to goms login server
- Convert the original ijl15.dll in maplestory to bytes
- Copy the bytes to ijl15_raw.h
- Modify OPT_ADDR_HOSTNAME in global.h to goms login server ip 
- Switch to Release x32 and build hook file
- Directly cover the hook file to the origin ijl15.dll in maplestory directory
- Or put the hook file into the assets directory of the launcher
## Launcher
The MapleStory launcher use for skipping original login method
- Upgrade go version to above 1.21
- Download package tool from [rsrc](https://github.com/akavel/rsrc/releases)
- Package manifest and icon to syso file:
- ```rsrc.exe -manifest launcher.manifest -ico ./assets/icon.ico -o assets.syso```
- To get rid of the cmd window:
- ```go build -ldflags="-s -w -H windowsgui"```
## Router
A http server that handles launcher login requests