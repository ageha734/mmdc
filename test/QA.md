# QA - E2E Test Cases

This document lists all E2E test cases that are covered by the automated tests in `test/e2e/`.

## Basic Options (basic_test.go)

### TestBasicRender_SVG

- **Description**: Basic SVG rendering with default options
- **Verification**: Input file read, SVG output, file existence check
- **Expected Behavior**: SVG file is generated successfully

### TestBasicRender_PNG

- **Description**: Basic PNG rendering
- **Verification**: Input file read, PNG output with valid magic bytes
- **Expected Behavior**: PNG file is generated successfully with correct format

### TestBasicRender_PDF

- **Description**: Basic PDF rendering
- **Verification**: Input file read, PDF output with valid magic bytes (%PDF)
- **Expected Behavior**: PDF file is generated successfully with correct format

### TestBasicRender_AutoDetectFormat

- **Description**: Auto-detection of output format from file extension
- **Verification**: Correct format output based on .svg, .png, .pdf extension
- **Expected Behavior**: Format is correctly determined from extension

### TestBasicRender_ComplexDiagram

- **Description**: Rendering of complex diagrams with subgraphs
- **Verification**: SVG output contains expected structure
- **Expected Behavior**: Complex diagram renders successfully

### TestBasicRender_Version

- **Description**: Version flag output
- **Verification**: Version string is displayed
- **Expected Behavior**: Version information is shown

### TestBasicRender_Help

- **Description**: Help flag output
- **Verification**: Help text contains all major options
- **Expected Behavior**: Help information is displayed

## Theme Options (theme_test.go)

### TestTheme_Default

- **Description**: Rendering with default theme
- **Verification**: SVG file is generated
- **Expected Behavior**: Default theme is applied

### TestTheme_Dark

- **Description**: Rendering with dark theme
- **Verification**: SVG file is generated
- **Expected Behavior**: Dark theme is applied

### TestTheme_Forest

- **Description**: Rendering with forest theme
- **Verification**: SVG file is generated
- **Expected Behavior**: Forest theme is applied

### TestTheme_Neutral

- **Description**: Rendering with neutral theme
- **Verification**: SVG file is generated
- **Expected Behavior**: Neutral theme is applied

### TestTheme_AllThemes

- **Description**: Rendering with all available themes
- **Verification**: SVG output for each theme
- **Expected Behavior**: All themes render successfully

### TestTheme_ComplexDiagram

- **Description**: Complex diagram with each theme
- **Verification**: SVG output for complex diagrams with each theme
- **Expected Behavior**: Complex diagrams render with all themes

## Dimension Options (dimensions_test.go)

### TestDimensions_Width

- **Description**: Custom width settings
- **Verification**: Rendering with various width values
- **Expected Behavior**: Width parameter is accepted

### TestDimensions_DefaultWidth

- **Description**: Default width (800px)
- **Verification**: Rendering without width parameter
- **Expected Behavior**: Default width is used

### TestDimensions_CustomDimensions

- **Description**: Custom width and height
- **Verification**: Rendering with -w and -H flags
- **Expected Behavior**: Custom dimensions are applied

### TestDimensions_Scale

- **Description**: Scale factor options
- **Verification**: Rendering with scale 1, 2, 3
- **Expected Behavior**: Scale parameter is accepted

### TestDimensions_PNGWithScale

- **Description**: PNG output with different scale factors
- **Verification**: PNG file size comparison between scales
- **Expected Behavior**: Scale affects PNG output

### TestDimensions_SmallDimensions

- **Description**: Small dimension values (200x150)
- **Verification**: Rendering with small dimensions
- **Expected Behavior**: Small dimensions work correctly

### TestDimensions_LargeDimensions

- **Description**: Large dimension values (3840x2160)
- **Verification**: Rendering with large dimensions
- **Expected Behavior**: Large dimensions work correctly

## Output Format (format_test.go)

### TestFormat_ExplicitSVG

- **Description**: Explicit SVG format with -e flag
- **Verification**: SVG output with format flag
- **Expected Behavior**: SVG is generated

### TestFormat_ExplicitPNG

- **Description**: Explicit PNG format with -e flag
- **Verification**: PNG output with valid magic bytes
- **Expected Behavior**: PNG is generated

### TestFormat_ExplicitPDF

- **Description**: Explicit PDF format with -e flag
- **Verification**: PDF output with valid magic bytes
- **Expected Behavior**: PDF is generated

### TestFormat_FormatOverridesExtension

- **Description**: Format flag overrides file extension
- **Verification**: -e flag takes precedence over extension
- **Expected Behavior**: Format flag determines output type

