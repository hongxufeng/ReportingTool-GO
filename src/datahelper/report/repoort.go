package report

import (
	"bytes"
	"datahelper/db"
	"fmt"
	"model"
	"utils/function"
	"utils/service"

	"github.com/beevik/etree"
)

type Param struct {
	TableConfig   model.TableConfig
	Uid           uint32
	Power         uint8 //用户判断得到的权限 暂时用0代表最高权限
	colName       []string
	ColConfigDict []model.ColumnConfig
}

func New(uid uint32, configFile string, tableID string) (param *Param, err error) {
	var ColConfigDict []model.ColumnConfig
	var ColName []string
	doc := etree.NewDocument()
	filename := "xml/" + configFile + ".xml"
	fmt.Println(filename)
	err = doc.ReadFromFile(filename)
	if err != nil {
		return
	}
	tableconfig := model.TableConfig{}
	if tableelement := doc.FindElements("./tables/table[@id='" + tableID + "']"); len(tableelement) > 0 {
		fmt.Println(len(tableelement))
		table := tableelement[0]
		defaultorder := table.SelectAttr("defaultorder")
		if defaultorder != nil {
			tableconfig.HasDefaultOrder = true
			tableconfig.DefaultOrder = defaultorder.Value
		}
		name := table.SelectAttr("name")
		if name == nil {
			err = service.NewError(service.ERR_XML_ATTRIBUTE_LACK, "XML配置属性name缺失哦！")
		} else {
			tableconfig.Name = name.Value
		}
		btncreatetext := table.SelectAttr("btn-create-text")
		if btncreatetext != nil {
			tableconfig.HasBtnCreateText = true
			tableconfig.BtnCreateText = btncreatetext.Value
		}
		//adminname := table.SelectAttr("adminname")
		//if adminname!=nil{
		//	tableconfig.HasAdminName=true
		//	tableconfig.AdminName=adminname.Value
		//}
		power := table.SelectAttr("power")
		if power != nil {
			tableconfig.HasPower = true
			tableconfig.Power, _ = function.StringToUint8(power.Value)
		}
		excel := table.SelectAttr("excel")
		if excel != nil {
			tableconfig.HasExcel = true
			tableconfig.Excel = excel.Value
		}
	} else {
		err = service.NewError(service.ERR_XML_TABLE_LACK, "您配置的XML表格是不在的啊！")
		return
	}
	path := "./tables/table[@id='" + tableID + "']/*"
	//fmt.Println(path)
	for _, elemnt := range doc.FindElements(path) {
		fmt.Printf("%s: %s\n", elemnt.Tag, elemnt.Text())
		if elemnt.Tag == "buttons" || elemnt.Tag == "pagerbuttons" {
			ColName = append(ColName, elemnt.Tag)
		}
		//赋值ColConfigDict
		cc := model.ColumnConfig{}
		cc.Tag = elemnt.Tag
		cc.Text = elemnt.Text()
		btnicon := elemnt.SelectAttr("btn-icon")
		btnfunc := elemnt.SelectAttr("btn-func")
		if btnicon != nil && btnfunc != nil {
			cc.HasBtn = true
			cc.BtnFunc = btnfunc.Value
			cc.BtnIcon = btnicon.Value
		}
		dateformat := elemnt.SelectAttr("dateformat")
		if dateformat != nil {
			cc.HasDateformat = true
			cc.DateFormat = dateformat.Value
		}
		checkbox := elemnt.SelectAttr("checkbox")
		if checkbox != nil && checkbox.Value == "true" {
			cc.ISCheckBox = true
		}
		defaultvalue := elemnt.SelectAttr("defaultvalue")
		if defaultvalue != nil {
			cc.HasDefaultValue = true
			cc.DefaultValue = defaultvalue.Value
		}
		power := elemnt.SelectAttr("power")
		if power != nil {
			cc.HasPower = true
			cc.Power, _ = function.StringToUint8(power.Value)
		}
		formatter := elemnt.SelectAttr("formatter")
		if formatter != nil {
			cc.HasFormatter = true
			cc.Formatter = formatter.Value
		}
		selector := elemnt.SelectAttr("selector")
		if selector != nil && selector.Value == "true" {
			cc.IsInSelector = true
		}
		formatterr := elemnt.SelectAttr("formatter-r")
		if formatterr != nil {
			cc.HasFormatterR = true
			cc.FormatterR = formatterr.Value
		}
		searchtype := elemnt.SelectAttr("search-type")
		if searchtype != nil {
			cc.HasSearchType = true
			cc.SearchType = searchtype.Value
		}
		selectorfuncagrs := elemnt.SelectAttr("selector-func-agrs")
		selectorfunc := elemnt.SelectAttr("selector-func")
		if selectorfunc != nil && selectorfuncagrs != nil {
			cc.IsInSelector = true
			cc.HasSelectorFunc = true
			cc.SelectorFunc = selectorfunc.Value
			cc.SelectorFuncAgrs = selectorfuncagrs.Value
		}
		// selectortext := elemnt.SelectAttr("selector-text")
		// if selectortext != nil {
		// 	cc.HasSelectorText = true
		// 	cc.SelectorText = selectortext.Value
		// }
		//linkto := elemnt.SelectAttr("linkto")
		//passedcol := elemnt.SelectAttr("passedcol")
		//if linkto!=nil&&passedcol!=nil{
		//	cc.HasLinkTo=true
		//	cc.LinkTo=linkto.Value
		//	cc.Passedcol =strings.Split(passedcol.Value,",")
		//	ignoredpassedcol := elemnt.SelectAttr("ignoredpassedcol")
		//	if param.Power==0&&ignoredpassedcol!=nil{
		//		ipd:=strings.Split(ignoredpassedcol.Value,",")
		//		for  no,_:=range ipd{
		//			cc.Passedcol=append(cc.Passedcol[:no],cc.Passedcol[no+1:]...)
		//		}
		//	}
		//}
		selectormulti := elemnt.SelectAttr("selector-multi")
		if selectormulti != nil && selectormulti.Value == "true" {
			cc.HasSelectorMulti = true
		}
		searchadv := elemnt.SelectAttr("search-adv")
		if searchadv != nil && searchadv.Value == "true" {
			cc.IsSearchAdv = true
		}
		navname := elemnt.SelectAttr("navname")
		if navname != nil {
			cc.HasNavName = true
			cc.NavName = navname.Value
		}
		searchbtnicon := elemnt.SelectAttr("searchbtnicon")
		if searchbtnicon != nil {
			cc.HasSearchBtnIcon = true
			cc.SearchBtnIcon = searchbtnicon.Value
		}
		searchbtnfunc := elemnt.SelectAttr("search-btn-func")
		if searchbtnfunc != nil {
			cc.HasSearchBtnFunc = true
			cc.SearchBtnFunc = searchbtnfunc.Value
		}
		search4admin := elemnt.SelectAttr("search4admin")
		if search4admin != nil && search4admin.Value == "true" {
			cc.Search4Admin = true
		}
		timetransfer := elemnt.SelectAttr("timetransfer")
		if timetransfer != nil {
			cc.HasTimeTransfer = true
			cc.TimeTransfer = timetransfer.Value
		}
		precision := elemnt.SelectAttr("precision")
		if precision != nil {
			cc.HasPrecision = true
			cc.Precision = precision.Value
		}
		visibility := elemnt.SelectAttr("visibility")
		if visibility != nil {
			cc.Visibility = visibility.Value
		}
		//percentageform := elemnt.SelectAttr("percentageform")
		//if percentageform!=nil&&percentageform.Value=="true"{
		//	cc.IsInPercentageform=true
		//}
		ColConfigDict = append(ColConfigDict, cc)
	}
	if len(ColConfigDict) == 0 {
		err = service.NewError(service.ERR_XML_ATTRIBUTE_LACK, "您至少需要配置一项XML中的列属性啊！")
		return
	}
	fmt.Println(ColConfigDict)
	fmt.Println(tableconfig)
	//根据uid判断权限
	ud, err := db.GetUserInfo(uid)
	param = &Param{tableconfig, uid, ud.Power, ColName, ColConfigDict}
	return
}

