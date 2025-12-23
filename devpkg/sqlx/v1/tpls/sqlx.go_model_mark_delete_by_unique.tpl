@def T
@def context.Context
@def UniqueSuffix
@def UniqueConds
@def session.For
@def MarkDeletionComment
@def frag.Fragment
@def builder.Update
@def builder.And
@def builder.CC
@def builder.Neq
@def driver.Value
@def DeletionMarker
--MarkDeletionByUnique
// MarkDeletionBy#UniqueSuffix# marks #T# as deleted
func (m *#T#) MarkDeletionBy#UniqueSuffix#(ctx #context.Context#) error {
	#DeletionMarker#

	deletion, modifications, v := m.SoftDeletion()
	cols := []#builder.Col#{T#T#.C(deletion)}
	for _, f := range modifications {
		cols = append(cols, T#T#.C(f))
	}

	conds := []#frag.Fragment#{
		#UniqueConds#
		#builder.CC#[#driver.Value#](T#T#.C(deletion)).AsCond(#builder.Neq#(v)),
	}

	_, err := #session.For#(ctx, m).Adaptor().Exec(
		ctx,
		#builder.Update#(T#T#).
			Set(T#T#.AssignmentFor(m, cols...)).
			Where(
				#builder.And#(conds...),
				#builder.Comment#(#MarkDeletionComment#),
			),
	)
	return err
}
