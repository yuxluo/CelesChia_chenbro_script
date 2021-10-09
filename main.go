package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func findPlotName(dir string) string {
	out, _ := exec.Command("ls").Output()
	s := strings.Split(string(out), "\n")
	for _, fileName := range s {
		if fileName[len(fileName) - 4:] == "plot" {
			return fileName
		}
	}
	return ""
}

func findEmpty() string {
	out, err := exec.Command("df").Output()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	//fmt.Printf("df output is \n%s\n", out)
	var dfOutput = string(out)


	s := strings.Split(dfOutput, "\n")
	for index, partition := range s {
		if index == 0 {
			continue
		}
		items := strings.Split(partition, " ")
		var withoutSpace []string
		for _, item := range items {
			if item != "" {
				withoutSpace = append(withoutSpace, item)
			}
		}
		intCapacity, err := strconv.Atoi(withoutSpace[1])
		intRemaining, err := strconv.Atoi(withoutSpace[3])
		if err != nil {

		}
		if intCapacity > 999999999 && intRemaining > 100999999 {
			return withoutSpace[5]
		}
	}
	return ""
}



func main() {
	var muji_ssd_path = "~/plots/"
	var nextEmptyDirve string
	fmt.Println("請確保母雞SSD已經mount好，不然要吃屎")

	for true {
		//	找到空盤
		nextEmptyDirve = findEmpty()
		if nextEmptyDirve == "" {
			fmt.Println("Either 吃屎了 or 全部盤都已經裝滿")
			break
		} else {
			for true {
				plotName := findPlotName(muji_ssd_path)
				if plotName == "" {
					fmt.Println("母雞沒有plot, 等待1分鐘....")
					time.Sleep(1 * time.Minute)
				} else {
					_, _ = exec.Command("mv", muji_ssd_path+plotName, nextEmptyDirve).Output()
					println("Transferred %s to %s", plotName, nextEmptyDirve)
					_, _ = exec.Command("rm", muji_ssd_path+plotName).Output()
					break
				}
			}
		}
		break
	}
	return
}