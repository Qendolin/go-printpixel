#include "blerp.h"
#include <stdio.h>
#include <string.h>
#include <math.h>

inline __m128 lerpFmaF(__m128 a, __m128 b, __m128 t)
{
	return _mm_fmadd_ps(t, b, _mm_fnmadd_ps(t, a, a));
}

inline __m128 lerpF(__m128 a, __m128 b, __m128 t)
{
	return _mm_add_ps(_mm_mul_ps(t, b), _mm_sub_ps(a, _mm_mul_ps(t, a)));
}

inline __m128 lerpF2(__m128 s, __m128 e, __m128 f)
{
	return _mm_add_ps(_mm_mul_ps(f, s), _mm_sub_ps(e, _mm_mul_ps(f, e)));
}

inline __m128 blerpFmaF(__m128 c00, __m128 c10, __m128 c01, __m128 c11, __m128 fx, __m128 fy)
{
	return lerpFmaF(
		lerpFmaF(c00, c10, fx),
		lerpFmaF(c01, c11, fx),
		fy);
}

inline __m128 blerpF(__m128 c00, __m128 c10, __m128 c01, __m128 c11, __m128 fx, __m128 fy)
{
	return lerpF(
		lerpF(c00, c10, fx),
		lerpF(c01, c11, fx),
		fy);
}

inline __m128 blerpF2(__m128 c00, __m128 c10, __m128 c01, __m128 c11, __m128 fx, __m128 fy)
{
	return lerpF2(
		lerpF2(c00, c10, fx),
		lerpF2(c01, c11, fx),
		fy);
}
// _mm512_mulhi_epu32
inline __m256i lerpI256(__m256i a, __m256i b, __m256i t)
{
	// return _mm256_add_epi32(a, _mm256_srli_epi64(_mm256_mul_epu32(_mm256_sub_epi32(b, a), t), 32));
	// const __m256i K_MASK = _mm256_setr_epi32(UINT32_MAX, 0, UINT32_MAX, 0, UINT32_MAX, 0, UINT32_MAX, 0);
	// return _mm256_add_epi32(a, _mm256_srli_epi64(_mm256_mul_epu32(_mm256_sub_epi64(_mm256_max_epu32(a, b), _mm256_min_epu32(a, b)), t), 32));
	__m256i ZERO = _mm256_setzero_si256();
	return _mm256_srli_epi64(_mm256_add_epi32(_mm256_mul_epu32(a, _mm256_sub_epi32(ZERO, t)), _mm256_mul_epu32(b, t)), 32);
}

inline __m256i blerpI256(__m256i c00, __m256i c10, __m256i c01, __m256i c11, __m256i fx, __m256i fy)
{
	return lerpI256(
		lerpI256(c00, c10, fx),
		lerpI256(c01, c11, fx),
		fy);
}

inline __m256i lerpI2562(__m256i a, __m256i b, __m256i t, __m256i ti)
{
	return _mm256_srli_epi64(_mm256_add_epi32(_mm256_mul_epu32(a, ti), _mm256_mul_epu32(b, t)), 32);
}

inline __m256i blerpI2562(__m256i c00, __m256i c10, __m256i c01, __m256i c11, __m256i fx, __m256i fxi, __m256i fy, __m256i fyi)
{
	return lerpI2562(
		lerpI2562(c00, c10, fx, fxi),
		lerpI2562(c01, c11, fx, fxi),
		fy, fyi);
}

inline __m128i lerpI128(__m128i a, __m128i b, __m128i t, __m128i ti)
{
	return _mm_add_epi32(_mm_mullo_epi32(a, ti), _mm_mullo_epi32(b, t));
}

inline __m128i blerpI128(__m128i c00, __m128i c10, __m128i c01, __m128i c11, __m128i fx, __m128i fxi, __m128i fy, __m128i fyi)
{
	return lerpI128(
		lerpI128(c00, c10, fx, fxi),
		lerpI128(c01, c11, fx, fxi),
		fy, fyi);
}

