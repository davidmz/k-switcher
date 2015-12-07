package main

type Transcoder struct {
	Directions []*TranscodeDirection
}

func NewTranscoder(A, B *KLayout) *Transcoder {
	return &Transcoder{
		Directions: []*TranscodeDirection{
			NewTranscodeDirection(A, B),
			NewTranscodeDirection(B, A),
		},
	}
}

func (t *Transcoder) Transcode(str string) (outs string, lastLayout *KLayout) {
	if str == "" {
		return
	}

	var (
		word, out []rune
		dir       *TranscodeDirection
	)

	for _, r := range str {
		runeDir, nDirs := (*TranscodeDirection)(nil), 0
		for _, d := range t.Directions {
			if d.CanTranscode(r) {
				nDirs++
				runeDir = d
			}
		}
		if nDirs != 1 {
			runeDir = nil
		}

		if dir == nil || runeDir == nil || dir == runeDir { // продолжаем то же слово
			if dir == nil {
				dir = runeDir
			}
			word = append(word, r)
		} else { // новое слово, с другим направление перекодировки
			// конвертируем и сохраняем имеющееся слово
			out = append(out, dir.Transcode(word)...)
			// стартуем новое слово
			word = []rune{r}
			dir = runeDir
		}
	}
	if dir != nil {
		// конвертируем и сохраняем имеющееся слово
		out = append(out, dir.Transcode(word)...)
		lastLayout = dir.Tgt
	} else {
		out = append(out, word...)
	}
	outs = string(out)
	return
}

// Перекодировка из раскладки А  в раскладку Б
type TranscodeDirection struct {
	Src, Tgt *KLayout
	// Таблица перекодирования
	Map map[rune]rune
}

func NewTranscodeDirection(A, B *KLayout) *TranscodeDirection {
	dir := &TranscodeDirection{
		Src: A,
		Tgt: B,
		Map: make(map[rune]rune),
	}
	rA, rB := []rune(A.Layout), []rune(B.Layout)
	for i, r := range rA {
		dir.Map[r] = rB[i]
	}
	return dir
}

func (t *TranscodeDirection) CanTranscode(r rune) bool {
	_, ok := t.Map[r]
	return ok
}

func (t *TranscodeDirection) Transcode(in []rune) (out []rune) {
	out = make([]rune, len(in))
	for i, r := range in {
		if r1, ok := t.Map[r]; ok {
			out[i] = r1
		} else {
			out[i] = r
		}
	}
	return
}
