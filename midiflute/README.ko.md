# midiflute - MIDI 키보드로 파이드 파이퍼의 피리 연주

[日本語](README.md) | [English](README.en.md) | **한국어**

**MIDI 키보드의 건반을 Tree of Savior 파이드 파이퍼의 피리 연주에 그대로 할당하는 애드온입니다.** 건반 1개를 누르면 음 1개가 납니다(건반 1개 = 음 1개). 화음과 검은 건반(#)도 지원합니다.

**⚠ 이 애드온만으로는 작동하지 않습니다.** **MIDI 키보드의 입력을 받는 `midiflute.exe`를 [Releases](https://github.com/8mada-8ma/TOSAddon/releases)에서 다운로드하여 실행해 두어야 합니다.** 애드온 매니저로 설치할 수 있는 것은 애드온(ipf)뿐입니다.

- **대상:** Tree of Savior 일본 서버 (Steam, "Tree of Savior (Japanese Ver.)")
- **필요한 것:** Windows (64비트), USB 연결 MIDI 키보드
- **필요한 exe 버전:** midiflute.exe v0.4.1 이상

> **📝 이 애드온과 midiflute.exe는 AI(Anthropic의 Claude)를 사용하여 만들었습니다.** 설계, 코드, 이 설명문의 대부분은 Claude가 작성했으며, 작성자가 자신의 환경에서 동작을 확인했습니다.

> **⚠ 동작 확인은 작성자의 환경(아래)에서만 했습니다.** 그 밖의 환경에서 작동하는지는 알 수 없습니다.
>
> | 항목 | 확인한 환경 |
> | --- | --- |
> | OS | Windows 11 (64비트, 일본어) |
> | MIDI 키보드 | KORG microKEY2 (25건반) |
> | 게임 | Tree of Savior 일본 서버 (Steam, 2026년 10월 기준) |
>
> **확인하지 않은 것:** Windows 10, 한국어·영어 등 일본어 이외의 Windows, 다른 제조사의 MIDI 키보드, 패드형 MIDI 기기, 다른 프로그램과 MIDI 기기를 동시에 사용할 수 없는 환경 등

## 할 수 있는 것

- MIDI 키보드의 건반이 게임 속 피리 음(3옥타브・36음)에 그대로 대응합니다.
- 건반 배치는 midiflute.exe 화면에서 확인・조정할 수 있습니다(1옥타브의 "도"로 할 건반을 고르기만 하면 됩니다).
- 3옥타브가 해금되지 않은 캐릭터에서는 3옥타브 건반이 소리 나지 않습니다.
- 연주 모드를 열었을 때 midiflute.exe가 실행되지 않았으면 화면과 채팅으로 알려 줍니다.

**할 수 없는 것(일부러 넣지 않았습니다):** MIDI 파일 자동 연주 등 자동으로 연주하는 기능은 없습니다. 연주하는 것은 여러분 자신입니다.

## 설치

### 1. 애드온 설치

애드온 매니저([MizukiBelhi/Addon-Manager](https://github.com/MizukiBelhi/Addon-Manager))로 **midiflute**를 설치하세요.

**수동으로 설치하는 경우**

1. [Releases](https://github.com/8mada-8ma/TOSAddon/releases)에서 `midiflute-v0.4.1.ipf`를 다운로드합니다.
2. 파일 이름을 아래 이름으로 바꿉니다. **⛄(눈사람)은 입력하기 어려우니 아래 이름을 복사해서 붙여 넣으세요.**

```
_midiflute-⛄-v0.4.1.ipf
```

   - 이름을 바꾸려면 파일을 오른쪽 클릭 → "이름 바꾸기"를 선택합니다.
   - 확장자(`.ipf`)가 표시되지 않는 경우에는 `.ipf`를 뺀 `_midiflute-⛄-v0.4.1` 부분만 붙여 넣으세요.
3. ToS 설치 폴더 안의 `data` 폴더에 넣습니다.
   - Steam에서 찾는 방법: Steam 라이브러리에서 Tree of Savior를 오른쪽 클릭 → "관리" → "로컬 파일 보기"

**이 이름 형식(`_이름-⛄-v버전.ipf`)이 아니면 게임이 읽지 않고 삭제할 수 있습니다.**

### 2. midiflute.exe 다운로드

1. [Releases](https://github.com/8mada-8ma/TOSAddon/releases)에서 `midiflute.exe`를 다운로드합니다.
2. **두는 위치는 어디든 상관없습니다**(바탕 화면이나 다운로드 폴더 그대로도 괜찮습니다).
   - ToS 폴더는 Steam 정보에서 자동으로 찾습니다. 찾지 못하면 화면의 "변경…"에서 ToS 폴더를 선택하세요.
3. 처음 실행할 때 Windows의 "Windows의 PC 보호" 화면이 나올 수 있습니다. 서명이 없는 프로그램에 나오는 경고입니다. "추가 정보" → "실행"으로 실행할 수 있습니다.
   - 걱정되는 경우 [소스 코드](exe/)를 공개하고 있으므로 내용을 확인하거나 직접 빌드할 수 있습니다(아래 "소스에서 빌드하기").

## 사용 방법

1. MIDI 키보드를 PC에 연결합니다.
2. `midiflute.exe`를 실행합니다. "MIDI:"가 "연결됨"이 되면 준비 완료입니다.
3. ToS에서 파이드 파이퍼 캐릭터를 휴식시키고 피리 연주 모드를 엽니다.
   - midiflute.exe의 "ToS:"가 "연주 모드 중 ♪"이 되면 연동된 것입니다.
4. 건반을 연주하세요.

### 건반 배치 맞추기

midiflute.exe 화면의 건반 그림에서 어느 건반이 어느 음이 되는지 확인할 수 있습니다.

| 색 | 의미 |
| --- | --- |
| 파랑 | 1옥타브 |
| 초록 | 2옥타브 |
| 주황 | 3옥타브 |
| 회색 | 소리 나지 않는 건반(범위 밖・3옥타브 미해금) |
| 빨강 | 누르고 있는 건반 |

- "건반으로 지정"을 누른 다음 1옥타브의 "도"로 할 건반을 누르면 그 건반이 기준이 됩니다.
- "±반음" "±1 옥타브"로 세밀하게 옮길 수도 있습니다.
- 키보드 본체의 옥타브 전환 버튼을 바꿨다면 다시 맞춰 주세요.

### 그 밖의 설정

- **MIDI 기기 / 연결・연결 해제:** MIDI 기기가 여러 개일 때 선택할 수 있습니다. "연결 해제"를 누르면 다른 프로그램(DAW 등)에 키보드를 넘겨줄 수 있습니다.
- **ToS가 앞에 있을 때만 연주 보내기:** 꺼져 있으면(기본값) ToS가 뒤에 있어도 연주됩니다.
- **言語 / Language:** 日本語・English・한국어를 전환할 수 있습니다. 게임 속 메시지도 여기서 고른 언어로 표시됩니다(처음에는 Windows의 표시 언어에 맞춥니다). **midiflute.exe를 한 번도 실행하지 않은 동안에는 게임 속 메시지가 일본어로 표시됩니다.**
- 최소화하면 작업 표시줄의 트레이(화면 오른쪽 아래)에 들어갑니다. 트레이 아이콘을 오른쪽 클릭 → "종료"로 종료합니다.
- midiflute.exe는 중복 실행되지 않습니다. 두 번째로 실행하면 이미 실행 중이라고 알려 줍니다.

## 문제가 생겼을 때

| 표시 | 대처 |
| --- | --- |
| 게임 속에서 "midiflute.exe가 실행되지 않았습니다" | midiflute.exe를 실행하세요. 실행하면 채팅에 "midiflute.exe와 연결되었습니다."가 나옵니다. |
| "ToS 폴더가 설정되지 않았습니다" | "변경…"에서 Tree of Savior (Japanese Ver.) 폴더를 선택하세요. |
| "ToS: 연동 안 됨"이 계속됨 | 애드온이 설치되어 있는지 확인하세요. 캐릭터로 로그인하면 연동됩니다. |
| "이 기기를 열 수 없습니다" | 다른 프로그램이 키보드를 사용 중일 수 있습니다. 그 프로그램을 닫거나 키보드 사용을 멈춰 주세요. |
| "다른 midiflute.exe가 동시에 실행 중인 것 같습니다" | 다른 위치나 이전 버전의 midiflute.exe가 실행 중입니다. 트레이를 확인하고 하나만 남겨 주세요. |
| 늦게 소리가 남 | 화면의 "지연"을 확인하세요. 작성자의 환경에서는 수십 ms 정도였습니다. ToS를 앞으로 가져오면 빨라질 수 있습니다. |

잘 작동하지 않으면 [Issues](https://github.com/8mada-8ma/TOSAddon/issues)로 연락해 주세요(한국어도 괜찮습니다). ToS 폴더에 `addons\midiflute\log.txt`(오류가 있을 때만 만들어집니다)가 있으면 첨부해 주세요.

## 제거

1. 애드온 매니저로 midiflute를 제거합니다(수동 설치의 경우 `data\_midiflute-⛄-vX.X.X.ipf`를 삭제).
2. `midiflute.exe`를 삭제합니다.
3. 다음 두 폴더를 삭제합니다.
   - ToS 폴더 안의 `addons\midiflute\`
   - `%APPDATA%\midiflute\`(탐색기 주소창에 입력하면 열립니다. midiflute.exe의 설정이 들어 있습니다)

midiflute.exe는 레지스트리에 쓰거나 시작 프로그램에 등록하지 않습니다.

## 작동 원리

```
MIDI 키보드 ──▶ midiflute.exe ──(파일)──▶ 애드온 ──▶ 게임의 피리 연주 함수
                 건반→음 변환        addons\midiflute\events.txt
```

- 애드온은 게임 속 연주 버튼을 눌렀을 때와 같은 함수(`Piedpiper.ReqPlayFluting` / `ReqStopFluting`)를 호출합니다. 키 입력을 위장하거나 게임 파일・메모리를 바꾸지 않습니다.
- midiflute.exe와 애드온은 ToS 폴더의 `addons\midiflute\` 안의 텍스트 파일로만 주고받습니다. 인터넷 통신은 하지 않습니다.
- 애드온이 연주 데이터를 읽는 것은 연주 모드가 열려 있는 동안뿐입니다.

## 소스에서 빌드하기

- 애드온: [`addon/midiflute/`](addon/midiflute/)의 Lua와 XML을 [`addon/ipfpack.py`](addon/ipfpack.py)로 ipf로 묶습니다.
- midiflute.exe: [`exe/`](exe/)를 Go 1.24로 빌드합니다.

```
cd exe
go install github.com/akavel/rsrc@latest
rsrc -manifest midiflute.manifest -arch amd64 -o rsrc_windows_amd64.syso
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o midiflute.exe .
```

## 이용약관에 대하여(꼭 읽어 주세요)

- 일본판 Tree of Savior [이용약관](https://treeofsavior.jp/page/ts.php) 제18조(금지 사항)에는 "외부 도구 또는 외부 프로그램의 이용, 개발 행위"(제32호) 등의 항목이 있습니다. **애드온이나 midiflute.exe의 사용이 이에 해당한다고 판단될 가능성이 있습니다.**
- 운영 측은 애드온을 공인하지 않았습니다. 이 애드온도 운영 측의 허가를 받은 것이 아닙니다.
- 이 애드온은 연주를 자동화하지 않으며, 한 번의 타건으로 음 하나를 낼 뿐입니다. 다만 이것이 문제없는지 판단하는 것은 운영 측입니다.
- **사용은 본인 책임입니다.** 운영 측의 지적이나 방침 발표가 있으면 그에 따라 주세요.

## 면책 사항

- 이 애드온과 midiflute.exe를 사용하여 생긴 어떠한 손해(계정 처분 포함)에 대해서도 작성자는 책임지지 않습니다.
- 개인이 만들어 배포하는 것이며 운영 측과는 전혀 관계없습니다.

## 연락처

- GitHub: [@8mada-8ma](https://github.com/8mada-8ma) (버그 신고는 [Issues](https://github.com/8mada-8ma/TOSAddon/issues)로)

## 업데이트 기록

- **v0.4.1**
  - 1초에 내는 음의 수에 상한(40음)을 두었습니다(비정상적인 연타를 서버로 보내지 않기 위해)
  - midiflute.exe가 실행되지 않았을 때는 연주 데이터를 읽지 않도록 했습니다
- **v0.4.0**
  - 영어・한국어 지원(midiflute.exe 화면과 게임 속 메시지)
  - 수동 설치 방법 추가
- **v0.3.0**
  - 첫 공개
