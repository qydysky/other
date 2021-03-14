package proxylist

import (
	// "flag"
	"bufio"
	"errors"
	// "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/app/router"
)

func Proxylist(){}

type ProxyList_item struct {
	Acce_file string
	Filename string
	Output string
	// Acce_file = flag.String("accefile", "", "acce.log location.")
	// //Sign    = flag.String("sign", "giaOut", "proxy sign string.")
	// Filename       = flag.String("file", "data/", "data save file.")
	// Output      = flag.String("o", "proxy.dat", "output file.")
}

type Entry struct {
	Type  string
	Value string
	Attrs []*router.Domain_Attribute
}

type List struct {
	Name  string
	Entry []Entry
}

type ParsedList struct {
	Name      string
	Inclusion map[string]bool
	Entry     []Entry
}

func (l *ParsedList) toProto() (*router.GeoSite, error) {
	site := &router.GeoSite{
		CountryCode: l.Name,
	}
	for _, entry := range l.Entry {
		switch entry.Type {
		case "domain":
			site.Domain = append(site.Domain, &router.Domain{
				Type:      router.Domain_Domain,
				Value:     entry.Value,
				Attribute: entry.Attrs,
			})
		case "regex":
			site.Domain = append(site.Domain, &router.Domain{
				Type:      router.Domain_Regex,
				Value:     entry.Value,
				Attribute: entry.Attrs,
			})
		case "keyword":
			site.Domain = append(site.Domain, &router.Domain{
				Type:      router.Domain_Plain,
				Value:     entry.Value,
				Attribute: entry.Attrs,
			})
		case "full":
			site.Domain = append(site.Domain, &router.Domain{
				Type:      router.Domain_Full,
				Value:     entry.Value,
				Attribute: entry.Attrs,
			})
		default:
			return nil, errors.New("unknown domain type: " + entry.Type)
		}
	}
	return site, nil
}

func removeComment(line string) string {
	idx := strings.Index(line, "#")
	if idx == -1 {
		return line
	}
	return strings.TrimSpace(line[:idx])
}

func parseDomain(domain string, entry *Entry) error {
	kv := strings.Split(domain, ":")
	if len(kv) == 1 {
		entry.Type = "domain"
		entry.Value = strings.ToLower(kv[0])
		return nil
	}

	if len(kv) == 2 {
		entry.Type = strings.ToLower(kv[0])
		entry.Value = strings.ToLower(kv[1])
		return nil
	}

	return errors.New("Invalid format: " + domain)
}

func parseAttribute(attr string) (router.Domain_Attribute, error) {
	var attribute router.Domain_Attribute
	if len(attr) == 0 || attr[0] != '@' {
		return attribute, errors.New("invalid attribute: " + attr)
	}

	attr = attr[0:]
	parts := strings.Split(attr, "=")
	if len(parts) == 1 {
		attribute.Key = strings.ToLower(parts[0])
		attribute.TypedValue = &router.Domain_Attribute_BoolValue{BoolValue: true}
	} else {
		attribute.Key = strings.ToLower(parts[0])
		intv, err := strconv.Atoi(parts[1])
		if err != nil {
			return attribute, errors.New("invalid attribute: " + attr + ": " + err.Error())
		}
		attribute.TypedValue = &router.Domain_Attribute_IntValue{IntValue: int64(intv)}
	}
	return attribute, nil
}

func parseEntry(line string) (Entry, error) {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")

	var entry Entry
	if len(parts) == 0 {
		return entry, errors.New("empty entry")
	}

	if err := parseDomain(parts[0], &entry); err != nil {
		return entry, err
	}

	for i := 1; i < len(parts); i++ {
		attr, err := parseAttribute(parts[i])
		if err != nil {
			return entry, err
		}
		entry.Attrs = append(entry.Attrs, &attr)
	}

	return entry, nil
}

