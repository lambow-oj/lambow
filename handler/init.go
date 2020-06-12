package handler

import "git/lambow-oj/lambow/cinex"

func Init() {
	opt := cinex.PostOpt{}

	cinex.RegPostHandler("/lambow/v1/ping", pingHandler, opt)
}
