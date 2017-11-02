package logcount

import (
	"sort"
	"os"
	"bufio"
	"strings"
	"flag"
	"strconv"
)
var MaxIPCount = flag.Int("n", 50, "количество IP-адресов")
//флаг для сортировки по IP-адресам, а не количеству запросов
var GroupByIp = flag.Bool("g", false, "сортировка по IP-адресам")
//флаг для вывода детализированного отчета по IP-адресам
var DetReport = flag.Bool("d", false, "детализированный отчет")
var StatusFlag = flag.Bool("s", false, "выводить статистику по кодам ответа")

type Request struct {
	Site, IP, Method, Uri, Status string
	ResponseBytes int64
}

type Requests []Request

type IpCount struct {
	IP string
	Count int
	StatusCode map[string]int
	StatusData map[string]int64
}


func (rr *Requests) SortByIP() []IpCount {

	ipMap := make(map[string]int)
	statusMap := make(map[string]map[string]int)
	statusDataMap := make(map[string]map[string]int64)

	for _, r := range *rr {
		ipMap[r.IP]++

		if statusMap[r.IP] == nil {
			statusMap[r.IP] = make(map[string]int)
			}
		statusMap[r.IP][r.Status]++

		if statusDataMap[r.IP] == nil {
			statusDataMap[r.IP] = make(map[string]int64)
			}
		statusDataMap[r.IP][r.Status] += r.ResponseBytes
		}

	var ipCnt []IpCount

	for ip, count := range ipMap {
		st := statusMap[ip]
		std := statusDataMap[ip]
		ipCnt = append(ipCnt, IpCount{ip,count, st, std} )
	}

	sort.Slice(ipCnt, func(i, j int) bool {
		return ipCnt[i].Count > ipCnt[j].Count
	})

	return ipCnt
}


//чтение из файла в массив строк
func ReadLines(path string) (Requests, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var rLog Requests
	var r Request
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" && strings.Count(scanner.Text(), " ") > 9 {
			line := strings.Split(scanner.Text(), " ")
			r.Site = line[0]
			r.IP = line[1]
			r.Method = strings.Trim(line[6], "\"")
			r.Uri = line[7]
			r.Status = line[9]

			r.ResponseBytes, err = strconv.ParseInt(line[10], 10, 64)
			if err != nil {
				r.ResponseBytes = 0
			}

			rLog = append(rLog, r)
		}
	}
	return rLog, scanner.Err()
}

/*func searchIp(ipr *[]RSorted, ip string) (index *int, ok bool) {
	for i, r := range *ipr {
		if r.IP == ip {
			return &i, true
		}
		continue
	}
	return nil, false

}

func searchRequestData(rData *[]RequestData, method string, url string, status string) (ok bool) {
	for i, d := range *rData {
		if d.Method == method && d.Url == url && d.Status == status {
			(*rData)[i].RCount++
			sort.Slice(*rData, func(i, j int) bool {
				return (*rData)[i].RCount > (*rData)[j].RCount
			})
			return true
		}
		continue
	}
	rd := RequestData{method, url, status, 1}
	*rData = append(*rData, rd)
	return false
}

//возвращает строку из n-пробелов в зависимости от длины строки параметра val
func InsSpace(elemlen int, maxelemlen int) string {
	space := ""
	for i := 0; i <= (maxelemlen - elemlen); i++ {
		space += " "
	}
	return space
}*/

