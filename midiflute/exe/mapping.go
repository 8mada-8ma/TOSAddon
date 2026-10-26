package main

// 鍵盤番号(MIDIノート番号)から、ToSの笛の音番号(fluting_scale.ies の ClassID)への変換

const (
	KeyOut    = iota // 範囲外
	KeyOct1          // 1オクターブ目
	KeyOct2          // 2オクターブ目
	KeyOct3          // 3オクターブ目(解放済み)
	KeyLocked        // 3オクターブ目(未解放)
)

// 白鍵 C D E F G A B / 黒鍵 C# D# F# G# A#
var naturalIdx = map[int]int{0: 0, 2: 1, 4: 2, 5: 3, 7: 4, 9: 5, 11: 6}
var sharpIdx = map[int]int{1: 0, 3: 1, 6: 2, 8: 3, 10: 4}
var naturalBase = [3]int{1, 8, 25}
var sharpBase = [3]int{15, 20, 32}

// noteToID は音番号と鍵盤の種類を返す。鳴らせない場合 id は 0。
func noteToID(note, base int, oct3 bool) (id int, kind int) {
	off := note - base
	if off < 0 || off > 35 {
		return 0, KeyOut
	}
	octave, semi := off/12, off%12
	if octave == 2 && !oct3 {
		return 0, KeyLocked
	}
	if i, ok := naturalIdx[semi]; ok {
		id = naturalBase[octave] + i
	} else {
		id = sharpBase[octave] + sharpIdx[semi]
	}
	return id, KeyOct1 + octave
}

func isBlack(note int) bool {
	_, ok := sharpIdx[((note%12)+12)%12]
	return ok
}

var noteNames = [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

// noteName は "C4" のような表記(中央のド=60=C4)
func noteName(note int) string {
	return noteNames[((note%12)+12)%12] + itoa(note/12-1)
}

func itoa(n int) string {
	if n < 0 {
		return "-" + itoa(-n)
	}
	if n < 10 {
		return string(rune('0' + n))
	}
	return itoa(n/10) + string(rune('0'+n%10))
}
