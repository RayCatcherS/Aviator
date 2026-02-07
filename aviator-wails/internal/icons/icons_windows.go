package icons

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"syscall"
	"unsafe"
)

var (
	shell32            = syscall.NewLazyDLL("shell32.dll")
	user32             = syscall.NewLazyDLL("user32.dll")
	gdi32              = syscall.NewLazyDLL("gdi32.dll")
	extractIconExW     = shell32.NewProc("ExtractIconExW")
	destroyIcon        = user32.NewProc("DestroyIcon")
	getIconInfo        = user32.NewProc("GetIconInfo")
	getDIBits          = gdi32.NewProc("GetDIBits")
	deleteDC           = gdi32.NewProc("DeleteDC")
	deleteObject       = gdi32.NewProc("DeleteObject")
	createCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
	getDC              = user32.NewProc("GetDC")
)

type ICONINFO struct {
	fIcon    uint32
	xHotspot uint32
	yHotspot uint32
	hbmMask  uintptr
	hbmColor uintptr
}

type BITMAP struct {
	bmType       int32
	bmWidth      int32
	bmHeight     int32
	bmWidthBytes int32
	bmPlanes     uint16
	bmBitsPixel  uint16
	bmBits       uintptr
}

type BITMAPINFOHEADER struct {
	biSize          uint32
	biWidth         int32
	biHeight        int32
	biPlanes        uint16
	biBitCount      uint16
	biCompression   uint32
	biSizeImage     uint32
	biXPelsPerMeter int32
	biYPelsPerMeter int32
	biClrUsed       uint32
	biClrImportant  uint32
}

// ExtractIconToBase64 extracts the icon from an exe file and returns it as base64 PNG
func ExtractIconToBase64(exePath string) (string, error) {
	// Try to extract large icon
	var hIcon uintptr
	path, err := syscall.UTF16PtrFromString(exePath)
	if err != nil {
		return "", err
	}

	// Extract one large icon
	ret, _, _ := extractIconExW.Call(
		uintptr(unsafe.Pointer(path)),
		0,                               // icon index
		uintptr(unsafe.Pointer(&hIcon)), // large icon
		0,                               // no small icon
		1,                               // extract 1 icon
	)

	if ret == 0 || hIcon == 0 {
		// Try extracting small icon if large fails
		ret, _, _ = extractIconExW.Call(
			uintptr(unsafe.Pointer(path)),
			0,
			0,
			uintptr(unsafe.Pointer(&hIcon)),
			1,
		)
		if ret == 0 || hIcon == 0 {
			return "", fmt.Errorf("failed to extract icon from %s", exePath)
		}
	}
	defer destroyIcon.Call(hIcon)

	// Convert icon to PNG
	img, err := iconToImage(hIcon)
	if err != nil {
		return "", err
	}

	// Encode to PNG and then base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func iconToImage(hIcon uintptr) (image.Image, error) {
	var info ICONINFO
	ret, _, _ := getIconInfo.Call(hIcon, uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return nil, fmt.Errorf("GetIconInfo failed")
	}
	defer deleteObject.Call(info.hbmColor)
	defer deleteObject.Call(info.hbmMask)

	// Get bitmap dimensions
	var bm BITMAP
	ret, _, _ = syscall.Syscall(
		gdi32.NewProc("GetObjectW").Addr(),
		3,
		info.hbmColor,
		unsafe.Sizeof(bm),
		uintptr(unsafe.Pointer(&bm)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("GetObject failed")
	}

	width := int(bm.bmWidth)
	height := int(bm.bmHeight)

	// Create DC
	hdcScreen, _, _ := getDC.Call(0)
	hdc, _, _ := createCompatibleDC.Call(hdcScreen)
	defer deleteDC.Call(hdc)

	// Prepare BITMAPINFO
	bi := BITMAPINFOHEADER{
		biSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
		biWidth:       int32(width),
		biHeight:      -int32(height), // top-down DIB
		biPlanes:      1,
		biBitCount:    32,
		biCompression: 0, // BI_RGB
	}

	// Allocate buffer for pixel data
	bufSize := width * height * 4
	pixels := make([]byte, bufSize)

	// Get DIB bits
	ret, _, _ = getDIBits.Call(
		hdc,
		info.hbmColor,
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&pixels[0])),
		uintptr(unsafe.Pointer(&bi)),
		0, // DIB_RGB_COLORS
	)
	if ret == 0 {
		return nil, fmt.Errorf("GetDIBits failed")
	}

	// Create RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Convert BGRA to RGBA
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (y*width + x) * 4
			b := pixels[idx]
			g := pixels[idx+1]
			r := pixels[idx+2]
			a := pixels[idx+3]

			imgIdx := y*img.Stride + x*4
			img.Pix[imgIdx] = r
			img.Pix[imgIdx+1] = g
			img.Pix[imgIdx+2] = b
			img.Pix[imgIdx+3] = a
		}
	}

	return img, nil
}
