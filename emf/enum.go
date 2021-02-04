package emf

// FormatSignature
const ENHMETA_SIGNATURE = 0x464D4520

// RecordType
const (
	EMR_HEADER                  = uint32(0x00000001)
	EMR_POLYBEZIER              = uint32(0x00000002)
	EMR_POLYGON                 = uint32(0x00000003)
	EMR_POLYLINE                = uint32(0x00000004)
	EMR_POLYBEZIERTO            = uint32(0x00000005)
	EMR_POLYLINETO              = uint32(0x00000006)
	EMR_POLYPOLYLINE            = uint32(0x00000007)
	EMR_POLYPOLYGON             = uint32(0x00000008)
	EMR_SETWINDOWEXTEX          = uint32(0x00000009)
	EMR_SETWINDOWORGEX          = uint32(0x0000000A)
	EMR_SETVIEWPORTEXTEX        = uint32(0x0000000B)
	EMR_SETVIEWPORTORGEX        = uint32(0x0000000C)
	EMR_SETBRUSHORGEX           = uint32(0x0000000D)
	EMR_EOF                     = uint32(0x0000000E)
	EMR_SETPIXELV               = uint32(0x0000000F)
	EMR_SETMAPPERFLAGS          = uint32(0x00000010)
	EMR_SETMAPMODE              = uint32(0x00000011)
	EMR_SETBKMODE               = uint32(0x00000012)
	EMR_SETPOLYFILLMODE         = uint32(0x00000013)
	EMR_SETROP2                 = uint32(0x00000014)
	EMR_SETSTRETCHBLTMODE       = uint32(0x00000015)
	EMR_SETTEXTALIGN            = uint32(0x00000016)
	EMR_SETCOLORADJUSTMENT      = uint32(0x00000017)
	EMR_SETTEXTCOLOR            = uint32(0x00000018)
	EMR_SETBKCOLOR              = uint32(0x00000019)
	EMR_OFFSETCLIPRGN           = uint32(0x0000001A)
	EMR_MOVETOEX                = uint32(0x0000001B)
	EMR_SETMETARGN              = uint32(0x0000001C)
	EMR_EXCLUDECLIPRECT         = uint32(0x0000001D)
	EMR_INTERSECTCLIPRECT       = uint32(0x0000001E)
	EMR_SCALEVIEWPORTEXTEX      = uint32(0x0000001F)
	EMR_SCALEWINDOWEXTEX        = uint32(0x00000020)
	EMR_SAVEDC                  = uint32(0x00000021)
	EMR_RESTOREDC               = uint32(0x00000022)
	EMR_SETWORLDTRANSFORM       = uint32(0x00000023)
	EMR_MODIFYWORLDTRANSFORM    = uint32(0x00000024)
	EMR_SELECTOBJECT            = uint32(0x00000025)
	EMR_CREATEPEN               = uint32(0x00000026)
	EMR_CREATEBRUSHINDIRECT     = uint32(0x00000027)
	EMR_DELETEOBJECT            = uint32(0x00000028)
	EMR_ANGLEARC                = uint32(0x00000029)
	EMR_ELLIPSE                 = uint32(0x0000002A)
	EMR_RECTANGLE               = uint32(0x0000002B)
	EMR_ROUNDRECT               = uint32(0x0000002C)
	EMR_ARC                     = uint32(0x0000002D)
	EMR_CHORD                   = uint32(0x0000002E)
	EMR_PIE                     = uint32(0x0000002F)
	EMR_SELECTPALETTE           = uint32(0x00000030)
	EMR_CREATEPALETTE           = uint32(0x00000031)
	EMR_SETPALETTEENTRIES       = uint32(0x00000032)
	EMR_RESIZEPALETTE           = uint32(0x00000033)
	EMR_REALIZEPALETTE          = uint32(0x00000034)
	EMR_EXTFLOODFILL            = uint32(0x00000035)
	EMR_LINETO                  = uint32(0x00000036)
	EMR_ARCTO                   = uint32(0x00000037)
	EMR_POLYDRAW                = uint32(0x00000038)
	EMR_SETARCDIRECTION         = uint32(0x00000039)
	EMR_SETMITERLIMIT           = uint32(0x0000003A)
	EMR_BEGINPATH               = uint32(0x0000003B)
	EMR_ENDPATH                 = uint32(0x0000003C)
	EMR_CLOSEFIGURE             = uint32(0x0000003D)
	EMR_FILLPATH                = uint32(0x0000003E)
	EMR_STROKEANDFILLPATH       = uint32(0x0000003F)
	EMR_STROKEPATH              = uint32(0x00000040)
	EMR_FLATTENPATH             = uint32(0x00000041)
	EMR_WIDENPATH               = uint32(0x00000042)
	EMR_SELECTCLIPPATH          = uint32(0x00000043)
	EMR_ABORTPATH               = uint32(0x00000044)
	EMR_COMMENT                 = uint32(0x00000046)
	EMR_FILLRGN                 = uint32(0x00000047)
	EMR_FRAMERGN                = uint32(0x00000048)
	EMR_INVERTRGN               = uint32(0x00000049)
	EMR_PAINTRGN                = uint32(0x0000004A)
	EMR_EXTSELECTCLIPRGN        = uint32(0x0000004B)
	EMR_BITBLT                  = uint32(0x0000004C)
	EMR_STRETCHBLT              = uint32(0x0000004D)
	EMR_MASKBLT                 = uint32(0x0000004E)
	EMR_PLGBLT                  = uint32(0x0000004F)
	EMR_SETDIBITSTODEVICE       = uint32(0x00000050)
	EMR_STRETCHDIBITS           = uint32(0x00000051)
	EMR_EXTCREATEFONTINDIRECTW  = uint32(0x00000052)
	EMR_EXTTEXTOUTA             = uint32(0x00000053)
	EMR_EXTTEXTOUTW             = uint32(0x00000054)
	EMR_POLYBEZIER16            = uint32(0x00000055)
	EMR_POLYGON16               = uint32(0x00000056)
	EMR_POLYLINE16              = uint32(0x00000057)
	EMR_POLYBEZIERTO16          = uint32(0x00000058)
	EMR_POLYLINETO16            = uint32(0x00000059)
	EMR_POLYPOLYLINE16          = uint32(0x0000005A)
	EMR_POLYPOLYGON16           = uint32(0x0000005B)
	EMR_POLYDRAW16              = uint32(0x0000005C)
	EMR_CREATEMONOBRUSH         = uint32(0x0000005D)
	EMR_CREATEDIBPATTERNBRUSHPT = uint32(0x0000005E)
	EMR_EXTCREATEPEN            = uint32(0x0000005F)
	EMR_POLYTEXTOUTA            = uint32(0x00000060)
	EMR_POLYTEXTOUTW            = uint32(0x00000061)
	EMR_SETICMMODE              = uint32(0x00000062)
	EMR_CREATECOLORSPACE        = uint32(0x00000063)
	EMR_SETCOLORSPACE           = uint32(0x00000064)
	EMR_DELETECOLORSPACE        = uint32(0x00000065)
	EMR_GLSRECORD               = uint32(0x00000066)
	EMR_GLSBOUNDEDRECORD        = uint32(0x00000067)
	EMR_PIXELFORMAT             = uint32(0x00000068)
	EMR_DRAWESCAPE              = uint32(0x00000069)
	EMR_EXTESCAPE               = uint32(0x0000006A)
	EMR_SMALLTEXTOUT            = uint32(0x0000006C)
	EMR_FORCEUFIMAPPING         = uint32(0x0000006D)
	EMR_NAMEDESCAPE             = uint32(0x0000006E)
	EMR_COLORCORRECTPALETTE     = uint32(0x0000006F)
	EMR_SETICMPROFILEA          = uint32(0x00000070)
	EMR_SETICMPROFILEW          = uint32(0x00000071)
	EMR_ALPHABLEND              = uint32(0x00000072)
	EMR_SETLAYOUT               = uint32(0x00000073)
	EMR_TRANSPARENTBLT          = uint32(0x00000074)
	EMR_GRADIENTFILL            = uint32(0x00000076)
	EMR_SETLINKEDUFIS           = uint32(0x00000077)
	EMR_SETTEXTJUSTIFICATION    = uint32(0x00000078)
	EMR_COLORMATCHTOTARGETW     = uint32(0x00000079)
	EMR_CREATECOLORSPACEW       = uint32(0x0000007A)
)

