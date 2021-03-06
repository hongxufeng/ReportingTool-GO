package report

import (
	"bytes"
	"datahelper/db"
	"datahelper/definition"
	"errors"
	"fmt"
	"model"
	"reflect"
	"strings"
	"time"
	"utils/function"
	"utils/service"
)

func Format(colConfig *model.ColumnConfig, cellValue string) (string, error) {
	if cellValue == "" {
		return "", nil
	}
	if colConfig.HasDateformat {
		time, e := time.Parse(colConfig.DateFormat, cellValue)
		if e == nil {
			cellValue = time.String()
		}
	} else if colConfig.HasTimeTransfer {
		switch colConfig.TimeTransfer {
		case "second":
			val, e := function.StringToInt(cellValue)
			if e == nil {
				cellValue = function.IntToString(val/(24*60*60)) + "日" + function.IntToString((val-(val/(24*60*60))*60*60*24)/3600) + "时" + function.IntToString((val-(val/(60*60))*60*60)/60) + "分" + function.IntToString(val-(val/60)*60) + "秒"
			}
			break
		}
	}
	fmt.Println(colConfig.HasFormatterR)
	if colConfig.HasFormatterR {
		cellValue = ReplacePlaceHolder(colConfig.FormatterR, cellValue)
	}
	var definitionData definition.Definition
	define := reflect.ValueOf(&definitionData)
	if colConfig.HasFormatter {
		//反射查找相应函数
		method := define.MethodByName(colConfig.Formatter)
		values := method.Call([]reflect.Value{reflect.ValueOf(*colConfig), reflect.ValueOf(cellValue)})
		fmt.Println(values)
		if len(values) != 2 {
			return cellValue, errors.New(fmt.Sprintf("method %s return value is not 2.", colConfig.Formatter))
		}
		switch x := values[1].Interface().(type) {
		case nil:
			cellValue = values[0].String()
		default:
			return cellValue, x.(error)
		}
	} else if colConfig.HasSelectorFunc {
		method := define.MethodByName(colConfig.SelectorFunc)
		values := method.Call([]reflect.Value{reflect.ValueOf(colConfig.SelectorFuncAgrs)})
		fmt.Println(values)
		if len(values) != 2 {
			return cellValue, errors.New(fmt.Sprintf("method %s return value is not 2.", colConfig.SelectorFunc))
		}
		switch x := values[1].Interface().(type) {
		case nil:
			switch y := values[0].Interface().(type) {
			default:
				selectordata := y.(map[string]string)
				fmt.Println(selectordata)
				cellValue = selectordata[cellValue]
			}
		default:
			return cellValue, x.(error)
		}
	}
	return cellValue, nil
}

func ReplacePlaceHolder(inStr string, cellValue string) string {
	outStr := inStr
	if cellValue == "" {
		outStr = strings.Replace(outStr, "$$::$$", "", -1)
	} else {
		outStr = strings.Replace(outStr, "$$::$$", cellValue, -1)
	}
	fmt.Println("outStr:", outStr)
	return outStr
}
func AppendWhere(req *service.HttpRequest, param *Param, buf *bytes.Buffer) error {
	hasWhere := false
	for _, colconfig := range param.ColConfigDict {
		var value string
		_ = req.GetParams(false, colconfig.Tag, &value)
		if len(value) == 0 {
			_ = req.GetParams(false, colconfig.Tag+"~~", &value)
			if len(value) == 0 {
				continue
			}
		}
		if hasWhere == false {
			hasWhere = true
			buf.WriteString(" where ")
		}
		buf.WriteString("(")
		queryList := strings.Split(value, "|")
		for i, v := range queryList {
			if i != 0 {
				buf.WriteString(" or ")
			}
			_, ok := function.StringToInt(v)
			buf.WriteString(colconfig.Tag)
			buf.WriteString("=")
			if ok == nil {
				buf.WriteString(v)
				buf.WriteString(" or ")
				buf.WriteString(colconfig.Tag)
				buf.WriteString("=")
			}
			buf.WriteString("\"")
			buf.WriteString(v)
			buf.WriteString("\"")
		}
		buf.WriteString(")")
		buf.WriteString(" and ")
	}
	if hasWhere {
		buf.Truncate(buf.Len() - 5)
	}
	return nil
}

