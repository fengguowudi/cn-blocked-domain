package main

import (
	"log"
	"strings"

	"github.com/Loyalsoldier/cn-blocked-domain/utils"
)

func buildTreeAndUnique(sortedDomainList []string) []string {
	tree := newList()
	remainList := make([]string, 0, len(sortedDomainList))

	for _, d := range sortedDomainList {
		parts := strings.Split(d, ".")
		l, i, e := tree.Insert(parts)

		if e != nil {
			log.Println(utils.Fatal("[Error]"), "check domain", utils.Info(d), "for redundancy.")
			continue
		}

		if !i {
			r := make([]string, 0, len(parts))

			for x := 0; x <= l; x++ {
				r = append(r, parts[x])
			}

			redundantStr := strings.Join(r, ".")
			log.Println("Found redundant domain:", utils.Info(d), "@", utils.Warning(redundantStr))
			continue
		}

		remainList = append(remainList, d)
	}

	return remainList
}