// StockObject
const (
	WHITE_BRUSH         = uint32(0x80000000)
	LTGRAY_BRUSH        = uint32(0x80000001)
	GRAY_BRUSH          = uint32(0x80000002)
	DKGRAY_BRUSH        = uint32(0x80000003)
	BLACK_BRUSH         = uint32(0x80000004)
	NULL_BRUSH          = uint32(0x80000005)
	WHITE_PEN           = uint32(0x80000006)
	BLACK_PEN           = uint32(0x80000007)
	NULL_PEN            = uint32(0x80000008)
	OEM_FIXED_FONT      = uint32(0x8000000A)
	ANSI_FIXED_FONT     = uint32(0x8000000B)
	ANSI_VAR_FONT       = uint32(0x8000000C)
	SYSTEM_FONT         = uint32(0x8000000D)
	DEVICE_DEFAULT_FONT = uint32(0x8000000E)
	DEFAULT_PALETTE     = uint32(0x8000000F)
	SYSTEM_FIXED_FONT   = uint32(0x80000010)
	DEFAULT_GUI_FONT    = uint32(0x80000011)
	DC_BRUSH            = uint32(0x80000012)
	DC_PEN              = uint32(0x80000013)
)

// BitCount
const (
	BI_BITCOUNT_0 = 0x0000
	BI_BITCOUNT_1 = 0x0001
	BI_BITCOUNT_2 = 0x0004
	BI_BITCOUNT_3 = 0x0008
	BI_BITCOUNT_4 = 0x0010
	BI_BITCOUNT_5 = 0x0018
	BI_BITCOUNT_6 = 0x0020
)

