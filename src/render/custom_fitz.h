#ifndef CUSTOM_FITZ_H
#define CUSTOM_FITZ_H

#include "mupdf/fitz.h"
#include "mupdf/fitz/color.h"
#include "mupdf/fitz/context.h"
#include "mupdf/fitz/glyph-cache.h"
#include "mupdf/fitz/path.h"
#include "mupdf/fitz/pixmap.h"

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
static inline void custom_drop_drawn_tiles(void *ctx_ptr, void *doc_ptr) {
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

// 渲染瓦片
static inline void custom_render_tile(void *ctx_ptr) {
  // fz_render_glyph_pixmap(fz_context * ctx, fz_font * font, int gid,
  // fz_matrix *ctm, const fz_irect *scissor, int aa)
}

#endif
