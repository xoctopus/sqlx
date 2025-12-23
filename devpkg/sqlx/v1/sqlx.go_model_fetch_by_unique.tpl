@def T
@def context.Context
@def UniqueSuffix
@def UniqueConds
@def UniqueFields
@def FetchComment
@def frag.Fragment
@def types.SoftDeletion
@def session.For
@def builder.CC
@def driver.Value
@def builder.Select
@def builder.Where
@def builder.Limit
@def builder.Comment
--FetchByUnique
// FetchBy#UniqueSuffix# fetch #T# by #UniqueFields#
func (m *#T#) FetchBy#UniqueSuffix#(ctx #context.Context#) error {
	conds := []#frag.Fragment#{
		#UniqueConds#
	}

	if x, ok := any(m).(#types.SoftDeletion#); ok {
		col, val := x.SoftDeletion()
		conds = append(
			conds,
			#builder.CC#[#driver.Value#](T#T#.C(col)).AsCond(builder.Eq(val)),
		)
	}

	rows, err := #session.For#(ctx, m).Adaptor().Query(
		ctx,
		#builder.Select#(nil).From(
			T#T#,
			#builder.Where#(builder.And(conds...)),
			#builder.Limit#(1),
			#builder.Comment#(#FetchComment#),
		),
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return #helper.Scan#(ctx, rows, m)
}
