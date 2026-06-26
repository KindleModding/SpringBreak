# SpringBreak

> [!NOTE]
> WinterBreak3, PenguinsBreak, Etc.

SpringBreak is a jailbreak for KT5, PW5(SE) on 5.19.2+, KT4/PW4 5.18.1.1.1.

Read the write-up [here](https://penguins184.xyz/blog/springbreak-jailbreak/).

## How It Works

1. Load modified store cache with upwards of 5-10 thousand nested folders to prevent cache deletion upon eject.
2. Schedule `SET.SCFG` document with ToDo (Kindle API) in order to set `winmgr.vibrancyMode.pref.path` to custom SH
3. Start `com.lab126.fts` via the AppMgr Kindle API seeing as it can communicate to `com.lab126.winmgr`, pointing to a custom URL. This, in turn, calls `vibrancyMode`, triggering the SH mandated in the ToDo document. (In this case, KindleModding's online jb.sh.)

## Build

Because MacOS scares me:

- Windows: `GOOS=windows GOARCH=amd64 go build -trimpath -o springbreak.exe installer.go`
- Darwin (Build from Source):
    - `GOOS=darwin GOARCH=arm64 go build -trimpath -o springbreak-darwin-arm64 installer.go`
    - `GOOS=darwin GOARCH=amd64 go build -trimpath -o springbreak-darwin-amd64 installer.go`
    - `lipo -create -output springbreak-darwin springbreak-darwin-arm64 springbreak-darwin-amd64`
    - `codesign --force --sign - ./springbreak-darwin` (Or, Developer Cert: Credit Fynn)
- Linux: `GOOS=linux GOARCH=amd64 go build -trimpath -o springbreak-linux installer.go`
