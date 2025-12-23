@def T
@def context.Context
@def UniqueSuffix
@def UniqueConds
@def UniqueFields
@def DeleteComment
@def frag.Fragment
@def types.SoftDeletion
@def session.For
@def builder.CC
@def driver.Value
@def builder.Delete
@def builder.Where
@def builder.Comment
--DeleteByUnique
// DeleteBy#UniqueSuffix# delete #T# recode by #UniqueFields#
func (m *#T#) DeleteBy#UniqueSuffix#(ctx #context.Context#) error {
	conds := []#frag.Fragment#{
		#UniqueConds#
	}

	_, err := #session.For#(ctx, m).Adaptor().Exec(
		ctx,
		#builder.Delete#().From(
			T#T#,
			#builder.Where#(#builder.And#(conds...)),
			#builder.Comment#(#DeleteComment#),
		),
	)
	return err
}
