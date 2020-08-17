package utils

import "math"

//分页方法，根据传递过来的当前页数，每页数，总数
func Paginator(page, pageNum int, nums int64) map[string]interface{} {

	var firstpage int //前一页地址
	var lastpage int  //后一页地址
	//根据nums总数，和pageNum每页数量 生成分页总数
	totalpages := int(math.Ceil(float64(nums) / float64(pageNum))) //page总数
	if page > totalpages {
		page = totalpages
	}
	if page <= 0 {
		page = 1
	}
	var num_pages []int
	switch {
	case page >= totalpages-5 && totalpages > 5: //最后5页
		start := totalpages - 5 + 1
		firstpage = page - 1
		lastpage = int(math.Min(float64(totalpages), float64(page+1)))
		num_pages = make([]int, 5)
		for i, _ := range num_pages {
			num_pages[i] = start + i
		}
	case page >= 3 && totalpages > 5:
		start := page - 3 + 1
		num_pages = make([]int, 5)
		firstpage = page - 3
		for i, _ := range num_pages {
			num_pages[i] = start + i
		}
		firstpage = page - 1
		lastpage = page + 1
	default:
		num_pages = make([]int, int(math.Min(5, float64(totalpages))))
		for i, _ := range num_pages {
			num_pages[i] = i + 1
		}
		firstpage = int(math.Max(float64(1), float64(page-1)))
		lastpage = page
	}
	paginatorMap := make(map[string]interface{})
	// 中间的页码数字
	paginatorMap["num_pages"] = num_pages
	// 总页数
	paginatorMap["totalpages"] = totalpages
	// 首页
	paginatorMap["firstpage"] = firstpage
	// 尾页
	paginatorMap["lastpage"] = lastpage
	// 当前页
	paginatorMap["currpage"] = page
	return paginatorMap
}
