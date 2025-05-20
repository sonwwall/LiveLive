package main

import (
	course "LiveLive/kitex_gen/livelive/course/courseservice"
	"log"
)

func main() {
	svr := course.NewServer(new(CourseServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
