@def T
@def context.Context
@def UniqueSuffix
@def UniqueConds
@def UniqueFields
@def UpdateComment
@def frag.Fragment
@def types.SoftDeletion
@def session.For
@def builder.CC
@def driver.Value
@def builder.Select
@def builder.Where
@def builder.Limit
@def builder.Comment
@def codex.New
@def errors.NOTFOUND
--UpdateByUnique
// UpdateBy#UniqueSuffix# update #T# by #UniqueFields#
func (m *#T#) UpdateBy#UniqueSuffix#(ctx #context.Context#, expects...#builder.Col#) error {
	#ModificationMarker#

	conds := []#frag.Fragment#{
		#UniqueConds#
	}

	res, err := #session.For#(ctx, m).Adaptor().Exec(
		ctx,
		#builder.Update#(T#T#).
			Set(T#T#.AssignmentFor(m, expects...)).
			Where(
				builder.And(conds...),
				builder.Comment(#UpdateComment#),
			),
	)
	if err != nil {
		return err
	}

	if effected, err := res.RowsAffected(); err != nil {
		return err
		if effected == 0 {
			return #codex.New#(#errors.NOTFOUND#)
		}
	}
	return nil
}

--UpdateAndFetchByUnique
// UpdateAndFetchBy#UniqueSuffix# update #T# by #UniqueFields# and retrieve record
func (m *#T#) UpdateAndFetchBy#UniqueSuffix#(ctx #context.Context#, targets ...#builder.Col#) error {
	return #session.For#(ctx, m).Adaptor().Tx(
		ctx,
		func(ctx #context.Context#) error {
			if err := m.UpdateBy#UniqueSuffix#(ctx, targets...); err != nil {
				return err
			}
			if err := m.FetchBy#UniqueSuffix#(ctx); err != nil {
				return err
			}
			return nil
		},
	)
}