func (param *Param) Table(req *service.HttpRequest, settings *model.Settings) (res map[string]interface{}, err error) {
	res = make(map[string]interface{}, 0)
	count, err := GetTableCount(param, "*")
	if err != nil {
		return
	}
	fmt.Println(count)
	query, err := BuildQuerySQL(req, param, settings)
	if err != nil {
		return
	}
	fmt.Println(query)
	//这里查询后可获得一个map
	// ret, err := db.FetchRows(query)
	// fmt.Println(len(ret))
	result, err := db.Query(query)
	if err != nil {
		return
	}
	defer result.Close()
	columns, _ := result.Columns()
	size := len(columns)
	var searchbuf, bodybuf, selectorbuf, conditionbuf, rowbuf bytes.Buffer
	err = GetTable(req, param, settings, result, size, &bodybuf, &searchbuf, &rowbuf, count)
	if err != nil {
		return
	}
	if settings.Style == model.Style_Table {
		err = BuildSelectorBar(req, param, size, &selectorbuf, &conditionbuf)
		if err != nil {
			return
		}
	}
	res["search"] = searchbuf.String()
	res["body"] = bodybuf.String()
	res["selector"] = selectorbuf.String()
	res["condition"] = conditionbuf.String()
	res["row"] = rowbuf.String()
	return
}

func (param *Param) SearchTree() (res map[string]interface{}, err error) {
	return
}

func (param *Param) LocateNode() (res map[string]interface{}, err error) {
	return
}

func (param *Param) CURD(req *service.HttpRequest, settings *model.CRUDSettings) (res map[string]interface{}, err error) {
	res = make(map[string]interface{}, 0)
	var bodybuf bytes.Buffer
	err = GetCURD(req, param, settings, &bodybuf)
	if err != nil {
		return
	}
	res["body"] = bodybuf.String()
	return
}