// https://stackoverflow.com/q/40144411/7448536
inline void storelo_4x32(uint32_t *dst, __m256i src)
{
	const __m256i K_PERM = _mm256_setr_epi32(0, 2, 4, 6, 1, 3, 5, 7);
	__m256i permuted = _mm256_permutevar8x32_epi32(src, K_PERM);
	__m128i lo128 = _mm256_extractf128_si256(permuted, 0);
	_mm_storeu_si128((__m128i *)dst, lo128);
}

inline void storelo_4x322(uint32_t *dst, __m256i src)
{
	const __m256i K_PERM = _mm256_setr_epi32(0, 2, 4, 6, 1, 3, 5, 7);
	__m256i permuted = _mm256_permutevar8x32_epi32(src, K_PERM);
	__m128i lo128 = _mm256_castsi256_si128(permuted);
	_mm_storeu_si128((__m128i *)dst, lo128);
}

inline __m256i loadlo_4x32(uint32_t *src)
{
	const __m256i K_GATHER = _mm256_setr_epi32(0, 0, 1, 1, 2, 2, 3, 3);
	return _mm256_i32gather_epi32(src, K_GATHER, sizeof(uint32_t));
}

inline __m256i loadlo_4x322(uint32_t *src)
{
	// Since high 128 bits are undefined this probably has undefined behavoir
	__m256i data = _mm256_castsi128_si256(_mm_loadu_si128((__m128i *)src));
	const __m256i K_PERM = _mm256_setr_epi32(0, 4, 1, 5, 2, 6, 3, 7);
	return _mm256_permutevar8x32_epi32(data, K_PERM);
}

inline __m128i loadlo_4x8(uint8_t *src)
{
	const __m128i ZERO = _mm_setzero_si128();
	__m128i xmm0 = _mm_cvtsi32_si128(*(const uint32_t *)src);
	xmm0 = _mm_unpacklo_epi8(xmm0, ZERO);
	return _mm_unpacklo_epi16(xmm0, ZERO);
}

void ScaleBlerpCFSimd(struct ValueFieldFC src, struct ValueFieldFC dst)
{
	float *sp = src.Data;
	float *dp = dst.Data;

	uint32_t send = src.Len;
	uint32_t srowstride = src.Width * 4;
	uint32_t drowstride = dst.Width * 4;

	uint32_t dTxF = UINT32_MAX / dst.Width;
	uint32_t dTyF = UINT32_MAX / dst.Height;
	uint32_t sTxF = UINT32_MAX / (src.Width - 1);
	uint32_t sTyF = UINT32_MAX / (src.Height - 1);

	float d2sxF = (float)(src.Width - 1) / (float)dst.Width;
	float d2syF = (float)(src.Height - 1) / (float)dst.Height;

	uint64_t dyF = 0, dyFend = 0, dy = 0;
	float fy = 0;

	for (int sy = 0; sy < send - srowstride; sy += srowstride)
	{
		float *syp = sp + sy;

		for (dyFend += sTyF; dyF < dyFend; dyF += dTyF)
		{
			float *dyp = dp + dy;

			uint64_t dxF = 0, dxFend = 0, dx = 0;
			float fx = 0;

			for (uint32_t sx = 0; sx < srowstride - 4; sx += 4)
			{
				float *sxp = syp + sx;

				__m128 c00 = _mm_loadu_ps(sxp);
				__m128 c10 = _mm_loadu_ps(sxp + 4);
				__m128 c01 = _mm_loadu_ps(sxp + srowstride);
				__m128 c11 = _mm_loadu_ps(sxp + srowstride + 4);

				for (dxFend += sTxF; dxF < dxFend; dxF += dTxF)
				{
					__m128 res = blerpFmaF(c00, c10, c01, c11, _mm_set_ps1(fx - (int)fx), _mm_set_ps1(fy - (int)fy));
					_mm_store_ps(dyp + dx, res);
					dx += 4;
					fx += d2sxF;
				}
			}

			dy += drowstride;
			fy += d2syF;
		}
	}
}