### TestFormat_ComplexDiagramAllFormats

- **Description**: Complex diagram in all formats
- **Verification**: SVG, PNG, PDF output validation
- **Expected Behavior**: All formats work for complex diagrams

### TestFormat_PDFWithFit

- **Description**: PDF output with pdfFit flag
- **Verification**: PDF generation with -f flag
- **Expected Behavior**: PDF is scaled to fit

## Markdown Processing (markdown_test.go)

### TestMarkdown_ProcessWithDiagrams

- **Description**: Markdown file with mermaid diagrams
- **Verification**: Output contains image references, no mermaid blocks
- **Expected Behavior**: Diagrams are replaced with images

### TestMarkdown_WithArtefactsDirectory

- **Description**: Artefacts directory creation
- **Verification**: Artefacts directory is created with SVG files
- **Expected Behavior**: Images are stored in artefacts directory

### TestMarkdown_WithPNGFormat

- **Description**: Markdown with PNG output
- **Verification**: PNG files in artefacts directory
- **Expected Behavior**: PNG images are generated

### TestMarkdown_PreservesNonMermaidContent

- **Description**: Non-mermaid content preservation
- **Verification**: Headers, text are preserved
- **Expected Behavior**: Non-mermaid content unchanged

### TestMarkdown_MultipleDiagrams

- **Description**: Multiple diagrams in one file
- **Verification**: All diagrams are rendered
- **Expected Behavior**: Multiple SVG files created

### TestMarkdown_ImageReferencesAreRelative

- **Description**: Relative path in image references
- **Verification**: Paths are relative, not absolute
- **Expected Behavior**: Relative paths are used

## Configuration (config_test.go)

### TestConfig_WithConfigFile

- **Description**: Mermaid config file
- **Verification**: Config file is applied
- **Expected Behavior**: Custom config is used

### TestConfig_WithCSSFile

- **Description**: Custom CSS file
- **Verification**: CSS file is applied
- **Expected Behavior**: Custom CSS is used

### TestConfig_WithBothConfigAndCSS

- **Description**: Config and CSS together
- **Verification**: Both files are applied
- **Expected Behavior**: Both config and CSS work together

### TestConfig_PuppeteerConfig

- **Description**: Puppeteer config file
- **Verification**: Puppeteer config handling
- **Expected Behavior**: Config is processed (may not be fully supported)

### TestConfig_NonExistentConfigFile

- **Description**: Non-existent config file error
- **Verification**: Error is returned
- **Expected Behavior**: Error for missing config

### TestConfig_NonExistentCSSFile

- **Description**: Non-existent CSS file handling
- **Verification**: Error or silent handling
- **Expected Behavior**: Appropriate error handling

### TestConfig_InvalidJSONConfig

- **Description**: Invalid JSON config file
- **Verification**: Parse error
- **Expected Behavior**: JSON parse error

### TestConfig_EmptyConfigFile

- **Description**: Empty JSON config file
- **Verification**: Empty config is accepted
- **Expected Behavior**: Empty config works

### TestConfig_ComplexDiagramWithConfig

- **Description**: Complex diagram with all config options
- **Verification**: Config, CSS, theme combined
- **Expected Behavior**: All options work together

## Advanced Options (advanced_test.go)

### TestAdvanced_SVGId

- **Description**: Custom SVG ID
- **Verification**: SVG contains custom ID
- **Expected Behavior**: Custom ID is applied

### TestAdvanced_BackgroundColor

- **Description**: Custom background colors
- **Verification**: Various background colors
- **Expected Behavior**: Background color is applied

### TestAdvanced_QuietMode

- **Description**: Quiet mode (-q)
- **Verification**: No log output
- **Expected Behavior**: Output is suppressed

### TestAdvanced_VerboseMode

- **Description**: Verbose mode (no -q)
- **Verification**: Log output present
- **Expected Behavior**: Output is shown

### TestAdvanced_PDFFit

- **Description**: PDF fit mode
- **Verification**: PDF with -f flag
- **Expected Behavior**: PDF is fit to content

### TestAdvanced_CombinedOptions

- **Description**: Multiple options combined
- **Verification**: Theme, dimensions, scale, colors, IDs
- **Expected Behavior**: All options work together

### TestAdvanced_LongSVGId

- **Description**: Long SVG ID string
- **Verification**: Long ID is preserved
- **Expected Behavior**: Long IDs work

### TestAdvanced_SpecialCharactersInSVGId

- **Description**: SVG ID with hyphens and numbers
- **Verification**: Valid XML ID characters
- **Expected Behavior**: Special characters work

