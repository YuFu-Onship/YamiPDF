#ifndef CUSTOM_FITZ_H
#define CUSTOM_FITZ_H

#include "mupdf/fitz.h"
#include "mupdf/fitz/color.h"
#include "mupdf/fitz/context.h"
#include "mupdf/fitz/device.h"
#include "mupdf/fitz/document.h"
#include "mupdf/fitz/geometry.h"
#include "mupdf/fitz/glyph-cache.h"
#include "mupdf/fitz/path.h"
#include "mupdf/fitz/pixmap.h"
#include "mupdf/fitz/types.h"

// 得到当前文档的页数
static inline int custom_get_page_count(void *ctx_ptr, void *doc_ptr) {
  fz_context *ctx = (fz_context *)ctx_ptr;
  fz_document *doc = (fz_document *)doc_ptr;

  if (!ctx || !doc) {
    return -1;
  }

  int count = 0;
  fz_try(ctx) { count = fz_count_pages(ctx, doc); }
  fz_catch(ctx) { return -1; }

  return count;
}

// 一次性释放瓦片内容
static inline void custom_drop_down_tiles(void *ctx_ptr, void *doc_ptr) {
  // 类型还原
  fz_context *ctx = (fz_context *)ctx_ptr; // (目标类型 *)强制类型转化
  fz_document *doc = (fz_document *)doc_ptr;

  // 检测是否为空
  if (!ctx || !doc) {
    return;
  }
  fz_drop_drawn_tiles_for_document(ctx, doc);
}

// 释放渲染路径
static inline void custom_drop_pixmap(void *ctx_ptr, void *pix_ptr) {
  fz_context *ctx = (fz_context *)ctx_ptr;
  fz_pixmap *path = (fz_pixmap *)pix_ptr;

  if (!ctx) {
    return;
  }

  fz_drop_pixmap(ctx, path);
}

// 得到指定dpi下的某一页的渲染尺寸
static inline void custom_pdf_get_page_size(void *ctx_ptr, void *doc_ptr,
                                            int page, float dpi) {
  fz_context *ctx = (fz_context *)ctx_ptr;
  fz_document *doc = (fz_document *)doc_ptr;
}

// 渲染瓦片
static inline void custom_render_full(void *ctx_ptr, void *doc_ptr,
                                      int page_num, float dpi, int tile_w,
                                      int tile_h, unsigned char *out_buffer) {
  fz_context *ctx = (fz_context *)ctx_ptr;
  fz_document *doc = (fz_document *)doc_ptr;

  fz_page *page = fz_load_page(ctx, doc, page_num);

  fz_rect bounds = fz_bound_page(ctx, page);
  fz_matrix ctm = fz_scale(dpi / 72, dpi / 72);

  bounds = fz_transform_rect(bounds, ctm);
  // fz_irect bbox = fz_round_rect(bounds);

  // printf("%d", );

  fz_irect bbox = fz_make_irect(bounds.x0, bounds.y0, bounds.x1, bounds.y1);

  fz_pixmap *pix =
      fz_new_pixmap_with_bbox(ctx, fz_device_rgb(ctx), bbox, NULL, 1);

  fz_device *device = fz_new_draw_device(ctx, ctm, pix);
  fz_enable_device_hints(ctx, device, FZ_NO_CACHE);

  fz_matrix drawMatrix = fz_identity;
  fz_run_page_contents(ctx, page, device, drawMatrix, NULL);

  unsigned char *pixels = fz_pixmap_samples(ctx, pix);
  if (pixels) {
    int stride = fz_pixmap_stride(ctx, pix);
    for (int y = 0; y < tile_h; y++) {
      unsigned char *dest = out_buffer + y * tile_w * 4;
      unsigned char *src = pixels + y * stride;
      memcpy(dest,      // 目标地址, 每行固定多出4字节的填充,不表示颜色
             src,       // 源地址
             tile_w * 4 // 复制的字节数,每个像素占用4个字节
      );
    }
  }

  fz_close_device(ctx, device);
  fz_clear_pixmap_with_value(ctx, pix, 0xff);

  fz_drop_device(ctx, device);
  fz_drop_pixmap(ctx, pix);
  fz_drop_page(ctx, page);
  // fzrender

  // fz_render_glyph_pixmap(fz_context * ctx, fz_font * font, int gid,
  // fz_matrix *ctm, const fz_irect *scissor, int aa)
}

#endif
