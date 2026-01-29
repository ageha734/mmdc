package domain

type MermaidDiagram struct {
	Code      string
	Index     int
	FullMatch string
}

type OutputFormat string

const (
	OutputFormatSVG OutputFormat = "svg"
	OutputFormatPNG OutputFormat = "png"
	OutputFormatPDF OutputFormat = "pdf"
)

func (f OutputFormat) String() string {
	return string(f)
}

func (f OutputFormat) IsValid() bool {
	return f == OutputFormatSVG || f == OutputFormatPNG || f == OutputFormatPDF
}