func BuildQuerySQL(req *service.HttpRequest, param *Param, settings *model.Settings) (sql string, err error) {
	if param.TableConfig.HasPower && param.Power >= param.TableConfig.Power {
		err = service.NewError(service.ERR_POWER_DENIED, "您的用户权限不足啊！")
		return
	}
	var buf bytes.Buffer
	buf.WriteString("select ")
	var size = len(param.ColConfigDict)
	if size == 0 {
		return buf.String(), service.NewError(service.ERR_XML_ATTRIBUTE_LACK, "您至少需要配置一项XML中的列属性啊！")
	}
	for i := 0; i < size; i++ {
		if param.ColConfigDict[i].Tag == "buttons" || param.ColConfigDict[i].Tag == "pagerbuttons" {
			continue
		}
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(param.ColConfigDict[i].Tag)
	}
	buf.WriteString(" from ")
	buf.WriteString(param.TableConfig.Name)

	err = AppendWhere(req, param, &buf)
	if err != nil {
		return
	}
	if settings.Order != "" {
		buf.WriteString(" order by ")
		buf.WriteString(settings.Order)
	} else if param.TableConfig.HasDefaultOrder {
		buf.WriteString(" order by ")
		buf.WriteString(param.TableConfig.DefaultOrder)
	}
	buf.WriteString(" limit ")
	buf.WriteString(function.IntToString(settings.Rows * (settings.Page - 1)))
	buf.WriteString(",")
	buf.WriteString(function.IntToString(settings.Rows*settings.Page - 1))
	sql = buf.String()
	return
}

func GetSelectQuery(param *Param, fields string) (query string, err error) {
	var buf bytes.Buffer
	buf.WriteString("select ")
	buf.WriteString(fields)
	buf.WriteString(" from ")
	buf.WriteString(param.TableConfig.Name)
	query = buf.String()
	fmt.Println(query)
	return
}

func GetTableCount(param *Param, fields string) (count int, err error) {
	query, _ := GetSelectQuery(param, "count("+fields+")")
	result, err := db.Query(query)
	if err != nil {
		return
	}
	defer result.Close()
	if result.Next() {
		err = result.Scan(&count)
	} else {
		err = service.NewError(service.ERR_MYSQL)
	}
	return
}

