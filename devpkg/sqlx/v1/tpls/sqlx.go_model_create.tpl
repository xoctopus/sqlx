@def T
@def context.Context
@def helper.CVsForInsertion
@def session.MustFor
@def builder.Insert
@def builder.Comment
--Create
// Create inserts #T# to database
func (m *#T#) Create(ctx #context.Context#) error {
	#CreationMarker#

	cols, values := #helper.CVsForInsertion#(m)
	_, err := #session.MustFor#(ctx, m).Adaptor().Exec(
		ctx,
		#builder.Insert#().Into(
			T#T#,
			#builder.Comment#(#CreateComment#),
		).Values(cols, values...),
	)
	return err
}

