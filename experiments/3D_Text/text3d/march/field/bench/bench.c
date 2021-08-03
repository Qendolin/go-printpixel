#include <time.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <immintrin.h>
#include <stdint.h>

struct ValueFieldIC
{
	uint32_t Len;
	uint32_t Width;
	uint32_t Height;
	uint32_t *Data;
};

static inline __m256i lerpI2562(__m256i a, __m256i b, __m256i t, __m256i ti)
{
	return _mm256_srli_epi64(_mm256_add_epi32(_mm256_mul_epu32(a, ti), _mm256_mul_epu32(b, t)), 32);
}

static inline __m256i blerpI2562(__m256i c00, __m256i c10, __m256i c01, __m256i c11, __m256i fx, __m256i fxi, __m256i fy, __m256i fyi)
{
	return lerpI2562(
		lerpI2562(c00, c10, fx, fxi),
		lerpI2562(c01, c11, fx, fxi),
		fy, fyi);
}

static inline void storelo_4x32(uint32_t *dst, __m256i src)
{
	const __m256i K_PERM = _mm256_setr_epi32(0, 2, 4, 6, 1, 3, 5, 7);
	__m256i permuted = _mm256_permutevar8x32_epi32(src, K_PERM);
	__m128i lo128 = _mm256_extractf128_si256(permuted, 0);
	_mm_storeu_si128((__m128i *)dst, lo128);
}

static inline __m256i loadlo_4x322(uint32_t *src)
{
	// Since high 128 bits are undefined this probably has undefined behavoir
	__m256i data = _mm256_castsi128_si256(_mm_loadu_si128((__m128i *)src));
	const __m256i K_PERM = _mm256_setr_epi32(0, 4, 1, 5, 2, 6, 3, 7);
	return _mm256_permutevar8x32_epi32(data, K_PERM);
}

static void ScaleBlerpCSimd2(struct ValueFieldIC src, struct ValueFieldIC dst)
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
					storelo_4x32((dyp + dx), res);

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

uint32_t rand32()
{
	uint32_t x = rand() & 0xff;
	x |= (rand() & 0xff) << 8;
	x |= (rand() & 0xff) << 16;
	x |= (rand() & 0xff) << 24;

	return x;
}

int main(void)
{
	struct ValueFieldIC src = {
		.Height = 37,
		.Width = 37,
		.Len = 37 * 37 * 4,
	};
	src.Data = (uint32_t *)calloc(src.Len, sizeof(uint32_t));
	if (!src.Data)
	{
		printf("calloc failed\n");
		exit(-1);
	}

	srand(1);
	for (int i = 0; i < src.Len; i++)
	{
		src.Data[i] = rand32();
	}

	struct ValueFieldIC dst = {
		.Height = 37 * 8,
		.Width = 37 * 8,
		.Len = 37 * 8 * 37 * 8 * 4,
	};
	dst.Data = (uint32_t *)calloc(dst.Len, sizeof(uint32_t));
	if (!dst.Data)
	{
		printf("calloc failed\n");
		exit(-1);
	}

	clock_t start_time = clock();
	for (int i = 0; i < 1000; i++)
	{
		ScaleBlerpCSimd2(src, dst);
	}
	double elapsed_time = (double)(clock() - start_time) / CLOCKS_PER_SEC;
	double ns_per_iter = 1e9 * (elapsed_time / 1000);
	printf("%.0f ns / op\n", elapsed_time);
}