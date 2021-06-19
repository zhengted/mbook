package models

type Relationship struct {
	RelationshipId int `orm:"pk;auto" json:"relationship_id"`
	MemberId       int `json:"member_id"`
	BookId         int `json:"book_id"`
	RoleId         int `json:"role_id"`
}

func (r *Relationship) TableName() string {
	return TNRelationship()
}

func NewRelationship() *Relationship {
	return &Relationship{}
}

func (m *Relationship) Select(bookId, memberId int) (*Relationship, error) {
	err := GetOrm("w").QueryTable(m.TableName()).Filter("book_id", bookId).Filter("member_id", memberId).One(m)
	return m, err
}

func (m *Relationship) SelectRoleId(bookId, memberId int) (int, error) {
	err := GetOrm("w").QueryTable(m.TableName()).Filter("book_id", bookId).Filter("member_id", memberId).One(m, "role_id")
	if err != nil {
		return 0, err
	}
	return m.RoleId, nil
}

func (m *Relationship) Insert() error {
	_, err := GetOrm("w").Insert(m)
	return err
}

func (m *Relationship) Update() error {
	_, err := GetOrm("w").Update(m)
	return err
}
