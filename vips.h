//#include <vips/vips.h>
//
//#define VIPS_TYPE_SOURCE_GO (vips_source_custom_get_type())
//#define VIPS_SOURCE_GO( obj ) \
//	(G_TYPE_CHECK_INSTANCE_CAST( (obj), \
//	VIPS_TYPE_SOURCE_GO, VipsSourceGo ))
//#define VIPS_SOURCE_GO_CLASS( klass ) \
//	(G_TYPE_CHECK_CLASS_CAST( (klass), \
//	VIPS_TYPE_SOURCE_GO, VipsSourceGoClass))
//#define VIPS_IS_SOURCE_GO( obj ) \
//	(G_TYPE_CHECK_INSTANCE_TYPE( (obj), VIPS_TYPE_SOURCE_GO ))
//#define VIPS_IS_SOURCE_GO_CLASS( klass ) \
//	(G_TYPE_CHECK_CLASS_TYPE( (klass), VIPS_TYPE_SOURCE_GO ))
//#define VIPS_SOURCE_GO_GET_CLASS( obj ) \
//	(G_TYPE_INSTANCE_GET_CLASS( (obj), \
//	VIPS_TYPE_SOURCE_GO, VipsSourceGoClass ))
//
//typedef struct _VipsSourceGo {
//	VipsSource parent_object;
//
//} VipsSourceGo;
//
//typedef struct _VipsSourceGoClass {
//	VipsSourceClass parent_class;
//
//	/* The action signals clients can use to implement read and seek.
//	 * We must use gint64 everywhere since there's no G_TYPE_SIZE.
//	 */
//
//	gint64 (*read)( VipsSourceGo *, void *, gint64 );
//	gint64 (*seek)( VipsSourceGo *, gint64, int );
//
//} VipsSourceGoClass;
//
//G_DEFINE_TYPE( VipsSourceGo, vips_source_go, VIPS_TYPE_SOURCE );