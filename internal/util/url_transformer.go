package util

import "net/url"

// BuildAbsoluteURL 函数接受站点域名和路径，返回绝对 URL
func BuildAbsoluteURL(base string, path string) (string, error) {
	// 解析基础 URL
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	// 解析路径
	relativeURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	// 如果路径是绝对路径，直接返回
	if relativeURL.IsAbs() {
		return relativeURL.String(), nil
	}

	// 将相对路径与基础 URL 合并
	absoluteURL := baseURL.ResolveReference(relativeURL)
	return absoluteURL.String(), nil
}
