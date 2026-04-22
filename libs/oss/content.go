package oss

var (
	// ImageContentTypes 图像类
	ImageContentTypes = []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/bmp",
		"image/webp",
		"image/svg+xml",
		"image/webp",
		"image/tiff",
		"image/vnd.microsoft.icon",
	}

	// VideoContentTypes 视频类
	VideoContentTypes = []string{
		"video/mp4",
		"video/mpeg",
		"video/ogg",
		"video/webm",
		"video/quicktime",
		"video/x-msvideo",
		"video/x-ms-wmv",
	}

	// DocumentContentTypes 文档类
	DocumentContentTypes = []string{
		"text/plain",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/html",
		"application/json",
	}

	// ZipContentTypes 压缩类
	ZipContentTypes = []string{
		"application/zip",
		"application/gzip",
		"application/x-tar",
		"application/x-rar-compressed",
	}
)