void ScaleBlerpCFSimd2(struct ValueFieldFC src, struct ValueFieldFC dst)
{
	float *sp = src.Data;
	float *dp = dst.Data;
	// __m128 d2sxF = _mm_set_ps1((float)(src.Width - 1) / (float)(dst.Width));
	float d2sxF = (float)(src.Width - 1) / (float)(dst.Width - 1);
	// uint64_t d2sxF = ((uint64_t)src.Width - 1) * (uint64_t)UINT32_MAX / (uint64_t)dst.Width;
	float s2dxF = (float)(dst.Width - 1) / (float)(src.Width - 1);
	// uint64_t s2dxF = ((uint64_t)(dst.Width) * UINT32_MAX) / (src.Width - 1);
	// __m128 d2syF = _mm_set_ps1((float)(src.Height - 1) / (float)(dst.Height));
	float d2syF = (float)(src.Height - 1) / (float)(dst.Height - 1);
	// uint64_t d2syF = ((uint64_t)src.Height - 1) * (uint64_t)UINT32_MAX / (uint64_t)dst.Height;
	float s2dyF = (float)(dst.Height - 1) / (float)(src.Height - 1);
	// uint64_t s2dyF = ((uint64_t)(dst.Height) * UINT32_MAX) / (src.Height - 1);
	// const __m128 ONE = _mm_set_ps1(1.0);

	// float dTx = 1 / (float)dst.Width;
	// float dTy = 1 / (float)dst.Height;
	// float sTx = 1 / (float)(src.Width - 1);
	// float sTy = 1 / (float)(src.Height - 1);
	// uint32_t dTxF = UINT32_MAX / dst.Width;
	// uint32_t dTyF = UINT32_MAX / dst.Height;
	// uint32_t sTxF = UINT32_MAX / (src.Width - 1);
	// uint32_t sTyF = UINT32_MAX / (src.Height - 1);

	uint32_t srowstride = src.Width * 4;
	uint32_t drowstride = dst.Width * 4;

	// __m128 fy = _mm_set_ps1(0);
	// float fy = 0;
	float dy = 0, dyEnd = 0;
	// uint64_t dy = 0, dyEnd = 0;
	uint64_t dyI = 0;

	for (uint32_t sy = 0; sy < src.Height - 1; sy++)
	{
		// __m128 fy = _mm_set_ps1(0);
		// __m128 fyi = _mm_set_ps1(1.0);
		float *syp = sp + sy * srowstride;

		// for (dyEnd += s2dyF; dy < dyEnd; dy += 0x100000000)
		for (dyEnd += s2dyF; dy < dyEnd; dy++)
		{
			// float *dyp = dp + (dy >> 32) * drowstride;
			float *dyp = dp + dyI;
			dyI += drowstride;
			// fy = (dyEnd - dy) / (float)s2dyF;
			// __m128 fx = _mm_set_ps1(0);
			// float fx = 0;

			float dx = 0, dxEnd = 0;
			// uint64_t dx = 0, dxEnd = 0;
			uint32_t dxI = 0;

			for (uint32_t sx = 0; sx < src.Width - 1; sx++)
			{
				// printf("x %zu\n", sx);
				//  __m128 fx = _mm_set_ps1(0);
				//  __m128 fxi = _mm_set_ps1(1.0);
				float *sxp = syp + sx * 4;

				__m128 c00 = _mm_loadu_ps(sxp);
				__m128 c10 = _mm_loadu_ps(sxp + 4);
				__m128 c01 = _mm_loadu_ps(sxp + srowstride);
				__m128 c11 = _mm_loadu_ps(sxp + srowstride + 4);

				// for (dxEnd += s2dxF; dx < dxEnd; dx += 0x100000000)
				// for (dxEnd++; dx < dxEnd; dx += d2sxF)
				for (dxEnd += s2dxF; dx < dxEnd; dx++)
				{
					// float *dpx = dyp + (dx >> 32) * 4;
					float *dpx = dyp + dxI;
					dxI += 4;
					// fx = (dxEnd - dx) / (float)s2dxF;

					// __m128 res = blerpF2(c00, c10, c01, c11, fx, fxi, fy, fyi);
					// __m128 res = blerpF(c00, c10, c01, c11, _mm_set_ps1(fx), _mm_set_ps1(fy));
					// __m128 res = blerpF(c00, c10, c01, c11, fx, fy);
					// __m128 res = blerpF(c00, c10, c01, c11, _mm_set_ps1(1 - (dxEnd - dx)), _mm_set_ps1(1 - (dyEnd - dy)));
					__m128 res = blerpF2(c00, c10, c01, c11, _mm_set_ps1((dxEnd - dx) * d2sxF), _mm_set_ps1((dyEnd - dy) * d2syF));
					_mm_store_ps(dpx, res);

					// fx = _mm_add_ps(fx, d2sxF);
					// fxi = _mm_sub_ps(ONE, fx);
					// fx += d2sxF;
					// fx = remainder(fx, 1.);
					// fx = fx - floorf(fx);
					// fx = fx - (int)fx;
					// fx = _mm_sub_ps(fx, _mm_round_ps(fx, _MM_FROUND_TRUNC));
				}
			}

			// fy = _mm_add_ps(fy, d2syF);
			// fyi = _mm_sub_ps(ONE, fy);
			// fy += d2syF;
			// fy = remainder(fy, 1.);
			// fy = fy - floorf(fy);
			// fy = fy - (int)fy;
			// fy = _mm_sub_ps(fy, _mm_round_ps(fy, _MM_FROUND_TRUNC));
		}
	}
}

