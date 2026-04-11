# SpringBreak

> [!NOTE]
> WinterBreak3, PenguinsBreak, Etc.

SpringBreak is a jailbreak for KT5, PW5(SE) on 5.19.2+.

## How It Works

1. Load modified store cache with upwards of 5-10 thousand nested folders to prevent cache deletion upon eject.
2. Schedule `SET.SCFG` document with ToDo (Kindle API) in order to set `winmgr.vibrancyMode.pref.path` to custom SH
3. Start `com.lab126.fts` via the AppMgr Kindle API seeing as it can communicate to `com.lab126.winmgr`, pointing to a custom URL. This, in turn, calls `vibrancyMode`, triggering the SH mandated in the ToDo document. (In this case, KindleModding's online jb.sh.)