func BuildTablePager(param *Param, settings *model.Settings, bodybuf *bytes.Buffer, count int, style string) (err error) {
	var x int
	if count%settings.Rows == 0 {
		x = 0
	} else {
		x = 1
	}
	totalpages := count/settings.Rows + x
	rowlist := strings.Split(settings.RowList, ",")
	start := (settings.Page-1)*settings.Rows + 1
	var end int
	if (settings.Page * settings.Rows) <= count {
		end = settings.Page * settings.Rows
	} else {
		end = count
	}
	bodybuf.WriteString("<div class=\"rt-pager-container\">")
	bodybuf.WriteString("<div class=\"rt-pager-buttons\">")
	if param.TableConfig.HasExcel && param.TableConfig.Excel == "true" {
		bodybuf.WriteString("<span class=\" rt-pager-export rt-pager-btn\"><span class=\"glyphicon glyphicon-export\" title=\"导出Excel\"></span>导出</span>")
	}
	if pagerbuttons := param.ColConfigDict[len(param.ColConfigDict)-1]; pagerbuttons.Tag == "pagerbuttons" {
		bodybuf.WriteString(pagerbuttons.Text)
	}
	bodybuf.WriteString("</div>")

	if style != model.Style_Tree {
		bodybuf.WriteString("<div class=\"rt-pager-controls\">")
		bodybuf.WriteString("&nbsp;<span class=\"glyphicon glyphicon-step-backward rt-pager-firstPage rt-pager-hover-color\"></span>")
		bodybuf.WriteString("&nbsp;<span class=\"glyphicon glyphicon-backward rt-pager-prevPage rt-pager-hover-color\"></span>")
		bodybuf.WriteString("&nbsp;<span class=\"pager-separator\"></span>&nbsp;")
		bodybuf.WriteString("第&nbsp;<input type=\"text\" class=\"rt-pager-page\" value=\"")
		bodybuf.WriteString(function.IntToString(settings.Page))
		bodybuf.WriteString("\"/>&nbsp;页，")
		bodybuf.WriteString("共&nbsp;<span class=\"rt-pager-totalPages\">")
		bodybuf.WriteString(function.IntToString(totalpages))
		bodybuf.WriteString("</span>&nbsp;页")
		bodybuf.WriteString("&nbsp;<span class=\"pager-separator\"></span>&nbsp;")
		bodybuf.WriteString("<span class=\"glyphicon glyphicon-forward rt-pager-nextPage rt-pager-hover-color\"></span>&nbsp;")
		bodybuf.WriteString("<span class=\"glyphicon glyphicon-step-forward rt-pager-lastPage rt-pager-hover-color\"></span>&nbsp;&nbsp;")
		bodybuf.WriteString("<select class=\"rt-pager-rowList\">")
		for _, v := range rowlist {
			bodybuf.WriteString("<option value=\"")
			bodybuf.WriteString(v)
			bodybuf.WriteString("\"")
			if i, _ := function.StringToInt(v); i == settings.Rows {
				bodybuf.WriteString(" selected")
			}
			bodybuf.WriteString(">")
			bodybuf.WriteString(v)
			bodybuf.WriteString("</option>")
		}
		bodybuf.WriteString("</select>")
		bodybuf.WriteString("</div>")
		bodybuf.WriteString("<div class=\"rt-pager-records\">第&nbsp;")
		bodybuf.WriteString(function.IntToString(start))
		bodybuf.WriteString(" - ")
		bodybuf.WriteString(function.IntToString(end))
		bodybuf.WriteString("&nbsp;条，")
		bodybuf.WriteString("共&nbsp;<span class=\"rt-pager-totalRecords\">")
		bodybuf.WriteString(function.IntToString(count))
		bodybuf.WriteString("</span>&nbsp;条</div>")
	}
	bodybuf.WriteString("</div>")
	return
}

