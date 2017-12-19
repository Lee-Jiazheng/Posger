package textrank

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"sort"
)



const SUMMARY_ADJUST = 0.85

func (self *Article) Summarizer() []string{
	const ITER = 10
	n := len(self.Sentences)
	similarityMartix := self.getSimilarityMatrix()

	PR := mat.NewDense(n, 1,  make([]float64, n))
	for c := 0; c < ITER; c++ {
		// PR = 0.15 + 0.85 * M * PR
		// and similarityMatrix is 0.85*M
		var mMulPR = new(mat.Dense)
		mMulPR.Mul(similarityMartix, PR)
		PR.Add(wholeMatrix(n, 1, 1 - SUMMARY_ADJUST), mMulPR)
		// if converagence, break, the threshold is 0.0001
	}

	resultMap := sortSentences(PR)
	result := make([]string, n)
	for i, v := range *resultMap {
		// assemble all the words, append to the result.
		result[i] = connectSliceWords(self.Sentences[v.Key].Words...)
	}
	return result
}

// n is the sentence length
func (self *Article)getSimilarityMatrix() (* mat.Dense) {
	n := len(self.Sentences)
	data := make([]float64, n * n)
	for i1, s1 := range self.Sentences {
		s1 := s1
		for i2, s2 := range self.Sentences {
			s2 := s2
			if i1 == i2 {
				// diagonal Set Zero
				data[i1*n + i2] = 0
			} else {
				data[i1*n + i2] = s1.BM25(&s2) * SUMMARY_ADJUST
			}
		}
	}
	return mat.NewDense(n, n, data)
}

func wholeMatrix(raw int, col int, num float64) *mat.Dense{
	data := make([]float64, raw, col)
	for i, _ := range data{
		data[i] = num
	}
	return mat.NewDense(raw, col, data)
}

func sortSentences(PR *mat.Dense) *ByPoint{
	// sort sentences points
	result_map := make(ByPoint, len(PR.RawMatrix().Data))
	for p, idf := range PR.RawMatrix().Data {
		result_map[p] = PointMap{p, idf}
	}
	sort.Sort(sort.Reverse(result_map))
	return &result_map
}

func (self *Sentence) BM25(s *Sentence) float64 {
	const k1, b, avsl = 2.0, 0.75, 5.0
	var score float64
	for _, word := range self.Words {
		var tf int
		for _, wordc := range s.Words {
			if word == wordc {
				tf += 1
			}
		}
		var B float64 = (1 - b) + b * (float64(len(s.Words)) / avsl)
		score += GetWordIDF(word) * (float64(tf*(k1+1)) / (float64(tf)+k1*B))
	}
	return score
}