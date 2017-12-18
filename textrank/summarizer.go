package textrank

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"sort"
)



const SUMMARY_ADJUST = 0.85

// n is the sentence length
func (self *Article)getSimilarityMatrix(n int) (* mat.Dense) {
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

func (self *Article) Summarizer() {
	const ITER = 10
	n := len(self.Sentences)
	similarityMartix := self.getSimilarityMatrix(n)

	PR := mat.NewDense(n, 1,  make([]float64, n))
	for c := 0; c < ITER; c++ {
		// PR = 0.15 + 0.85 * M * PR
		// and similarityMatrix is 0.85*M
		var mMulPR = new(mat.Dense)
		mMulPR.Mul(similarityMartix, PR)
		PR.Add(wholeMatrix(n, 1, (1 - SUMMARY_ADJUST)), mMulPR)

		// if converagence, break, the threshold is 0.0001
	}
	// sort sentences points
	result := PR.RawMatrix().Data
	result_map := make(ByPoint, len(result))

	for p, idf := range result {
		result_map[p] = PointMap{p, idf}
	}
	sort.Sort(sort.Reverse(result_map))
	for _, v := range result_map {
		fmt.Print(v)
		//fmt.Print(self.Sentences[v.Key].content + fmt.Sprintf("%.3f", v.Value) + "\n")
	}
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