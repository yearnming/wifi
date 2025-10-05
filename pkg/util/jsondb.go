package util

import (
	"encoding/json"
	"os"
	"sync"
)

// ===== 对外 API =====

type DB struct {
	path string
	mu   sync.RWMutex
	data map[string]*Entry // SSID -> 条目
}

type Entry struct {
	Failed   map[string]struct{} `json:"failed"`
	Abnormal map[string]struct{} `json:"abnormal"`
	Success  map[string]struct{} `json:"success"`
}

// New 创建实例（未加载文件）
func New(path string) *DB {
	return &DB{
		path: path,
		data: make(map[string]*Entry),
	}
}

// Load 读取三分类 JSON（空文件则视为空库）
// func (db *DB) Load() error {
// 	db.mu.Lock()
// 	defer db.mu.Unlock()

// 	f, err := os.Open(db.path)
// 	if os.IsNotExist(err) {
// 		return nil // 空库
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

//		// 磁盘结构：map[string]*Entry
//		raw := make(map[string]*Entry)
//		if err := json.NewDecoder(f).Decode(&raw); err != nil {
//			return err
//		}
//		db.data = raw
//		return nil
//	}
func (db *DB) Load() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	type fileEntry struct {
		Failed   []string `json:"failed"`
		Abnormal []string `json:"abnormal"`
		Success  []string `json:"success"`
	}
	raw := make(map[string]*fileEntry)

	f, err := os.Open(db.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return err
	}

	// 磁盘切片 -> 内存 map
	db.data = make(map[string]*Entry, len(raw))
	for ssid, fe := range raw {
		db.data[ssid] = &Entry{
			Failed:   sliceToSet(fe.Failed),
			Abnormal: sliceToSet(fe.Abnormal),
			Success:  sliceToSet(fe.Success),
		}
	}
	return nil
}

func sliceToSet(ss []string) map[string]struct{} {
	m := make(map[string]struct{}, len(ss))
	for _, v := range ss {
		m[v] = struct{}{}
	}
	return m
}

// Save 整表覆盖写盘（程序退出时调用一次）
// func (db *DB) Save() error {
// 	db.mu.RLock()
// 	defer db.mu.RUnlock()

// 	f, err := os.Create(db.path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

//		enc := json.NewEncoder(f)
//		enc.SetIndent("", "  ")
//		return enc.Encode(db.data)
//	}
//
// 内存不变，磁盘用 []string
type fileEntry struct { // 仅落盘结构
	Failed   []string `json:"failed"`
	Abnormal []string `json:"abnormal"`
	Success  []string `json:"success"`
}

func (db *DB) Save() error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// 内存 map -> 磁盘切片
	raw := make(map[string]*fileEntry, len(db.data))
	for ssid, ent := range db.data {
		raw[ssid] = &fileEntry{
			Failed:   keys(ent.Failed),
			Abnormal: keys(ent.Abnormal),
			Success:  keys(ent.Success),
		}
	}

	f, err := os.Create(db.path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(raw)
}

func keys(m map[string]struct{}) []string {
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

// 三分类追加（内存级别，不触盘）
func (db *DB) AddFail(ssid, password string) {
	db.add(ssid, password, "failed")
}
func (db *DB) AddAbnormal(ssid, password string) {
	db.add(ssid, password, "abnormal")
}
func (db *DB) AddSuccess(ssid, password string) {
	db.add(ssid, password, "success")
}

// FilterFresh 返回“从未出现过”的密码（三分类都没试过）
func (db *DB) FilterFresh(ssid string, dict []string) []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	ent := db.entry(ssid)
	var out []string
	for _, p := range dict {
		if !ent.has(p) {
			out = append(out, p)
		}
	}
	return out
}

// FilterFailed 去掉 failed+abnormal，保留 success 也可选
func (db *DB) FilterFailed(ssid string, dict []string) []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	ent := db.entry(ssid)
	var out []string
	for _, p := range dict {
		if _, ok := ent.Failed[p]; ok {
			continue
		}
		if _, ok := ent.Abnormal[p]; ok {
			continue
		}
		out = append(out, p)
	}
	return out
}

// ===== 内部辅助 =====

func (db *DB) add(ssid, password, category string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	ent := db.entry(ssid)
	switch category {
	case "failed":
		ent.Failed[password] = struct{}{}
	case "abnormal":
		ent.Abnormal[password] = struct{}{}
	case "success":
		ent.Success[password] = struct{}{}
	}
}

// 获取或创建 Entry
func (db *DB) entry(ssid string) *Entry {
	if db.data[ssid] == nil {
		db.data[ssid] = &Entry{
			Failed:   make(map[string]struct{}),
			Abnormal: make(map[string]struct{}),
			Success:  make(map[string]struct{}),
		}
	}
	return db.data[ssid]
}

// 判断密码是否在任何分类出现过
func (e *Entry) has(password string) bool {
	_, ok1 := e.Failed[password]
	_, ok2 := e.Abnormal[password]
	_, ok3 := e.Success[password]
	return ok1 || ok2 || ok3
}

// GetSuccess 返回该 SSID 所有成功密码（切片拷贝）
func (db *DB) GetSuccess(ssid string) []string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	ent := db.entry(ssid)
	out := make([]string, 0, len(ent.Success))
	for p := range ent.Success {
		out = append(out, p)
	}
	return out
}

// MoveSuccessToFailed 把一条成功密码降级为失败（并删除 success）
func (db *DB) MoveSuccessToFailed(ssid, pwd string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	ent := db.entry(ssid)
	delete(ent.Success, pwd)
	ent.Failed[pwd] = struct{}{}
}