func DetectPath(path string) (string, error) {
	arrPath := strings.Split(path, string(filepath.ListSeparator))
	for _, content := range arrPath {
		fullPath := filepath.Join(content, "src", "github.com", "v2ray", "domain-list-community", "data")
		_, err := os.Stat(fullPath)
		if err == nil || os.IsExist(err) {
			return fullPath, nil
		}
	}
	err := errors.New("No file found in GOPATH")
	return "", err
}

func Load(path string) (*List, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	list := &List{
		Name: strings.ToUpper(filepath.Base(path)),
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeComment(line)
		if len(line) == 0 {
			continue
		}
		entry, err := parseEntry(line)
		if err != nil {
			return nil, err
		}
		list.Entry = append(list.Entry, entry)
	}

	return list, nil
}

func ParseList(list *List, ref map[string]*List) (*ParsedList, error) {
	pl := &ParsedList{
		Name:      list.Name,
		Inclusion: make(map[string]bool),
	}
	entryList := list.Entry
	for {
		newEntryList := make([]Entry, 0, len(entryList))
		hasInclude := false
		for _, entry := range entryList {
			if entry.Type == "include" {
				if pl.Inclusion[entry.Value] {
					continue
				}
				refName := strings.ToUpper(entry.Value)
				pl.Inclusion[refName] = true
				r := ref[refName]
				if r == nil {
					return nil, errors.New(entry.Value + " not found.")
				}
				newEntryList = append(newEntryList, r.Entry...)
				hasInclude = true
			} else {
				newEntryList = append(newEntryList, entry)
			}
		}
		entryList = newEntryList
		if !hasInclude {
			break
		}
	}
	pl.Entry = entryList

	return pl, nil
}

func Main(p ProxyList_item) error {
	Acce_file := p.Acce_file
	Filename := p.Filename
	Output := p.Output
	dir := Filename
	if Filename == "" {Filename = ""}
	if Output == "" {Output = "proxy.dat"}

	var exist func(string) bool = func (s string) bool {
		_, err := os.Stat(s)
		return err == nil || os.IsExist(err)
	}
	if Acce_file != ""||exist(Acce_file) {
		Main_proxy(Main_proxy_type {
			Acce_file:Acce_file,
			Sign:"gia",
			Filename:Filename+"giaOut",
			Discard:true,
		})
		Main_proxy(Main_proxy_type {
			Acce_file:Acce_file,
			Sign:"fast",
			Filename:Filename+"fastOut",
			Discard:true,
		})
	}
	NewPath(Filename)
	ref := make(map[string]*List)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		list, err := Load(path)
		if err != nil {
			return err
		}
		ref[list.Name] = list
		return nil
	})
	if err != nil {
		return err
	}
	protoList := new(router.GeoSiteList)
	for _, list := range ref {
		pl, err := ParseList(list, ref)
		if err != nil {
			return err
		}
		site, err := pl.toProto()
		if err != nil {
			return err
		}
		protoList.Entry = append(protoList.Entry, site)
	}

	protoBytes, err := proto.Marshal(protoList)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(Output, protoBytes, 0777); err != nil {
		return err
	}
	return nil
}

func NewPath(filename string) error {
	var newpath func(string) error = func (filename string)error{
		/*
			如果filename路径不存在，就新建它
		*/	
		var exist func(string) bool = func (s string) bool {
			_, err := os.Stat(s)
			return err == nil || os.IsExist(err)
		}
	
		for i:=0;true;{
			a := strings.Index(filename[i:],"/")
			if a == -1 {break}
			if a == 0 {a = 1}//bug fix 当绝对路径时开头的/导致问题
			i=i+a+1
			if !exist(filename[:i-1]) {
				err := os.Mkdir(filename[:i-1], os.ModePerm)
				if err != nil {return err}
			}
		}
		
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fd,err:=os.Create(filename)
			if err != nil {
				return err
			}else{
				fd.Close()
			}
		}
		return nil
	}
	return newpath(filename)

}