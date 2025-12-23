@def T
@def context.Context
@def builder.SqlCondition
@def builder.Additions
@def builder.Col
@def frag.Fragment
@def builder.ColsOf
@def SoftDeletionCondition
@def builder.Where
@def builder.And
@def builder.Comment
@def ListComment
@def session.For
@def builder.Select
@def helper.Scan
--List
// List fetch #T# datalist with condition and additions
func (m *#T#) List(ctx #context.Context#, cond #builder.SqlCondition#, adds #builder.Additions#, expects ...#builder.Col#) ([]#T#, error) {
	cols := #frag.Fragment#(nil)
	if len(expects) > 0 {
		cols = #builder.ColsOf#(expects...)
	}
	conds := []#frag.Fragment#{cond}
	#SoftDeletionCondition#

	adds = append(
		adds,
		#builder.Where#(#builder.And#(conds...)),
		#builder.Comment#(#ListComment#),
	)

	rows, err := #session.For#(ctx, m).Adaptor().Query(
		ctx,
		#builder.Select#(cols).From(T#T#, adds...),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := new([]#T#)
	if err = #helper.Scan#(ctx, rows, res); err != nil {
		return nil, err
	}
	return *res, nil
}
