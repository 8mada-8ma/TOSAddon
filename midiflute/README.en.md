# midiflute - Play the Pied Piper flute with a MIDI keyboard

[日本語](README.md) | **English** | [한국어](README.ko.md)

**An addon that maps the keys of your MIDI keyboard directly to the Pied Piper flute in Tree of Savior.** One key plays one note (one key = one note). Chords and black keys (sharps) are supported.

**⚠ The addon alone does not work.** **You must download `midiflute.exe`, which receives input from your MIDI keyboard, from [Releases](https://github.com/8mada-8ma/TOSAddon/releases) and keep it running.** The Addon Manager can only install the addon (ipf).

- **Target:** Tree of Savior Japanese server (Steam, "Tree of Savior (Japanese Ver.)")
- **Requirements:** Windows (64-bit), a USB MIDI keyboard
- **Required exe version:** midiflute.exe v0.4.1 or later

> **📝 This addon and midiflute.exe were made with AI (Anthropic's Claude).** Claude wrote most of the design, code, and this document. The author tested it in their own environment.

> **⚠ It has only been tested in the author's environment (below).** It is unknown whether it works anywhere else.
>
> | Item | Tested environment |
> | --- | --- |
> | OS | Windows 11 (64-bit, Japanese) |
> | MIDI keyboard | KORG microKEY2 (25 keys) |
> | Game | Tree of Savior Japanese server (Steam, as of October 2026) |
>
> **Not tested:** Windows 10, non-Japanese Windows (English, Korean, etc.), MIDI keyboards from other makers, pad-type MIDI devices, systems where a MIDI device cannot be shared between programs, and so on.

## Features

- The keys of your MIDI keyboard map directly to the in-game flute notes (3 octaves, 36 notes).
- You can check and adjust the key mapping in the midiflute.exe window (just choose which key is C of the 1st octave).
- On characters that have not unlocked the 3rd octave, 3rd-octave keys make no sound.
- If midiflute.exe is not running when you open playing mode, you are notified on screen and in chat.

**Not included (on purpose):** There is no automatic playing, such as playing MIDI files. You play it yourself.

## Installation

### 1. Install the addon

Install **midiflute** with the Addon Manager ([MizukiBelhi/Addon-Manager](https://github.com/MizukiBelhi/Addon-Manager)).

**Manual installation**

1. Download `midiflute-v0.4.1.ipf` from [Releases](https://github.com/8mada-8ma/TOSAddon/releases).
2. Rename the file to the name below. **⛄ (snowman) is hard to type, so copy and paste it.**

```
_midiflute-⛄-v0.4.1.ipf
```

   - To rename, right-click the file → "Rename".
   - If file extensions are hidden, paste only `_midiflute-⛄-v0.4.1` (without `.ipf`).
3. Put it in the `data` folder inside the ToS installation folder.
   - To find it in Steam: right-click Tree of Savior in your library → "Manage" → "Browse local files".

**If the name does not have this form (`_name-⛄-vVersion.ipf`), the game may not load it and may delete it.**

### 2. Download midiflute.exe

1. Download `midiflute.exe` from [Releases](https://github.com/8mada-8ma/TOSAddon/releases).
2. **You can put it anywhere** (your desktop or Downloads folder is fine).
   - The ToS folder is found automatically from Steam. If it is not found, choose the ToS folder with "Change…" in the window.
3. On first launch, Windows may show "Windows protected your PC". This warning appears for software without a signature. Click "More info" → "Run anyway".
   - If you are worried, the [source code](exe/) is public, so you can read it or build it yourself (see "Building from source" below).

## Usage

1. Connect your MIDI keyboard to the PC.
2. Start `midiflute.exe`. It is ready when "MIDI:" shows "Connected".
3. In ToS, let your Pied Piper character rest and open the flute playing mode.
   - When "ToS:" in midiflute.exe shows "Playing mode ♪", it is linked.
4. Play the keys.

### Adjusting the key mapping

The keyboard picture in the midiflute.exe window shows which key plays which note.

| Color | Meaning |
| --- | --- |
| Blue | 1st octave |
| Green | 2nd octave |
| Orange | 3rd octave |
| Gray | Silent keys (out of range / 3rd octave locked) |
| Red | Keys being pressed |

- Click "Set by key", then press the key you want as C of the 1st octave. That key becomes the base.
- You can also fine-tune with "±1 semitone" and "±1 octave".
- If you change the octave with the buttons on your keyboard, set it again.

### Other settings

- **MIDI device / Connect・Disconnect:** Choose a device if you have more than one. "Disconnect" releases the keyboard for other software (such as a DAW).
- **Send notes only while ToS is in the foreground:** When off (default), notes are played even if ToS is in the background.
- **言語 / Language:** Switch between 日本語, English, and 한국어. In-game messages use the language chosen here (the first time, it follows your Windows display language). **Until midiflute.exe has been started once, in-game messages are shown in Japanese.**
- Minimizing the window puts it in the system tray (bottom right of the screen). Right-click the tray icon → "Quit" to exit.
- midiflute.exe does not start twice. If you start a second one, you are told it is already running.

## Troubleshooting

| Message | What to do |
| --- | --- |
| "midiflute.exe is not running" in the game | Start midiflute.exe. "Connected to midiflute.exe." appears in chat. |
| "The ToS folder is not set" | Choose the Tree of Savior (Japanese Ver.) folder with "Change…". |
| "ToS: Not linked" stays | Check that the addon is installed. It links after you log in with a character. |
| "Could not open this device" | Another program may be using the keyboard. Close it or stop it from using the keyboard. |
| "Another midiflute.exe seems to be running" | Another copy or an older version of midiflute.exe is running. Check the system tray and keep only one. |
| Notes are delayed | Check "Latency" in the window. In the author's environment it was a few tens of ms. Bringing ToS to the front may make it faster. |

If something does not work, please report it in [Issues](https://github.com/8mada-8ma/TOSAddon/issues) (English is fine). If `addons\midiflute\log.txt` exists in the ToS folder (it is created only when an error occurs), please attach it.

## Uninstall

1. Uninstall midiflute with the Addon Manager (if installed manually, delete `data\_midiflute-⛄-vX.X.X.ipf`).
2. Delete `midiflute.exe`.
3. Delete these two folders:
   - `addons\midiflute\` inside the ToS folder
   - `%APPDATA%\midiflute\` (type it into the Explorer address bar to open it; it contains midiflute.exe's settings)

midiflute.exe does not write to the registry or register itself to start with Windows.

## How it works

```
MIDI keyboard ──▶ midiflute.exe ──(files)──▶ addon ──▶ the game's flute functions
                   key → note           addons\midiflute\events.txt
```

- The addon calls the same functions that run when you click the in-game playing buttons (`Piedpiper.ReqPlayFluting` / `ReqStopFluting`). It does not fake key input or modify game files or memory.
- midiflute.exe and the addon communicate only through text files in `addons\midiflute\` inside the ToS folder. No internet communication is used.
- The addon reads the playing data only while playing mode is open.

## Building from source

- Addon: the Lua and XML in [`addon/midiflute/`](addon/midiflute/) are packed into an ipf with [`addon/ipfpack.py`](addon/ipfpack.py).
- midiflute.exe: [`exe/`](exe/) is built with Go 1.24.

```
cd exe
go install github.com/akavel/rsrc@latest
rsrc -manifest midiflute.manifest -arch amd64 -o rsrc_windows_amd64.syso
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o midiflute.exe .
```

## About the Terms of Service (please read)

- Article 18 (Prohibited acts) of the Japanese Tree of Savior [Terms of Service](https://treeofsavior.jp/page/ts.php) includes "using or developing external tools or external programs" (item 32). **Using addons or midiflute.exe may be judged to fall under this.**
- The operator does not officially approve addons. This addon is not approved by the operator either.
- This addon does not automate playing; one keystroke plays one note. However, it is the operator who decides whether this is acceptable.
- **Use at your own risk.** If the operator issues a notice or policy, follow it.

## Disclaimer

- The author is not responsible for any damage caused by using this addon or midiflute.exe (including action taken against your account).
- This is made and distributed by an individual and has nothing to do with the operator.

## Contact

- GitHub: [@8mada-8ma](https://github.com/8mada-8ma) (bug reports: [Issues](https://github.com/8mada-8ma/TOSAddon/issues))

## Changelog

- **v0.4.1**
  - Added a limit of 40 notes per second (so abnormal rapid input is not sent to the server)
  - The addon no longer reads playing data while midiflute.exe is not running
- **v0.4.0**
  - English and Korean support (midiflute.exe window and in-game messages)
  - Added manual installation steps
- **v0.3.0**
  - First release
