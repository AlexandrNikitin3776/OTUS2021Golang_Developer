package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("readme example", func(t *testing.T) {
		inputText := "cat and dog, one dog,two cats and one man"
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"and",
				"dog",
				"one",
				"cat",
				"cats",
				"man",
				"two",
			}
			require.Equal(t, expected, Top10(inputText))
		} else {
			expected := []string{
				"and",
				"one",
				"cat",
				"cats",
				"dog,",
				"dog,two",
				"man",
			}
			require.Equal(t, expected, Top10(inputText))
		}
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}

func TestGetWordsFromText(t *testing.T) {
	tableTests := []struct {
		name     string
		text     string
		expected []string
	}{
		{name: "empty", text: "", expected: nil},
		{name: "punctuation", text: "!", expected: nil},
		{name: "one letter words", text: "К В Н", expected: []string{"К", "В", "Н"}},
		{name: "one word", text: "one", expected: []string{"one"}},
		{name: "two words", text: "one two", expected: []string{"one", "two"}},
		{name: "three cyryllyc words", text: "Раз два три", expected: []string{"Раз", "два", "три"}},
		{
			name:     "with hyphen and dash",
			text:     "Какой-либо - это не какой-то",
			expected: []string{"Какой-либо", "это", "не", "какой-то"},
		},
		{
			name:     "with hyphens",
			text:     "тра-та-та-та-та туда-сюда-оттуда",
			expected: []string{"тра-та-та-та-та", "туда-сюда-оттуда"},
		},
		{
			name:     "with punctuation",
			text:     "\"Василий Косяков\" - человек,пароход,имя и фамилия!",
			expected: []string{"Василий", "Косяков", "человек", "пароход", "имя", "и", "фамилия"},
		},
		{name: "with numbers", text: "АК-47 99-го года", expected: []string{"АК-47", "99-го", "года"}},
	}

	for _, tc := range tableTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := getWordsFromText(tc.text)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestMakeFrequencyMap(t *testing.T) {
	tableTests := []struct {
		name     string
		words    []string
		expected map[string]int
	}{
		{name: "empty", words: nil, expected: nil},
		{
			name:     "different word forms",
			words:    []string{"нога", "ногу", "ноги"},
			expected: map[string]int{"нога": 1, "ногу": 1, "ноги": 1},
		},
		{name: "same word forms", words: []string{"Нога", "нога"}, expected: map[string]int{"нога": 2}},
	}

	for _, tc := range tableTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := makeFrequencyMap(tc.words)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestGetKeysFromMap(t *testing.T) {
	tableTests := []struct {
		name     string
		inputMap map[string]int
		expected []string
	}{
		{name: "empty", inputMap: nil, expected: nil},
		{name: "one key", inputMap: map[string]int{"так": 100500}, expected: []string{"так"}},
		{
			name:     "many keys",
			inputMap: map[string]int{"так": 100500, "так-так": 100500, "туда": 1, "сюда": 2},
			expected: []string{"так", "так-так", "туда", "сюда"},
		},
	}

	for _, tc := range tableTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := getKeysFromMap(tc.inputMap)
			require.Subset(t, tc.expected, got)
			require.Subset(t, got, tc.expected)
		})
	}
}

func TestMinOfTwo(t *testing.T) {
	tableTests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{name: "equals", a: 0, b: 0, expected: 0},
		{name: "a is min", a: 100, b: 100500, expected: 100},
		{name: "b is min", a: -100, b: -100500, expected: -100500},
	}

	for _, tc := range tableTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := minOfTwo(tc.a, tc.b)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestSortSliceByMapValues(t *testing.T) {
	tableTests := []struct {
		name     string
		sl       []string
		mp       map[string]int
		expected []string
	}{
		{name: "empty", sl: nil, mp: nil, expected: nil},
		{
			name:     "one word slice",
			sl:       []string{"ё"},
			mp:       map[string]int{"ё": 1},
			expected: []string{"ё"},
		},
		{
			name:     "lexicographical slice",
			sl:       []string{"нога", "ноги", "ногу", "ногой", "ногами"},
			mp:       map[string]int{"нога": 1, "ноги": 1, "ногу": 1, "ногой": 1, "ногами": 1},
			expected: []string{"нога", "ногами", "ноги", "ногой", "ногу"},
		},
	}

	for _, tc := range tableTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			sortSliceByMapValues(tc.sl, tc.mp)
			require.Equal(t, tc.expected, tc.sl)
		})
	}
}
