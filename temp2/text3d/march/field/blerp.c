#include "blerp.h"
#include <stdio.h>
#include <immintrin.h>

uint32_t lerpC32(uint32_t s, uint32_t e, uint64_t f)
{
	return s + ((f * (e - s)) >> 32);
}

uint64_t lerpC64A(uint64_t s, uint64_t e, uint64_t f)
{
	return (s + ((f * (e - s)) >> 32)) & UINT32_MAX;
}

uint32_t lerpC32B(uint32_t s, uint32_t e, uint64_t f)
{
	return s + (f * (e - s)) >> 32;
}

uint32_t blerpC32B(uint32_t c00, uint32_t c10, uint32_t c01, uint32_t c11, uint64_t fx, uint64_t fy)
{
	uint64_t fxi = UINT32_MAX - fx;
	uint64_t fyi = UINT32_MAX - fy;
	uint64_t a = (fx * fy) >> 32;
	uint64_t b = (fy * fxi) >> 32;
	uint64_t c = (fx * fyi) >> 32;
	uint64_t d = (fxi * fyi) >> 32;
	return (c11 * a + c01 * b + c10 * c + c00 * d) >> 32;
}

uint32_t blerpC32simd(uint32_t c00, uint32_t c10, uint32_t c01, uint32_t c11, uint64_t fx, uint64_t fy)
{
	uint64_t fxi = UINT32_MAX - fx;
	uint64_t fyi = UINT32_MAX - fy;
	__m256i lhs = _mm256_set_epi32(fx, 0, fy, 0, fx, 0, fxi, 0);
	__m256i rhs = _mm256_set_epi32(fy, 0, fxi, 0, fyi, 0, fyi, 0);
	__m256i coeffs = _mm256_mul_epu32(lhs, rhs);
	coeffs = _mm256_srli_epi64(coeffs, 32);
	__m256i points = _mm256_set_epi32(c11, 0, c01, 0, c10, 0, c00, 0);
	__m256i res = _mm256_mul_epu32(points, coeffs);
	res = _mm256_add_epi32(res, _mm256_srli_epi64(res, 8));
	res = _mm256_add_epi32(res, _mm256_srli_epi64(res, 4));
	return _mm256_extract_epi32(res, 0);
}

uint32_t blerpC32C(uint32_t c00, uint32_t c10, uint32_t c01, uint32_t c11, uint64_t fx, uint64_t fy)
{
	return c00 - (c00 * fx) >> 32 - (c00 * fy) >> 32 + (c10 * fx) >> 32 + (c01 * fy) >> 32 + (c00 * fx * fy) >> 32 - (c01 * fx * fy) >> 32 - (c10 * fx * fy) >> 32 + (c11 * fx * fy) >> 32;
}

uint32_t blerpC32(uint32_t c00, uint32_t c10, uint32_t c01, uint32_t c11, uint32_t tx, uint32_t ty)
{
	return (uint32_t)(lerpC64A(
		lerpC64A((uint64_t)c00, (uint64_t)c10, (uint64_t)tx),
		lerpC64A((uint64_t)c01, (uint64_t)c11, (uint64_t)tx),
		(uint64_t)ty));
}

uint32_t blerpC64(uint32_t c00, uint32_t c10, uint32_t c01, uint32_t c11, uint64_t tx, uint64_t ty)
{
	return (uint32_t)(lerpC32B(
		lerpC32B((uint32_t)c00, (uint32_t)c10, tx),
		lerpC32B((uint32_t)c01, (uint32_t)c11, tx),
		ty));
}

void ScaleBlerpC(struct ValueFieldIC src, struct ValueFieldIC dst)
{
	uint32_t *sp = src.Data;
	uint32_t *dp = dst.Data;

	uint32_t send = src.Len;
	uint32_t srowstride = src.Width * 3;
	uint32_t drowstride = dst.Width * 3;

	uint32_t dTxF = UINT32_MAX / dst.Width;
	uint32_t dTyF = UINT32_MAX / dst.Height;
	uint32_t sTxF = UINT32_MAX / (src.Width - 1);
	uint32_t sTyF = UINT32_MAX / (src.Height - 1);

	uint32_t d2sxF = (uint32_t)(((uint64_t)(src.Width - 1) * UINT32_MAX) / dst.Width);
	uint32_t d2syF = (uint32_t)(((uint64_t)(src.Height - 1) * UINT32_MAX) / dst.Height);

	uint64_t dyF = 0, dyFend = 0, dy = 0;
	uint32_t fy = 0;

	for (uint32_t sy = 0; sy < send - srowstride; sy += srowstride)
	{
		uint32_t *syp = sp + sy;

		for (dyFend += sTyF; dyF < dyFend; dyF += dTyF)
		{
			uint32_t *dyp = dp + dy;

			uint64_t dxF = 0, dxFend = 0, dx = 0;
			uint32_t fx = 0;

			for (uint32_t sx = 0; sx < srowstride - 3; sx += 3)
			{
				uint32_t *sxp = syp + sx;

				uint32_t r00 = *(sxp + 0);
				uint32_t g00 = *(sxp + 1);
				uint32_t b00 = *(sxp + 2);
				uint32_t r10 = *(sxp + 3);
				uint32_t g10 = *(sxp + 4);
				uint32_t b10 = *(sxp + 5);
				uint32_t r01 = *(sxp + 0 + srowstride);
				uint32_t g01 = *(sxp + 1 + srowstride);
				uint32_t b01 = *(sxp + 2 + srowstride);
				uint32_t r11 = *(sxp + 3 + srowstride);
				uint32_t g11 = *(sxp + 4 + srowstride);
				uint32_t b11 = *(sxp + 5 + srowstride);

				for (dxFend += sTxF; dxF < dxFend; dxF += dTxF)
				{
					*(dyp + dx + 0) = blerpC32(r00, r10, r01, r11, (uint64_t)fx, (uint64_t)fy);
					*(dyp + dx + 1) = blerpC32(g00, g10, g01, g11, (uint64_t)fx, (uint64_t)fy);
					*(dyp + dx + 2) = blerpC32(b00, b10, b01, b11, (uint64_t)fx, (uint64_t)fy);

					dx += 3;
					fx += d2sxF;
				}
			}

			dy += drowstride;
			fy += d2syF;
		}
	}
}
