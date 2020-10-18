package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/common/base"
	bridageModels "weiXinBot/app/bridage/models"
	"weiXinBot/app/main/models"
)

// KeyWordsController ...
type KeyWordsController struct {
	base.BaseController
}

// Post ...
// @Title Post
// @Description create KeyWords
// @Param	body		body 	models.KeyWords	true		"body for Bots content"
// @Success 201 {int} models.KeyWords
// @Failure 403 body is empty
// @router / [post]
func (c *KeyWordsController) Post() {
	var err error
	var v *bridageModels.KeyWords
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: v}
		} else {
			c.Data["json"] = common.RestResult{Code: 0, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		return
	}
	_, err = models.AddKeyWord(v)
}

// GetOne ...
// @Title Get One
// @Description get KeyWords by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.KeyWords
// @Failure 403 :id is empty
// @router /:id [get]
func (c *KeyWordsController) GetOne() {
	var v *bridageModels.KeyWords
	var err error
	idstr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idstr, 0, 64)
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: v}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	v, err = models.GetKeyWordByID(id)
}

// GetOneKeyWord 根据类别和planid获取具体关键词回复
// @router /category [get]
func (c *KeyWordsController) GetOneKeyWord() {
	var v *bridageModels.KeyWords
	var err error
	typeStr := c.GetString("type")
	typeID, _ := strconv.Atoi(typeStr)
	planIDStr := c.GetString("planid")
	planID, _ := strconv.ParseInt(planIDStr, 0, 64)
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: v}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	v, err = models.GetKeyWordByTypeAndPlan(planID, typeID)
}

// GetAll ...
// @Title Get All
// @Description get KeyWords
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.KeyWords
// @Failure 403
// @router / [get]
func (c *KeyWordsController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query []*common.QueryConditon
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k|type:v|v|v,k|type:v|v|v  其中Type可以没有,默认值是 MultiText
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") { // 分割多个查询key
			qcondtion := new(common.QueryConditon)
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key:value pair," + cond)
				c.ServeJSON()
				return
			}
			kInit, vInit := kv[0], kv[1]         // 初始分割查询key和value（备注，value是多个用|分割）
			keyType := strings.Split(kInit, "|") // 解析key中的type信息
			if len(keyType) == 2 {
				qcondtion.QueryKey = keyType[0]
				qcondtion.QueryType = keyType[1]
			} else if len(keyType) == 1 {
				qcondtion.QueryKey = keyType[0]
				qcondtion.QueryType = common.MultiText
			} else {
				c.Data["json"] = errors.New("Error: invalid query key|type format," + kInit)
				c.ServeJSON()
				return
			}
			qcondtion.QueryValues = strings.Split(vInit, "|") // 解析出values信息
			qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
			query = append(query, qcondtion)
		}
	}
	l, count, err := models.GetAllKeyWord(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
	} else {
		c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: struct {
			Items interface{}
			Total int64
		}{
			Items: l,
			Total: count,
		}}
	}
	c.ServeJSON()
}

// Put ...
// @router /:id [put]
func (c *KeyWordsController) Put() {
	var err error
	defer func() {
		if err != nil {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		} else {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
		}
		c.ServeJSON()
	}()
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var v = bridageModels.KeyWords{ID: id}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		return
	}
	err = models.UpdateKeyWordByID(&v)
}

// Delete ...
// @router /:id [delete]
func (c *KeyWordsController) Delete() {
	var err error
	defer func() {
		if err != nil {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		} else {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
		}
		c.ServeJSON()
	}()
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	err = models.DeleteKeyWordByID(id)
}

// DeleteList ...
// @router /deletelist [delete]
func (c *KeyWordsController) DeleteList() {
	var idslice []interface{}
	if ids := c.GetString("ids"); ids != "" {
		s := strings.Split(ids, ",")
		for _, v := range s {
			idslice = append(idslice, v)
		}
	}
	if err := models.MultiDeleteKeyWordByIDs(idslice); err == nil {
		c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
	} else {
		c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
	}
	c.ServeJSON()
}
