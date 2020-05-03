#ifndef HAVE_GO_H
#define HAVE_GO_H

#include "source.h"

#define INT_TO_GBOOLEAN(bool) (bool > 0 ? TRUE : FALSE)

#endif

int
vips_thumbnail_bridge(VipsSourceGo *source, VipsImage **out, int width, int height, int no_rotate, int crop) {
	if (crop) {
		return vips_thumbnail_source(source, out, width,
			"height", height,
			"no_rotate", INT_TO_GBOOLEAN(no_rotate),
			"crop", VIPS_INTERESTING_CENTRE,
			NULL
		);
	}
	return vips_thumbnail_source(source, out, width,
		"height", height,
		"no_rotate", INT_TO_GBOOLEAN(no_rotate),
		NULL
	);
}