// BackgroundMode
const (
	TRANSPARENT = 0x0001
	OPAQUE      = 0x0002
)

// PenStyle
const (
	PS_COSMETIC      = uint32(0x00000000)
	PS_ENDCAP_ROUND  = uint32(0x00000000)
	PS_JOIN_ROUND    = uint32(0x00000000)
	PS_SOLID         = uint32(0x00000000)
	PS_DASH          = uint32(0x00000001)
	PS_DOT           = uint32(0x00000002)
	PS_DASHDOT       = uint32(0x00000003)
	PS_DASHDOTDOT    = uint32(0x00000004)
	PS_NULL          = uint32(0x00000005)
	PS_INSIDEFRAME   = uint32(0x00000006)
	PS_USERSTYLE     = uint32(0x00000007)
	PS_ALTERNATE     = uint32(0x00000008)
	PS_ENDCAP_SQUARE = uint32(0x00000100)
	PS_ENDCAP_FLAT   = uint32(0x00000200)
	PS_JOIN_BEVEL    = uint32(0x00001000)
	PS_JOIN_MITER    = uint32(0x00002000)
	PS_GEOMETRIC     = uint32(0x00010000)
)

//ModifyWorldTransformMode
const (
	MWT_IDENTITY      = 0x01
	MWT_LEFTMULTIPLY  = 0x02
	MWT_RIGHTMULTIPLY = 0x03
	MWT_SET           = 0x04
)

// PolygonFillMode
const (
	ALTERNATE = 0x01
	WINDING   = 0x02
)

// ExtTextOutOptions
const (
	ETO_OPAQUE            = uint32(0x00000002)
	ETO_CLIPPED           = uint32(0x00000004)
	ETO_GLYPH_INDEX       = uint32(0x00000010)
	ETO_RTLREADING        = uint32(0x00000080)
	ETO_NO_RECT           = uint32(0x00000100)
	ETO_SMALL_CHARS       = uint32(0x00000200)
	ETO_NUMERICSLOCAL     = uint32(0x00000400)
	ETO_NUMERICSLATIN     = uint32(0x00000800)
	ETO_IGNORELANGUAGE    = uint32(0x00001000)
	ETO_PDY               = uint32(0x00002000)
	ETO_REVERSE_INDEX_MAP = uint32(0x00010000)
)

