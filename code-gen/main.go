package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	. "github.com/dave/jennifer/jen"
)

const (
	ROOT    = "D:/go_workspace/src"
	SRV     = "servcie"
	CTRL    = "controller"
	DAO     = "domain"
	PERSIST = "persistence"
)

func main() {
	NewP("calais.gaodun.com", "glive", "pear").Gen()
}

//代码
type Srv struct {
}
type Ctrl struct {
}
type Dao struct {
}

type P struct {
	projectName string
	domain      string
	prefix      string
	controller  *Ctrl
	service     *Srv
	dao         *Dao
}

func NewP(projectName string, domain string, prefix string) *P {
	return &P{projectName: projectName, domain: domain, prefix: prefix}
}

func (p *P) Gen() {
	p.genDao()
	p.genService()
	p.genController()
}

func (p *P) genService() {
	var err error
	srvFile := NewFile(SRV)
	srvFileWriter, err := p.tryCreteFile(fmt.Sprintf("%s/service/%sService.go", p.domain, p.prefix))
	srvName := fmt.Sprintf("%s%s", Capitalize(p.prefix), Capitalize(SRV))
	if err == nil && srvFileWriter != nil {
		srvFile.Type().Id(srvName).Interface()
		srvFile.Render(srvFileWriter)
	}
	srvImplFile := NewFile(SRV)
	srvImplFileWriter, err := p.tryCreteFile(fmt.Sprintf("%s/service/%sServiceImpl.go", p.domain, p.prefix))
	if err == nil && srvImplFileWriter != nil {
		daoName := fmt.Sprintf("%sDao", Capitalize(p.prefix))
		srvImplFile.Type().Id(fmt.Sprintf("%sImpl", srvName)).Struct(
			Id(daoName).Qual(fmt.Sprintf("%s/%s/domain", p.projectName, p.domain), daoName),
		)
		srvImplFile.Render(srvImplFileWriter)
	}
}

func (p *P) genController() {
	//使用NewFilePath取代NewFile，引入包名的好处，可以保证当前包下添加的引用自动去除包名。
	ctrlFile := NewFilePath(fmt.Sprintf("%s/%s", p.projectName, CTRL))
	ctrlFileWriter, err := p.tryCreteFile(fmt.Sprintf("controller/%s.go", p.prefix))
	if err == nil && ctrlFileWriter != nil {
		srvName := fmt.Sprintf("%sService", Capitalize(p.prefix))
		ctrlFile.Type().Id(fmt.Sprintf("%sApi", Capitalize(p.prefix))).Struct(
			Id("").Qual(fmt.Sprintf("%s/%s", p.projectName, CTRL), "Base"),
			Id(srvName).Qual(fmt.Sprintf("%s/%s/service", p.projectName, p.domain), srvName),
		)
		ctrlFile.Render(ctrlFileWriter)
	}

}
func (p *P) genDao() {
	daoFile := NewFile(DAO)
	daoFileWriter, err := p.tryCreteFile(fmt.Sprintf("%s/domain/%sDao.go", p.domain, p.prefix))
	if err == nil && daoFileWriter != nil {
		daoName := fmt.Sprintf("%sDao", Capitalize(p.prefix))
		daoFile.Type().Id(daoName).Interface()
		daoFile.Render(daoFileWriter)
	}
	dbalFile := NewFile(PERSIST)
	dbalFileWriter, err := p.tryCreteFile(fmt.Sprintf("%s/infrastructure/persistence/dbal/%sDAODBAL.go", p.domain, p.prefix))
	if err == nil && daoFileWriter != nil {
		dbalName := fmt.Sprintf("%sDAODBAL", Capitalize(p.prefix))
		engineName := "Engine"
		dbalFile.Type().Id(dbalName).Struct(
			Id("Table").Op("*").Qual(fmt.Sprintf("%s/%s/domain", p.projectName, p.domain), Capitalize(p.prefix)),
			Id(engineName).Op("*").Qual("github.com/go-xorm/xorm", engineName),
		)
		dbalFile.Render(dbalFileWriter)
	}
}

func (p *P) tryCreteFile(filePath string) (*os.File, error) {
	filePath = fmt.Sprintf("%s/%s/%s", ROOT, p.projectName, filePath)
	var file *os.File
	var err error
	if checkFileIsExist(filePath) { //如果文件存在
		return nil, nil
	} else {
		basepath := path.Dir(filePath)
		if os.MkdirAll(basepath, 0666) != nil {
			err = errors.New("创建目录失败！")
		}
		file, err = os.Create(filePath) //创建文件
	}
	check(err)
	return file, nil
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
