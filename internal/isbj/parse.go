package isbj

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type Einrichtung struct {
	Name            string `json:"name"`
	ElternBetreuung int    `json:"eltern_betreuung"`
	ElternEssen     int    `json:"eltern_essen"`
	ZuschlagQM      int    `json:"zuschlag_qm"`
	ZuschlagMSS     int    `json:"zuschlag_mss"`
	ZuschlagNDH     int    `json:"zuschlag_ndh"`
	ZuschlagSPH     int    `json:"zuschlag_sph"`
	Summe           int    `json:"summe"`
}

func (e Einrichtung) String() string {
	return fmt.Sprintf("Einrichtung: %s: ∑ %.2f € (QM: %.2f €, MSS: %.2f €, NDH: %.2f €, SPH: %.2f €)",
		e.Name,
		float64(e.Summe)/100,
		float64(e.ZuschlagQM)/100,
		float64(e.ZuschlagMSS)/100,
		float64(e.ZuschlagNDH)/100,
		float64(e.ZuschlagSPH)/100,
	)
}

type Abrechnung struct {
	VertragsBuchung  int `json:"vertrags_buchung"`
	KorrekturBuchung int `json:"korrektur_buchung"`
}

func (a Abrechnung) String() string {
	return fmt.Sprintf("Abrechnung: %.2f € (Vertrag), %.2f € (Korrektur)",
		float64(a.VertragsBuchung)/100,
		float64(a.KorrekturBuchung)/100,
	)
}

type Kind struct {
	Gutscheinnummer  string `json:"gutscheinnummer"`
	Name             string `json:"name"`
	Geburtsdatum     string `json:"geburtsdatum"`
	QM               string `json:"qm"`
	MSS              string `json:"mss"`
	HS               string `json:"hs"`
	SPH              string `json:"sph"`
	Registrierdatum  string `json:"registrierdatum"`
	Betreuungsumfang string `json:"betreuungsumfang"`
	Bezirk           int64  `json:"bezirk"`
	Basisentgeld     int    `json:"basisentgeld"`
	AbzugOM          int    `json:"abzug_om"`
	ElternBetreuung  int    `json:"eltern_betreuung"`
	ElternEssen      int    `json:"eltern_essen"`
	BuT              int    `json:"but"`
	AnteilBezirk     int    `json:"anteil_bezirk"`
	ZuschlagQM       int    `json:"zuschlag_qm"`
	ZuschlagMSS      int    `json:"zuschlag_mss"`
	ZuschlagNDH      int    `json:"zuschlag_ndh"`
	ZuschlagSPH      int    `json:"zuschlag_sph"`
	Summe            int    `json:"summe"`
}

func (k Kind) String() string {
	return fmt.Sprintf("%s (%s, geb. %s): %.2f €", k.Name, k.Gutscheinnummer, k.Geburtsdatum, float64(k.Summe)/100)
}

type Vertrag struct {
	Kinder []Kind `json:"kinder"`
}

func (v Vertrag) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Vertrag mit %d Kindern:\n", len(v.Kinder)))
	for i, k := range v.Kinder {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, k.String()))
	}
	return sb.String()
}

type SenatsabrechnungOutput struct {
	Einrichtung  *Einrichtung `json:"einrichtung"`
	Abrechnung   *Abrechnung  `json:"abrechnung"`
	Vertrag      *Vertrag     `json:"vertrag"`
	BillingMonth time.Time    `json:"billing_month"`
}

func cellAsString(f *excelize.File, sheet, cell string) string {
	val, err := f.GetCellValue(sheet, cell)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(val)
}

func cellAsInt64(f *excelize.File, sheet, cell string) (int64, error) {
	val := cellAsString(f, sheet, cell)
	if val == "" {
		return 0, nil
	}
	return strconv.ParseInt(val, 10, 64)
}

func cellAsCents(f *excelize.File, sheet, cell string) (int, error) {
	val := cellAsString(f, sheet, cell)
	if val == "" {
		return 0, nil
	}
	// numbers use US format from excelize - like 1,234.56
	val = strings.ReplaceAll(val, ",", "")
	f64, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("sheet %q cell %s: %w", sheet, cell, err)
	}
	return int(math.Round(f64 * 100)), nil
}

