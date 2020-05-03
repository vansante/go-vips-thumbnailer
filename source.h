#include <vips/vips.h>

typedef struct _VipsSourceGo {
	VipsSource parent_object;

	int id;
} VipsSourceGo;

typedef struct _VipsSourceGoClass {
	VipsSourceClass parent_class;

	gint64 (*read)( VipsSourceGo *, void *, gint64 );
	gint64 (*seek)( VipsSourceGo *, gint64, int );

} VipsSourceGoClass;

GType vips_source_go_get_type( void );
VipsSourceGo *vips_source_go_new( int id );

static void vips_source_go_class_init ( VipsSourceGoClass *class );
static void vips_source_go_init ( VipsSourceGo *source_go );

