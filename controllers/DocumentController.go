package controllers

import "mbook/models"

type DocumentController struct {
	BaseController
}

// 获取图书内容并判断权限
func (c DocumentController) getBookData(identify, token string) *models.BookData {
	// TODO: Complete me!
	return nil

}
