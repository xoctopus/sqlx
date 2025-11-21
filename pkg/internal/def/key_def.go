package def

// ParseKeyDefine parses key define
// eg:
//
//	| Kind         | Name[,Method]      | Fields[,Option]            | Comments                                  |
//	| :---         | :---               | :----                      | :----                                     |
//	| idx          | idx_name,BTREE     | Name                       | 普通索引, 使用btree                       |
//	| index        | idx_name,GIST      | Geo,gist_trgm_ops          | 位置索引, 使用GIST, 开启gist_trgm_ops选项 |
//	| unique_index | idx_name           | OrgID,NULLS,FIRST MemberID | 组织ID+成员ID唯一索引 OrgID排序空在后     |
//	| u_idx        | idx_name           | OrgID MemberID,NULLS,FIRST | 组织ID+成员ID唯一索引 MemberID排序空在前  |
//	| primary      |                    | ID                         | 主键                                      |
//	| pk           |                    | ID                         | 主键                                      |
func ParseKeyDefine(def string) *KeyDefine {
	return nil
}

type KeyDefine struct {
	Kind    string
	Name    string
	Using   string
	Comment string
	Options []KeyColumnOption
}

type KeyColumnOption struct {
	Name    string
	Options []string
}
