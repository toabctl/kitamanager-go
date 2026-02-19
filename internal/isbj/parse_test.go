package isbj

import (
	"os"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

const testFile = "testdata/Abrechnung_11-25_0770_anonymized.xlsx"

func TestOpenSenatsabrechnung(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	if f == nil {
		t.Fatal("OpenSenatsabrechnung() returned nil file")
	}
}

func TestOpenSenatsabrechnungNonExistent(t *testing.T) {
	_, err := OpenSenatsabrechnung("nonexistent.xlsx")
	if err == nil {
		t.Error("OpenSenatsabrechnung() expected error for non-existent file")
	}
}

func TestParseEinrichtung(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	eu, err := ParseEinrichtung(f)
	if err != nil {
		t.Fatalf("ParseEinrichtung() error = %v", err)
	}

	if eu == nil {
		t.Fatal("ParseEinrichtung() returned nil")
	}

	// Check name
	if eu.Name != "Kindertagesstätte Sternenstaub" {
		t.Errorf("ParseEinrichtung() Name = %q, expected %q", eu.Name, "Kindertagesstätte Sternenstaub")
	}

	// Check summe is positive
	if eu.Summe <= 0 {
		t.Errorf("ParseEinrichtung() Summe = %d, expected positive value", eu.Summe)
	}

	// Check SPH surcharge (1656.80 EUR = 165680 cents)
	expectedSPH := 165680
	if eu.ZuschlagSPH != expectedSPH {
		t.Errorf("ParseEinrichtung() ZuschlagSPH = %d, expected %d", eu.ZuschlagSPH, expectedSPH)
	}

	// Check Summe (49406.53 EUR = 4940653 cents)
	expectedSumme := 4940653
	if eu.Summe != expectedSumme {
		t.Errorf("ParseEinrichtung() Summe = %d, expected %d", eu.Summe, expectedSumme)
	}
}

func TestParseAbrechnung(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	ar, err := ParseAbrechnung(f)
	if err != nil {
		t.Fatalf("ParseAbrechnung() error = %v", err)
	}

	if ar == nil {
		t.Fatal("ParseAbrechnung() returned nil")
	}

	// VertragsBuchung: 49406.53 EUR = 4940653 cents
	if ar.VertragsBuchung != 4940653 {
		t.Errorf("ParseAbrechnung() VertragsBuchung = %d, expected %d", ar.VertragsBuchung, 4940653)
	}

	// KorrekturBuchung: 18230.85 EUR = 1823085 cents
	if ar.KorrekturBuchung != 1823085 {
		t.Errorf("ParseAbrechnung() KorrekturBuchung = %d, expected %d", ar.KorrekturBuchung, 1823085)
	}
}

func TestParseVertrag(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	v, err := ParseVertrag(f)
	if err != nil {
		t.Fatalf("ParseVertrag() error = %v", err)
	}

	if v == nil {
		t.Fatal("ParseVertrag() returned nil")
	}

	// Check we have exactly 40 children
	expectedCount := 40
	if len(v.Kinder) != expectedCount {
		t.Errorf("ParseVertrag() got %d children, expected %d", len(v.Kinder), expectedCount)
	}
}

func TestParseVertragKinderDetails(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	v, err := ParseVertrag(f)
	if err != nil {
		t.Fatalf("ParseVertrag() error = %v", err)
	}

	if len(v.Kinder) == 0 {
		t.Fatal("ParseVertrag() returned no children")
	}

	// Check first child has required fields
	first := v.Kinder[0]

	if first.Gutscheinnummer == "" {
		t.Error("First child has empty Gutscheinnummer")
	}

	if !voucherPattern.MatchString(first.Gutscheinnummer) {
		t.Errorf("First child Gutscheinnummer %q doesn't match pattern GB-XXXXXXXXXXX-XX", first.Gutscheinnummer)
	}

	if first.Name == "" {
		t.Error("First child has empty Name")
	}

	if first.Geburtsdatum == "" {
		t.Error("First child has empty Geburtsdatum")
	}

	if first.Betreuungsumfang == "" {
		t.Error("First child has empty Betreuungsumfang")
	}

	if first.Bezirk == 0 {
		t.Error("First child has zero Bezirk")
	}

	// Check QM/MSS are valid values
	validQMMSS := map[string]bool{"ja": true, "nein": true, "": true}
	if !validQMMSS[first.QM] {
		t.Errorf("First child QM = %q, expected 'ja', 'nein', or empty", first.QM)
	}
	if !validQMMSS[first.MSS] {
		t.Errorf("First child MSS = %q, expected 'ja', 'nein', or empty", first.MSS)
	}

	// Check Betreuungsumfang is valid
	validScope := map[string]bool{"teilzeit": true, "ganztags": true, "erweitert": true}
	if !validScope[first.Betreuungsumfang] {
		t.Errorf("First child Betreuungsumfang = %q, expected 'teilzeit', 'ganztags', or 'erweitert'", first.Betreuungsumfang)
	}

	// Check first child Summe (899.87 EUR = 89987 cents)
	if first.Summe != 89987 {
		t.Errorf("First child Summe = %d, expected %d", first.Summe, 89987)
	}
}

func TestParseVertragAllKinderHaveVoucher(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	v, err := ParseVertrag(f)
	if err != nil {
		t.Fatalf("ParseVertrag() error = %v", err)
	}

	for i, kind := range v.Kinder {
		if !voucherPattern.MatchString(kind.Gutscheinnummer) {
			t.Errorf("Child %d: Gutscheinnummer %q doesn't match pattern", i, kind.Gutscheinnummer)
		}
	}
}

func TestParseVertragAllKinderHavePositiveSumme(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	v, err := ParseVertrag(f)
	if err != nil {
		t.Fatalf("ParseVertrag() error = %v", err)
	}

	for i, kind := range v.Kinder {
		if kind.Summe <= 0 {
			t.Errorf("Child %d (%s): Summe = %d, expected positive value", i, kind.Name, kind.Summe)
		}
	}
}

func TestParseKind(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	// Row 9 should be the first child
	kind, err := ParseKind(f, 9)
	if err != nil {
		t.Fatalf("ParseKind() error = %v", err)
	}

	if kind == nil {
		t.Fatal("ParseKind() returned nil")
	}

	if kind.Gutscheinnummer == "" {
		t.Error("ParseKind() Gutscheinnummer is empty")
	}

	if kind.Name == "" {
		t.Error("ParseKind() Name is empty")
	}
}

func TestEinrichtungString(t *testing.T) {
	eu := Einrichtung{
		Name:       "Test Kita",
		Summe:      100050, // 1000.50 EUR
		ZuschlagQM: 10000,  // 100.00 EUR
	}

	s := eu.String()

	if s == "" {
		t.Error("Einrichtung.String() returned empty string")
	}

	if !strings.Contains(s, "Test Kita") {
		t.Errorf("Einrichtung.String() = %q, expected to contain 'Test Kita'", s)
	}

	if !strings.Contains(s, "1000.50") {
		t.Errorf("Einrichtung.String() = %q, expected to contain '1000.50'", s)
	}
}

func TestAbrechnungString(t *testing.T) {
	ar := Abrechnung{
		VertragsBuchung:  500000, // 5000.00 EUR
		KorrekturBuchung: 100000, // 1000.00 EUR
	}

	s := ar.String()

	if s == "" {
		t.Error("Abrechnung.String() returned empty string")
	}

	if !strings.Contains(s, "5000.00") {
		t.Errorf("Abrechnung.String() = %q, expected to contain '5000.00'", s)
	}
	if !strings.Contains(s, "1000.00") {
		t.Errorf("Abrechnung.String() = %q, expected to contain '1000.00'", s)
	}
}

func TestKindString(t *testing.T) {
	kind := Kind{
		Gutscheinnummer: "GB-12345678901-02",
		Name:            "Löwe, Tiger",
		Geburtsdatum:    "05.20",
		Summe:           89987, // 899.87 EUR
	}

	s := kind.String()

	if s == "" {
		t.Error("Kind.String() returned empty string")
	}

	if !strings.Contains(s, "Löwe, Tiger") {
		t.Errorf("Kind.String() = %q, expected to contain 'Löwe, Tiger'", s)
	}

	if !strings.Contains(s, "GB-12345678901-02") {
		t.Errorf("Kind.String() = %q, expected to contain voucher number", s)
	}

	if !strings.Contains(s, "05.20") {
		t.Errorf("Kind.String() = %q, expected to contain birth date", s)
	}

	if !strings.Contains(s, "899.87") {
		t.Errorf("Kind.String() = %q, expected to contain '899.87'", s)
	}
}

func TestVertragString(t *testing.T) {
	v := Vertrag{
		Kinder: []Kind{
			{Name: "Kind 1", Gutscheinnummer: "GB-11111111111-01", Geburtsdatum: "01.20", Summe: 10000},
			{Name: "Kind 2", Gutscheinnummer: "GB-22222222222-02", Geburtsdatum: "02.21", Summe: 20000},
		},
	}

	s := v.String()

	if s == "" {
		t.Error("Vertrag.String() returned empty string")
	}

	if !strings.Contains(s, "2 Kindern") {
		t.Errorf("Vertrag.String() = %q, expected to contain '2 Kindern'", s)
	}

	if !strings.Contains(s, "Kind 1") {
		t.Errorf("Vertrag.String() = %q, expected to contain 'Kind 1'", s)
	}
	if !strings.Contains(s, "Kind 2") {
		t.Errorf("Vertrag.String() = %q, expected to contain 'Kind 2'", s)
	}
}

func TestCellAsString(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	// Test reading a known cell (E9 should be a voucher number)
	val := cellAsString(f, SheetVertrag, "E9")
	if val == "" {
		t.Error("cellAsString() returned empty for E9")
	}

	if !voucherPattern.MatchString(val) {
		t.Errorf("cellAsString() E9 = %q, expected voucher pattern", val)
	}
}

func TestCellAsInt64(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	// Test reading Bezirk (O9)
	val, err := cellAsInt64(f, SheetVertrag, "O9")
	if err != nil {
		t.Fatalf("cellAsInt64() error = %v", err)
	}

	// Bezirk should be a valid Berlin district (1-12)
	if val < 1 || val > 12 {
		t.Errorf("cellAsInt64() O9 = %d, expected Berlin district 1-12", val)
	}
}

func TestCellAsInt64Empty(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	// Test reading an empty cell (should return 0, no error)
	val, err := cellAsInt64(f, SheetVertrag, "ZZ999")
	if err != nil {
		t.Fatalf("cellAsInt64() error for empty cell = %v", err)
	}

	if val != 0 {
		t.Errorf("cellAsInt64() for empty cell = %d, expected 0", val)
	}
}

func TestCellAsCents(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	// Test reading Summe (Z9)
	val, err := cellAsCents(f, SheetVertrag, "Z9")
	if err != nil {
		t.Fatalf("cellAsCents() error = %v", err)
	}

	// Should be positive
	if val <= 0 {
		t.Errorf("cellAsCents() Z9 = %d, expected positive value", val)
	}
}

func TestCellAsCentsEmpty(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	// Test reading an empty cell
	val, err := cellAsCents(f, SheetVertrag, "ZZ999")
	if err != nil {
		t.Fatalf("cellAsCents() error = %v", err)
	}

	// Should return zero
	if val != 0 {
		t.Errorf("cellAsCents() for empty cell = %d, expected 0", val)
	}
}

func TestVoucherPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"GB-12345678901-02", true},
		{"GB-00000000000-00", true},
		{"GB-99999999999-99", true},
		{"", false},
		{"GB-1234567890-02", false},   // too few digits
		{"GB-123456789012-02", false}, // too many digits
		{"GB-12345678901-2", false},   // suffix too short
		{"GB-12345678901-123", false}, // suffix too long
		{"AB-12345678901-02", false},  // wrong prefix
		{"GB12345678901-02", false},   // missing dash
		{"Abrechnungssumme", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := voucherPattern.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("voucherPattern.MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSenatsabrechnungOutputStruct(t *testing.T) {
	f, err := OpenSenatsabrechnung(testFile)
	if err != nil {
		t.Fatalf("OpenSenatsabrechnung() error = %v", err)
	}
	defer func() { _ = f.Close() }()

	eu, err := ParseEinrichtung(f)
	if err != nil {
		t.Fatalf("ParseEinrichtung() error = %v", err)
	}

	ar, err := ParseAbrechnung(f)
	if err != nil {
		t.Fatalf("ParseAbrechnung() error = %v", err)
	}

	v, err := ParseVertrag(f)
	if err != nil {
		t.Fatalf("ParseVertrag() error = %v", err)
	}

	output := SenatsabrechnungOutput{
		Einrichtung: eu,
		Abrechnung:  ar,
		Vertrag:     v,
	}

	if output.Einrichtung == nil {
		t.Error("SenatsabrechnungOutput.Einrichtung is nil")
	}
	if output.Abrechnung == nil {
		t.Error("SenatsabrechnungOutput.Abrechnung is nil")
	}
	if output.Vertrag == nil {
		t.Error("SenatsabrechnungOutput.Vertrag is nil")
	}
}

// --- Malformed-input tests using in-memory Excel files ---

func TestParseVertragMissingSheet(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	_, err := ParseVertrag(f)
	if err != nil {
		t.Fatalf("ParseVertrag() unexpected error for missing sheet: %v", err)
	}
}

func TestParseVertragEmptySheet(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	_, _ = f.NewSheet(SheetVertrag)

	v, err := ParseVertrag(f)
	if err != nil {
		t.Fatalf("ParseVertrag() error = %v", err)
	}
	if len(v.Kinder) != 0 {
		t.Errorf("ParseVertrag() got %d children, expected 0 for empty sheet", len(v.Kinder))
	}
}

func TestCellAsCentsRequiredInvalidValue(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheet := "Sheet1"
	_ = f.SetCellValue(sheet, "A1", "not-a-number")

	_, err := cellAsCentsRequired(f, sheet, "A1")
	if err == nil {
		t.Fatal("cellAsCentsRequired() expected error for non-numeric value")
	}
	if !strings.Contains(err.Error(), "Sheet1") {
		t.Errorf("error %q should contain sheet name", err.Error())
	}
	if !strings.Contains(err.Error(), "A1") {
		t.Errorf("error %q should contain cell reference", err.Error())
	}
}

func TestCellAsCentsInvalidValue(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	sheet := "Sheet1"
	_ = f.SetCellValue(sheet, "B2", "abc")

	_, err := cellAsCents(f, sheet, "B2")
	if err == nil {
		t.Fatal("cellAsCents() expected error for non-numeric value")
	}
	if !strings.Contains(err.Error(), "B2") {
		t.Errorf("error %q should contain cell reference", err.Error())
	}
}

func TestParseEinrichtungMissingHeader(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	_, _ = f.NewSheet(SheetEinrichtung)

	_, err := ParseEinrichtung(f)
	if err == nil {
		t.Fatal("ParseEinrichtung() expected error for missing headers")
	}
	if !strings.Contains(err.Error(), SheetEinrichtung) {
		t.Errorf("error %q should contain sheet name", err.Error())
	}
	if !strings.Contains(err.Error(), "Einrichtungsname") {
		t.Errorf("error %q should name the missing header", err.Error())
	}
}

func TestParseAbrechnungMissingHeader(t *testing.T) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	_, _ = f.NewSheet(SheetAbrechnung)

	_, err := ParseAbrechnung(f)
	if err == nil {
		t.Fatal("ParseAbrechnung() expected error for missing headers")
	}
	if !strings.Contains(err.Error(), SheetAbrechnung) {
		t.Errorf("error %q should contain sheet name", err.Error())
	}
}

