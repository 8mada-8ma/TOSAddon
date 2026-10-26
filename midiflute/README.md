# midiflute - MIDIキーボードでパイドパイパーの笛を演奏

**日本語** | [English](README.en.md) | [한국어](README.ko.md)

**MIDIキーボードの鍵盤を、Tree of Savior のパイドパイパーの笛の演奏にそのまま割り当てるアドオンです。** 鍵盤を1つ押すと1つの音が鳴ります(1鍵=1音)。和音や黒鍵(#)にも対応しています。

**⚠ このアドオンだけでは動きません。** MIDIキーボードの入力を受け取る **`midiflute.exe` を [Releases](https://github.com/8mada-8ma/TOSAddon/releases) からダウンロードし、起動しておく必要があります。** アドオンマネージャーでインストールできるのはアドオン(ipf)だけです。

---

**English:** Plays the Pied Piper flute directly from a MIDI keyboard (one key = one note). **The addon alone does not work. You must download `midiflute.exe` from [Releases](https://github.com/8mada-8ma/TOSAddon/releases) and keep it running.** (Japanese server only.) → [English README](README.en.md)

**한국어:** MIDI 키보드의 건반을 파이드 파이퍼의 피리 연주에 그대로 할당하는 애드온입니다(건반 1개 = 음 1개). **이 애드온만으로는 작동하지 않습니다. [Releases](https://github.com/8mada-8ma/TOSAddon/releases)에서 `midiflute.exe`를 다운로드하여 실행해 두어야 합니다.** (일본 서버 전용) → [한국어 README](README.ko.md)

---

- **対象:** 日本版 Tree of Savior(Steam)
- **必要なもの:** Windows(64bit)、USB接続のMIDIキーボード
- **必要なexeの版:** midiflute.exe v0.4.1 以降

> **📝 このアドオンと midiflute.exe は、AI(Anthropic の Claude)を使って作成しました。** 設計・コード・この説明文の大部分は Claude が書き、作者が自分の環境で動作を確認しています。

> **⚠ 動作確認は、作者の環境(下記)でしか行っていません。** それ以外の環境で動くかどうかは分かりません。
>
> | 項目 | 確認した環境 |
> | --- | --- |
> | OS | Windows 11(64bit) |
> | MIDIキーボード | KORG microKEY2(25鍵) |
> | ゲーム | 日本版 Tree of Savior(Steam、2026年10月時点) |
>
> **確認していないもの:** Windows 10、英語版などの日本語以外の Windows、ほかのメーカーの MIDIキーボード、パッド型の MIDI 機器、ほかのソフトと MIDI 機器を同時に使えない環境 など

## できること

- MIDIキーボードの鍵盤が、ゲーム内の笛の音(3オクターブ・36音)にそのまま対応します
- 鍵盤の割り当ては midiflute.exe の画面で確認・調整できます(1オクターブ目の「ド」にする鍵盤を選ぶだけ)
- 3オクターブ目が解放されていないキャラクターでは、3オクターブ目の鍵盤は鳴りません
- 演奏モードを開いたときに midiflute.exe が起動していなければ、画面とチャットでお知らせします

**できないこと(あえて入れていません):** MIDIファイルの自動演奏など、自動で演奏する機能はありません。演奏するのはあなた自身です。

## インストール

### 1. アドオンを入れる

アドオンマネージャー([MizukiBelhi/Addon-Manager](https://github.com/MizukiBelhi/Addon-Manager))で **midiflute** をインストールしてください。

**手動で入れる場合**

1. [Releases](https://github.com/8mada-8ma/TOSAddon/releases) から `midiflute-v0.4.1.ipf` をダウンロードします。
2. ファイル名を、次の名前に変えます。**⛄(雪だるま)は入力しにくいので、下の名前をコピーして貼り付けてください。**

```
_midiflute-⛄-v0.4.1.ipf
```

   - 名前を変えるには、ファイルを右クリック →「名前の変更」(Windows 11 では右クリックメニューのアイコン)を選びます。
   - 拡張子(`.ipf`)が表示されていない場合は、`.ipf` を除いた `_midiflute-⛄-v0.4.1` の部分だけを貼り付けてください。
3. ToSのインストールフォルダの中の `data` フォルダに入れます。
   - Steamでの探し方: Steamのライブラリで Tree of Savior を右クリック →「管理」→「ローカルファイルを閲覧」

**この名前の形(`_名前-⛄-v版.ipf`)でないと、ゲームに読み込まれず、削除されることがあります。**

### 2. midiflute.exe をダウンロードする

1. [Releases](https://github.com/8mada-8ma/TOSAddon/releases) から `midiflute.exe` をダウンロードします。
2. **置き場所はどこでもかまいません**(デスクトップやダウンロードフォルダのままでOK)。
   - ToSのフォルダは、Steamの情報から自動で探します。見つからないときは、画面の「変更…」でToSのフォルダを選んでください。
3. 初回起動時に、Windowsの「WindowsによってPCが保護されました」という画面が出ることがあります。これは、署名(発行元の証明)の無いソフトに出る警告です。「詳細情報」→「実行」で起動できます。
   - 心配な場合は、[ソースコード](exe/)を公開しているので、中身を確認したり、自分でビルドしたりできます(下の「ソースからビルドする」)。

## 使い方

1. MIDIキーボードをPCにつなぎます。
2. `midiflute.exe` を起動します。「MIDI:」が「接続中」になればOKです。
3. ToSでパイドパイパーのキャラクターを休憩させ、笛の演奏モードを開きます。
   - midiflute.exe の「ToS:」が「演奏モード中 ♪」になれば連携できています。
4. 鍵盤を弾いてください。

### 鍵盤の割り当てを合わせる

midiflute.exe の画面の鍵盤の図で、どの鍵盤がどの音になるか確認できます。

| 色 | 意味 |
| --- | --- |
| 青 | 1オクターブ目 |
| 緑 | 2オクターブ目 |
| 橙 | 3オクターブ目 |
| 灰 | 鳴らない鍵盤(範囲外・3オクターブ目が未解放) |
| 赤 | 押している鍵盤 |

- 「鍵盤で指定」を押してから、1オクターブ目の「ド」にしたい鍵盤を押すと、そこが基準になります。
- 「±半音」「±1オクターブ」で細かく動かすこともできます。
- キーボード本体のオクターブ切替ボタンを動かしたときは、もう一度合わせてください。

### そのほかの設定

- **MIDI機器 / 接続・切断:** 複数のMIDI機器があるときは選べます。「切断」を押すと、ほかのソフト(DAWなど)にキーボードを譲れます。
- **ToSが前面にあるときだけ演奏を送る:** オフ(初期値)だと、ToSが裏にあっても演奏されます。
- **言語 / Language:** 日本語・English・한국어 を切り替えられます。ゲーム内のメッセージも、ここで選んだ言語になります(最初は Windows の表示言語に合わせます)。
- 最小化するとタスクトレイ(画面右下)に入ります。トレイのアイコンを右クリック →「終了」で終了します。
- midiflute.exe は二重には起動しません。2つ目を起動すると、起動中であることをお知らせします。

## 困ったとき

| 表示 | 対処 |
| --- | --- |
| ゲーム内で「midiflute.exe が起動していません」 | midiflute.exe を起動してください。起動するとチャットに「接続しました」と出ます。 |
| 「ToSフォルダが未設定です」 | 「変更…」から、Tree of Savior (Japanese Ver.) のフォルダを選んでください。 |
| 「ToS: 未連携」のまま | アドオンが入っているか確認してください。キャラクターでログインすると連携します。 |
| 「この機器は開けませんでした」 | ほかのソフトがキーボードを使っている可能性があります。そのソフトを閉じるか、キーボードの使用を止めてください。 |
| 「ほかの midiflute.exe が同時に動いているようです」 | 別の場所や古い版の midiflute.exe が動いています。タスクトレイを確認し、1つだけにしてください。 |
| 遅れて鳴る | 画面の「遅延」を確認してください。作者の環境では数十ms程度でした。ToSを前面にすると速くなることがあります。 |

うまく動かないときは、[Issues](https://github.com/8mada-8ma/TOSAddon/issues) で連絡してください。ToSフォルダの `addons\midiflute\log.txt`(エラーがあったときだけ作られます)があれば添えてください。

## アンインストール

1. アドオンマネージャーで midiflute をアンインストールします(手動の場合は `data\_midiflute-⛄-vX.X.X.ipf` を削除)。
2. `midiflute.exe` を削除します。
3. 次の2つのフォルダを削除します。
   - ToSフォルダの中の `addons\midiflute\`
   - `%APPDATA%\midiflute\`(エクスプローラーのアドレス欄に入力すると開けます。midiflute.exe の設定が入っています)

midiflute.exe は、レジストリへの書き込みやスタートアップへの登録をしません。

## しくみ

```
MIDIキーボード ──▶ midiflute.exe ──(ファイル)──▶ アドオン ──▶ ゲームの笛の演奏関数
                    鍵盤→音の変換      addons\midiflute\events.txt
```

- アドオンは、ゲーム内の演奏ボタンを押したときと同じ関数(`Piedpiper.ReqPlayFluting` / `ReqStopFluting`)を呼んでいます。キー入力の偽装や、ゲームのファイル・メモリの書き換えはしていません。
- midiflute.exe とアドオンは、ToSフォルダの `addons\midiflute\` の中のテキストファイルだけでやり取りします。インターネット通信は行いません。
  - `events.txt`(exe → アドオン): 鳴らす・止める音
  - `ack.txt`(アドオン → exe): 処理済みの位置
  - `status.txt`(アドオン → exe): 演奏モード中か、3オクターブ目が解放済みか
  - `helper.txt`(exe → アドオン): exe が起動している証拠
  - `lang.txt`(exe → アドオン): メッセージの言語
- アドオンが演奏データを読むのは、演奏モードが開いている間だけです。

## ソースからビルドする

- アドオン: [`addon/midiflute/`](addon/midiflute/) の Lua と XML を、[`addon/ipfpack.py`](addon/ipfpack.py) で ipf にまとめています。
- midiflute.exe: [`exe/`](exe/) を Go 1.24 でビルドしています。

```
cd exe
go install github.com/akavel/rsrc@latest
rsrc -manifest midiflute.manifest -arch amd64 -o rsrc_windows_amd64.syso
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o midiflute.exe .
```

## 利用規約について(必ずお読みください)

- 日本版 Tree of Savior の[利用規約](https://treeofsavior.jp/page/ts.php)第18条(禁止事項)には、「外部ツールまたは外部プログラムの利用、開発行為」(第32号)などの項目があります。**アドオンや midiflute.exe の使用がこれに当たると判断される可能性はあります。**
- 運営はアドオンを公認していません。このアドオンも、運営の許可を得たものではありません。
- このアドオンは、演奏を自動化するものではなく、1回の打鍵で1つの音を鳴らすだけです。ただし、それで問題ないかどうかを判断するのは運営です。
- **使用は自己責任でお願いします。** 運営から指摘や方針の発表があった場合は、そちらに従ってください。

## 免責事項

- このアドオンと midiflute.exe を使用したことによって生じたいかなる損害(アカウントへの処分を含みます)についても、作者は責任を負いません。
- 個人が作成・配布しているものであり、運営とは一切関係ありません。

## 連絡先

- GitHub: [@8mada-8ma](https://github.com/8mada-8ma)(不具合報告は [Issues](https://github.com/8mada-8ma/TOSAddon/issues) へ)

## 更新履歴

- **v0.4.1**
  - 1秒あたりに鳴らす音の数に上限(40音)を設けました(異常な連打をサーバーへ送らないため)
  - midiflute.exe が動いていないときは、演奏データを読まないようにしました
- **v0.4.0**
  - 英語・韓国語に対応(midiflute.exe の画面とゲーム内メッセージ)
  - 手動インストールの手順を追加
- **v0.3.0**
  - 初回公開
