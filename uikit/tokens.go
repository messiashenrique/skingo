package uikit

// DesignTokens defines the central design token system for Skingo UI components.
// All color, spacing, typography, and layout values should be sourced from here.
type DesignTokens struct {
	// Color Tokens
	Colors ColorTokens

	// Spacing Tokens
	Spacing SpacingTokens

	// Typography Tokens
	Typography TypographyTokens

	// Border Tokens
	Border BorderTokens

	// Shadow Tokens
	Shadows ShadowTokens

	// Component-specific tokens
	Components ComponentTokens
}

// ColorTokens defines the color palette
type ColorTokens struct {
	// Semantic Colors
	Primary   string // Primary action color
	Secondary string // Secondary action color
	Success   string // Success/positive state
	Warning   string // Warning state
	Error     string // Error/danger state
	Info      string // Information state

	// Neutral colors
	Background string // Main background
	Surface    string // Card/container background
	Border     string // Border color
	Text       string // Primary text color
	TextMuted  string // Muted/secondary text color

	// Variants for each semantic color
	PrimaryLight   string // Light variant
	SecondaryLight string
	SuccessLight   string
	WarningLight   string
	ErrorLight     string
	InfoLight      string

	// Outline variants
	PrimaryOutline   string
	SecondaryOutline string
	SuccessOutline   string
	WarningOutline   string
	ErrorOutline     string
	InfoOutline      string
}

// SpacingTokens defines spacing scale
type SpacingTokens struct {
	XS  string // Extra small (0.25rem)
	SM  string // Small (0.5rem)
	MD  string // Medium (1rem)
	LG  string // Large (1.5rem)
	XL  string // Extra large (2rem)
	XXL string // 2x extra large (3rem)
}

// TypographyTokens defines typography settings
type TypographyTokens struct {
	FontFamily string
	FontSize   FontSizeTokens
	LineHeight LineHeightTokens
	FontWeight FontWeightTokens
}

// FontSizeTokens defines font size scale
type FontSizeTokens struct {
	XS  string // 0.75rem
	SM  string // 0.875rem
	MD  string // 1rem
	LG  string // 1.125rem
	XL  string // 1.25rem
	XXL string // 1.5rem
}

// LineHeightTokens defines line height values
type LineHeightTokens struct {
	Tight   string // 1.2
	Normal  string // 1.5
	Relaxed string // 1.75
}

// FontWeightTokens defines font weight values
type FontWeightTokens struct {
	Regular  string // 400
	Medium   string // 500
	Semibold string // 600
	Bold     string // 700
}

// BorderTokens defines border styling
type BorderTokens struct {
	Radius RadiusTokens
	Width  WidthTokens
}

// RadiusTokens defines border radius values
type RadiusTokens struct {
	None string // 0
	SM   string // 0.25rem
	MD   string // 0.5rem
	LG   string // 0.75rem
	XL   string // 1rem
}

// WidthTokens defines border width values
type WidthTokens struct {
	Thin  string // 1px
	Base  string // 2px
	Thick string // 3px
}

// ShadowTokens defines shadow values
type ShadowTokens struct {
	SM    string
	MD    string
	LG    string
	XL    string
	Inset string
}

// ComponentTokens defines component-specific token overrides
type ComponentTokens struct {
	Button ButtonTokens
	Input  InputTokens
	Card   CardTokens
	Badge  BadgeTokens
	Info   InfoTokens
}

// ButtonTokens defines button-specific tokens
type ButtonTokens struct {
	PaddingVertical    string
	PaddingHorizontal  string
	MinHeight          string
	FontSize           string
	BorderRadius       string
	TransitionDuration string
}

// InputTokens defines input-specific tokens
type InputTokens struct {
	PaddingVertical   string
	PaddingHorizontal string
	MinHeight         string
	FontSize          string
	BorderRadius      string
	BorderWidth       string
}

// CardTokens defines card-specific tokens
type CardTokens struct {
	Padding      string
	BorderRadius string
	BorderWidth  string
	ShadowColor  string
}

// BadgeTokens defines badge-specific tokens
type BadgeTokens struct {
	PaddingVertical   string
	PaddingHorizontal string
	FontSize          string
	BorderRadius      string
}

// InfoTokens defines info/alert-specific tokens
type InfoTokens struct {
	Padding      string
	BorderRadius string
	BorderLeft   string
	BorderWidth  string
	IconSize     string
}