void ScaleBlerpCFull(struct ValueFieldFC src, struct ValueFieldFC dst)
{
	float *sp = src.Data;
	float *dp = dst.Data;
	float d2sxF = (float)(src.Width) / (float)(dst.Width);
	float s2dxF = (float)(dst.Width) / (float)(src.Width);
	float d2syF = (float)(src.Height) / (float)(dst.Height);
	float s2dyF = (float)(dst.Height) / (float)(src.Height);

	uint32_t sRowStride = src.Width * 4;
	uint32_t dRowStride = dst.Width * 4;

	float dy = 0, dyEnd = 0;
	uint64_t offsetX = s2dxF == 1 ? 0 : (uint64_t)(0.5 * s2dxF) - 1;
	uint64_t offsetY = s2dyF == 1 ? 0 : (uint64_t)(0.5 * s2dyF) - 1;
	uint64_t dyI = offsetX * 4 + offsetY * dRowStride;

	// Inner
	for (uint32_t sy = 0; sy < src.Height - 1; sy++)
	{
		float *syp = sp + sy * sRowStride;

		for (dyEnd += s2dyF; dy <= dyEnd; dy++)
		{
			float *dyp = dp + dyI;
			dyI += dRowStride;

			float dx = 0, dxEnd = 0;
			uint32_t dxI = 0;

			for (uint32_t sx = 0; sx < src.Width - 1; sx++)
			{
				float *sxp = syp + sx * 4;

				__m128 c00 = _mm_loadu_ps(sxp);
				__m128 c10 = _mm_loadu_ps(sxp + 4);
				__m128 c01 = _mm_loadu_ps(sxp + sRowStride);
				__m128 c11 = _mm_loadu_ps(sxp + sRowStride + 4);

				for (dxEnd += s2dxF; dx <= dxEnd; dx++)
				{
					float *dpx = dyp + dxI;
					dxI += 4;

					__m128 res = blerpF2(c00, c10, c01, c11, _mm_set_ps1((dxEnd - dx) * d2sxF), _mm_set_ps1((dyEnd - dy) * d2syF));
					_mm_store_ps(dpx, res);
				}
			}
		}
	}

	uint32_t margin_x_start = floorf(s2dxF * 0.5);
	uint32_t margin_x_end = ceilf(s2dxF * 0.5);

	uint32_t margin_y_start = floorf(s2dyF * 0.5);
	uint32_t margin_y_end = ceilf(s2dyF * 0.5);

	dy = 0, dyEnd = 0;
	dyI = dRowStride * (margin_y_start == 0 ? 0 : margin_y_start - 1);
	// Interpolate left & right edges vertically
	for (uint32_t sy = 0; sy < src.Height - 1; sy++)
	{
		float *syp = sp + sy * sRowStride;

		__m128 c00 = _mm_loadu_ps(syp);
		__m128 c10 = _mm_loadu_ps(syp + sRowStride - 4);
		__m128 c01 = _mm_loadu_ps(syp + sRowStride);
		__m128 c11 = _mm_loadu_ps(syp + sRowStride + sRowStride - 4);
		for (dyEnd += s2dyF; dy <= dyEnd; dy++)
		{
			float *dyp = dp + dyI;
			dyI += dRowStride;

			__m128 res00 = lerpF2(c00, c01, _mm_set_ps1((dyEnd - dy) * d2syF));
			__m128 res10 = lerpF2(c10, c11, _mm_set_ps1((dyEnd - dy) * d2syF));

			// Left edge
			for (uint32_t dx = 0; dx + 1 < margin_x_start; dx++)
			{
				_mm_store_ps(dyp + dx * 4, res00);
			}

			// Right edge
			for (uint32_t dx = margin_x_end; dx > 0; dx--)
			{
				_mm_store_ps(dyp + dRowStride - dx * 4, res10);
			}
		}
	}

	float dx = 0, dxEnd = 0;
	uint64_t dxI = 4 * (margin_x_start == 0 ? 0 : margin_x_start - 1);
	// Interpolate top & bottom edges horizontally
	for (uint32_t sx = 0; sx < src.Width - 1; sx++)
	{
		float *sxp = sp + sx * 4;

		__m128 c00 = _mm_loadu_ps(sxp);
		__m128 c10 = _mm_loadu_ps(sxp + 4);
		__m128 c01 = _mm_loadu_ps(sxp + (src.Width - 1) * sRowStride);
		__m128 c11 = _mm_loadu_ps(sxp + (src.Width - 1) * sRowStride + 4);
		for (dxEnd += s2dxF; dx <= dxEnd; dx++)
		{
			float *dxp = dp + dxI;
			dxI += 4;

			__m128 res00 = lerpF2(c00, c10, _mm_set_ps1((dxEnd - dx) * d2sxF));
			__m128 res01 = lerpF2(c01, c11, _mm_set_ps1((dxEnd - dx) * d2sxF));

			// Top edge
			for (uint32_t dy = 0; dy + 1 < margin_y_start; dy++)
			{
				_mm_store_ps(dxp + dy * dRowStride, res00);
			}

			// Bottom edge
			for (uint32_t dy = margin_y_end; dy > 0; dy--)
			{
				_mm_store_ps(dxp + dst.Height * dRowStride - dy * dRowStride, res01);
			}
		}
	}

	// Corner colors
	__m128 c00 = _mm_loadu_ps(sp);
	__m128 c10 = _mm_loadu_ps(sp + sRowStride - 4);
	__m128 c01 = _mm_loadu_ps(sp + sRowStride * (src.Height - 1));
	__m128 c11 = _mm_loadu_ps(sp + sRowStride * src.Height - 4);

	float *d00 = dp;
	float *d10 = dp + dRowStride;
	float *d01 = dp + dst.Height * dRowStride;
	float *d11 = dp + dst.Height * dRowStride + dRowStride;
	// Fill top left corner
	for (uint32_t dy = 0; dy < margin_y_start; dy++)
		for (uint32_t dx = 0; dx < margin_x_start; dx++)
		{
			_mm_store_ps(d00 + dy * dRowStride + dx * 4, c00);
		}

	// Fill top right corner
	for (uint32_t dy = 0; dy < margin_y_start; dy++)
		for (uint32_t dx = margin_x_end; dx > 0; dx--)
		{
			_mm_store_ps(d10 + dy * dRowStride - dx * 4, c10);
		}

	// Fill bottom right corner
	for (uint32_t dy = margin_y_end; dy > 0; dy--)
		for (uint32_t dx = 0; dx < margin_x_start; dx++)
		{
			_mm_store_ps(d01 - dy * dRowStride + dx * 4, c01);
		}

	// Fill bottom right corner
	for (uint32_t dy = margin_y_end; dy > 0; dy--)
		for (uint32_t dx = margin_x_end; dx > 0; dx--)
		{
			_mm_store_ps(d11 - dy * dRowStride - dx * 4, c11);
		}
}

