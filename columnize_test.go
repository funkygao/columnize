package columnize

import (
	"fmt"
	"testing"

	crand "crypto/rand"
)

func TestSimpleFormatWithChinese(t *testing.T) {
	lines := []string{"CPU%|IP"}
	lines = append(lines, fmt.Sprintf("%s|%s", "阿尔法项目", "10.187.131.227"))
	lines = append(lines, fmt.Sprintf("%s|%s", "ide.cc.polaris.jd.com", "10.191.92.130"))

	println(Format(lines, &Config{Glue: "  ", Empty: " "}))
	println()
	fmt.Printf("%-21s %s\n", "阿尔法项目", "abc")
	fmt.Printf("%-21s %s\n", "ide.cc.polaris.jd.com", "cm")
}

func TestRuneLen(t *testing.T) {
	if runeLen("ab") != 2 {
		t.Fatalf("ab should 2")
	}
	if runeLen("我") != 2 {
		t.Fatalf("我 should 2")
	}
	if runeLen("我b") != 3 {
		t.Fatalf("我b should 3")
	}
}

func TestListOfStringsInput(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyLinesOutput(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestLeadingSpacePreserved(t *testing.T) {
	input := []string{
		"| Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "   Column B  Column C\n"
	expected += "x  y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestColumnWidthCalculator(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"Longer than A | Longer than B | Longer than C",
		"short | short | short",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A       Column B       Column C\n"
	expected += "Longer than A  Longer than B  Longer than C\n"
	expected += "short          short          short"

	if output != expected {
		printableProof := fmt.Sprintf("\nGot:      %+q", output)
		printableProof += fmt.Sprintf("\nExpected: %+q", expected)
		t.Fatalf("\n%s", printableProof)
	}
}

func TestColumnWidthCalculatorNonASCII(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"⌘⌘⌘⌘⌘⌘⌘⌘ | Longer than B | Longer than C",
		"short | short | short",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B       Column C\n"
	expected += "⌘⌘⌘⌘⌘⌘⌘⌘  Longer than B  Longer than C\n"
	expected += "short     short          short"

	if output != expected {
		printableProof := fmt.Sprintf("\nGot:      %+q", output)
		printableProof += fmt.Sprintf("\nExpected: %+q", expected)
		t.Fatalf("\n%s", printableProof)
	}
}

func BenchmarkColumnWidthCalculator(b *testing.B) {
	// Generate the input
	input := []string{
		"UUID A | UUID B | UUID C | Column D | Column E",
	}

	format := "%s|%s|%s|%s"
	short := "short"

	uuid := func() string {
		buf := make([]byte, 16)
		if _, err := crand.Read(buf); err != nil {
			panic(fmt.Errorf("failed to read random bytes: %v", err))
		}

		return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
			buf[0:4],
			buf[4:6],
			buf[6:8],
			buf[8:10],
			buf[10:16])
	}

	for i := 0; i < 1000; i++ {
		l := fmt.Sprintf(format, uuid()[:8], uuid()[:12], uuid(), short, short)
		input = append(input, l)
	}

	config := DefaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Format(input, config)
	}
}

func TestVariedInputSpacing(t *testing.T) {
	input := []string{
		"Column A       |Column B|    Column C",
		"x|y|          z",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestUnmatchedColumnCounts(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"Value A | Value B",
		"Value A | Value B | Value C | Value D",
	}

	config := DefaultConfig()
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "Value A   Value B\n"
	expected += "Value A   Value B   Value C   Value D"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternateDelimiter(t *testing.T) {
	input := []string{
		"Column | A % Column | B % Column | C",
		"Value A % Value B % Value C",
	}

	config := DefaultConfig()
	config.Delim = "%"
	output := Format(input, config)

	expected := "Column | A  Column | B  Column | C\n"
	expected += "Value A     Value B     Value C"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternateSpacingString(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	config.Glue = "    "
	output := Format(input, config)

	expected := "Column A    Column B    Column C\n"
	expected += "x           y           z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestSimpleFormatBasic(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	output := SimpleFormat(input)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestAlternatePrefixString(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := DefaultConfig()
	config.Prefix = "  "
	output := Format(input, config)

	expected := "  Column A  Column B  Column C\n"
	expected += "  x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyFieldReplacement(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | | z",
	}

	config := DefaultConfig()
	config.Empty = "<none>"
	output := Format(input, config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         <none>    z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestEmptyConfigValues(t *testing.T) {
	input := []string{
		"Column A | Column B | Column C",
		"x | y | z",
	}

	config := Config{}
	output := Format(input, &config)

	expected := "Column A  Column B  Column C\n"
	expected += "x         y         z"

	if output != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", expected, output)
	}
}

func TestMergeConfig(t *testing.T) {
	conf1 := &Config{Delim: "a", Glue: "a", Prefix: "a", Empty: "a"}
	conf2 := &Config{Delim: "b", Glue: "b", Prefix: "b", Empty: "b"}
	conf3 := &Config{Delim: "c", Prefix: "c"}

	m := MergeConfig(conf1, conf2)
	if m.Delim != "b" || m.Glue != "b" || m.Prefix != "b" || m.Empty != "b" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, conf3)
	if m.Delim != "c" || m.Glue != "a" || m.Prefix != "c" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, nil)
	if m.Delim != "a" || m.Glue != "a" || m.Prefix != "a" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}

	m = MergeConfig(conf1, &Config{})
	if m.Delim != "a" || m.Glue != "a" || m.Prefix != "a" || m.Empty != "a" {
		t.Fatalf("bad: %#v", m)
	}
}