// LightTheme returns the default light theme tokens
func LightTheme() *DesignTokens {
	return &DesignTokens{
		Colors: ColorTokens{
			Primary:          "#0f172a", // Slate-900 (shadcn-ui inspired)
			Secondary:        "#334155", // Slate-700
			Success:          "#059669", // Emerald-600
			Warning:          "#d97706", // Amber-600
			Error:            "#e11d48", // Rose-600
			Info:             "#0284c7", // Sky-600
			Background:       "#ffffff",
			Surface:          "#f9fafb",
			Border:           "#e5e7eb",
			Text:             "#0f172a",
			TextMuted:        "#64748b",
			PrimaryLight:     "#f1f5f9", // Slate-100
			SecondaryLight:   "#cbd5e1", // Slate-200
			SuccessLight:     "#d1fae5", // Emerald-100
			WarningLight:     "#fed7aa", // Amber-100
			ErrorLight:       "#ffe4e6", // Rose-100
			InfoLight:        "#e0f2fe", // Sky-100
			PrimaryOutline:   "#0f172a",
			SecondaryOutline: "#334155",
			SuccessOutline:   "#059669",
			WarningOutline:   "#d97706",
			ErrorOutline:     "#e11d48",
			InfoOutline:      "#0284c7",
		},
		Spacing: SpacingTokens{
			XS:  "0.25rem",
			SM:  "0.5rem",
			MD:  "1rem",
			LG:  "1.5rem",
			XL:  "2rem",
			XXL: "3rem",
		},
		Typography: TypographyTokens{
			FontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif",
			FontSize: FontSizeTokens{
				XS:  "0.75rem",
				SM:  "0.875rem",
				MD:  "1rem",
				LG:  "1.125rem",
				XL:  "1.25rem",
				XXL: "1.5rem",
			},
			LineHeight: LineHeightTokens{
				Tight:   "1.2",
				Normal:  "1.5",
				Relaxed: "1.75",
			},
			FontWeight: FontWeightTokens{
				Regular:  "400",
				Medium:   "500",
				Semibold: "600",
				Bold:     "700",
			},
		},
		Border: BorderTokens{
			Radius: RadiusTokens{
				None: "0",
				SM:   "0.25rem",
				MD:   "0.5rem",
				LG:   "0.75rem",
				XL:   "1rem",
			},
			Width: WidthTokens{
				Thin:  "1px",
				Base:  "2px",
				Thick: "3px",
			},
		},
		Shadows: ShadowTokens{
			SM:    "0 1px 2px 0 rgba(0, 0, 0, 0.05)",
			MD:    "0 4px 6px -1px rgba(0, 0, 0, 0.1)",
			LG:    "0 10px 15px -3px rgba(0, 0, 0, 0.1)",
			XL:    "0 20px 25px -5px rgba(0, 0, 0, 0.1)",
			Inset: "inset 0 2px 4px 0 rgba(0, 0, 0, 0.05)",
		},
		Components: ComponentTokens{
			Button: ButtonTokens{
				PaddingVertical:    "0.5rem",
				PaddingHorizontal:  "1rem",
				MinHeight:          "2.5rem",
				FontSize:           "1rem",
				BorderRadius:       "0.5rem",
				TransitionDuration: "150ms",
			},
			Input: InputTokens{
				PaddingVertical:   "0.5rem",
				PaddingHorizontal: "0.75rem",
				MinHeight:         "2.5rem",
				FontSize:          "1rem",
				BorderRadius:      "0.375rem",
				BorderWidth:       "1px",
			},
			Card: CardTokens{
				Padding:      "1.5rem",
				BorderRadius: "0.5rem",
				BorderWidth:  "1px",
				ShadowColor:  "0 1px 3px 0 rgba(0, 0, 0, 0.1)",
			},
			Badge: BadgeTokens{
				PaddingVertical:   "0.25rem",
				PaddingHorizontal: "0.75rem",
				FontSize:          "0.875rem",
				BorderRadius:      "9999px",
			},
			Info: InfoTokens{
				Padding:      "1rem",
				BorderRadius: "0.375rem",
				BorderLeft:   "4px",
				BorderWidth:  "1px",
				IconSize:     "1.25rem",
			},
		},
	}
}

// DarkTheme returns a dark mode theme tokens
func DarkTheme() *DesignTokens {
	tokens := LightTheme()
	// Override colors for dark mode (shadcn-ui inspired)
	tokens.Colors = ColorTokens{
		Primary:          "#f1f5f9", // Slate-100
		Secondary:        "#cbd5e1", // Slate-200
		Success:          "#6ee7b7", // Emerald-400
		Warning:          "#fcd34d", // Amber-300
		Error:            "#fb7185", // Rose-400
		Info:             "#38bdf8", // Sky-400
		Background:       "#020617", // Slate-950
		Surface:          "#0f172a", // Slate-900
		Border:           "#1e293b", // Slate-800
		Text:             "#f1f5f9", // Slate-100
		TextMuted:        "#94a3b8", // Slate-400
		PrimaryLight:     "#1e293b", // Slate-800
		SecondaryLight:   "#334155", // Slate-700
		SuccessLight:     "#064e3b", // Emerald-900
		WarningLight:     "#78350f", // Amber-900
		ErrorLight:       "#881337", // Rose-900
		InfoLight:        "#0c4a6e", // Sky-900
		PrimaryOutline:   "#e2e8f0",
		SecondaryOutline: "#cbd5e1",
		SuccessOutline:   "#6ee7b7",
		WarningOutline:   "#fcd34d",
		ErrorOutline:     "#fb7185",
		InfoOutline:      "#38bdf8",
	}
	return tokens
}

// GetTheme returns the design tokens for the given theme name
func GetTheme(name string) *DesignTokens {
	switch name {
	case "dark":
		return DarkTheme()
	case "light":
		fallthrough
	default:
		return LightTheme()
	}
}
