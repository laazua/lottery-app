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
	IsRandomMode bool  // true:随机注, false:胆拖
	FrontDare    []int // 前区胆码
	FrontDrag    []int // 前区拖码
	BackDare     []int // 后区胆码
	BackDrag     []int // 后区拖码
}

// 获取随机数生成器
func getRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 随机生成一注
func GenerateRandomTicket() LotteryResult {
	r := getRand()

	// 生成前区
	front := generateRandomNumbers(r, FrontMin, FrontMax, FrontCount)
	// 生成后区
	back := generateRandomNumbers(r, BackMin, BackMax, BackCount)

	return LotteryResult{
		FrontNumbers: front,
		BackNumbers:  back,
		IsRandomMode: true,
	}
}

// 随机生成胆拖号码
// frontDareCount: 前区胆码个数(1-4)
// frontDragCount: 前区拖码个数(至少需要FrontCount-frontDareCount个)
// backDareCount: 后区胆码个数(1)
// backDragCount: 后区拖码个数(至少需要BackCount-backDareCount个)
func GenerateRandomDanTuo(frontDareCount, frontDragCount, backDareCount, backDragCount int) (LotteryResult, error) {
	// 验证规则
	if frontDareCount < 1 || frontDareCount > 4 {
		return LotteryResult{}, strconv.ErrSyntax
	}
	if backDareCount != 1 {
		return LotteryResult{}, strconv.ErrSyntax
	}
	if frontDareCount+frontDragCount < FrontCount {
		return LotteryResult{}, strconv.ErrSyntax
	}
	if backDareCount+backDragCount < BackCount {
		return LotteryResult{}, strconv.ErrSyntax
	}

	r := getRand()

	// 生成所有可用数字池
	frontPool := make([]int, FrontMax-FrontMin+1)
	for i := range frontPool {
		frontPool[i] = FrontMin + i
	}

	backPool := make([]int, BackMax-BackMin+1)
	for i := range backPool {
		backPool[i] = BackMin + i
	}

	// 随机打乱
	r.Shuffle(len(frontPool), func(i, j int) {
		frontPool[i], frontPool[j] = frontPool[j], frontPool[i]
	})
	r.Shuffle(len(backPool), func(i, j int) {
		backPool[i], backPool[j] = backPool[j], backPool[i]
	})

	// 选择胆码
	frontDare := frontPool[:frontDareCount]
	backDare := backPool[:backDareCount]

	// 从剩余号码中选择拖码
	remainingFront := frontPool[frontDareCount:]
	remainingBack := backPool[backDareCount:]

	frontDrag := remainingFront[:frontDragCount]
	backDrag := remainingBack[:backDragCount]

	// 合并成完整号码
	frontNumbers := make([]int, len(frontDare))
	copy(frontNumbers, frontDare)
	frontNumbers = append(frontNumbers, frontDrag...)
	sort.Ints(frontNumbers)

	backNumbers := make([]int, len(backDare))
	copy(backNumbers, backDare)
	backNumbers = append(backNumbers, backDrag...)
	sort.Ints(backNumbers)

	// 对胆码和拖码分别排序
	sort.Ints(frontDare)
	sort.Ints(frontDrag)
	sort.Ints(backDare)
	sort.Ints(backDrag)

	return LotteryResult{
		FrontNumbers: frontNumbers,
		BackNumbers:  backNumbers,
		IsRandomMode: false,
		FrontDare:    frontDare,
		FrontDrag:    frontDrag,
		BackDare:     backDare,
		BackDrag:     backDrag,
	}, nil
}

// 生成不重复的随机数
func generateRandomNumbers(r *rand.Rand, min, max, count int) []int {
	pool := make([]int, max-min+1)
	for i := range pool {
		pool[i] = min + i
	}

	r.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	result := pool[:count]
	sort.Ints(result)
	return result
}

// 格式化显示（普通随机模式）
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

// 格式化显示胆拖模式
func FormatDanTuoNumbers(result LotteryResult) string {
	str := "前区:\n"

	// 显示胆码
	str += "   胆码: "
	for i, num := range result.FrontDare {
		if i > 0 {
			str += "  "
		}
		str += formatNumber(num)
	}

	// 显示拖码
	str += "\n   拖码: "
	for i, num := range result.FrontDrag {
		if i > 0 {
			str += "  "
		}
		str += formatNumber(num)
	}

	// 后区
	str += "\n后区:\n"
	str += "   胆码: "
	for i, num := range result.BackDare {
		if i > 0 {
			str += "  "
		}
		str += formatNumber(num)
	}

	str += "\n   拖码: "
	for i, num := range result.BackDrag {
		if i > 0 {
			str += "  "
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