func TestParseFromReader(t *testing.T) {
	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("os.Open() error = %v", err)
	}
	defer file.Close()

	output, err := ParseFromReader(file)
	if err != nil {
		t.Fatalf("ParseFromReader() error = %v", err)
	}

	if output == nil {
		t.Fatal("ParseFromReader() returned nil")
	}

	if output.Einrichtung == nil {
		t.Error("ParseFromReader() Einrichtung is nil")
	}
	if output.Abrechnung == nil {
		t.Error("ParseFromReader() Abrechnung is nil")
	}
	if output.Vertrag == nil {
		t.Error("ParseFromReader() Vertrag is nil")
	}

	// Verify same values as file-based parsing
	if output.Einrichtung.Name != "Kindertagesstätte Sternenstaub" {
		t.Errorf("ParseFromReader() Einrichtung.Name = %q, expected %q", output.Einrichtung.Name, "Kindertagesstätte Sternenstaub")
	}
	if len(output.Vertrag.Kinder) != 40 {
		t.Errorf("ParseFromReader() got %d children, expected 40", len(output.Vertrag.Kinder))
	}
	if output.Einrichtung.ZuschlagSPH != 165680 {
		t.Errorf("ParseFromReader() ZuschlagSPH = %d, expected %d", output.Einrichtung.ZuschlagSPH, 165680)
	}
}
