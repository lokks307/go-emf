package emf

import "bytes"

// map of readers for records
var records = map[uint32]func(*bytes.Reader, uint32) (Recorder, error){
	EMR_HEADER:                  readHeaderRecord,
	EMR_POLYBEZIER:              nil,
	EMR_POLYGON:                 nil,
	EMR_POLYLINE:                nil,
	EMR_POLYBEZIERTO:            nil,
	EMR_POLYLINETO:              nil,
	EMR_POLYPOLYLINE:            nil,
	EMR_POLYPOLYGON:             nil,
	EMR_SETWINDOWEXTEX:          readSetWindowExtExRecord,
	EMR_SETWINDOWORGEX:          readSetWindowOrgExRecord,
	EMR_SETVIEWPORTEXTEX:        readSetViewportExtExRecord,
	EMR_SETVIEWPORTORGEX:        readSetViewportOrgExRecord,
	EMR_SETBRUSHORGEX:           readSetBrushOrgExRecord,
	EMR_EOF:                     readEofRecord,
	EMR_SETPIXELV:               readSetPixelvRecord,
	EMR_SETMAPPERFLAGS:          readSetMapperFlagsRecord,
	EMR_SETMAPMODE:              readSetMapModeRecord,
	EMR_SETBKMODE:               readSetBkModeRecord,
	EMR_SETPOLYFILLMODE:         readSetPolyfillModeRecord,
	EMR_SETROP2:                 readSetROP2Record,
	EMR_SETSTRETCHBLTMODE:       readSetStretchBltModeRecord,
	EMR_SETTEXTALIGN:            readSetTextAlignRecord,
	EMR_SETCOLORADJUSTMENT:      nil,
	EMR_SETTEXTCOLOR:            readSetTextColorRecord,
	EMR_SETBKCOLOR:              readSetBkColorRecord,
	EMR_OFFSETCLIPRGN:           readOffSetClipRgnRecord,
	EMR_MOVETOEX:                readMoveToExRecord,
	EMR_SETMETARGN:              readSetMetaRgnRecord,
	EMR_EXCLUDECLIPRECT:         nil,
	EMR_INTERSECTCLIPRECT:       readIntersectClipRectRecord,
	EMR_SCALEVIEWPORTEXTEX:      nil,
	EMR_SCALEWINDOWEXTEX:        readScaleWindowExtExRecord,
	EMR_SAVEDC:                  readSaveDCRecord,
	EMR_RESTOREDC:               readRestoreDCRecord,
	EMR_SETWORLDTRANSFORM:       readSetWorldTransformRecord,
	EMR_MODIFYWORLDTRANSFORM:    readModifyWorldTransformRecord,
	EMR_SELECTOBJECT:            readSelectObjectRecord,
	EMR_CREATEPEN:               readCreatePenRecord,
	EMR_CREATEBRUSHINDIRECT:     readCreateBrushIndirectRecord,
	EMR_DELETEOBJECT:            readDeleteObjectRecord,
	EMR_ANGLEARC:                nil,
	EMR_ELLIPSE:                 nil,
	EMR_RECTANGLE:               readRectangleRecord,
	EMR_ROUNDRECT:               nil,
	EMR_ARC:                     readArcRecord,
	EMR_CHORD:                   nil,
	EMR_PIE:                     nil,
	EMR_SELECTPALETTE:           readSelectPaletteRecord,
	EMR_CREATEPALETTE:           readCreatePaletteRecord,
	EMR_SETPALETTEENTRIES:       nil,
	EMR_RESIZEPALETTE:           nil,
	EMR_REALIZEPALETTE:          nil,
	EMR_EXTFLOODFILL:            nil,
	EMR_LINETO:                  readLineToRecord,
	EMR_ARCTO:                   nil,
	EMR_POLYDRAW:                nil,
	EMR_SETARCDIRECTION:         nil,
	EMR_SETMITERLIMIT:           readSetMiterLimitRecord,
	EMR_BEGINPATH:               readBeginPathRecord,
	EMR_ENDPATH:                 readEndPathRecord,
	EMR_CLOSEFIGURE:             readCloseFigureRecord,
	EMR_FILLPATH:                readFillPathRecord,
	EMR_STROKEANDFILLPATH:       readStrokeAndFillPathRecord,
	EMR_STROKEPATH:              readStrokePathRecord,
	EMR_FLATTENPATH:             nil,
	EMR_WIDENPATH:               nil,
	EMR_SELECTCLIPPATH:          readSelectClipPathRecord,
	EMR_ABORTPATH:               readAbortPathRecord,
	EMR_COMMENT:                 readCommentRecord,
	EMR_FILLRGN:                 readFillRgnRecord,
	EMR_FRAMERGN:                nil,
	EMR_INVERTRGN:               nil,
	EMR_PAINTRGN:                nil,
	EMR_EXTSELECTCLIPRGN:        readExtSelectClipRgnRecord,
	EMR_BITBLT:                  readBitBltRecord,
	EMR_STRETCHBLT:              readStretchBltRecord,
	EMR_MASKBLT:                 readMaskBltRecord,
	EMR_PLGBLT:                  nil,
	EMR_SETDIBITSTODEVICE:       readSetDIBitsToDeviceRecord,
	EMR_STRETCHDIBITS:           readStretchDIBitsRecord,
	EMR_EXTCREATEFONTINDIRECTW:  readExtCreateFontIndirectWRecord,
	EMR_EXTTEXTOUTA:             nil,
	EMR_EXTTEXTOUTW:             readExtTextOutWRecord,
	EMR_POLYBEZIER16:            readPolyBezier16Record,
	EMR_POLYGON16:               readPolygon16Record,
	EMR_POLYLINE16:              readPolyLine16Record,
	EMR_POLYBEZIERTO16:          readPolyBezierTo16Record,
	EMR_POLYLINETO16:            readPolyLineTo16Record,
	EMR_POLYPOLYLINE16:          nil,
	EMR_POLYPOLYGON16:           readPolyPolygon16Record,
	EMR_POLYDRAW16:              nil,
	EMR_CREATEMONOBRUSH:         nil,
	EMR_CREATEDIBPATTERNBRUSHPT: nil,
	EMR_EXTCREATEPEN:            readExtCreatePenRecord,
	EMR_POLYTEXTOUTA:            nil,
	EMR_POLYTEXTOUTW:            nil,
	EMR_SETICMMODE:              readSetICMModeRecord,
	EMR_CREATECOLORSPACE:        nil,
	EMR_SETCOLORSPACE:           nil,
	EMR_DELETECOLORSPACE:        nil,
	EMR_GLSRECORD:               nil,
	EMR_GLSBOUNDEDRECORD:        nil,
	EMR_PIXELFORMAT:             nil,
	EMR_DRAWESCAPE:              nil,
	EMR_EXTESCAPE:               nil,
	EMR_SMALLTEXTOUT:            nil,
	EMR_FORCEUFIMAPPING:         nil,
	EMR_NAMEDESCAPE:             nil,
	EMR_COLORCORRECTPALETTE:     nil,
	EMR_SETICMPROFILEA:          nil,
	EMR_SETICMPROFILEW:          nil,
	EMR_ALPHABLEND:              nil,
	EMR_SETLAYOUT:               nil,
	EMR_TRANSPARENTBLT:          nil,
	EMR_GRADIENTFILL:            nil,
	EMR_SETLINKEDUFIS:           nil,
	EMR_SETTEXTJUSTIFICATION:    readSetTextJustificationRecord,
	EMR_COLORMATCHTOTARGETW:     nil,
	EMR_CREATECOLORSPACEW:       nil,
}