func cellAsCentsRequired(f *excelize.File, sheetName string, cell string) (int, error) {
	cellVal, err := f.GetCellValue(sheetName, cell)
	if err != nil {
		return 0, fmt.Errorf("sheet %q cell %s: %w", sheetName, cell, err)
	}
	// numbers use US format from excelize - like 1,234.56
	cellVal = strings.ReplaceAll(cellVal, ",", "")
	cellVal = strings.TrimSpace(cellVal)
	f64, err := strconv.ParseFloat(cellVal, 64)
	if err != nil {
		return 0, fmt.Errorf("sheet %q cell %s: %w", sheetName, cell, err)
	}
	return int(math.Round(f64 * 100)), nil
}

func OpenSenatsabrechnung(dateiname string) (*excelize.File, error) {
	f, err := excelize.OpenFile(dateiname)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// ParseFromReader parses a Senatsabrechnung from an io.Reader (e.g., HTTP upload body).
func ParseFromReader(r io.Reader) (*SenatsabrechnungOutput, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("opening excel from reader: %w", err)
	}
	defer func() { _ = f.Close() }()

	eu, err := ParseEinrichtung(f)
	if err != nil {
		return nil, fmt.Errorf("parsing Einrichtung: %w", err)
	}

	ar, err := ParseAbrechnung(f)
	if err != nil {
		return nil, fmt.Errorf("parsing Abrechnung: %w", err)
	}

	v, err := ParseVertrag(f)
	if err != nil {
		return nil, fmt.Errorf("parsing Vertrag: %w", err)
	}

	billingMonth, err := ParseBillingMonth(f)
	if err != nil {
		return nil, fmt.Errorf("parsing billing month: %w", err)
	}

	return &SenatsabrechnungOutput{
		Einrichtung:  eu,
		Abrechnung:   ar,
		Vertrag:      v,
		BillingMonth: billingMonth,
	}, nil
}

// ParseBillingMonth extracts the billing month from the Vertragsübersicht sheet.
// The layout is: "Trägernummer" label in one cell, with the billing month in the
// next column one row below (e.g., label at Y2, billing month at Z3).
// It also checks the cell directly below as a fallback.
func ParseBillingMonth(f *excelize.File) (time.Time, error) {
	col, row, err := getValueByHeaderName(f, SheetVertrag, "Trägernummer")
	if err != nil {
		return time.Time{}, fmt.Errorf("finding Trägernummer header: %w", err)
	}

	// Primary location: next column, one row below the label
	colNum, err := excelize.ColumnNameToNumber(col)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing column name %q: %w", col, err)
	}
	nextCol, err := excelize.ColumnNumberToName(colNum + 1)
	if err != nil {
		return time.Time{}, fmt.Errorf("computing next column: %w", err)
	}

	dateStr := cellAsString(f, SheetVertrag, fmt.Sprintf("%s%d", nextCol, row+1))

	// Fallback: directly below the label (some anonymized/variant files)
	if dateStr == "" {
		dateStr = cellAsString(f, SheetVertrag, fmt.Sprintf("%s%d", col, row+1))
	}

	if dateStr == "" {
		return time.Time{}, fmt.Errorf("billing month cell near Trägernummer is empty")
	}

	t, err := time.Parse("01/06", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing billing month %q: expected MM/YY format: %w", dateStr, err)
	}

	return t, nil
}

func getValueByHeaderName(f *excelize.File, sheetName string, header string) (string, int, error) {
	result, err := f.SearchSheet(sheetName, regexp.QuoteMeta(header), true)
	if err != nil {
		return "", 0, err
	}
	if len(result) != 1 {
		return "", 0, fmt.Errorf("sheet %q: header %q found %d times, expected 1", sheetName, header, len(result))
	}
	col, row, err := excelize.SplitCellName(result[0])
	if err != nil {
		return "", 0, err
	}
	return col, row, nil
}

func getEinrichtungName(f *excelize.File) (string, error) {
	col, row, err := getValueByHeaderName(f, SheetEinrichtung, "Einrichtungsname")
	if err != nil {
		return "", err
	}
	val, err := f.GetCellValue(SheetEinrichtung, fmt.Sprintf("%s%d", col, row+2))
	if err != nil {
		return "", err
	}
	return val, nil
}

func getEinrichtungZuschlagQM(f *excelize.File) (int, error) {
	col, row, err := getValueByHeaderName(f, SheetEinrichtung, "QM")
	if err != nil {
		return 0, err
	}
	coords := fmt.Sprintf("%s%d", col, row+1)
	return cellAsCentsRequired(f, SheetEinrichtung, coords)
}

