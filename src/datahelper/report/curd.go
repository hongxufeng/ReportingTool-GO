package report

import (
	"bytes"
	"model"
	"utils/service"
)

func GetCURD(req *service.HttpRequest, param *Param, settings *model.CRUDSettings, bodybuf *bytes.Buffer) (err error) {
	bodybuf.WriteString("<form class=\"form-horizontal\">")
	err = BuildCURDBody(param, settings, bodybuf)
	if err != nil {
		return
	}
	bodybuf.WriteString("</form>")
	return
}
func BuildCURDBody(param *Param, settings *model.CRUDSettings, bodybuf *bytes.Buffer) (err error) {
	for _, colConfig := range param.ColConfigDict {
		if colConfig.Tag == "buttons" || colConfig.Tag == "pagerbuttons" {
			continue
		}
		bodybuf.WriteString("<div class=\"form-group\">")
		bodybuf.WriteString("<label class=\"col-sm-3 control-label\">")
		bodybuf.WriteString(colConfig.Text)
		bodybuf.WriteString("&nbsp;&nbsp;<span class=\"rt-glyphicon-color\">:</span></label>")
		bodybuf.WriteString("<div class=\"col-sm-6\">")
		bodybuf.WriteString("<input name=\"")
		bodybuf.WriteString(colConfig.Tag)
		bodybuf.WriteString("\" type=\"text\" class=\"form-control rt-form-control\" placeholder=\"")
		bodybuf.WriteString(colConfig.Text)
		bodybuf.WriteString("\">")
		bodybuf.WriteString("</div>")
		bodybuf.WriteString("</div>")
	}
	bodybuf.WriteString("<div class=\"form-group\">")
	bodybuf.WriteString("<div class=\"col-sm-offset-3  col-sm-6\">")
	if settings.Cmd == model.Cmd_Create {
		bodybuf.WriteString("<input type=\"submit\" class=\"btn btn-primary form-control\" value=\"")
		if param.TableConfig.HasBtnCreateText {
			bodybuf.WriteString(param.TableConfig.BtnCreateText)
		} else {
			bodybuf.WriteString("新&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;增")
		}
		bodybuf.WriteString("\" />")
	}
	return
}
