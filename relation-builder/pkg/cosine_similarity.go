package pkg

import "math"

func CosineSimilarity(vec1, vec2 []float64) float64 {
	var normAB, normA, normB float64
	for i := 0; i < len(vec1); i++ {
		normAB += vec1[i] * vec2[i]
		normA += vec1[i] * vec1[i]
		normB += vec2[i] * vec2[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return normAB / (math.Sqrt(normA) * math.Sqrt(normB))
}