### TestAdvanced_AllDiagramTypes

- **Description**: All Mermaid diagram types
- **Verification**: flowchart, sequence, class, state, er, gantt, pie, gitgraph
- **Expected Behavior**: All diagram types render

## Error Handling (error_test.go)

### TestError_NoInputFile

- **Description**: No input file specified
- **Verification**: Appropriate error or stdin read
- **Expected Behavior**: Error or stdin behavior

### TestError_NonExistentInputFile

- **Description**: Non-existent input file
- **Verification**: File not found error
- **Expected Behavior**: Error is returned

### TestError_InvalidMermaidSyntax

- **Description**: Invalid mermaid syntax
- **Verification**: Parse error
- **Expected Behavior**: Syntax error is returned

### TestError_InvalidOutputDirectory

- **Description**: Non-existent output directory
- **Verification**: Directory error
- **Expected Behavior**: Directory error is returned

### TestError_InvalidOutputFormat

- **Description**: Invalid output format
- **Verification**: Format error
- **Expected Behavior**: Format error is returned

### TestError_EmptyInputFile

- **Description**: Empty input file
- **Verification**: Error or empty handling
- **Expected Behavior**: Appropriate error handling

### TestError_WhitespaceOnlyInput

- **Description**: Whitespace-only input
- **Verification**: Error or empty handling
- **Expected Behavior**: Appropriate error handling

### TestError_ReadOnlyOutputDirectory

- **Description**: Read-only output directory
- **Verification**: Permission error
- **Expected Behavior**: Permission error is returned

### TestError_InvalidNegativeWidth

- **Description**: Negative width value
- **Verification**: Error or handling
- **Expected Behavior**: Appropriate handling

### TestError_InvalidNegativeHeight

- **Description**: Negative height value
- **Verification**: Error or handling
- **Expected Behavior**: Appropriate handling

### TestError_InvalidScale

- **Description**: Invalid scale value
- **Verification**: Error or handling
- **Expected Behavior**: Appropriate handling

### TestError_MarkdownWithInvalidDiagram

- **Description**: Markdown with invalid mermaid
- **Verification**: Render error
- **Expected Behavior**: Error is returned

### TestError_BinaryInputFile

- **Description**: Binary file as input
- **Verification**: Parse error
- **Expected Behavior**: Error is returned

### TestError_VeryLargeInputFile

- **Description**: Very large diagram
- **Verification**: Timeout or completion
- **Expected Behavior**: Handles large input

## Output Validation (output_validation_test.go)

### TestOutputValidation_SVGStructure

- **Description**: SVG XML structure
- **Verification**: Valid XML, root element is svg
- **Expected Behavior**: Valid SVG XML

### TestOutputValidation_SVGContainsExpectedElements

- **Description**: SVG contains expected elements
- **Verification**: svg tags present
- **Expected Behavior**: SVG structure is correct

### TestOutputValidation_SVGWithCustomId

- **Description**: Custom ID in SVG output
- **Verification**: Custom ID present in output
- **Expected Behavior**: Custom ID is in SVG

### TestOutputValidation_PNGMagicBytes

- **Description**: PNG file magic bytes
- **Verification**: 0x89 0x50 0x4E 0x47 0x0D 0x0A 0x1A 0x0A
- **Expected Behavior**: Valid PNG header

### TestOutputValidation_PNGMinimumSize

- **Description**: PNG file minimum size
- **Verification**: File size > 100 bytes
- **Expected Behavior**: PNG is not empty

### TestOutputValidation_PDFMagicBytes

- **Description**: PDF file magic bytes
- **Verification**: Starts with %PDF-
- **Expected Behavior**: Valid PDF header

### TestOutputValidation_PDFMinimumSize

- **Description**: PDF file minimum size
- **Verification**: File size > 100 bytes
- **Expected Behavior**: PDF is not empty

### TestOutputValidation_PDFContainsEOFMarker

- **Description**: PDF EOF marker
- **Verification**: Contains %%EOF
- **Expected Behavior**: Valid PDF trailer

### TestOutputValidation_ComplexDiagramSVG

- **Description**: Complex diagram SVG structure
- **Verification**: Contains group elements
- **Expected Behavior**: Complex SVG is valid

### TestOutputValidation_MarkdownOutput

- **Description**: Markdown output validation
- **Verification**: No mermaid blocks, has image refs, SVGs exist
- **Expected Behavior**: Markdown is processed correctly

### TestOutputValidation_FileSizeConsistency

- **Description**: Consistent output file size
- **Verification**: Multiple renders produce similar sizes
- **Expected Behavior**: Output is deterministic
