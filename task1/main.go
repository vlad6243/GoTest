package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Holiday struct for json ...
type Holiday struct {
	Date      string
	LocalName string
	Name      string
}

// Weekennd struct for json ...
type LongWeekend struct {
	StartDate string
	EndDate   string
	DayCount  int
}

func main() {

	const URLHOLIDAY string = "https://date.nager.at/api/v2/NextPublicHolidays/Ua"
	// const URLWEEKENDS string = "https://date.nager.at/Api/v2/LongWeekend/2020/UA"

	// var holidays = []Holiday{
	// 	{
	// 		Date:      "2020-08-24",
	// 		Name:      "some",
	// 		LocalName: "something",
	// 	},
	// }
	var holidays []Holiday

	var weekends []LongWeekend

	if err := json.Unmarshal([]byte(getHoliday(URLHOLIDAY)), &holidays); err != nil {
		fmt.Println(err)
	}

	weekends = getLongWeekends(holidays)

	// if err := json.Unmarshal([]byte(getHoliday(URLWEEKENDS)), &weekends); err != nil {
	// 	fmt.Println(err)
	// }

	t := time.Now()
	if holidays[0].Date == t.Format("2006-01-02") {

		fmt.Println("Today holiday " + holidays[0].Name)

	} else {

		flag := false
		//parse holiday
		ph, _ := time.Parse("2006-01-02", holidays[0].Date)

		for i := 0; i < len(weekends); i++ {

			flag = true

			if (holidays[0].Date == weekends[i].StartDate) || (holidays[0].Date == weekends[i].EndDate) {

				flag = false

				//parse weekends start
				pws, _ := time.Parse("2006-01-02", weekends[i].StartDate)
				//parse weekends end
				pwe, _ := time.Parse("2006-01-02", weekends[i].EndDate)

				fmt.Println("The next holiday is " + holidays[0].Name + ", " + ph.Format("Jan-02") + ",")
				fmt.Printf("and the weekend will last %d days:", weekends[i].DayCount)
				fmt.Print(pws.Format("Jan-02") + " - " + pwe.Format("Jan-02"))

				break
			}
		}

		if flag {
			fmt.Println("Next holiday " + holidays[0].Name + ", " + ph.Format("Jan-02"))
		}

	}

}

func getHoliday(url string) string {

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", url, nil,
	)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(body)

}

func getLongWeekends(holidays []Holiday) []LongWeekend {

	var weekends []LongWeekend

	for i := 0; i < len(holidays); i++ {

		d, _ := time.Parse("2006-01-02", holidays[i].Date)

		if d.Weekday() == 1 {

			newWeekend := LongWeekend{
				StartDate: d.AddDate(0, 0, -2).Format("2006-01-02"),
				EndDate:   holidays[i].Date,
				DayCount:  3,
			}

			weekends = append(weekends, newWeekend)
		}

		if d.Weekday() == 2 {

			newWeekend := LongWeekend{
				StartDate: d.AddDate(0, 0, -3).Format("2006-01-02"),
				EndDate:   holidays[i].Date,
				DayCount:  3,
			}

			weekends = append(weekends, newWeekend)
		}

		if d.Weekday() == 4 {

			newWeekend := LongWeekend{
				StartDate: holidays[i].Date,
				EndDate:   d.AddDate(0, 0, 3).Format("2006-01-02"),
				DayCount:  4,
			}

			weekends = append(weekends, newWeekend)
		}

		if d.Weekday() == 5 {

			newWeekend := LongWeekend{
				StartDate: holidays[i].Date,
				EndDate:   d.AddDate(0, 0, 2).Format("2006-01-02"),
				DayCount:  3,
			}

			weekends = append(weekends, newWeekend)
		}

	}

	return weekends
}
