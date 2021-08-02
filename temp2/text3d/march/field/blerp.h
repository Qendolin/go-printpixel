#ifndef BLERP_H
#define BLERP_H

#include <stdint.h>
#include <immintrin.h>
#include <emmintrin.h>

struct ValueFieldIC
{
	uint32_t Len;
	uint32_t Width;
	uint32_t Height;
	uint32_t *Data;
};

struct ValueFieldI8C
{
	uint32_t Len;
	uint32_t Width;
	uint32_t Height;
	uint8_t *Data;
};

struct ValueFieldFC
{
	uint32_t Len;
	uint32_t Width;
	uint32_t Height;
	float *Data;
};

extern uint32_t lerpC32(uint32_t s, uint32_t e, uint64_t f);
extern uint64_t lerpC64A(uint64_t s, uint64_t e, uint64_t f);
extern __m128 lerpFmaF(__m128 a, __m128 b, __m128 t);
extern __m128 lerpF(__m128 a, __m128 b, __m128 t);
extern __m128 lerpF2(__m128 a, __m128 b, __m128 t);
extern __m128i lerpI128(__m128i a, __m128i b, __m128i t, __m128i ti);
extern __m256i lerpI256(__m256i a, __m256i b, __m256i t);
extern __m256i lerpI2562(__m256i a, __m256i b, __m256i t, __m256i ti);

extern uint32_t blerpC32(uint32_t c00, uint32_t c10, uint32_t c01, uint32_t c11, uint32_t tx, uint32_t ty);
extern uint32_t blerpC64(uint32_t c00, uint32_t c10, uint32_t c01, uint32_t c11, uint64_t tx, uint64_t ty);
extern __m128 blerpFmaF(__m128 c00, __m128 c10, __m128 c01, __m128 c11, __m128 fx, __m128 fy);
extern __m128 blerpF(__m128 c00, __m128 c10, __m128 c01, __m128 c11, __m128 fx, __m128 fy);
extern __m128 blerpF2(__m128 c00, __m128 c10, __m128 c01, __m128 c11, __m128 fx, __m128 fy);
extern __m128i blerpI128(__m128i c00, __m128i c10, __m128i c01, __m128i c11, __m128i fx, __m128i fxi, __m128i fy, __m128i fyi);
extern __m256i blerpI256(__m256i c00, __m256i c10, __m256i c01, __m256i c11, __m256i fx, __m256i fy);
extern __m256i blerpI2562(__m256i c00, __m256i c10, __m256i c01, __m256i c11, __m256i fx, __m256i fxi, __m256i fy, __m256i fyi);

extern __m256i loadlo_4x32(uint32_t *src);
extern __m256i loadlo_4x322(uint32_t *src);
extern __m128i loadlo_4x8(uint8_t *src);

extern void storelo_4x32(uint32_t *dst, __m256i src);
extern void storelo_4x322(uint32_t *dst, __m256i src);

extern void ScaleBlerpC(struct ValueFieldIC src, struct ValueFieldIC dst);
extern void ScaleBlerpCFSimd(struct ValueFieldFC src, struct ValueFieldFC dst);
extern void ScaleBlerpCFSimd2(struct ValueFieldFC src, struct ValueFieldFC dst);
extern void ScaleBlerpCFull(struct ValueFieldFC src, struct ValueFieldFC dst);
extern void ScaleBlerpCSimd(struct ValueFieldIC src, struct ValueFieldIC dst);
extern void ScaleBlerpCSimd2(struct ValueFieldIC src, struct ValueFieldIC dst);
extern void ScaleBlerpCSimd3(struct ValueFieldI8C src, struct ValueFieldIC dst);

#endif