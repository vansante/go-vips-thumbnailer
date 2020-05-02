/*
#define VIPS_DEBUG
 */

#ifdef HAVE_CONFIG_H
#include <config.h>
#endif /*HAVE_CONFIG_H*/
#include <vips/intl.h>

#include <stdio.h>
#include <stdlib.h>
#ifdef HAVE_UNISTD_H
#include <unistd.h>
#endif /*HAVE_UNISTD_H*/
#include <string.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

#include <vips/vips.h>
#include <vips/debug.h>

#include "source.h"


static gint64
vips_source_go_read_real ( VipsSource *source, void *data, size_t length )
{
	VIPS_DEBUG_MSG( "vips_source_go_read:\n" );

    	return( 0 );
}

static gint64
vips_source_go_seek_real ( VipsSource *source, gint64 offset, int whence )
{
	VIPS_DEBUG_MSG( "vips_source_go_seek:\n" );

	return( -1 );
}

static void
vips_source_go_class_init( VipsSourceGoClass *class )
{
	VipsObjectClass *object_class = VIPS_OBJECT_CLASS( class );
	VipsSourceClass *source_class = VIPS_SOURCE_CLASS( class );

	object_class->nickname = "go source";
	object_class->description = _( "Go source" );

	source_class->read = vips_source_go_read_real;
	source_class->seek = vips_source_go_seek_real;
}

static void
vips_source_go_init( VipsSourceGo *source_custom )
{
}
