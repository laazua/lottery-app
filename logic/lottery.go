package logic

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// 大乐透规则：前区1-35选5个，后区1-12选2个
const (
	FrontMin   = 1
	FrontMax   = 35
	FrontCount = 5
	BackMin    = 1
	BackMax    = 12
	BackCount  = 2
)

type LotteryResult struct {
	FrontNumbers []int
	BackNumbers  []int
	IsRandomMode bool // true:随机注, false:胆拖
}

// 随机生成一注
func GenerateRandomTicket() LotteryResult {
	rand.Seed(time.Now().UnixNano())

	// 生成前区
	front := generateRandomNumbers(FrontMin, FrontMax, FrontCount)
	// 生成后区
	back := generateRandomNumbers(BackMin, BackMax, BackCount)

	return LotteryResult{
		FrontNumbers: front,
		BackNumbers:  back,
		IsRandomMode: true,
	}
}

// 胆拖方式生成
// frontDare: 前区胆码，frontDrag: 前区拖码
// backDare: 后区胆码，backDrag: 后区拖码
func GenerateDanTuoTicket(frontDare, frontDrag, backDare, backDrag []int) []LotteryResult {
	var results []LotteryResult

	// 验证胆码数量（前区胆码1-4个，后区胆码1个）
	if len(frontDare) < 1 || len(frontDare) > 4 {
		return results
	}
	if len(backDare) != 1 {
		return results
	}

	// 计算需要选择的号码数量
	frontNeed := FrontCount - len(frontDare)
	backNeed := BackCount - len(backDare)

	// 从前区拖码中选择frontNeed个号码的组合
	frontCombos := combinations(frontDrag, frontNeed)
	// 从后区拖码中选择backNeed个号码的组合
	backCombos := combinations(backDrag, backNeed)

	// 组合所有可能
	for _, frontCombo := range frontCombos {
		for _, backCombo := range backCombos {
			// 合并胆码和拖码
			frontNumbers := make([]int, len(frontDare))
			copy(frontNumbers, frontDare)
			frontNumbers = append(frontNumbers, frontCombo...)
			sort.Ints(frontNumbers)

			backNumbers := make([]int, len(backDare))
			copy(backNumbers, backDare)
			backNumbers = append(backNumbers, backCombo...)
			sort.Ints(backNumbers)

			results = append(results, LotteryResult{
				FrontNumbers: frontNumbers,
				BackNumbers:  backNumbers,
				IsRandomMode: false,
			})
		}
	}

	return results
}

// 生成不重复的随机数
func generateRandomNumbers(min, max, count int) []int {
	pool := make([]int, max-min+1)
	for i := range pool {
		pool[i] = min + i
	}

	rand.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	result := pool[:count]
	sort.Ints(result)
	return result
}

// 组合数计算（C(n,m)）
func combinations(nums []int, k int) [][]int {
	result := [][]int{}
	if k == 0 {
		return result
	}
	if k > len(nums) {
		return result
	}

	var backtrack func(start int, current []int)
	backtrack = func(start int, current []int) {
		if len(current) == k {
			temp := make([]int, len(current))
			copy(temp, current)
			result = append(result, temp)
			return
		}
		for i := start; i < len(nums); i++ {
			backtrack(i+1, append(current, nums[i]))
		}
	}

	backtrack(0, []int{})
	return result
}

// 格式化显示
func FormatNumbers(front, back []int) string {
	str := "前区: "
	for i, num := range front {
		if i > 0 {
			str += " "
		}
		str += formatNumber(num)
	}
	str += "  后区: "
	for i, num := range back {
		if i > 0 {
			str += " "
		}
		str += formatNumber(num)
	}
	return str
}

func formatNumber(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
