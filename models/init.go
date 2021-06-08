package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(
		new(Category),
		new(Book),
		new(BookCategory),
		new(Attachment),
		new(Document),
		new(Member),
		new(Relationship),
	)

}

/*
* Table Names
* */

func TNCategory() string {
	return "md_category"
}

func TNBookCategory() string {
	return "md_book_category"
}

func TNBook() string {
	return "md_books"
}

func TNDocuments() string {
	return "md_documents"
}
func TNDocumentStore() string {
	return "md_document_store"
}

func TNAttachment() string {
	return "md_attachment"
}

func TNRelationship() string {
	return "md_relationship"
}

func TNMembers() string {
	return "md_members"
}

func TNCollection() string {
	return "md_star"
}

func TNFans() string {
	return "md_fans"
}

func TNComments() string {
	return "md_comments"
}

func TNScore() string {
	return "md_score"
}