func getEinrichtungZuschlagMSS(f *excelize.File) (int, error) {
	col, row, err := getValueByHeaderName(f, SheetEinrichtung, "MSS")
	if err != nil {
		return 0, err
	}
	coords := fmt.Sprintf("%s%d", col, row+1)
	return cellAsCentsRequired(f, SheetEinrichtung, coords)
}

func getEinrichtungZuschlagNDH(f *excelize.File) (int, error) {
	col, row, err := getValueByHeaderName(f, SheetEinrichtung, "ndH")
	if err != nil {
		return 0, err
	}
	coords := fmt.Sprintf("%s%d", col, row+1)
	return cellAsCentsRequired(f, SheetEinrichtung, coords)
}

func getEinrichtungZuschlagSPH(f *excelize.File) (int, error) {
	col, row, err := getValueByHeaderName(f, SheetEinrichtung, "SpH")
	if err != nil {
		return 0, err
	}
	coords := fmt.Sprintf("%s%d", col, row+1)
	return cellAsCentsRequired(f, SheetEinrichtung, coords)
}

func getEinrichtungSumme(f *excelize.File) (int, error) {
	col, row, err := getValueByHeaderName(f, SheetEinrichtung, "Abrechn.")
	if err != nil {
		return 0, err
	}
	coords := fmt.Sprintf("%s%d", col, row+2)
	return cellAsCentsRequired(f, SheetEinrichtung, coords)
}

func ParseEinrichtung(f *excelize.File) (*Einrichtung, error) {
	eu := &Einrichtung{}
	name, err := getEinrichtungName(f)
	if err != nil {
		return nil, err
	}
	eu.Name = name

	zuschlagQM, err := getEinrichtungZuschlagQM(f)
	if err != nil {
		return nil, err
	}
	eu.ZuschlagQM = zuschlagQM

	zuschlagMSS, err := getEinrichtungZuschlagMSS(f)
	if err != nil {
		return nil, err
	}
	eu.ZuschlagMSS = zuschlagMSS

	zuschlagNDH, err := getEinrichtungZuschlagNDH(f)
	if err != nil {
		return nil, err
	}
	eu.ZuschlagNDH = zuschlagNDH

	zuschlagSPH, err := getEinrichtungZuschlagSPH(f)
	if err != nil {
		return nil, err
	}
	eu.ZuschlagSPH = zuschlagSPH

	summe, err := getEinrichtungSumme(f)
	if err != nil {
		return nil, err
	}
	eu.Summe = summe

	return eu, nil
}

func ParseAbrechnung(f *excelize.File) (*Abrechnung, error) {
	ar := &Abrechnung{}
	// Get the row where the sum is
	_, rowSum, err := getValueByHeaderName(f, SheetAbrechnung, "Summe:")
	if err != nil {
		return nil, err
	}

	// Get Vertragsbuchung
	col, _, err := getValueByHeaderName(f, SheetAbrechnung, "Vertragsbuchung")
	if err != nil {
		return nil, err
	}
	coords := fmt.Sprintf("%s%d", col, rowSum)
	data, err := cellAsCentsRequired(f, SheetAbrechnung, coords)
	if err != nil {
		return nil, err
	}
	ar.VertragsBuchung = data

	// Get Korrekturbuchung
	col, _, err = getValueByHeaderName(f, SheetAbrechnung, "Korrekturbuchung")
	if err != nil {
		return nil, err
	}
	coords = fmt.Sprintf("%s%d", col, rowSum)
	data, err = cellAsCentsRequired(f, SheetAbrechnung, coords)
	if err != nil {
		return nil, err
	}
	ar.KorrekturBuchung = data
	return ar, nil
}

// voucherPattern matches voucher numbers like GB-XXXXXXXXXXX-XX
var voucherPattern = regexp.MustCompile(`^GB-\d{11}-\d{2}$`)

const (
	SheetVertrag     = "Vertragsübersicht"
	SheetAbrechnung  = "Abrechnungsübersicht"
	SheetEinrichtung = "Einrichtungsübersicht"
	VertragStartRow  = 9
)