// Compression
const (
	BI_RGB       = 0x0000
	BI_RLE8      = 0x0001
	BI_RLE4      = 0x0002
	BI_BITFIELDS = 0x0003
	BI_JPEG      = 0x0004
	BI_PNG       = 0x0005
	BI_CMYK      = 0x000B
	BI_CMYKRLE8  = 0x000C
	BI_CMYKRLE4  = 0x000D
)

// MapMode
const (
	MM_TEXT        = 0x01
	MM_LOMETRIC    = 0x02
	MM_HIMETRIC    = 0x03
	MM_LOENGLISH   = 0x04
	MM_HIENGLISH   = 0x05
	MM_TWIPS       = 0x06
	MM_ISOTROPIC   = 0x07
	MM_ANISOTROPIC = 0x08
)

// DIBColors
const (
	DIB_RGB_COLORS  = 0x00
	DIB_PAL_COLORS  = 0x01
	DIB_PAL_INDICES = 0x02
)

// StretchMode
const (
	STRETCH_ANDSCANS    = 0x01
	STRETCH_ORSCANS     = 0x02
	STRETCH_DELETESCANS = 0x03
	STRETCH_HALFTONE    = 0x04
)

// RegionMode
const (
	RGN_AND  = 0x01
	RGN_OR   = 0x02
	RGN_XOR  = 0x03
	RGN_DIFF = 0x04
	RGN_COPY = 0x05
)

// BrushStyle
const (
	BS_SOLID         = 0x0000
	BS_NULL          = 0x0001
	BS_HATCHED       = 0x0002
	BS_PATTERN       = 0x0003
	BS_INDEXED       = 0x0004
	BS_DIBPATTERN    = 0x0005
	BS_DIBPATTERNPT  = 0x0006
	BS_PATTERN8X8    = 0x0007
	BS_DIBPATTERN8X8 = 0x0008
	BS_MONOPATTERN   = 0x0009
)

// FamilyFont
const (
	FF_DONTCARE   = 0x00
	FF_ROMAN      = 0x01
	FF_SWISS      = 0x02
	FF_MODERN     = 0x03
	FF_SCRIPT     = 0x04
	FF_DECORATIVE = 0x05
)

// PitchFont
const (
	DEFAULT_PITCH  = 0
	FIXED_PITCH    = 1
	VARIABLE_PITCH = 2
)

// CharacterSet
const (
	ANSI_CHARSET        = uint32(0x00000000)
	DEFAULT_CHARSET     = uint32(0x00000001)
	SYMBOL_CHARSET      = uint32(0x00000002)
	MAC_CHARSET         = uint32(0x0000004D)
	SHIFTJIS_CHARSET    = uint32(0x00000080)
	HANGUL_CHARSET      = uint32(0x00000081)
	JOHAB_CHARSET       = uint32(0x00000082)
	GB2312_CHARSET      = uint32(0x00000086)
	CHINESEBIG5_CHARSET = uint32(0x00000088)
	GREEK_CHARSET       = uint32(0x000000A1)
	TURKISH_CHARSET     = uint32(0x000000A2)
	VIETNAMESE_CHARSET  = uint32(0x000000A3)
	HEBREW_CHARSET      = uint32(0x000000B1)
	ARABIC_CHARSET      = uint32(0x000000B2)
	BALTIC_CHARSET      = uint32(0x000000BA)
	RUSSIAN_CHARSET     = uint32(0x000000CC)
	THAI_CHARSET        = uint32(0x000000DE)
	EASTEUROPE_CHARSET  = uint32(0x000000EE)
	OEM_CHARSET         = uint32(0x000000FF)
)

// TextAlignmentMode [MS-WMF]

const (
	TA_NOUPDATECP = 0x0000
	TA_LEFT       = 0x0000
	TA_TOP        = 0x0000
	TA_UPDATECP   = 0x0001
	TA_RIGHT      = 0x0002
	TA_CENTER     = 0x0003
	TA_BOTTOM     = 0x0008
	TA_BASELINE   = 0x0018
	TA_RTLREADING = 0x0100
)
