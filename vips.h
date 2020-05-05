#ifndef HAVE_GO_H
#define HAVE_GO_H

#define INT_TO_GBOOLEAN(bool) (bool > 0 ? TRUE : FALSE)

#include <vips/vips.h>

VipsImage * vips_new_image_bridge(VipsSourceCustom *source) {
	return vips_image_new_from_source( (VipsSource*) source, "", NULL);
}

int vips_thumbnail_bridge(VipsSourceCustom *source, VipsImage **out, int width, int height, int no_rotate, int crop) {
	if (crop) {
		return vips_thumbnail_source( (VipsSource*) source, out, width,
			"height", height,
			"no_rotate", INT_TO_GBOOLEAN(no_rotate),
			"crop", VIPS_INTERESTING_CENTRE,
			NULL
		);
	}
	return vips_thumbnail_source( (VipsSource*) source, out, width,
		"height", height,
		"no_rotate", INT_TO_GBOOLEAN(no_rotate),
		NULL
	);
}

int vips_jpegsave_bridge(VipsImage *in, VipsTargetCustom *target, int strip, int quality, int interlace) {
	return vips_jpegsave_target(in, (VipsTarget*) target,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"optimize_coding", TRUE,
		"interlace", INT_TO_GBOOLEAN(interlace),
		NULL
	);
}

int vips_pngsave_bridge(VipsImage *in, VipsTargetCustom *target, int strip, int compression, int quality, int interlace) {
	return vips_pngsave_target(in, (VipsTarget*) target,
		"strip", INT_TO_GBOOLEAN(strip),
		"compression", compression,
		"interlace", INT_TO_GBOOLEAN(interlace),
		"filter", VIPS_FOREIGN_PNG_FILTER_ALL,
		NULL
	);
}

int vips_webpsave_bridge(VipsImage *in, VipsTargetCustom *target, int strip, int quality, int lossless) {
	return vips_webpsave_target(in, (VipsTarget*) target,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"lossless", INT_TO_GBOOLEAN(lossless),
		NULL
	);
}

// FIXME: Not supported yet in VIPS.
//int vips_tiffsave_bridge(VipsImage *in, VipsTargetCustom *target) {
//	return vips_tiffsave_target(in, (VipsTarget*) target, NULL);
//}

// FIXME: Not supported yet in VIPS.
//int vips_heifsave_bridge(VipsImage *in, VipsTargetCustom *target, int strip, int quality, int lossless) {
//	return vips_heifsave_target(in, (VipsTarget*) target,
//		"strip", INT_TO_GBOOLEAN(strip),
//		"Q", quality,
//		"lossless", INT_TO_GBOOLEAN(lossless),
//		NULL
//	);
//}

#endif