void print128_num(__m256i var)
{
	uint32_t val[8];
	memcpy(val, &var, sizeof(val));
	printf("%llu %llu %llu %llu %llu %llu %llu %llu\n",
		   val[0], val[1], val[2], val[3], val[4], val[5],
		   val[6], val[7]);
}

void ScaleBlerpCSimd(struct ValueFieldIC src, struct ValueFieldIC dst)
{
	uint32_t *sp = src.Data;
	uint32_t *dp = dst.Data;

	uint32_t send = src.Len;
	uint32_t srowstride = src.Width * 4;
	uint32_t drowstride = dst.Width * 4;

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

			for (uint32_t sx = 0; sx < srowstride - 4; sx += 4)
			{
				uint32_t *sxp = syp + sx;

				__m256i c00 = loadlo_4x32(sxp);
				__m256i c10 = loadlo_4x32(sxp + 4);
				__m256i c01 = loadlo_4x32(sxp + srowstride);
				__m256i c11 = loadlo_4x32(sxp + srowstride + 4);

				for (dxFend += sTxF; dxF < dxFend; dxF += dTxF)
				{
					__m256i res = blerpI256(c00, c10, c01, c11, _mm256_set1_epi32(fx), _mm256_set1_epi32(fy));
					storelo_4x32(dyp + dx, res);

					dx += 4;
					fx += d2sxF;
				}
			}

			dy += drowstride;
			fy += d2syF;
		}
	}
}

