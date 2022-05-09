package search

import (
	"fmt"
	"math"
	"time"
	"trap-search/db"
	"trap-search/twitter"
)

func FromCreator(userId string) {
	token := twitter.GetBearerToken()
	timeLine := twitter.GetTimeLine(token, userId, map[string]string{
		"max_results": "100",
		"exclude":     "retweets,replies",
		"expansions":  "attachments.media_keys",
	})
	tweets := []twitter.Tweet{}
	for _, val := range timeLine.Data {
		if val.Attach != nil {
			tweets = append(tweets, val)
		}
	}

	fmt.Printf("%+v\n", tweets)
	for _, val := range tweets {
		var likes []string
		for {
			var err error
			likes, err = twitter.GetLikes(token, val.Id, map[string]string{
				"max_results": "100",
			})
			if err == nil || len(likes) >= 5000 {
				break
			}
			fmt.Println(err.Error())
			time.Sleep(30 * time.Second)
		}
		coef := 0.0
		users := []db.User{}
		for _, v := range likes {
			var u db.User
			user := []db.User{}
			db.Select(&user, "select * from user where id = ?", v)
			if len(user) == 0 {
				u = db.User{
					Id:          v,
					Name:        "hoge",
					Coef:        0.5,
					CreatorCoef: 0.1,
				}
				db.Push("insert into user values(?,?,?,?)", u.Id, u.Name, u.Coef, u.CreatorCoef)
				coef += 0.4999999999999
			} else {
				u = user[0]
				coef += u.Coef
			}
			users = append(users, u)
		}
		fmt.Print(val.Text + ":")
		if coef/float64(len(likes)) >= 0.5 {
			fmt.Println("This is trap")
			for _, v := range users {
				db.Push("UPDATE user SET trap_coef=? WHERE id = ?", math.Min(v.Coef+0.01, 1.0), v.Id)
			}
		} else {
			fmt.Println("This is not trap")
			for _, v := range users {
				db.Push("UPDATE user SET trap_coef=? WHERE id = ?", math.Max(v.Coef-0.01, 0.0), v.Id)
			}
		}
	}
}
