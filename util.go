package oss

import (
	"net/url"
	"path"
	"strings"

	"github.com/prismgo/framework/filesystem"
)

// normalizeKey 把业务层传入的相对路径规范为统一的 slash 形式。
func normalizeKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, "\\", "/")
	key = strings.TrimPrefix(key, "/")
	key = path.Clean("/" + key)
	key = strings.TrimPrefix(key, "/")
	if key == "." {
		return ""
	}
	return key
}

// normalizeDir 保证目录前缀统一以单个斜杠结尾。
func normalizeDir(key string) string {
	key = normalizeKey(key)
	if key == "" {
		return ""
	}
	return strings.TrimSuffix(key, "/") + "/"
}

// joinKey 把目录前缀和文件名拼接成磁盘内的相对路径。
func joinKey(prefix, key string) string {
	prefix = normalizeDir(prefix)
	key = normalizeKey(key)
	switch {
	case prefix == "":
		return key
	case key == "":
		return strings.TrimSuffix(prefix, "/")
	default:
		return prefix + key
	}
}

// joinURL 把 URL 前缀与相对路径安全拼接为可访问地址。
func joinURL(base string, key string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	key = normalizeKey(key)
	if base == "" {
		return "/" + key
	}
	if key == "" {
		return base
	}
	segments := strings.Split(key, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return base + "/" + strings.Join(segments, "/")
}

// ensureVisibility 把可见性值归一化，并在非法值时回退到默认值。
func ensureVisibility(visibility, fallback string) string {
	switch strings.ToLower(strings.TrimSpace(visibility)) {
	case filesystem.VisibilityPublic:
		return filesystem.VisibilityPublic
	case filesystem.VisibilityPrivate:
		return filesystem.VisibilityPrivate
	default:
		return strings.ToLower(strings.TrimSpace(fallback))
	}
}

// rejectEmptyDirectory 在破坏性目录操作前拒绝空目录，避免把根目录当作目标删除。
func rejectEmptyDirectory(dir string) error {
	if normalizeDir(dir) == "" {
		return filesystem.ErrEmptyDirectory
	}
	return nil
}