void ScaleBlerpCSimd2(struct ValueFieldIC src, struct ValueFieldIC dst)
{
	uint32_t *sp = src.Data;
	uint32_t *dp = dst.Data;

	uint32_t send = src.Len;
	uint32_t srowstride = src.Width * 4;
	uint32_t drowstride = dst.Width * 4;

	uint32_t dTxF = UINT32_MAX / dst.Width;
	uint32_t dTyF = UINT32_MAX / dst.Height;
	uint32_t sTxF = UINT32_MAX / (src.Width - 1);
	uint32_t sTyF = UINT32_MAX / (src.Height - 1);

	const __m256i ZERO = _mm256_setzero_si256();
	__m256i d2sxF = _mm256_set1_epi32((uint32_t)(((uint64_t)(src.Width - 1) * UINT32_MAX) / dst.Width));
	__m256i d2syF = _mm256_set1_epi32((uint32_t)(((uint64_t)(src.Height - 1) * UINT32_MAX) / dst.Height));

	uint64_t dyF = 0, dyFend = 0, dy = 0;
	__m256i fy = _mm256_setzero_si256(), fyi = _mm256_set1_epi32(UINT32_MAX);

	for (uint32_t sy = 0; sy < send - srowstride; sy += srowstride)
	{
		uint32_t *syp = sp + sy;

		for (dyFend += sTyF; dyF < dyFend; dyF += dTyF)
		{
			uint32_t *dyp = dp + dy;

			uint64_t dxF = 0, dxFend = 0, dx = 0;
			__m256i fx = _mm256_setzero_si256(), fxi = _mm256_set1_epi32(UINT32_MAX);

			for (uint32_t sx = 0; sx < srowstride - 4; sx += 4)
			{
				uint32_t *sxp = syp + sx;

				__m256i c00 = loadlo_4x322(sxp);
				__m256i c10 = loadlo_4x322(sxp + 4);
				__m256i c01 = loadlo_4x322(sxp + srowstride);
				__m256i c11 = loadlo_4x322(sxp + srowstride + 4);

				for (dxFend += sTxF; dxF < dxFend; dxF += dTxF)
				{
					__m256i res = blerpI2562(c00, c10, c01, c11, fx, fxi, fy, fyi);
					storelo_4x322(dyp + dx, res);

					dx += 4;
					fx = _mm256_add_epi32(fx, d2sxF);
					fxi = _mm256_sub_epi32(ZERO, fx);
				}
			}

			dy += drowstride;
			fy = _mm256_add_epi32(fy, d2syF);
			fyi = _mm256_sub_epi32(ZERO, fy);
		}
	}
}