func BuildSelectorBar(req *service.HttpRequest, param *Param, size int, selectorbuf *bytes.Buffer, conditionbuf *bytes.Buffer) (err error) {
	var selectordata map[string]string
	var definitionData definition.Definition
	define := reflect.ValueOf(&definitionData)
	for i := 0; i < size; i++ {
		selectordata = make(map[string]string, 0)
		if !param.ColConfigDict[i].IsInSelector {
			continue
		}
		if param.ColConfigDict[i].HasSelectorFunc == true {
			method := define.MethodByName(param.ColConfigDict[i].SelectorFunc)
			values := method.Call([]reflect.Value{reflect.ValueOf(param.ColConfigDict[i].SelectorFuncAgrs)})
			fmt.Println(values)
			if len(values) != 2 {
				err = errors.New(fmt.Sprintf("method %s return value is not 2.", param.ColConfigDict[i].SelectorFunc))
			}
			switch x := values[1].Interface().(type) {
			case nil:
				switch y := values[0].Interface().(type) {
				case nil:
					return errors.New(fmt.Sprintf("method %s second return value is nil.", param.ColConfigDict[i].SelectorFunc))
				default:
					selectordata = y.(map[string]string)
				}
			default:
				return x.(error)
			}
		} else {
			being, _ := db.GetSelectorCachePersistence(param.TableConfig.Name, param.ColConfigDict[i].Tag)
			if being == true {
				selectordata, _ = db.HGetSelectorBarCache(param.TableConfig.Name, param.ColConfigDict[i].Tag)
				fmt.Println(being, selectordata)
			} else {
				_ = db.SetSelectorCachePersistence(param.TableConfig.Name, param.ColConfigDict[i].Tag)
				fmt.Println("nothing")
				query, _ := GetSelectQuery(param, "distinct("+param.ColConfigDict[i].Tag+")")
				result, e := db.Query(query)
				if e != nil {
					return e
				}
				defer result.Close()
				for j := 0; result.Next(); j++ {
					var value string
					if e = result.Scan(&value); e != nil {
						fmt.Println("BuildSelectorBar:", e.Error())
						return e
					}
					fmt.Println(value)
					selectordata[value] = value
					e := db.HSetSelectorBarCache(param.TableConfig.Name, param.ColConfigDict[i].Tag, value, value)
					if e != nil {
						return e
					}
				}
			}
		}
		var value string
		e := req.GetParams(false, param.ColConfigDict[i].Tag, &value)
		if e == nil {
			if value == "" {
				continue
			}
			originValue := strings.Split(value, "|")
			var valueText []string
			for _, v := range originValue {
				sd, ok := selectordata[v]
				if !ok {
					continue
				}
				valueText = append(valueText, sd)
			}
			conditionbuf.WriteString("<div data-value=\"")

			conditionbuf.WriteString(param.ColConfigDict[i].Tag)
			conditionbuf.WriteString("\">")
			conditionbuf.WriteString(param.ColConfigDict[i].Text)
			conditionbuf.WriteString(" : ")
			conditionbuf.WriteString(strings.Join(valueText, "、"))
			conditionbuf.WriteString("<span class=\"glyphicon glyphicon-remove rt-glyphicon-color rt-condition-remove\"></div>")
			conditionbuf.WriteString("</div>")
			continue
		}
		selectorbuf.WriteString("<div class=\"rt-selector-folder\">")
		selectorbuf.WriteString("<div class=\"rt-selector-key\" data-value=\"")
		selectorbuf.WriteString(param.ColConfigDict[i].Tag)
		selectorbuf.WriteString("\">")
		selectorbuf.WriteString(param.ColConfigDict[i].Text)
		selectorbuf.WriteString("：</div>")
		selectorbuf.WriteString("<div class=\"rt-selector-value\">")
		selectorbuf.WriteString("<ul class=\"rt-selector-list\">")
		for k, v := range selectordata {
			selectorbuf.WriteString("<li data-value=\"")
			selectorbuf.WriteString(k)
			selectorbuf.WriteString("\"><span class=\"rt-selector-list-text\"><span class=\"glyphicon glyphicon-unchecked\"></span>")
			selectorbuf.WriteString(v)
			selectorbuf.WriteString("</span></li>")
		}
		selectorbuf.WriteString("</ul>")
		selectorbuf.WriteString("</div>")
		if param.ColConfigDict[i].HasSelectorMulti {
			selectorbuf.WriteString("<div class=\"rt-multiselect-btns\">")
			selectorbuf.WriteString("<button class=\"btn btn-primary btn-xs rt-multiselect-ok\">确&nbsp;&nbsp;定</button><button class=\"btn btn-default btn-xs rt-multiselect-cancel\">取&nbsp;&nbsp;消</button>")
			selectorbuf.WriteString("</div>")
		}
		selectorbuf.WriteString("<div class=\"rt-selector-btns\">")
		selectorbuf.WriteString("<span class=\"rt-selector-selectmore\"><span class\"rt-selectmore-txt\">更多</span><span class=\"glyphicon glyphicon-chevron-down rt-glyphicon-color\"></span></span>")
		if param.ColConfigDict[i].HasSelectorMulti {
			selectorbuf.WriteString("<span class=\"rt-selector-multiselect\">多选<span class=\"glyphicon glyphicon-plus rt-glyphicon-color\"></span></span>")
		}
		selectorbuf.WriteString("</div>")
		selectorbuf.WriteString("</div>")
	}
	return nil
}