func ParseKind(f *excelize.File, row int) (*Kind, error) {
	kind := &Kind{}

	// Column E: GutscheinNr.
	kind.Gutscheinnummer = cellAsString(f, SheetVertrag, fmt.Sprintf("E%d", row))
	// Column F: Name des Kindes
	kind.Name = cellAsString(f, SheetVertrag, fmt.Sprintf("F%d", row))
	// Column G: Geb.-datum (MM.YY format)
	kind.Geburtsdatum = cellAsString(f, SheetVertrag, fmt.Sprintf("G%d", row))
	// Column I: QM ("ja"/"nein")
	kind.QM = cellAsString(f, SheetVertrag, fmt.Sprintf("I%d", row))
	// Column J: MSS ("ja"/"nein")
	kind.MSS = cellAsString(f, SheetVertrag, fmt.Sprintf("J%d", row))
	// Column K: Hs ("D"/"N"/etc)
	kind.HS = cellAsString(f, SheetVertrag, fmt.Sprintf("K%d", row))
	// Column L: SpH ("N"/etc)
	kind.SPH = cellAsString(f, SheetVertrag, fmt.Sprintf("L%d", row))
	// Column M: Registr. ab (MM.YY)
	kind.Registrierdatum = cellAsString(f, SheetVertrag, fmt.Sprintf("M%d", row))
	// Column N: Betreu.-umfang
	kind.Betreuungsumfang = cellAsString(f, SheetVertrag, fmt.Sprintf("N%d", row))

	// Column O: Be-zirk (int)
	bezirk, err := cellAsInt64(f, SheetVertrag, fmt.Sprintf("O%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing Bezirk at row %d: %w", row, err)
	}
	if bezirk < 1 || bezirk > 12 {
		return nil, fmt.Errorf("row %d: Bezirk %d is not a valid Berlin district (1-12)", row, bezirk)
	}
	kind.Bezirk = bezirk

	switch kind.Betreuungsumfang {
	case "teilzeit", "ganztags", "erweitert":
		// valid
	default:
		return nil, fmt.Errorf("row %d: Betreuungsumfang %q is not valid (expected \"teilzeit\", \"ganztags\", or \"erweitert\")", row, kind.Betreuungsumfang)
	}

	// Column P: Basis-entgelt (decimal)
	kind.Basisentgeld, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("P%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing Basisentgeld at row %d: %w", row, err)
	}
	// Column Q: Abzug o.M. (decimal)
	kind.AbzugOM, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("Q%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing AbzugOM at row %d: %w", row, err)
	}
	// Column R: Anteil Eltern (Betreuung)
	kind.ElternBetreuung, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("R%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing ElternBetreuung at row %d: %w", row, err)
	}
	// Column S: Anteil Eltern (Essen)
	kind.ElternEssen, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("S%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing ElternEssen at row %d: %w", row, err)
	}
	// Column T: BuT
	kind.BuT, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("T%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing BuT at row %d: %w", row, err)
	}
	// Column U: Anteil Bezirk
	kind.AnteilBezirk, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("U%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing AnteilBezirk at row %d: %w", row, err)
	}
	// Column V: Zuschlag QM
	kind.ZuschlagQM, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("V%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing ZuschlagQM at row %d: %w", row, err)
	}
	// Column W: Zuschlag MSS
	kind.ZuschlagMSS, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("W%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing ZuschlagMSS at row %d: %w", row, err)
	}
	// Column X: Zuschlag ndH
	kind.ZuschlagNDH, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("X%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing ZuschlagNDH at row %d: %w", row, err)
	}
	// Column Y: SpH
	kind.ZuschlagSPH, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("Y%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing ZuschlagSPH at row %d: %w", row, err)
	}
	// Column Z: Abrechn.-summe
	kind.Summe, err = cellAsCents(f, SheetVertrag, fmt.Sprintf("Z%d", row))
	if err != nil {
		return nil, fmt.Errorf("parsing Summe at row %d: %w", row, err)
	}

	return kind, nil
}

func ParseVertrag(f *excelize.File) (*Vertrag, error) {
	vertrag := &Vertrag{
		Kinder: []Kind{},
	}

	row := VertragStartRow

	for {
		// Check column E for voucher number
		voucher := cellAsString(f, SheetVertrag, fmt.Sprintf("E%d", row))

		// Stop if empty or contains "Abrechnungssumme"
		if voucher == "" || strings.Contains(voucher, "Abrechnungssumme") {
			break
		}

		// Validate voucher pattern
		if !voucherPattern.MatchString(voucher) {
			break
		}

		kind, err := ParseKind(f, row)
		if err != nil {
			return nil, fmt.Errorf("parsing kind at row %d: %w", row, err)
		}

		vertrag.Kinder = append(vertrag.Kinder, *kind)
		row++
	}

	return vertrag, nil
}