void ScaleBlerpCSimd3(struct ValueFieldI8C src, struct ValueFieldIC dst)
{
	uint8_t *sp = src.Data;
	uint32_t *dp = dst.Data;

	uint32_t send = src.Len;
	uint32_t srowstride = src.Width * 4;
	uint32_t drowstride = dst.Width * 4;

	const uint32_t UINT12_MAX = (1 << 12) - 1;

	uint32_t dTxF = UINT12_MAX / dst.Width;
	uint32_t dTyF = UINT12_MAX / dst.Height;
	uint32_t sTxF = UINT12_MAX / (src.Width - 1);
	uint32_t sTyF = UINT12_MAX / (src.Height - 1);

	const __m128i ZERO = _mm_setzero_si128();
	__m128i d2sxF = _mm_set1_epi32((uint32_t)(((uint64_t)(src.Width - 1) * UINT12_MAX) / dst.Width));
	__m128i d2syF = _mm_set1_epi32((uint32_t)(((uint64_t)(src.Height - 1) * UINT12_MAX) / dst.Height));

	uint64_t dyF = 0, dyFend = 0, dy = 0;
	__m128i fy = _mm_setzero_si128(), fyi = _mm_set1_epi32(UINT12_MAX);

	for (uint32_t sy = 0; sy < send - srowstride; sy += srowstride)
	{
		uint8_t *syp = sp + sy;

		for (dyFend += sTyF; dyF < dyFend; dyF += dTyF)
		{
			uint32_t *dyp = dp + dy;

			uint64_t dxF = 0, dxFend = 0, dx = 0;
			__m128i fx = _mm_setzero_si128(), fxi = _mm_set1_epi32(UINT12_MAX);

			for (uint32_t sx = 0; sx < srowstride - 4; sx += 4)
			{
				uint8_t *sxp = syp + sx;

				__m128i c00 = loadlo_4x8(sxp);
				__m128i c10 = loadlo_4x8(sxp + 4);
				__m128i c01 = loadlo_4x8(sxp + srowstride);
				__m128i c11 = loadlo_4x8(sxp + srowstride + 4);

				for (dxFend += sTxF; dxF < dxFend; dxF += dTxF)
				{
					__m128i res = blerpI128(c00, c10, c01, c11, fx, fxi, fy, fyi);
					_mm_storeu_si128((__m128i *)(dyp + dx), res);

					dx += 4;
					fx = _mm_add_epi32(fx, d2sxF);
					fxi = _mm_sub_epi32(ZERO, fx);
				}
			}

			dy += drowstride;
			fy = _mm_add_epi32(fy, d2syF);
			fyi = _mm_sub_epi32(ZERO, fy);
		}
	}
}